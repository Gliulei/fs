## ADDED Requirements

### Requirement: Configuration import command
The system SHALL provide an `import` command that extracts and imports host configuration files from a packaged archive.

#### Scenario: Import configurations from archive
- **WHEN** user executes `fs import config-archive.tar.gz`
- **THEN** system extracts and imports all YAML configuration files to ~/.fs/hosts/

#### Scenario: Import with conflict handling
- **WHEN** user executes `fs import config-archive.tar.gz` and conflicts exist
- **THEN** system prompts user to choose conflict resolution strategy (overwrite/skip/rename)

#### Scenario: Import with force flag
- **WHEN** user executes `fs import config-archive.tar.gz --force`
- **THEN** system overwrites existing configurations without prompting

#### Scenario: Import with skip flag
- **WHEN** user executes `fs import config-archive.tar.gz --skip-existing`
- **THEN** system skips importing configurations that already exist

### Requirement: Archive validation
The system SHALL validate imported archive files before extraction:
- Verify file format is tar.gz
- Check archive integrity
- Validate required file structure
- Confirm VERSION compatibility

#### Scenario: Archive format validation
- **WHEN** user attempts to import non-tar.gz file
- **THEN** system reports error and rejects the import

#### Scenario: Archive structure validation
- **WHEN** user attempts to import archive with missing required files
- **THEN** system reports validation error with specific missing components

#### Scenario: Version compatibility check
- **WHEN** user attempts to import archive created with incompatible fs version
- **THEN** system warns user about version mismatch but allows proceeding

### Requirement: Configuration file validation
The system SHALL validate individual configuration files before importing:
- Parse YAML format
- Validate required fields
- Check for valid SSH configuration structure

#### Scenario: YAML format validation
- **WHEN** archive contains invalid YAML file
- **THEN** system reports the specific file as invalid and skips import

#### Scenario: Required fields validation
- **WHEN** configuration file is missing required fields
- **THEN** system reports validation error with missing field names

### Requirement: Conflict resolution options
The system SHALL provide the following conflict resolution strategies:
- overwrite: Replace existing configuration with imported one
- skip: Keep existing configuration, skip imported one
- rename: Save imported configuration with new name
- abort: Stop import process entirely

#### Scenario: Interactive conflict resolution
- **WHEN** conflicts are detected during import
- **THEN** system presents user with resolution options and waits for selection

#### Scenario: Batch conflict resolution with flags
- **WHEN** user provides conflict resolution flags
- **THEN** system applies specified strategy to all conflicts without prompting

### Requirement: Import reporting
The system SHALL provide detailed import operation reporting:
- Number of configurations successfully imported
- Number of configurations skipped due to conflicts
- Number of configurations that failed validation
- List of imported configuration names

#### Scenario: Import completion report
- **WHEN** import operation completes
- **THEN** system displays summary of import results including success/failure counts