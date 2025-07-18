# {{PROJECT_NAME}} - Recent Changes

*This document is automatically updated based on git history and merged pull requests.*

## Recent Updates

{{#each RECENT_CHANGES}}
### {{date}} - {{type}}
**{{title}}**
{{#if description}}
{{description}}
{{/if}}

**Changes:**
{{#each changes}}
- {{this}}
{{/each}}

**Impact:** {{impact}}
**Author:** {{author}}
{{#if pullRequest}}
**Pull Request:** [#{{pullRequest.number}}]({{pullRequest.url}})
{{/if}}

---
{{/each}}

## Significant Changes (Last 30 days)

### New Features
{{#each NEW_FEATURES}}
- **{{title}}** ({{date}}) - {{description}}
{{/each}}

### Bug Fixes
{{#each BUG_FIXES}}
- **{{title}}** ({{date}}) - {{description}}
{{/each}}

### Breaking Changes
{{#each BREAKING_CHANGES}}
- **{{title}}** ({{date}}) - {{description}}
  - **Migration Guide:** {{migrationGuide}}
{{/each}}

### Dependencies
{{#each DEPENDENCY_UPDATES}}
- **{{package}}** updated from {{fromVersion}} to {{toVersion}} ({{date}})
{{/each}}

## Change Statistics
- **Total commits:** {{TOTAL_COMMITS}}
- **Contributors:** {{CONTRIBUTORS_COUNT}}
- **Files changed:** {{FILES_CHANGED}}
- **Lines added:** {{LINES_ADDED}}
- **Lines removed:** {{LINES_REMOVED}}

## Upcoming Changes
{{#each UPCOMING_CHANGES}}
- **{{title}}** - {{description}} (Expected: {{expectedDate}})
{{/each}}

---
*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*
