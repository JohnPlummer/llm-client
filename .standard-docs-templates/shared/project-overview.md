# {{PROJECT_NAME}} - Project Overview

## Description
{{PROJECT_DESCRIPTION}}

## Purpose & Goals
<!-- What problem does this project solve? What are its main objectives? -->

## High-Level Architecture
<!-- Describe the overall system architecture, main components, and how they interact -->

## Technology Stack
{{#each FRAMEWORKS}}
- {{this}}
{{/each}}

## Key Features
<!-- List the main features and capabilities of this project -->

## Project Structure
```
{{PROJECT_ROOT}}/
├── {{#each MAIN_DIRECTORIES}}
├── {{this}}/
{{/each}}
```

## Getting Started
For development setup instructions, see [development-setup.md](./development-setup.md).

## Related Documentation
- [Development Setup](./development-setup.md)
- [Key Components](./key-components.md)
- [API Reference](./api-reference.md) {{#if HAS_API}}{{/if}}
- [Deployment Guide](./deployment-guide.md)
- [Troubleshooting](./troubleshooting.md)

## Contact & Support
<!-- Project maintainers, support channels, etc. -->

---
*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*
