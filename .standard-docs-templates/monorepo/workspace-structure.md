# {{PROJECT_NAME}} - Workspace Structure

## Monorepo Overview
**Monorepo Tool:** {{MONOREPO_TOOL}}
**Package Manager:** {{PACKAGE_MANAGER}}
**Total Packages:** {{TOTAL_PACKAGES}}

## Workspace Layout
```
{{PROJECT_ROOT}}/
├── packages/
{{#each PACKAGES}}
│   ├── {{name}}/                    # {{type}} - {{description}}
{{/each}}
├── apps/
{{#each APPS}}
│   ├── {{name}}/                    # {{type}} - {{description}}
{{/each}}
├── tools/
{{#each TOOLS}}
│   ├── {{name}}/                    # {{description}}
{{/each}}
├── docs/                            # Shared documentation
└── scripts/                         # Build and utility scripts
```

## Workspaces

### Libraries/Packages
{{#each PACKAGES}}
#### {{name}}
**Type:** {{type}}
**Path:** `{{path}}`
**Description:** {{description}}
**Version:** {{version}}

**Dependencies:**
{{#each dependencies}}
- {{name}} ({{version}})
{{/each}}

**Used By:**
{{#each consumers}}
- {{name}}
{{/each}}

**Documentation:** [{{name}} README]({{path}}/README.md)

---
{{/each}}

### Applications
{{#each APPS}}
#### {{name}}
**Type:** {{type}}
**Path:** `{{path}}`
**Description:** {{description}}
**Version:** {{version}}

**Dependencies:**
{{#each dependencies}}
- {{name}} ({{version}})
{{/each}}

**Environment:** {{environment}}
**Build Command:** `{{buildCommand}}`
**Start Command:** `{{startCommand}}`

**Documentation:** [{{name}} README]({{path}}/README.md)

---
{{/each}}

### Tools & Utilities
{{#each TOOLS}}
#### {{name}}
**Path:** `{{path}}`
**Purpose:** {{purpose}}
**Usage:** `{{usage}}`

---
{{/each}}

## Dependency Management

### Shared Dependencies
{{#each SHARED_DEPENDENCIES}}
- **{{name}}:** {{version}} - {{description}}
{{/each}}

### Dependency Graph
<!-- Include dependency graph visualization -->

### Version Management
{{#if HAS_VERSION_MANAGEMENT}}
**Strategy:** {{VERSION_STRATEGY}}
**Sync Tool:** {{SYNC_TOOL}}

**Version Constraints:**
{{#each VERSION_CONSTRAINTS}}
- {{package}}: {{constraint}}
{{/each}}
{{/if}}

## Build System

### Build Order
{{#each BUILD_ORDER}}
{{order}}. {{package}} ({{buildTime}})
{{/each}}

### Build Commands
```bash
# Build all packages
{{BUILD_ALL_COMMAND}}

# Build specific package
{{BUILD_PACKAGE_COMMAND}}

# Build with dependencies
{{BUILD_WITH_DEPS_COMMAND}}
```

### Caching
{{#if HAS_BUILD_CACHING}}
**Cache Type:** {{CACHE_TYPE}}
**Cache Location:** {{CACHE_LOCATION}}

**Cache Commands:**
```bash
# Clear cache
{{CLEAR_CACHE_COMMAND}}

# Cache status
{{CACHE_STATUS_COMMAND}}
```
{{/if}}

## Development Workflow

### Getting Started
```bash
# Install all dependencies
{{INSTALL_COMMAND}}

# Build all packages
{{BUILD_ALL_COMMAND}}

# Run all tests
{{TEST_ALL_COMMAND}}
```

### Working with Packages
```bash
# Add dependency to package
{{ADD_DEPENDENCY_COMMAND}}

# Run command in package
{{RUN_IN_PACKAGE_COMMAND}}

# Run command in all packages
{{RUN_IN_ALL_COMMAND}}
```

### Linking Packages
{{#if HAS_PACKAGE_LINKING}}
```bash
# Link packages
{{LINK_COMMAND}}

# Unlink packages
{{UNLINK_COMMAND}}
```
{{/if}}

## Testing Strategy

### Test Structure
{{#each TEST_STRUCTURE}}
- **{{level}}:** {{description}}
{{/each}}

### Test Commands
```bash
# Run all tests
{{TEST_ALL_COMMAND}}

# Run tests for specific package
{{TEST_PACKAGE_COMMAND}}

# Run tests with coverage
{{TEST_COVERAGE_COMMAND}}
```

### Test Coverage
{{#if HAS_TEST_COVERAGE}}
**Coverage Tool:** {{COVERAGE_TOOL}}
**Coverage Threshold:** {{COVERAGE_THRESHOLD}}
**Coverage Report:** {{COVERAGE_REPORT_LOCATION}}
{{/if}}

## Code Quality

### Linting
{{#if HAS_LINTING}}
**Linter:** {{LINTER}}
**Config:** `{{LINT_CONFIG}}`

```bash
# Lint all packages
{{LINT_ALL_COMMAND}}

# Fix linting issues
{{LINT_FIX_COMMAND}}
```
{{/if}}

### Code Formatting
{{#if HAS_FORMATTING}}
**Formatter:** {{FORMATTER}}
**Config:** `{{FORMAT_CONFIG}}`

```bash
# Format all code
{{FORMAT_COMMAND}}

# Check formatting
{{FORMAT_CHECK_COMMAND}}
```
{{/if}}

### Type Checking
{{#if HAS_TYPE_CHECKING}}
**Type Checker:** {{TYPE_CHECKER}}
**Config:** `{{TYPE_CONFIG}}`

```bash
# Type check all packages
{{TYPE_CHECK_COMMAND}}
```
{{/if}}

## Release Management

### Versioning Strategy
{{#if HAS_VERSIONING}}
**Strategy:** {{VERSIONING_STRATEGY}}
**Version Tool:** {{VERSION_TOOL}}

**Release Commands:**
```bash
# Create new version
{{VERSION_COMMAND}}

# Publish packages
{{PUBLISH_COMMAND}}
```
{{/if}}

### Changelog
{{#if HAS_CHANGELOG}}
**Changelog Tool:** {{CHANGELOG_TOOL}}
**Changelog Format:** {{CHANGELOG_FORMAT}}

**Changelog Commands:**
```bash
# Generate changelog
{{CHANGELOG_COMMAND}}
```
{{/if}}

## CI/CD Integration

### Pipeline Structure
{{#if HAS_CI_CD}}
**CI/CD Platform:** {{CI_CD_PLATFORM}}
**Config File:** `{{CI_CD_CONFIG}}`

**Pipeline Stages:**
{{#each CI_CD_STAGES}}
- **{{name}}:** {{description}}
{{/each}}
{{/if}}

### Deployment
{{#if HAS_DEPLOYMENT}}
**Deployment Strategy:** {{DEPLOYMENT_STRATEGY}}

**Deployment Environments:**
{{#each DEPLOYMENT_ENVIRONMENTS}}
- **{{name}}:** {{description}}
{{/each}}
{{/if}}

## Performance Optimization

### Bundle Analysis
{{#if HAS_BUNDLE_ANALYSIS}}
```bash
# Analyze bundle size
{{BUNDLE_ANALYSIS_COMMAND}}
```
{{/if}}

### Build Performance
{{#each BUILD_OPTIMIZATIONS}}
- **{{technique}}:** {{description}}
{{/each}}

## Troubleshooting

### Common Issues
{{#each COMMON_ISSUES}}
#### {{title}}
**Problem:** {{problem}}
**Solution:** {{solution}}
{{/each}}

### Debug Commands
```bash
# Debug dependency resolution
{{DEBUG_DEPS_COMMAND}}

# Debug build issues
{{DEBUG_BUILD_COMMAND}}
```

---
*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*