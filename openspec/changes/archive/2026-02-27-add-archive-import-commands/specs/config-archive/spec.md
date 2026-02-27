## ADDED Requirements

### Requirement: Configuration archive command
The system SHALL provide an `archive` command that creates a compressed package containing all host configuration files.

#### Scenario: Archive all configurations
- **WHEN** user executes `fs archive`
- **THEN** system creates a tar.gz file containing all YAML configuration files from ~/.fs/hosts/

#### Scenario: Archive with custom output path
- **WHEN** user executes `fs archive --output /path/to/custom.tar.gz`
- **THEN** system creates the archive file at the specified location

#### Scenario: Archive with compression level
- **WHEN** user executes `fs archive --compression-level 9`
- **THEN** system creates archive with specified compression level (1-9)

### Requirement: Archive file structure
The system SHALL create archive files with the following structure:
```
fs-config-archive/
├── config/
│   ├── host1.yaml
│   ├── host2.yaml
│   └── ...
├── metadata.json
└── VERSION
```

#### Scenario: Archive file structure validation
- **WHEN** archive command completes successfully
- **THEN** the generated tar.gz file contains the specified directory structure

### Requirement: Metadata inclusion
The system SHALL include metadata.json file in the archive containing:
- Creation timestamp
- fs version used for archiving
- Number of configuration files included
- Host configuration names list

#### Scenario: Metadata file generation
- **WHEN** archive command executes
- **THEN** metadata.json contains all required information fields

### Requirement: Version tracking
The system SHALL include a VERSION file containing the fs tool version.

#### Scenario: Version file creation
- **WHEN** archive is created
- **THEN** VERSION file contains the current fs version string