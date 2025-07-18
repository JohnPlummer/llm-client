# {{PROJECT_NAME}} - Key Components

## Architecture Overview
<!-- High-level component diagram or description -->

## Core Components

{{#each COMPONENTS}}
### {{name}}
**Location:** `{{path}}`
**Purpose:** {{purpose}}

{{#if description}}
{{description}}
{{/if}}

**Key Files:**
{{#each files}}
- `{{path}}` - {{description}}
{{/each}}

**Dependencies:**
{{#each dependencies}}
- {{this}}
{{/each}}

**Used By:**
{{#each usedBy}}
- {{this}}
{{/each}}

---
{{/each}}

## Component Relationships
<!-- Describe how components interact with each other -->

## Data Flow
<!-- Describe how data flows through the system -->

## Extension Points
<!-- Areas where new functionality can be added -->

## Performance Considerations
<!-- Important performance characteristics of key components -->

---
*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*
