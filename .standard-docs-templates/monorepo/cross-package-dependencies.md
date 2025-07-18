# {{PROJECT_NAME}} - Cross-Package Dependencies

## Dependency Overview

This document describes how packages within the monorepo depend on each other and external dependencies.

## Dependency Graph

<!-- Include visual dependency graph -->

## Internal Dependencies

### Package Relationships

{{#each INTERNAL_DEPENDENCIES}}
#### {{package}}

**Type:** {{type}}
**Depends On:**

{{#each dependencies}}
- **{{name}}** ({{version}}) - {{relationship}}
{{/each}}

**Used By:**

{{#each consumers}}
- **{{name}}** - {{usage}}
{{/each}}

**Circular Dependencies:** {{circularDependencies}}

---
{{/each}}

### Dependency Layers

{{#each DEPENDENCY_LAYERS}}
#### Layer {{level}}: {{name}}

**Description:** {{description}}
**Packages:**

{{#each packages}}
- {{name}} - {{description}}
{{/each}}

**Dependencies:** {{dependencies}}
**Can depend on:** {{canDependOn}}

---
{{/each}}

## External Dependencies

### Shared Dependencies

{{#each SHARED_EXTERNAL_DEPENDENCIES}}
#### {{name}}

**Version:** {{version}}
**Type:** {{type}}
**Used By:**

{{#each usedBy}}
- {{package}} ({{version}})
{{/each}}

**Reason for sharing:** {{reason}}

---
{{/each}}

### Package-Specific Dependencies

{{#each PACKAGE_SPECIFIC_DEPENDENCIES}}
#### {{package}}

**Dependencies:**

{{#each dependencies}}
- **{{name}}** ({{version}}) - {{purpose}}
{{/each}}

---
{{/each}}

## Dependency Management

### Version Synchronization

{{#if HAS_VERSION_SYNC}}
**Strategy:** {{VERSION_SYNC_STRATEGY}}
**Tool:** {{VERSION_SYNC_TOOL}}

**Synchronized Packages:**

{{#each SYNCHRONIZED_PACKAGES}}
- {{name}}: {{versions}}
{{/each}}

{{/if}}

### Hoisting Strategy

{{#if HAS_HOISTING}}
**Hoisting Method:** {{HOISTING_METHOD}}
**Hoisted Dependencies:**

{{#each HOISTED_DEPENDENCIES}}
- {{name}} ({{version}})
{{/each}}

{{/if}}

### Workspace Protocols

{{#if HAS_WORKSPACE_PROTOCOLS}}
**Protocol:** {{WORKSPACE_PROTOCOL}}

**Usage Examples:**

```json
{
  "dependencies": {
    "{{EXAMPLE_PACKAGE}}": "workspace:*"
  }
}
```

{{/if}}

## Development Workflow

### Adding Dependencies

```bash
# Add dependency to specific package
{{ADD_DEPENDENCY_COMMAND}}

# Add workspace dependency
{{ADD_WORKSPACE_DEPENDENCY_COMMAND}}

# Add shared dependency
{{ADD_SHARED_DEPENDENCY_COMMAND}}
```

### Updating Dependencies

```bash
# Update all dependencies
{{UPDATE_ALL_COMMAND}}

# Update specific dependency
{{UPDATE_SPECIFIC_COMMAND}}

# Update workspace dependencies
{{UPDATE_WORKSPACE_COMMAND}}
```

### Removing Dependencies

```bash
# Remove dependency from package
{{REMOVE_DEPENDENCY_COMMAND}}

# Remove shared dependency
{{REMOVE_SHARED_DEPENDENCY_COMMAND}}
```

## Dependency Validation

### Validation Rules

{{#each VALIDATION_RULES}}
- **{{rule}}:** {{description}}
{{/each}}

### Validation Commands

```bash
# Check dependency consistency
{{VALIDATE_DEPENDENCIES_COMMAND}}

# Check for circular dependencies
{{CHECK_CIRCULAR_COMMAND}}

# Check for duplicate dependencies
{{CHECK_DUPLICATES_COMMAND}}
```

## Build Dependencies

### Build Order

{{#each BUILD_ORDER}}
{{order}}. **{{package}}** - {{buildTime}}
   - Dependencies: {{dependencies}}
   - Affects: {{affects}}
{{/each}}

### Incremental Builds

{{#if HAS_INCREMENTAL_BUILDS}}
**Tool:** {{INCREMENTAL_BUILD_TOOL}}
**Strategy:** {{INCREMENTAL_STRATEGY}}

**Commands:**

```bash
# Build changed packages
{{BUILD_CHANGED_COMMAND}}

# Build affected packages
{{BUILD_AFFECTED_COMMAND}}
```

{{/if}}

## Testing Dependencies

### Test Dependencies

{{#each TEST_DEPENDENCIES}}
#### {{package}}

**Test Dependencies:**

{{#each dependencies}}
- **{{name}}** ({{version}}) - {{purpose}}
{{/each}}

**Shared Test Utilities:**

{{#each sharedTestUtils}}
- {{name}} - {{description}}
{{/each}}

---
{{/each}}

### Test Execution Order

{{#if HAS_TEST_ORDERING}}

{{#each TEST_ORDER}}
{{order}}. **{{package}}** - {{testTime}}
   - Prerequisites: {{prerequisites}}
   - Test Type: {{testType}}
{{/each}}

{{/if}}

## Performance Considerations

### Bundle Impact

{{#each BUNDLE_IMPACT}}
#### {{package}}

**Bundle Size:** {{size}}
**Dependencies Contributing to Size:**

{{#each largeDependencies}}
- {{name}}: {{size}}
{{/each}}

**Optimization Opportunities:**

{{#each optimizations}}
- {{technique}}: {{description}}
{{/each}}

---
{{/each}}

### Installation Performance

{{#each INSTALLATION_METRICS}}
- **{{metric}}:** {{value}}
{{/each}}

## Troubleshooting

### Common Dependency Issues

{{#each COMMON_DEPENDENCY_ISSUES}}
#### {{title}}

**Symptoms:** {{symptoms}}
**Causes:** {{causes}}
**Solutions:**

{{#each solutions}}
- {{solution}}
{{/each}}

---
{{/each}}

### Debug Commands

```bash
# Debug dependency resolution
{{DEBUG_RESOLUTION_COMMAND}}

# Show dependency tree
{{SHOW_TREE_COMMAND}}

# Find dependency conflicts
{{FIND_CONFLICTS_COMMAND}}
```

## Migration Guide

### Dependency Updates

{{#each MIGRATION_GUIDES}}
#### {{title}}

**From:** {{fromVersion}}
**To:** {{toVersion}}

**Breaking Changes:**

{{#each breakingChanges}}
- {{change}}
{{/each}}

**Migration Steps:**

{{#each steps}}
{{step}}. {{description}}
{{/each}}

---
{{/each}}

---
*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*

