-- Create the userprofile table
CREATE TABLE user_profiles (
    user_profile_id         VARCHAR(255) NOT NULL,                  -- References keycloak_user_id in the users table
    first_name              VARCHAR(100) NOT NULL,                  -- First name of the user
    last_name               VARCHAR(100) NOT NULL,                  -- Last name of the user
    full_name               VARCHAR(200) NOT NULL,                  -- Full name of the user
    phone_number            VARCHAR(15) NOT NULL,                   -- Phone number of the user
    adhaar_no               BIGINT NOT ,                 -- Aadhaar number (unique constraint)
    gender                  VARCHAR(10) ,                   -- Gender of the user
    guardian                VARCHAR(100),                           -- Guardian's name (nullable)
    guardian_type           VARCHAR(50),                            -- Type of guardian (nullable)
    date_of_birth           DATE,                                   -- Date of birth (nullable)
    department              VARCHAR(100),                           -- Department (nullable)
    designation             VARCHAR(100),                           -- Designation (nullable)
    work_location           VARCHAR(200),                           -- Work location (nullable)
    profile_picture         TEXT,                                   -- Profile picture URL or data (nullable)
    relationship_to_property VARCHAR(50) ,                  -- Relationship to the property
    ownership_share         DOUBLE PRECISION  CHECK (ownership_share >= 0 AND ownership_share <= 100), -- Ownership share (0-100%)
    is_primary_owner        BOOLEAN  DEFAULT FALSE,          -- Indicates if the user is the primary owner
    is_verified             BOOLEAN DEFAULT FALSE,                  -- Indicates if the user is verified
    address_id              UUID,                                   -- Foreign key to the address table (nullable)
    PRIMARY KEY (user_profile_id),                                  -- Primary key on the `user_profile_id` column
    CONSTRAINT userprofile_adhaar_no_key UNIQUE (adhaar_no),        -- Unique constraint on Aadhaar number
    CONSTRAINT userprofile_id_fkey FOREIGN KEY (user_profile_id) REFERENCES users(keycloak_user_id), -- Foreign key to the users table
    CONSTRAINT userprofile_address_id_fkey FOREIGN KEY (address_id) REFERENCES address(id) -- Foreign key to the address table
);