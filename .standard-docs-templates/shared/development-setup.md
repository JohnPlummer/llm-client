# {{PROJECT_NAME}} - Development Setup

## Prerequisites
{{#each PREREQUISITES}}
- {{this}}
{{/each}}

## Environment Setup

### 1. Clone the Repository
```bash
git clone {{REPOSITORY_URL}}
cd {{PROJECT_NAME}}
```

### 2. Install Dependencies
{{#if PACKAGE_MANAGER}}
```bash
{{PACKAGE_MANAGER}} install
```
{{/if}}

### 3. Environment Configuration
{{#if HAS_ENV_FILE}}
Copy the example environment file and configure:
```bash
cp .env.example .env
```

Required environment variables:
{{#each ENV_VARIABLES}}
- `{{name}}`: {{description}}
{{/each}}
{{/if}}

### 4. Database Setup
{{#if HAS_DATABASE}}
<!-- Database setup instructions -->
{{/if}}

### 5. Start Development Server
```bash
{{START_COMMAND}}
```

## Development Workflow

### Running Tests
```bash
{{TEST_COMMAND}}
```

### Code Quality
{{#if HAS_LINTER}}
```bash
{{LINT_COMMAND}}
```
{{/if}}

{{#if HAS_FORMATTER}}
```bash
{{FORMAT_COMMAND}}
```
{{/if}}

### Building for Production
```bash
{{BUILD_COMMAND}}
```

## IDE Configuration
{{#if HAS_VSCODE_CONFIG}}
This project includes VS Code configuration in `.vscode/`. Recommended extensions:
{{#each VSCODE_EXTENSIONS}}
- {{this}}
{{/each}}
{{/if}}

## Troubleshooting
For common development issues, see [troubleshooting.md](./troubleshooting.md).

## Contributing
<!-- Link to contributing guidelines -->

---
*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*
