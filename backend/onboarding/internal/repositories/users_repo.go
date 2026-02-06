package repositories

import (
	"context"
	"fmt"
	"property-tax-onboarding/internal/constants"
	"property-tax-onboarding/internal/errors"
	"property-tax-onboarding/internal/models"
	"property-tax-onboarding/pkg/logger"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// PostgreSQLUserRepository implements UserRepository interface for unified users table
type PostgreSQLUserRepository struct {
	dbd *gorm.DB
}

// NewPostgreSQLUserRepository creates a new PostgreSQL user repository
func NewPostgreSQLUserRepository(dbd *gorm.DB) UserRepository {
	return &PostgreSQLUserRepository{dbd: dbd}
}

// GetUserProfileByUserID retrieves a user profile by user ID.
//
// Parameters:
//   - userID: The unique identifier of the user whose profile is to be retrieved.
//
// Returns:
//   - A pointer to the UserProfile object if found.
//   - An error if the operation fails or the user profile is not found.
//
// Description:
// This method queries the database to retrieve the user profile associated with the given user ID.
// If the user profile is not found, it returns an error. The method uses centralized error handling
// and logs the operation at various stages for traceability.
func (r *PostgreSQLUserRepository) GetUserProfileByUserID(ctx context.Context, userID string) (*models.UserProfile, error) {
	logger.Info(constants.LogRetrieveUserProfileStart, "userID", userID)

	// Create an empty UserProfile object to hold the result
	var profile models.UserProfile

	if err := r.dbd.WithContext(ctx).Where("user_profile_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error(constants.ErrUserProfileNotFound, "userID", userID)
			return nil, errors.NewRepositoryError(fmt.Sprintf("%s: %s", constants.ErrUserProfileNotFound, userID))
		}
		logger.Error(constants.ErrRetrieveUserProfileFailed, "error", err, "userID", userID)
		return nil, errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrRetrieveUserProfileFailed, err))
	}

	logger.Info(constants.LogRetrieveUserProfileSuccess, "userID", userID)
	return &profile, nil
}

// GetByKeycloakUserID retrieves a user by Keycloak user ID.
//
// Parameters:
//   - keycloakUserID: The unique identifier of the user in Keycloak.
//
// Returns:
//   - A pointer to the User object if found.
//   - An error if the operation fails or the user is not found.
//
// Description:
// This method queries the database to retrieve the user associated with the given Keycloak user ID.
// If the user is not found, it returns an error. The method uses centralized error handling
// and logs the operation at various stages for traceability.
func (r *PostgreSQLUserRepository) GetByKeycloakUserID(ctx context.Context, keycloakUserID string) (*models.User, error) {
	logger.Info(constants.LogRetrieveUserStart, "keycloak_user_id", keycloakUserID)

	// Create an empty User object to hold the result
	var user models.User

	// Use GORM to find the user by Keycloak user ID
	if err := r.dbd.WithContext(ctx).Where("keycloak_user_id = ?", keycloakUserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error(constants.ErrUserNotFound, "keycloak_user_id", keycloakUserID)
			return nil, errors.NewRepositoryError(fmt.Sprintf("%s: %s", constants.ErrUserNotFound, keycloakUserID))
		}
		logger.Error(constants.ErrRetrieveUserFailed, "error", err, "keycloak_user_id", keycloakUserID)
		return nil, errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrRetrieveUserFailed, err))
	}

	logger.Info(constants.LogRetrieveUserSuccess, "keycloak_user_id", keycloakUserID)
	return &user, nil
}

// GetByRole retrieves all users by role with pagination.
//
// Parameters:
//   - role: The role of the users to retrieve.
//   - limit: The maximum number of users to retrieve.
//   - offset: The starting point for pagination.
//
// Returns:
//   - A slice of pointers to User objects if found.
//   - An error if the operation fails or no users are found.
//
// Description:
// This method queries the database to retrieve users associated with the given role.
// It supports pagination using the limit and offset parameters. If no users are found,
// or if an error occurs during the query, it returns an appropriate error. The method
// uses centralized error handling and logs the operation at various stages for traceability.
func (r *PostgreSQLUserRepository) GetByRole(ctx context.Context, role models.UserRole, limit, offset int) ([]*models.User, error) {
	logger.Info(constants.LogGetUsersByRoleStart, "role", role)

	var users []*models.User
	// Use GORM to find users by role with limit and offset
	if err := r.dbd.WithContext(ctx).Where("role = ? AND deleted = false", role).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		logger.Error(constants.ErrGetUsersByRole, "error", err, "role", role)
		return nil, errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrGetUsersByRole, err))
	}
	// Check if no users were found
	if len(users) == 0 {
		logger.Warn(constants.ErrNoUsersFound, "role", role)
		return nil, errors.NewRepositoryError(constants.ErrNoUsersFound)
	}

	logger.Info(constants.LogGetUsersByRoleSuccess, "role", role, "returnedCount", len(users))
	return users, nil
}

// Delete removes a user and related data from users, userprofile, and addresses tables.
//
// Parameters:
//   - keycloakUserID: The unique identifier of the user in Keycloak.
//
// Returns:
//   - An error if the operation fails during any step of the deletion process.
//
// Description:
// This method deletes a user and their related data from the database. It first checks if the user exists.
// If the user exists, it uses a GORM transaction to ensure atomicity while deleting the user profile,
// associated address, and the user record. If any step fails, the transaction is rolled back, and an
// appropriate error is returned. The method uses centralized error handling and logs the operation at
// various stages for traceability.
func (r *PostgreSQLUserRepository) Delete(ctx context.Context, keycloakUserID string) error {
	logger.Info(constants.LogDeleteUserStart, "keycloak_user_id", keycloakUserID)

	// First check if user exists
	var user models.User
	if err := r.dbd.Where("keycloak_user_id = ?", keycloakUserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error(constants.ErrUserNotFound, "keycloak_user_id", keycloakUserID)
			return errors.NewRepositoryError(fmt.Sprintf("%s: %s", constants.ErrUserNotFound, keycloakUserID))
		}
		logger.Error(constants.ErrCheckUserExistenceFailed, "error", err, "keycloak_user_id", keycloakUserID)
		return errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrCheckUserExistenceFailed, err))
	}

	// Check if user is already inactive
	if user.Deleted {
		logger.Warn(constants.LogUserAlreadyDeleted, "keycloak_user_id", keycloakUserID)
		return errors.NewRepositoryError(fmt.Sprintf("User is already deleted: %s", keycloakUserID))
	}

	// Perform soft delete by setting deleted to true
	logger.Info(constants.LogSoftDeleteUserStart, "keycloak_user_id", keycloakUserID)
	if err := r.dbd.Model(&models.User{}).
		Where("keycloak_user_id = ?", keycloakUserID).
		Update("deleted", true).Error; err != nil {
		logger.Error(constants.ErrSoftDeleteUserFailed, "error", err, "keycloak_user_id", keycloakUserID)
		return errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrSoftDeleteUserFailed, err))
	}

	logger.Info(constants.LogSoftDeleteUserSuccess, "keycloak_user_id", keycloakUserID)
	return nil
}

// GetAddressByID retrieves an address by its ID.
//
// Parameters:
//   - addressID: The unique identifier of the address to retrieve.
//
// Returns:
//   - A pointer to the AddressDB object if found.
//   - An error if the operation fails or the address is not found.
//
// Description:
// This method queries the database to retrieve the address associated with the given address ID.
// If the address is not found, it returns an error. The method uses centralized error handling
// and logs the operation at various stages for traceability.
func (r *PostgreSQLUserRepository) GetAddressByID(ctx context.Context, addressID string) (*models.AddressDB, error) {
	logger.Info(constants.LogRetrieveAddressStart, "addressID", addressID)

	var address models.AddressDB
	// Use GORM to find the address by ID
	if err := r.dbd.WithContext(ctx).Where("id = ?", addressID).First(&address).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Error(constants.ErrAddressNotFound, "addressID", addressID)
			return nil, errors.NewRepositoryError(fmt.Sprintf("%s: %s", constants.ErrAddressNotFound, addressID))
		}
		logger.Error(constants.ErrRetrieveAddressFailed, "error", err, "addressID", addressID)
		return nil, errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrRetrieveAddressFailed, err))
	}

	logger.Info(constants.LogRetrieveAddressSuccess, "addressID", addressID)
	return &address, nil
}

// Count returns the total number of users in the database.
//
// Returns:
//   - The total number of users as an int64.
//   - An error if the operation fails.
//
// Description:
// This method uses GORM to count the total number of users in the database. If an error occurs
// during the counting process, it returns an appropriate error message.
func (r *PostgreSQLUserRepository) Count() (int64, error) {
	logger.Info(constants.LogCountUsersStart)

	var count int64
	// Use GORM to count the number of users
	if err := r.dbd.Model(&models.User{}).Count(&count).Error; err != nil {
		logger.Error(constants.ErrCountUsersFailed, "error", err)
		return 0, errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrCountUsersFailed, err))
	}

	logger.Info(constants.LogCountUsersSuccess, "count", count)
	return count, nil
}

// UpdateUser updates user basic information including preferred language.
//
// Parameters:
//   - keycloakUserID: The unique identifier of the user in Keycloak.
//   - updateReq: A pointer to the UpdateUserRequest object containing the updated user information.
//
// Returns:
//   - An error if the operation fails.
//
// Description:
// This method updates the basic information of a user in the database, including their preferred language.
// If the update operation fails, it returns an appropriate error message. The method uses centralized
// error handling and logs the operation at various stages for traceability.
func (r *PostgreSQLUserRepository) UpdateUser(ctx context.Context, keycloakUserID string, updateReq *models.UpdateUserRequest) error {
	logger.Info(constants.LogUpdateUserStartt, "keycloak_user_id", keycloakUserID)

	// Use GORM to update the user
	if err := r.dbd.WithContext(ctx).Model(&models.User{}).Where("keycloak_user_id = ? AND is_active= true", keycloakUserID).Updates(models.User{
		Email:             updateReq.Email,
		Role:              models.UserRole(updateReq.Role),
		IsActive:          updateReq.IsActive,
		PreferredLanguage: updateReq.PreferredLanguage,
		EndDate:           updateReq.EndDate,
	}).Error; err != nil {
		logger.Error(constants.ErrUpdateUserFailed, "error", err, "keycloak_user_id", keycloakUserID)
		return errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrUpdateUserFailed, err))
	}

	logger.Info(constants.LogUpdateUserSuccess, "keycloak_user_id", keycloakUserID)
	return nil
}

// UpdateUserProfileComplete updates complete user profile.
//
// Parameters:
//   - keycloakUserID: The unique identifier of the user in Keycloak.
//   - profile: A pointer to the UpdateUserProfile object containing the updated profile information.
//
// Returns:
//   - An error if the operation fails.
//
// Description:
// This method updates the complete user profile in the database. If the user profile does not have an
// associated address, it creates a new address record and updates the profile with the new AddressID.
// If an address already exists, it updates the existing address record. The method uses centralized
// error handling and logs the operation at various stages for traceability.
func (r *PostgreSQLUserRepository) UpdateUserProfileComplete(ctx context.Context, keycloakUserID string, profile *models.UpdateUserProfile) error {
	logger.Info(constants.LogUpdateUserProfileStart, "keycloak_user_id", keycloakUserID)

	// Update the user profile
	if err := r.dbd.WithContext(ctx).Model(&models.UserProfile{}).Where("user_profile_id = ?", keycloakUserID).Updates(profile).Error; err != nil {
		logger.Error(constants.ErrUpdateUserProfileFailed, "error", err, "keycloak_user_id", keycloakUserID)
		return errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrUpdateUserProfileFailed, err))
	}

	// Retrieve the AddressID from the user_profiles table
	var userProfile models.UserProfile
	if err := r.dbd.WithContext(ctx).Select("address_id").Where("user_profile_id = ?", keycloakUserID).First(&userProfile).Error; err != nil {
		logger.Error(constants.ErrRetrieveAddressIDFailed, "error", err, "keycloak_user_id", keycloakUserID)
		return errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrRetrieveAddressIDFailed, err))
	}

	// If AddressID is NULL, create a new address record
	if userProfile.AddressID == nil && profile.Address != nil {
		logger.Info(constants.LogCreateNewAddressStart, "keycloak_user_id", keycloakUserID)

		// Generate a new UUID for the address ID
		newAddress := profile.Address
		newAddress.ID = uuid.New().String() // Generate a valid UUID

		// Create a new address record
		if err := r.dbd.WithContext(ctx).Create(newAddress).Error; err != nil {
			logger.Error(constants.ErrCreateNewAddressFailed, "error", err, "keycloak_user_id", keycloakUserID)
			return errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrCreateNewAddressFailed, err))
		}

		// Update the user_profiles table with the new AddressID
		if err := r.dbd.WithContext(ctx).Model(&models.UserProfile{}).Where("user_profile_id = ?", keycloakUserID).Update("address_id", newAddress.ID).Error; err != nil {
			logger.Error(constants.ErrUpdateAddressIDFailed, "error", err, "keycloak_user_id", keycloakUserID)
			return errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrUpdateAddressIDFailed, err))
		}

		logger.Info(constants.LogCreateNewAddressSuccess, "address_id", newAddress.ID, "keycloak_user_id", keycloakUserID)
	} else if userProfile.AddressID != nil && profile.Address != nil {

		// If AddressID exists, update the existing address record
		logger.Info(constants.LogUpdateExistingAddressStart, "address_id", *userProfile.AddressID)
		if err := r.dbd.WithContext(ctx).Model(&models.AddressDB{}).Where("id = ?", *userProfile.AddressID).Updates(profile.Address).Error; err != nil {
			logger.Error(constants.ErrUpdateExistingAddressFailed, "error", err, "address_id", *userProfile.AddressID)
			return errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrUpdateExistingAddressFailed, err))
		}
		logger.Info(constants.LogUpdateExistingAddressSuccess, "address_id", *userProfile.AddressID)
	}

	logger.Info(constants.LogUpdateUserProfileSuccess, "keycloak_user_id", keycloakUserID)
	return nil
}

// GetAllUsersWithFilters retrieves users with optional filters and pagination.
//
// Parameters:
//   - filters: A UserFilters object containing optional filters for the query.
//   - limit: The maximum number of users to retrieve.
//   - offset: The starting point for pagination.
//
// Returns:
//   - A slice of pointers to User objects if found.
//   - The total count of users matching the filters.
//   - An error if the operation fails.
//
// Description:
// This method queries the database to retrieve users based on the provided filters. It supports
// pagination using the limit and offset parameters. If an error occurs during the query or counting
// process, it returns an appropriate error message. The method uses centralized error handling and
// logs the operation at various stages for traceability.
func (r *PostgreSQLUserRepository) GetAllUsersWithFilters(ctx context.Context, filters models.UserFilters, limit, offset int) ([]*models.User, int64, error) {
	logger.Info(constants.LogGetUsersWithFiltersStartRetrive)

	var users []*models.User
	var totalCount int64
	query := r.dbd.WithContext(ctx).Model(&models.User{}).Where("deleted = false")

	// Apply filters
	if filters.Role != nil {
		query = query.WithContext(ctx).Where("role IN ?", filters.Role)
	}
	if filters.IsActive != nil {
		query = query.WithContext(ctx).Where("is_active = ?", *filters.IsActive)
	}
	if filters.Username != nil && *filters.Username != "" {
		query = query.WithContext(ctx).Where("username = ?", *filters.Username)
	}
	if filters.Email != nil && *filters.Email != "" {
		query = query.WithContext(ctx).Where("email ILIKE ?", "%"+*filters.Email+"%")
	}

	if filters.Ward != nil {
		query = query.WithContext(ctx).Joins("JOIN \"DIGIT3\".\"zone_mapping\" ON \"DIGIT3\".\"zone_mapping\".user_id = \"DIGIT3\".\"users\".keycloak_user_id").
			Where("\"DIGIT3\".\"zone_mapping\".ward && ?", pq.Array(filters.Ward))
	}

	if filters.PhoneNumber != nil {
		query = query.WithContext(ctx).Joins("JOIN \"DIGIT3\".\"user_profiles\" ON \"DIGIT3\".\"users\".keycloak_user_id = \"DIGIT3\".\"user_profiles\".user_profile_id").
			Where("\"DIGIT3\".\"user_profiles\".phone_number = ?", *filters.PhoneNumber)
	}

	// Count total records
	if err := query.WithContext(ctx).Count(&totalCount).Error; err != nil {
		logger.Error(constants.ErrCountUsersWithFiltersFailed, "error", err)
		return nil, 0, errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrCountUsersWithFiltersFailed, err))
	}

	// Apply pagination and retrieve users
	if err := query.WithContext(ctx).Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		logger.Error(constants.ErrRetrieveUsersWithFiltersFailed, "error", err)
		return nil, 0, errors.NewRepositoryError(fmt.Sprintf("%s: %v", constants.ErrRetrieveUsersWithFiltersFailed, err))
	}
	logger.Info(constants.LogGetUsersWithFiltersSuccess, "totalCount", totalCount, "returnedCount", len(users), "limit", limit, "offset", offset)
	return users, totalCount, nil
}

// UpdateUserIsActive updates the isActive field for a user in the database.
//
// Parameters:
//   - ctx: The context for managing request-scoped values.
//   - userID: The unique identifier of the user in Keycloak.
//   - isActive: A boolean value indicating the new `isActive` status.
//
// Returns:
//   - An error if the operation fails, or nil if the update is successful.
//
// Description:
// This method updates the `is_active` field for a user in the database using the provided `userID`.
// It ensures that the user's active status in the database is consistent with their status in Keycloak.
func (repo *PostgreSQLUserRepository) UpdateUserIsActive(ctx context.Context, userID string, isActive bool) error {
	// Use GORM to update the isActive field
	if err := repo.dbd.WithContext(ctx).
		Model(&models.User{}).
		Where("keycloak_user_id = ?", userID).
		Update("is_active", isActive).Error; err != nil {
		return fmt.Errorf("failed to update isActive for user %s: %w", userID, err)
	}
	return nil
}

// GetUsersByStartDate retrieves users whose `start_date` matches the given date and are currently inactive.
//
// Parameters:
//   - ctx: The context for managing request-scoped values.
//   - date: The date to match against the `start_date` field (format: "YYYY-MM-DD").
//
// Returns:
//   - A slice of `User` objects representing the users to be activated.
//   - An error if the database query fails.
//
// Description:
// This method queries the database to find users whose `start_date` matches the provided date
// and whose `is_active` status is `false`.
func (repo *PostgreSQLUserRepository) GetUsersByStartDate(ctx context.Context, date string) ([]models.User, error) {
	var users []models.User
	err := repo.dbd.WithContext(ctx).
		Where("start_date = ? AND is_active = false", date).
		Find(&users).Error
	return users, err
}

// GetUsersByEndDate retrieves users whose `end_date` is less than or equal to the given date and are currently active.
//
// Parameters:
//   - ctx: The context for managing request-scoped values.
//   - date: The date to match against the `end_date` field (format: "YYYY-MM-DD").
//
// Returns:
//   - A slice of `User` objects representing the users to be deactivated.
//   - An error if the database query fails.
//
// Description:
// This method queries the database to find users whose `end_date` is less than or equal to the provided date
// and whose `is_active` status is `true`.
func (repo *PostgreSQLUserRepository) GetUsersByEndDate(ctx context.Context, date string) ([]models.User, error) {
	var users []models.User
	err := repo.dbd.WithContext(ctx).
		Where("end_date <= ? AND is_active = true", date).
		Find(&users).Error
	return users, err
}

// GetUsersWithFutureStartDate retrieves users whose `start_date` is after the given date and are currently active.
//
// Parameters:
//   - ctx: The context for managing request-scoped values.
//   - currentDate: The current date to compare against the `start_date` field (format: "YYYY-MM-DD").
//
// Returns:
//   - A slice of `User` objects representing the users to be deactivated temporarily.
//   - An error if the database query fails.
//
// Description:
// This method queries the database to find users whose `start_date` is after the provided date
// and whose `is_active` status is `true`. These users are candidates for temporary deactivation.
func (repo *PostgreSQLUserRepository) GetUsersWithFutureStartDate(ctx context.Context, currentDate string) ([]models.User, error) {
	var users []models.User
	err := repo.dbd.WithContext(ctx).
		Where("start_date > ? AND is_active = true", currentDate).
		Find(&users).Error
	return users, err
}

func (repo *PostgreSQLUserRepository) GetUserCounts(ctx context.Context) (int64, int64, int64, error) {
	var totalUsers, activeUsers, fieldAgents int64

	// Count total users (excluding deleted)
	if err := repo.dbd.WithContext(ctx).Model(&models.User{}).Where("deleted = ?", false).Count(&totalUsers).Error; err != nil {
		logger.Error("Failed to count total users", "error", err)
		return 0, 0, 0, err
	}

	// Count active users
	if err := repo.dbd.WithContext(ctx).Model(&models.User{}).Where("deleted = ? AND is_active = ?", false, true).Count(&activeUsers).Error; err != nil {
		logger.Error("Failed to count active users", "error", err)
		return 0, 0, 0, err
	}

	// Count field agents
	if err := repo.dbd.WithContext(ctx).Model(&models.User{}).Where("deleted = ? AND role = ?", false, models.RoleAgent).Count(&fieldAgents).Error; err != nil {
		logger.Error("Failed to count field agents", "error", err)
		return 0, 0, 0, err
	}
	logger.Info("User counts retrieved successfully", "totalUsers", totalUsers, "activeUsers", activeUsers, "fieldAgents", fieldAgents)
	return totalUsers, activeUsers, fieldAgents, nil
}
