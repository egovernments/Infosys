package services

import (
	"context"
	"encoding/json"
	workflow "enumeration/internal/clients"
	"enumeration/internal/config"
	"enumeration/internal/constants"
	"enumeration/internal/dto"
	"enumeration/internal/models"
	"enumeration/internal/repositories"
	"enumeration/pkg/logger"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ ApplicationService = (*applicationService)(nil)

// applicationService handles business logic for applications
type applicationService struct {
	appRepo        repositories.ApplicationRepository
	ownerRepo      repositories.PropertyOwnerRepository
	workflowClient *workflow.Client
	cfg            *config.Config
	applicationlog ApplicationLogService
}

// NewApplicationService creates a new application service
func NewApplicationService(appRepo repositories.ApplicationRepository, ownerRepo repositories.PropertyOwnerRepository, workflowClient *workflow.Client, cfg *config.Config, applicationlog ApplicationLogService) ApplicationService {
	return &applicationService{
		appRepo:        appRepo,
		ownerRepo:      ownerRepo,
		workflowClient: workflowClient,
		cfg:            cfg,
		applicationlog: applicationlog,
	}
}
// Create creates a new application with property owners
func (s *applicationService) Create(ctx context.Context, tenantID string, citizenID string, req *dto.CreateApplicationRequest) (*models.Application, error) {
	
	if req == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrValidation)
	}
	// Generate unique application number
	applicationNo := s.generateApplicationNumber()

	process, err := s.workflowClient.GetProcess(ctx, tenantID, "PROP_ENUM")
	if err != nil {

		return nil, fmt.Errorf("PROP_ENUM workflow process not found - please ensure workflow is pre-created: %w", err)
	}

	noOfDays, err := workflow.FetchNoOfDays(s.cfg.MDMSURL, tenantID)
	logger.Info("No of days fetched from MDMS: ", noOfDays)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch noOfDays from MDMS: %w", err)
	}
	req.DueDate = time.Now().AddDate(0, 0, noOfDays)
	req.Priority = constants.PriorityLow

	// Validate AssesseeID exists in users table
	if err := s.appRepo.CheckUserExists(ctx, req.AssesseeID.String()); err != nil {
		return nil, err // Error already has proper context from repository
	}

	// Create application entity
	application := &models.Application{

		ApplicationNo:      applicationNo,
		PropertyID:         req.PropertyID,
		Priority:           req.Priority,
		IsDraft:            req.IsDraft,
		DueDate:            req.DueDate,
		Status:             constants.StatusInitiated,
		TenantID:           tenantID,
		AppliedBy:          req.AppliedBy,
		AssesseeID:         req.AssesseeID.String(),
		WorkflowInstanceID: process.ID,
	}

	err = s.appRepo.Create(ctx, application)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrApplicationCreationFailed+": %w", err)
	} else {
		logger.Info("Application created with ID: ", application.ID)
	}

	// Load the complete application with Property data
	fullApplication, err := s.appRepo.GetByID(ctx, application.ID) // Add ctx parameter
	if err != nil {
		return nil, fmt.Errorf("failed to load created application: %w", err)
	}
	username := ctx.Value("user").(string)
	appLog := dto.CreateApplicationLogRequest{
		ApplicationID: fullApplication.ID,
		Action:        "Application Created",
		Actor:         ctx.Value("role").(string),
		PerformedBy:   username,
		Metadata:      map[string]interface{}{},
		Comments:      `Application applied by ` + fullApplication.AppliedBy,
	}
	s.applicationlog.Create(ctx, &appLog)

	return fullApplication, nil
}

// GetByID retrieves application by ID
func (s *applicationService) GetByID(ctx context.Context, id uuid.UUID) (*models.Application, error) {
	application, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.ErrApplicationNotFound+":%w", err)
		}
		return nil, err
	}
	logger.Info("Fetched application with ID: ", application.ID)
	return application, nil
}

// GetByApplicationNo retrieves application by application number
func (s *applicationService) GetByApplicationNo(ctx context.Context, applicationNo string) (*models.Application, error) {
	application, err := s.appRepo.GetByApplicationNo(ctx, applicationNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.ErrApplicationNotFound)
		}
		return nil, err
	}
	logger.Info("Fetched application with Application No: ", application.ApplicationNo)
	return application, nil
}

// Update updates an existing application
func (s *applicationService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateApplicationRequest) (*models.Application, error) {
	// Check if application exists
	if req == nil {
		return nil, fmt.Errorf("%w: request is nil", ErrValidation)
	}
	existing, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf(constants.ErrApplicationNotFound)
		}
		return nil, err
	}

	// Update fields
	if req.Priority != "" {
		existing.Priority = req.Priority
	}
	if req.DueDate != nil {
		existing.DueDate = *req.DueDate
	}
	if req.AssignedAgent != nil {
		existing.AssignedAgent = req.AssignedAgent
	}
	if req.AppliedBy != "" {
		existing.AppliedBy = req.AppliedBy
	}
	if req.IsDraft != nil {
		existing.IsDraft = *req.IsDraft
	}
	 if req.ImportantNote != nil {
        existing.ImportantNote = *req.ImportantNote
    }

	// Update application
	if err := s.appRepo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}
	logger.Info("Updated application with ID: ", existing.ID)

	return existing, nil
}

// UpdateStatus updates application status
func (s *applicationService) UpdateStatus(ctx context.Context, id uuid.UUID, req *dto.UpdateStatusRequest) error {
	// Check if application exists
	if req == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	existing, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(constants.ErrApplicationNotFound)
		}
		return err
	}

	// Update status
	existing.Status = req.Status

	if err := s.appRepo.Update(ctx, existing); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// AssignAgent assigns an agent to application
func (s *applicationService) AssignAgent(ctx context.Context, id uuid.UUID, req *dto.ActionRequest, tenantID string) error {
	// Check if application exists
	if req == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}
	existing, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(constants.ErrApplicationNotFound)
		}
		return err
	}
	logger.Info("Fetched application with ID: ", existing.ID)

	err = s.appRepo.CheckUserExists(ctx, req.AgentID.String())
	if err != nil {
		return err
	}
	// Check if workflow process exists
	if existing.WorkflowInstanceID == "" {
		return fmt.Errorf("workflow process not initialized for application %s", id)
	}
	// Assign agent
	if req.AgentID == uuid.Nil {
		return fmt.Errorf("agent ID is required for assignment")
	}
	existing.AssignedAgent = &req.AgentID

	existing.Status = constants.StatusAssigned
	if err := s.appRepo.Update(ctx, existing); err != nil {
		return fmt.Errorf("failed to assign agent: %w", err)
	}
	logger.Info("Assigned agent with ID: ", req.AgentID, " to application ID: ", existing.ID)
	// Skip workflow transition creation for re-assign action
	if req.Action == "re-assign" {
		username, _ := workflow.GetUsername(req.AgentID.String())
		appLog := dto.CreateApplicationLogRequest{
			ApplicationID: existing.ID,
			Action:        "Application Re-assigned to Agent",
			Actor: ctx.Value("role").(string),
			PerformedBy:   ctx.Value("user").(string),
			Metadata:      map[string]interface{}{},
			Comments:      "application Re-assigned to agent : " + username,
		}
		s.applicationlog.Create(ctx, &appLog)
		return nil
	}
	// Create workflow transition
	transitionReq := &workflow.TransitionRequest{
		ProcessID: existing.WorkflowInstanceID,
		EntityID:  existing.ApplicationNo,
		Action:    "Assign to Agent",
		Attributes: map[string]interface{}{
			"roles":         []string{constants.RoleServiceManager},
			"jurisdiction":  []string{"Punjab.Amritsar"},
			"assignedAgent": []string{req.AgentID.String()},
			"comments":      []string{req.Comments},
		},
	}
	_, err = s.workflowClient.CreateTransition(ctx, "pb.amritsar", transitionReq)
	if err != nil {
		return fmt.Errorf("failed to create workflow transition: %w", err)
	}
	username, _ := workflow.GetUsername(req.AgentID.String())
	appLog := dto.CreateApplicationLogRequest{
		ApplicationID: existing.ID,
		Action:        "Application Assigned to Agent",
		PerformedBy:   ctx.Value("user").(string),
		Actor:         ctx.Value("role").(string),
		Metadata:      map[string]interface{}{},
		Comments:      "application assigned to agent : " + username,
	}
	s.applicationlog.Create(ctx, &appLog)

	return nil
}

// Delete deletes application by ID
func (s *applicationService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if application exists
	_, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf(constants.ErrApplicationNotFound)
		}
		return err
	}
	logger.Info("Deleting application with ID: ", id)
	return s.appRepo.Delete(ctx, id)
}

// List retrieves all applications with pagination
func (s *applicationService) List(ctx context.Context, UserID string, tenantID string, Role string, verify string, page, size int) ([]models.Application, int64, error) {
	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 {
		size = constants.DefaultSize
	}

	var applications []models.Application
	var total int64
	var err error

	switch Role {
	case constants.RoleServiceManager:

		applications, total, err = s.appRepo.GetByTenantIDAndStatus(ctx, tenantID, verify, page, size)

	case constants.RoleAgent:
		if verify == "" {
			logger.Error("Agent List: verify parameter is missing")
			return []models.Application{}, 0, errors.New("verify parameter is missing")
		}
		applications, total, err = s.appRepo.GetByAssignedAgent(ctx, UserID, verify, page, size)

	case constants.RoleCommissioner:

		applications, total, err = s.appRepo.GetByTenantIDAndStatus(ctx, tenantID, verify, page, size)
	default:
		logger.Error("Invalid role provided")
		return []models.Application{}, 0, errors.New("invalid role provided")

	}
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch applications: %w", err)
	}
	// currentTime := time.Now()
	// for i := range applications {
	// 	createdTime := applications[i].CreatedAt
	// 	dueDate := applications[i].DueDate

	// 	// Skip priority calc when dates are missing or invalid to avoid divide-by-zero
	// 	if createdTime.IsZero() || dueDate.IsZero() {
	// 		// leave priority as-is (or set default) and continue
	// 		continue
	// 	}

	// 	// Calculate progress

	// 	progress := s.CalculateProgress(ctx, currentTime, createdTime, dueDate)
	// 	logger.Info("Calculated progress for application ID: ", applications[i].ID, " Progress: ", progress)

	// 	// Assign priority based on progress
	// 	switch {
	// 	case progress < 0.5:
	// 		applications[i].Priority = constants.PriorityLow
	// 		if err := s.appRepo.Update(ctx, &applications[i]); err != nil {
	// 			logger.Error("failed to update application priority for %s: %v", applications[i].ID, err)
	// 		}
	// 	case progress >= 0.5 && progress < 0.75:
	// 		applications[i].Priority = constants.PriorityMedium
	// 		if err := s.appRepo.Update(ctx, &applications[i]); err != nil {
	// 			logger.Error("failed to update application priority for %s: %v", applications[i].ID, err)
	// 		}
	// 	default:
	// 		applications[i].Priority = constants.PriorityHigh
	// 		if err := s.appRepo.Update(ctx, &applications[i]); err != nil {
	// 			logger.Error("failed to update application priority for %s: %v", applications[i].ID, err)
	// 		}
	// 	}
	// }

	return applications, total, err
}

func (s *applicationService) CalculateProgress(ctx context.Context, currentTime time.Time, createdTime time.Time, dueDate time.Time) float64 {

	totalDuration := dueDate.Sub(createdTime).Hours()
	elapsedDuration := currentTime.Sub(createdTime).Hours()
	if totalDuration <= 0 {
		return 1.0 // If due date is before or same as created date, consider progress as complete
	}
	progress := elapsedDuration / totalDuration
	return progress

}
func (s *applicationService) ApproveApplication(ctx context.Context, tenantID, commissionerID, applicationID string, req *dto.ActionRequest) error {

	if req == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}
	appID, err := uuid.Parse(applicationID)
	if err != nil {
		return fmt.Errorf("invalid application ID format: %w", err)
	}
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	transitionReq := &workflow.TransitionRequest{
		ProcessID: app.WorkflowInstanceID,
		EntityID:  app.ApplicationNo,
		Action:    "Approve Application",
		Attributes: map[string]interface{}{
			"roles":            []string{constants.RoleCommissioner},
			"jurisdiction":     []string{"Punjab.Amritsar"},
			"approvalComments": []string{req.Comments},
		},
	}

	_, err = s.workflowClient.CreateTransition(ctx, "pb.amritsar", transitionReq)

	if err != nil {
		return fmt.Errorf("failed to create workflow transition: %w", err)
	}

	newState := constants.StatusApproved
	if !req.Approved {
		newState = constants.StatusRejected
	}

	app.Status = newState

	if err := s.appRepo.Update(ctx, app); err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}
	logger.Info("Updated application with ID: ", app.ID, " to status: ", newState)
	username, _ := workflow.GetUsername(commissionerID)
	appLog := dto.CreateApplicationLogRequest{
		ApplicationID: app.ID,
		Action:        "Application Approved",
		PerformedBy:   ctx.Value("user").(string),
		Actor: ctx.Value("role").(string),
		Metadata:      map[string]interface{}{},
		Comments:      "application " + newState + " by commissioner : " + username,
	}
	s.applicationlog.Create(ctx, &appLog)
	return nil
}

func (s *applicationService) VerifyApplication(ctx context.Context, tenantID, serviceManagerID, applicationID string, req *dto.ActionRequest) error {
	if req == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}
	appID, err := uuid.Parse(applicationID)
	if err != nil {
		return fmt.Errorf("invalid application ID format: %w", err)
	}
	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	// Check if workflow process exists
	if app.WorkflowInstanceID == "" {
		return fmt.Errorf("workflow process not initialized for application %s", applicationID)
	}

	newState := constants.StatusAuditVerified

	transitionReq := &workflow.TransitionRequest{
		ProcessID: app.WorkflowInstanceID,
		EntityID:  app.ApplicationNo,
		Action:    "Audit Verification",
		Attributes: map[string]interface{}{
			"roles":         []string{constants.RoleServiceManager},
			"jurisdiction":  []string{"Punjab.Amritsar"},
			"auditComments": []string{req.Comments},
		},
	}

	_, err = s.workflowClient.CreateTransition(ctx, "pb.amritsar", transitionReq)
	if err != nil {
		return fmt.Errorf("failed to create workflow transition: %w", err)
	}
	app.Status = newState

	if err := s.appRepo.Update(ctx, app); err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}
	logger.Info("Updated application with ID: ", app.ID, " to status: ", newState)
	username, _ := workflow.GetUsername(serviceManagerID)
	appLog := dto.CreateApplicationLogRequest{
		ApplicationID: app.ID,
		Action:        "Application Audit Verified",
		PerformedBy:   ctx.Value("user").(string),
		Metadata:      map[string]interface{}{},
		Actor: ctx.Value("role").(string),
		Comments:      "application " + newState + " by service manager : " + username,
	}
	s.applicationlog.Create(ctx, &appLog)
	return nil
}

func (s *applicationService) VerifyApplicationByAgent(ctx context.Context, tenantID, agentID, applicationID string, req *dto.ActionRequest) error {
	if req == nil {
		return fmt.Errorf("%w: request is nil", ErrValidation)
	}

	appID, err := uuid.Parse(applicationID)
	if err != nil {
		return fmt.Errorf("invalid application ID format: %w", err)
	}

	app, err := s.appRepo.GetByID(ctx, appID)
	if err != nil {
		return fmt.Errorf("application not found: %w", err)
	}

	transitionReq := &workflow.TransitionRequest{
		ProcessID: app.WorkflowInstanceID,
		EntityID:  app.ApplicationNo,
		Action:    "Submit for Verification",
		Attributes: map[string]interface{}{
			"roles":              []string{constants.RoleAgent},
			"jurisdiction":       []string{"Punjab.Amritsar"},
			"verificationReport": []string{req.Comments},
		},
	}

	_, err = s.workflowClient.CreateTransition(ctx, "pb.amritsar", transitionReq)
	if err != nil {
		return fmt.Errorf("failed to create workflow transition: %w", err)
	}

	app.Status = constants.StatusVerified

	if err := s.appRepo.Update(ctx, app); err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}
	logger.Info("Updated application with ID: ", app.ID, " to status: ", constants.StatusVerified)
	username, _ := workflow.GetUsername(agentID)
	appLog := dto.CreateApplicationLogRequest{
		ApplicationID: app.ID,
		Action:        "Application Verified",
		PerformedBy:   ctx.Value("user").(string),
		Metadata:      map[string]interface{}{},
		Actor: ctx.Value("role").(string),
		Comments:      "application " + constants.StatusVerified + " by agent : " + username,
	}
	s.applicationlog.Create(ctx, &appLog)
	return nil
}

// Search retrieves applications based on search criteria
func (s *applicationService) Search(ctx context.Context, criteria *dto.ApplicationSearchCriteria, page, size int) ([]*models.Application, int64, error) {
	if criteria == nil {
		// If your repo accepts nil criteria, you can allow it; otherwise return validation
		return nil, 0, fmt.Errorf("%w: search criteria is nil", ErrValidation)
	}

	if page < 0 {
		page = constants.DefaultPage
	}
	if size <= 0 {
		size = constants.DefaultSize
	}
	logger.Info("Searching applications with criteria: ", criteria)
	return s.appRepo.Search(ctx, criteria, page, size)
}

// generateApplicationNo generates a unique application number
func (s *applicationService) generateApplicationNumber() string {
	// Generate application number by calling external ID generation service
	reqBody := `{
		"templateId": "applId",
		"variables": {
			"ORG": "APPL"
		}
	}`

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(s.cfg.IDGenURL, "application/json", strings.NewReader(reqBody))
	if err != nil {
		// fallback to local generation if service fails
		return fmt.Sprintf("PROP-%d-%s", time.Now().Unix(), uuid.New().String()[:8])
	}
	defer resp.Body.Close()

	var result struct {
		ID string `json:"id"`
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("PROP-%d-%s", time.Now().Unix(), uuid.New().String()[:8])
	}
	if err := json.Unmarshal(body, &result); err != nil || result.ID == "" {
		return fmt.Sprintf("PROP-%d-%s", time.Now().Unix(), uuid.New().String()[:8])
	}
	return result.ID
}
