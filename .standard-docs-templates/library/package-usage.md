# {{PROJECT_NAME}} - Package Usage

## Installation

### Package Manager

{{#each PACKAGE_MANAGERS}}
#### {{name}}

```bash
{{installCommand}}
```

{{/each}}

### CDN

{{#if HAS_CDN}}

```html
<script src="{{CDN_URL}}"></script>
```

{{/if}}

## Getting Started

### Basic Usage

```javascript
{{BASIC_USAGE_EXAMPLE}}
```

### Configuration

{{#if HAS_CONFIG}}

```javascript
{{CONFIG_EXAMPLE}}
```

**Configuration Options:**

{{#each CONFIG_OPTIONS}}
- `{{name}}` ({{type}}) - {{description}} {{#if defaultValue}}(Default: {{defaultValue}}){{/if}}
{{/each}}

{{/if}}

## API Reference

### Core API

{{#each CORE_API}}
#### {{name}}

{{#if description}}
{{description}}
{{/if}}

**Signature:**
```javascript
{{signature}}
```

**Parameters:**
{{#each parameters}}
- `{{name}}` ({{type}}) {{#if required}}**Required**{{/if}} - {{description}}
{{/each}}

**Returns:**
{{returnType}} - {{returnDescription}}

**Example:**
```javascript
{{example}}
```

---
{{/each}}

### Utility Functions
{{#each UTILITY_FUNCTIONS}}
#### {{name}}
{{#if description}}
{{description}}
{{/if}}

**Signature:**
```javascript
{{signature}}
```

**Example:**
```javascript
{{example}}
```

---
{{/each}}

## Advanced Usage

### Plugins
{{#if HAS_PLUGINS}}
{{#each PLUGINS}}
#### {{name}}
{{description}}

**Installation:**
```bash
{{installCommand}}
```

**Usage:**
```javascript
{{usage}}
```
{{/each}}
{{/if}}

### Middleware
{{#if HAS_MIDDLEWARE}}
{{#each MIDDLEWARE}}
#### {{name}}
{{description}}

**Usage:**
```javascript
{{usage}}
```
{{/each}}
{{/if}}

### Hooks
{{#if HAS_HOOKS}}
{{#each HOOKS}}
#### {{name}}
{{description}}

**Usage:**
```javascript
{{usage}}
```
{{/each}}
{{/if}}

## Framework Integration

### React
{{#if HAS_REACT_INTEGRATION}}
```jsx
{{REACT_EXAMPLE}}
```
{{/if}}

### Vue
{{#if HAS_VUE_INTEGRATION}}
```vue
{{VUE_EXAMPLE}}
```
{{/if}}

### Angular
{{#if HAS_ANGULAR_INTEGRATION}}
```typescript
{{ANGULAR_EXAMPLE}}
```
{{/if}}

### Node.js
{{#if HAS_NODE_INTEGRATION}}
```javascript
{{NODE_EXAMPLE}}
```
{{/if}}

## Examples

### Common Use Cases
{{#each COMMON_EXAMPLES}}
#### {{title}}
{{description}}

```javascript
{{code}}
```
{{/each}}

### Complete Examples
{{#each COMPLETE_EXAMPLES}}
#### {{title}}
{{description}}

**Files:**
{{#each files}}
- `{{filename}}`:
```{{language}}
{{content}}
```
{{/each}}
{{/each}}

## TypeScript Support

### Type Definitions
{{#if HAS_TYPESCRIPT}}
Type definitions are included in the package.

**Main Types:**
{{#each MAIN_TYPES}}
- `{{name}}` - {{description}}
{{/each}}

**Generic Types:**
{{#each GENERIC_TYPES}}
- `{{name}}<{{generics}}>` - {{description}}
{{/each}}
{{/if}}

## Error Handling

### Error Types
{{#each ERROR_TYPES}}
#### {{name}}
{{description}}

**Properties:**
{{#each properties}}
- `{{name}}`: {{type}} - {{description}}
{{/each}}

**Example:**
```javascript
{{example}}
```
{{/each}}

### Error Handling Patterns
```javascript
{{ERROR_HANDLING_EXAMPLE}}
```

## Performance

### Optimization Tips
{{#each PERFORMANCE_TIPS}}
- **{{title}}:** {{description}}
{{/each}}

### Bundle Size
{{#if HAS_BUNDLE_INFO}}
- **Minified:** {{MINIFIED_SIZE}}
- **Gzipped:** {{GZIPPED_SIZE}}
- **Tree-shakeable:** {{TREE_SHAKEABLE}}
{{/if}}

## Browser Support
{{#if HAS_BROWSER_SUPPORT}}
{{#each BROWSER_SUPPORT}}
- **{{browser}}:** {{version}}
{{/each}}
{{/if}}

## Migration Guide

{{#each MIGRATION_GUIDES}}
### From {{fromVersion}} to {{toVersion}}

**Breaking Changes:**
{{#each breakingChanges}}
- {{description}}
{{/each}}

**Migration Steps:**
{{#each steps}}
{{step}}. {{description}}
{{#if code}}
```javascript
{{code}}
```
{{/if}}
{{/each}}
{{/each}}

## Contributing

### Development Setup
```bash
{{DEV_SETUP_COMMAND}}
```

### Building
```bash
{{BUILD_COMMAND}}
```

### Testing
```bash
{{TEST_COMMAND}}
```

### Submitting Issues
{{ISSUE_GUIDELINES}}

### Pull Requests
{{PR_GUIDELINES}}

---
*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*
