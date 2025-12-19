// Package constants holds validation rule constants for the service
package constants

import "time"

// Validation rule constants for input and business logic
const (
	// String Length Limits
	MinNameLength        = 2
	MaxNameLength        = 100
	MinDescriptionLength = 10
	MaxDescriptionLength = 500
	MaxComplexNameLength = 200

	// Time Limits
	MinDueDateDays       = 1
	MaxDueDateDays       = 365
	ApplicationTimeout   = 30 * time.Second
	DatabaseQueryTimeout = 15 * time.Second

	// Business Rules
	MaxPropertiesPerApplication = 1
	MinPropertyOwners           = 1
	MaxPropertyOwners           = 10

	// Coordinate Validation Ranges
	MinLatitude  = -90.0
	MaxLatitude  = 90.0
	MinLongitude = -180.0
	MaxLongitude = 180.0

	// Geographic ranges for India (approximate bounds)
	IndiaMinLatitude  = 8.0
	IndiaMaxLatitude  = 37.0
	IndiaMinLongitude = 68.0
	IndiaMaxLongitude = 97.5

	// Geographic ranges for Karnataka, India (approximate bounds)
	KarnatakaMinLatitude  = 11.5
	KarnatakaMaxLatitude  = 18.5
	KarnatakaMinLongitude = 74.0
	KarnatakaMaxLongitude = 78.5
)
