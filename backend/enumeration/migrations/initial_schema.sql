
-- -- Enums

-- DROP SCHEMA "DIGIT3";

CREATE OR REPLACE FUNCTION uuid_generate_v4()
 RETURNS uuid
 LANGUAGE c
 STRICT
AS '$libdir/uuid-ossp', $function$uuid_generate_v4$function$
;
CREATE SCHEMA "DIGIT3" AUTHORIZATION root;

-- DROP TYPE "DIGIT3"."amenity_type_enum";

CREATE TYPE "DIGIT3"."amenity_type_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."digit3_action_enum";

CREATE TYPE "DIGIT3"."digit3_action_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."document_type_enum";

CREATE TYPE "DIGIT3"."document_type_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."gender_enum";

CREATE TYPE "DIGIT3"."gender_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."gis_source_enum";

CREATE TYPE "DIGIT3"."gis_source_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."gis_type_enum";

CREATE TYPE "DIGIT3"."gis_type_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."guardian_type_enum";

CREATE TYPE "DIGIT3"."guardian_type_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."log_action_enum";

CREATE TYPE "DIGIT3"."log_action_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."ownership_type_enum";

CREATE TYPE "DIGIT3"."ownership_type_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."priority_enum";

CREATE TYPE "DIGIT3"."priority_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."property_type_enum";

CREATE TYPE "DIGIT3"."property_type_enum" AS ENUM (
);

-- DROP TYPE "DIGIT3"."relationship_property_enum";

CREATE TYPE "DIGIT3"."relationship_property_enum" AS ENUM (
);
-- "DIGIT3".address definition

-- Drop table

-- DROP TABLE "DIGIT3".address;

CREATE TABLE "DIGIT3".address (
	id uuid DEFAULT "DIGIT3".gen_random_uuid() NOT NULL,
	address_line1 varchar(200) NOT NULL,
	address_line2 varchar(200) NULL,
	city varchar(100) NOT NULL,
	state varchar(100) NOT NULL,
	pin_code varchar(10) NOT NULL,
	CONSTRAINT address_pkey PRIMARY KEY (id)
);


-- "DIGIT3".backup_match definition

-- Drop table

-- DROP TABLE "DIGIT3".backup_match;

CREATE TABLE "DIGIT3".backup_match (
	count int8 NULL
);


-- "DIGIT3"."document" definition

-- Drop table

-- DROP TABLE "DIGIT3"."document";

CREATE TABLE "DIGIT3"."document" (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	property_id uuid NOT NULL,
	document_type varchar(100) NOT NULL,
	document_name varchar(255) NOT NULL,
	file_store_id varchar(200) NULL,
	upload_date timestamptz DEFAULT now() NOT NULL,
	"action" varchar(100) DEFAULT 'PENDING'::character varying NOT NULL,
	uploaded_by varchar(200) NULL,
	"size" varchar(50) NULL,
	CONSTRAINT chk_document_name_not_empty CHECK ((char_length((COALESCE(document_name, ''::character varying))::text) > 0)),
	CONSTRAINT chk_document_type_not_empty CHECK ((char_length((COALESCE(document_type, ''::character varying))::text) > 0)),
	CONSTRAINT document_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_document_file_store_id ON "DIGIT3".document USING btree (file_store_id);
CREATE INDEX idx_document_name ON "DIGIT3".document USING btree (document_name);
CREATE INDEX idx_document_property_id ON "DIGIT3".document USING btree (property_id);
CREATE INDEX idx_document_type ON "DIGIT3".document USING btree (document_type);


-- "DIGIT3".properties definition

-- Drop table

-- DROP TABLE "DIGIT3".properties;

CREATE TABLE "DIGIT3".properties (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	property_no varchar(100) NOT NULL,
	ownership_type "DIGIT3"."ownership_type_enum" NULL,
	property_type "DIGIT3"."property_type_enum" NULL,
	complex_name varchar(200) NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT properties_pkey PRIMARY KEY (id),
	CONSTRAINT properties_property_no_key UNIQUE (property_no)
);


-- "DIGIT3".users definition

-- Drop table

-- DROP TABLE "DIGIT3".users;

CREATE TABLE "DIGIT3".users (
	keycloak_user_id varchar(255) NOT NULL,
	username varchar(255) NOT NULL,
	email varchar(255) NOT NULL,
	"role" varchar(50) NOT NULL,
	is_active bool DEFAULT true NOT NULL,
	preferred_language varchar(50) NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	created_by varchar(255) NULL,
	updated_by varchar(255) NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (keycloak_user_id),
	CONSTRAINT users_username_key UNIQUE (username)
);


-- "DIGIT3".with_ward definition

-- Drop table

-- DROP TABLE "DIGIT3".with_ward;

CREATE TABLE "DIGIT3".with_ward (
	count int8 NULL
);


-- "DIGIT3".additional_property_details definition

-- Drop table

-- DROP TABLE "DIGIT3".additional_property_details;

CREATE TABLE "DIGIT3".additional_property_details (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	property_id uuid NOT NULL,
	field_name varchar(100) NOT NULL,
	field_value text NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT additional_property_details_pkey PRIMARY KEY (id),
	CONSTRAINT fk_additional_property_details_property FOREIGN KEY (property_id) REFERENCES "DIGIT3".properties(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- "DIGIT3".amenities definition

-- Drop table

-- DROP TABLE "DIGIT3".amenities;

CREATE TABLE "DIGIT3".amenities (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	property_id uuid NOT NULL,
	"type" _text NOT NULL,
	description text NULL,
	expiry_date date NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT amenities_pkey PRIMARY KEY (id),
	CONSTRAINT uk_amenities_property_id UNIQUE (property_id),
	CONSTRAINT fk_amenities_property FOREIGN KEY (property_id) REFERENCES "DIGIT3".properties(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_amenities_property_id ON "DIGIT3".amenities USING btree (property_id);


-- "DIGIT3".applications definition

-- Drop table

-- DROP TABLE "DIGIT3".applications;

CREATE TABLE "DIGIT3".applications (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	application_no varchar(50) NULL,
	property_id uuid NOT NULL,
	priority varchar(10) NULL,
	due_date date NULL,
	assigned_agent varchar(255) NULL,
	status varchar(50) NULL,
	workflow_instance_id varchar(100) NULL,
	applied_by varchar(200) NOT NULL,
	assessee_id varchar(255) NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	tenant_id text NULL,
	is_draft bool DEFAULT false NULL,
	CONSTRAINT applications_application_no_key UNIQUE (application_no),
	CONSTRAINT applications_pkey PRIMARY KEY (id),
	CONSTRAINT applications_priority_check CHECK (((priority)::text = ANY ((ARRAY['LOW'::character varying, 'MEDIUM'::character varying, 'HIGH'::character varying, 'NULL'::character varying])::text[]))),
	CONSTRAINT fk_application_property FOREIGN KEY (property_id) REFERENCES "DIGIT3".properties(id) ON DELETE RESTRICT ON UPDATE CASCADE,
	CONSTRAINT fk_applications_assessee FOREIGN KEY (assessee_id) REFERENCES "DIGIT3".users(keycloak_user_id)
);


-- "DIGIT3".assessment_details definition

-- Drop table

-- DROP TABLE "DIGIT3".assessment_details;

CREATE TABLE "DIGIT3".assessment_details (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	property_id uuid NOT NULL,
	reason_of_creation varchar(200) NULL,
	occupancy_certificate_number varchar(100) NULL,
	occupancy_certificate_date date NULL,
	extend_of_site varchar(200) NULL,
	is_land_underneath_building bool DEFAULT false NULL,
	is_unspecified_share bool DEFAULT false NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT assessment_details_pkey PRIMARY KEY (id),
	CONSTRAINT assessment_details_property_id_key UNIQUE (property_id),
	CONSTRAINT fk_assessment_details_property FOREIGN KEY (property_id) REFERENCES "DIGIT3".properties(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- "DIGIT3".construction_details definition

-- Drop table

-- DROP TABLE "DIGIT3".construction_details;

CREATE TABLE "DIGIT3".construction_details (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	property_id uuid NOT NULL,
	floor_type varchar(100) NULL,
	wall_type varchar(100) NULL,
	roof_type varchar(100) NULL,
	wood_type varchar(100) NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT construction_details_pkey PRIMARY KEY (id),
	CONSTRAINT construction_details_property_id_key UNIQUE (property_id),
	CONSTRAINT fk_construction_details_property FOREIGN KEY (property_id) REFERENCES "DIGIT3".properties(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- "DIGIT3".floor_details definition

-- Drop table

-- DROP TABLE "DIGIT3".floor_details;

CREATE TABLE "DIGIT3".floor_details (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	construction_details_id uuid NOT NULL,
	floor_no int4 NOT NULL,
	classification varchar(100) NULL,
	nature_of_usage varchar(100) NULL,
	firm_name varchar(200) NULL,
	occupancy_type varchar(100) NULL,
	occupancy_name varchar(200) NULL,
	construction_date date NULL,
	effective_from_date date NULL,
	unstructured_land varchar(200) NULL,
	length_ft numeric(10, 2) NULL,
	breadth_ft numeric(10, 2) NULL,
	plinth_area_sq_ft numeric(10, 2) NULL,
	building_permission_no varchar(100) NULL,
	floor_details_entered bool DEFAULT false NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT floor_details_pkey PRIMARY KEY (id),
	CONSTRAINT fk_floor_details_construction FOREIGN KEY (construction_details_id) REFERENCES "DIGIT3".construction_details(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_floor_details_construction_id ON "DIGIT3".floor_details USING btree (construction_details_id);


-- "DIGIT3".gis_data definition

-- Drop table

-- DROP TABLE "DIGIT3".gis_data;

CREATE TABLE "DIGIT3".gis_data (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	property_id uuid NOT NULL,
	"source" "DIGIT3"."gis_source_enum" NOT NULL,
	"type" "DIGIT3"."gis_type_enum" NOT NULL,
	entity_type varchar(100) NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT gis_data_pkey PRIMARY KEY (id),
	CONSTRAINT gis_data_property_id_key UNIQUE (property_id),
	CONSTRAINT fk_gis_data_property FOREIGN KEY (property_id) REFERENCES "DIGIT3".properties(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- Add latitude and longitude columns to gis_data table
ALTER TABLE "DIGIT3".gis_data 
ADD COLUMN latitude NUMERIC(10, 8) NULL,
ADD COLUMN longitude NUMERIC(11, 8) NULL;

-- Add check constraints for valid coordinate ranges
ALTER TABLE "DIGIT3".gis_data
ADD CONSTRAINT chk_gis_data_latitude CHECK (latitude >= -90.0 AND latitude <= 90.0),
ADD CONSTRAINT chk_gis_data_longitude CHECK (longitude >= -180.0 AND longitude <= 180.0);

-- Add comment for documentation
COMMENT ON COLUMN "DIGIT3".gis_data.latitude IS 'Latitude coordinate (-90 to 90)';
COMMENT ON COLUMN "DIGIT3".gis_data.longitude IS 'Longitude coordinate (-180 to 180)';

-- Optional: Add index for coordinate-based queries
CREATE INDEX idx_gis_data_coordinates ON "DIGIT3".gis_data(latitude, longitude);


-- "DIGIT3".igrs definition

-- Drop table

-- DROP TABLE "DIGIT3".igrs;

CREATE TABLE "DIGIT3".igrs (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	property_id uuid NULL,
	habitation varchar(200) NOT NULL,
	igrs_ward varchar(100) NULL,
	igrs_locality varchar(100) NULL,
	igrs_block varchar(100) NULL,
	door_no_from varchar(50) NULL,
	door_no_to varchar(50) NULL,
	igrs_classification varchar(100) NULL,
	built_up_area_pct numeric(7, 2) NULL,
	front_setback numeric(8, 2) NULL,
	rear_setback numeric(8, 2) NULL,
	side_setback numeric(8, 2) NULL,
	total_plinth_area numeric(10, 2) NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	CONSTRAINT igrs_pkey PRIMARY KEY (id),
	CONSTRAINT igrs_property_id_key UNIQUE (property_id),
	CONSTRAINT fk_igrs_property FOREIGN KEY (property_id) REFERENCES "DIGIT3".properties(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- "DIGIT3".property_addresses definition

-- Drop table

-- DROP TABLE "DIGIT3".property_addresses;

CREATE TABLE "DIGIT3".property_addresses (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	property_id uuid NOT NULL,
	locality varchar(200) NULL,
	zone_no varchar(50) NULL,
	ward_no varchar(50) NULL,
	block_no varchar(50) NULL,
	street varchar(200) NULL,
	election_ward varchar(50) NULL,
	secretariat_ward varchar(50) NULL,
	pin_code int4 NULL,
	different_correspondence_address bool DEFAULT false NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT property_addresses_pkey PRIMARY KEY (id),
	CONSTRAINT property_addresses_property_id_key UNIQUE (property_id),
	CONSTRAINT fk_property_address_property FOREIGN KEY (property_id) REFERENCES "DIGIT3".properties(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- "DIGIT3".property_owner definition

-- Drop table

-- DROP TABLE "DIGIT3".property_owner;

CREATE TABLE "DIGIT3".property_owner (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	property_id uuid NOT NULL,
	adhaar_no int8 NOT NULL,
	"name" varchar(200) NOT NULL,
	contact_no varchar(15) NOT NULL,
	email varchar(100) NULL,
	gender varchar(10) NOT NULL,
	guardian varchar(200) NULL,
	guardian_type varchar(10) NULL,
	relationship_to_property varchar(20) NULL,
	ownership_share numeric(5, 2) DEFAULT 0 NULL,
	is_primary_owner bool DEFAULT false NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT property_owner_gender_check CHECK (((gender)::text = ANY ((ARRAY['MALE'::character varying, 'FEMALE'::character varying, 'OTHER'::character varying])::text[]))),
	CONSTRAINT property_owner_guardian_type_check CHECK (((guardian_type)::text = ANY ((ARRAY['FATHER'::character varying, 'MOTHER'::character varying, 'SPOUSE'::character varying, 'OTHER'::character varying])::text[]))),
	CONSTRAINT property_owner_pkey PRIMARY KEY (id),
	CONSTRAINT property_owner_relationship_to_property_check CHECK (((relationship_to_property)::text = ANY ((ARRAY['OWNER'::character varying, 'CO_OWNER'::character varying, 'JOINT_OWNER'::character varying, 'LEGAL_HEIR'::character varying, 'POWER_OF_ATTORNEY'::character varying, 'OTHER'::character varying])::text[]))),
	CONSTRAINT fk_property_owner_property FOREIGN KEY (property_id) REFERENCES "DIGIT3".properties(id) ON DELETE CASCADE ON UPDATE CASCADE
);


-- "DIGIT3".user_profiles definition

-- Drop table

-- DROP TABLE "DIGIT3".user_profiles;

CREATE TABLE "DIGIT3".user_profiles (
	user_profile_id varchar(255) NOT NULL,
	first_name varchar(100) NOT NULL,
	last_name varchar(100) NOT NULL,
	full_name varchar(200) NOT NULL,
	phone_number varchar(15) NOT NULL,
	adhaar_no int8 NOT NULL,
	gender varchar(10) NOT NULL,
	guardian varchar(100) NULL,
	guardian_type varchar(50) NULL,
	date_of_birth date NULL,
	department varchar(100) NULL,
	designation varchar(100) NULL,
	work_location varchar(200) NULL,
	profile_picture text NULL,
	relationship_to_property varchar(50) DEFAULT ''::character varying NOT NULL,
	ownership_share float8 DEFAULT 0 NOT NULL,
	is_primary_owner bool DEFAULT false NOT NULL,
	is_verified bool DEFAULT false NULL,
	address_id uuid NULL,
	CONSTRAINT userprofile_adhaar_no_key UNIQUE (adhaar_no),
	CONSTRAINT userprofile_ownership_share_check CHECK (((ownership_share >= (0)::double precision) AND (ownership_share <= (100)::double precision))),
	CONSTRAINT userprofile_pkey PRIMARY KEY (user_profile_id),
	CONSTRAINT userprofile_address_id_fkey FOREIGN KEY (address_id) REFERENCES "DIGIT3".address(id) ON DELETE CASCADE,
	CONSTRAINT userprofile_id_fkey FOREIGN KEY (user_profile_id) REFERENCES "DIGIT3".users(keycloak_user_id)
);


-- "DIGIT3".zone_mapping definition

-- Drop table

-- DROP TABLE "DIGIT3".zone_mapping;

CREATE TABLE "DIGIT3".zone_mapping (
	user_id varchar(255) NOT NULL,
	"zone" varchar(50) NOT NULL,
	ward _text NOT NULL,
	created_at timestamp DEFAULT now() NULL,
	updated_at timestamp DEFAULT now() NULL,
	CONSTRAINT zone_mapping_pkey PRIMARY KEY (user_id, zone),
	CONSTRAINT zone_mapping_user_id_fkey FOREIGN KEY (user_id) REFERENCES "DIGIT3".users(keycloak_user_id) ON DELETE CASCADE
);
CREATE INDEX idx_zone_mapping_user ON "DIGIT3".zone_mapping USING btree (user_id);
CREATE INDEX idx_zone_mapping_ward_gin ON "DIGIT3".zone_mapping USING gin (ward);
CREATE INDEX idx_zone_mapping_zone ON "DIGIT3".zone_mapping USING btree (zone);


-- "DIGIT3".application_logs definition

-- Drop table

-- DROP TABLE "DIGIT3".application_logs;

CREATE TABLE "DIGIT3".application_logs (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	"action" "DIGIT3"."digit3_action_enum" NULL,
	performed_by varchar(200) NOT NULL,
	performed_date timestamp NOT NULL,
	"comments" text NULL,
	metadata json NULL,
	file_store_id uuid NULL,
	application_id uuid NOT NULL,
	created_at timestamp DEFAULT now() NULL,
	CONSTRAINT application_logs_pkey PRIMARY KEY (id),
	CONSTRAINT fk_application FOREIGN KEY (application_id) REFERENCES "DIGIT3".applications(id)
);
CREATE INDEX idx_application_logs_application_id ON "DIGIT3".application_logs USING btree (application_id);
CREATE INDEX idx_application_logs_file_store_id ON "DIGIT3".application_logs USING btree (file_store_id);
CREATE INDEX idx_application_logs_performed_date ON "DIGIT3".application_logs USING btree (performed_date);


-- "DIGIT3".coordinates definition

-- Drop table

-- DROP TABLE "DIGIT3".coordinates;

CREATE TABLE "DIGIT3".coordinates (
	id uuid DEFAULT "DIGIT3".uuid_generate_v4() NOT NULL,
	gis_data_id uuid NOT NULL,
	latitude numeric(10, 8) NOT NULL,
	longitude numeric(11, 8) NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT coordinates_pkey PRIMARY KEY (id),
	CONSTRAINT fk_coordinates_gis_data FOREIGN KEY (gis_data_id) REFERENCES "DIGIT3".gis_data(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX idx_coordinates_gis_data_id ON "DIGIT3".coordinates USING btree (gis_data_id);


-- DROP FUNCTION "DIGIT3".uuid_generate_v4();
