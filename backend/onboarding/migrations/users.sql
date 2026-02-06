CREATE TABLE users (
    keycloak_user_id    VARCHAR(255) NOT NULL PRIMARY KEY, -- Unique identifier for the user
    username            VARCHAR(255) NOT NULL UNIQUE,      -- Username must be unique
    email               VARCHAR(255) NOT NULL UNIQUE,      -- Email must be unique
    role                VARCHAR(50) NOT NULL,             -- Role of the user (e.g., AGENT, CITIZEN, etc.)
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,    -- Indicates if the user is active
    preferred_language  VARCHAR(50),                      -- Preferred language of the user (nullable)
    created_at          TIMESTAMP NOT NULL DEFAULT NOW(), -- Timestamp when the record was created
    updated_at          TIMESTAMP NOT NULL DEFAULT NOW(), -- Timestamp when the record was last updated
    created_by          VARCHAR(255),                     -- User who created the record (nullable)
    updated_by          VARCHAR(255),                     -- User who last updated the record (nullable)
    start_date          DATE,                             -- start date of the user
    end_date            DATE                              -- end_date of the user
);