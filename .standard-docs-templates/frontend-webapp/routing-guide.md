# {{PROJECT_NAME}} - Routing Guide

## Router Configuration

**Router Library:** {{ROUTER_LIBRARY}}
**Configuration File:** `{{ROUTER_CONFIG_FILE}}`

## Route Structure

### Public Routes

{{#each PUBLIC_ROUTES}}
- **{{path}}** - {{description}}
  - Component: `{{component}}`
  - Layout: `{{layout}}`
{{/each}}

### Protected Routes

{{#each PROTECTED_ROUTES}}
- **{{path}}** - {{description}}
  - Component: `{{component}}`
  - Layout: `{{layout}}`
  - Required Role: {{requiredRole}}
  - Redirect if unauthorized: {{unauthorizedRedirect}}
{{/each}}

### API Routes

{{#if HAS_API_ROUTES}}
{{#each API_ROUTES}}
- **{{path}}** - {{description}}
  - Handler: `{{handler}}`
  - Method: {{method}}
{{/each}}

{{/if}}

## Route Parameters

### Dynamic Routes

{{#each DYNAMIC_ROUTES}}
- **{{path}}** - {{description}}
  - Parameters: {{parameters}}
  - Example: `{{example}}`
{{/each}}

### Query Parameters

{{#each QUERY_PARAMS}}
- **{{name}}** - {{description}}
  - Type: {{type}}
  - Required: {{required}}
{{/each}}

## Navigation

### Navigation Components

{{#each NAVIGATION_COMPONENTS}}

#### {{name}}

**Location:** `{{path}}`
**Purpose:** {{purpose}}

**Props:**

{{#each props}}
- `{{name}}` ({{type}}) - {{description}}
{{/each}}

{{/each}}

### Programmatic Navigation

```jsx
{{PROGRAMMATIC_NAVIGATION_EXAMPLE}}
```

## Route Guards

### Authentication Guard

{{#if HAS_AUTH_GUARD}}
**Location:** `{{AUTH_GUARD_PATH}}`

Protected routes require authentication. Unauthenticated users are redirected to: `{{AUTH_REDIRECT}}`
{{/if}}

### Role-Based Access

{{#if HAS_ROLE_BASED_ACCESS}}
**Location:** `{{ROLE_GUARD_PATH}}`

Role-based access control is implemented for the following routes:

{{#each ROLE_PROTECTED_ROUTES}}
- **{{path}}** - Requires: {{requiredRoles}}
{{/each}}

{{/if}}

## Route Layouts

### Layout Components

{{#each LAYOUTS}}

#### {{name}}

**Location:** `{{path}}`
**Used by:** {{usedBy}}

**Features:**

{{#each features}}
- {{this}}
{{/each}}

{{/each}}

## Route Optimization

### Code Splitting

{{#if HAS_CODE_SPLITTING}}
Routes are dynamically imported for better performance:

```jsx
{{CODE_SPLITTING_EXAMPLE}}
```
{{/if}}

### Preloading

{{#if HAS_PRELOADING}}
Critical routes are preloaded:

{{#each PRELOADED_ROUTES}}
- {{path}}
{{/each}}

{{/if}}

## Error Handling

### Error Boundaries

{{#if HAS_ERROR_BOUNDARIES}}
**Location:** `{{ERROR_BOUNDARY_PATH}}`

Error boundaries catch routing errors and display fallback UI.
{{/if}}

### 404 Handling

**Not Found Component:** `{{NOT_FOUND_COMPONENT}}`
**Fallback Route:** `{{FALLBACK_ROUTE}}`

## Breadcrumbs

{{#if HAS_BREADCRUMBS}}
**Breadcrumb Component:** `{{BREADCRUMB_COMPONENT}}`

Breadcrumbs are automatically generated based on route hierarchy.
{{/if}}

## Testing Routes

### Route Testing

{{#if HAS_ROUTE_TESTS}}
**Test Location:** `{{ROUTE_TEST_PATH}}`

Test each route for:

- Proper rendering
- Access control
- Parameter handling
- Error scenarios

{{/if}}

---

*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*
