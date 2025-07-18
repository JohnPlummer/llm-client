# {{PROJECT_NAME}} - Component Library

## Overview

This document describes the reusable components available in the {{PROJECT_NAME}} application.

## Component Categories

### UI Components

{{#each UI_COMPONENTS}}

#### {{name}}

**Location:** `{{path}}`
**Purpose:** {{purpose}}

**Props:**

{{#each props}}
- `{{name}}` ({{type}}) {{#if required}}**Required**{{/if}} - {{description}}
{{/each}}

**Usage:**

```jsx
{{usage}}
```

**Styling:**

- CSS Class: `{{cssClass}}`
- Variants: {{variants}}

---

{{/each}}

### Form Components

{{#each FORM_COMPONENTS}}

#### {{name}}

**Location:** `{{path}}`
**Purpose:** {{purpose}}

**Props:**

{{#each props}}
- `{{name}}` ({{type}}) {{#if required}}**Required**{{/if}} - {{description}}
{{/each}}

**Validation:**

{{#each validation}}
- {{rule}}: {{message}}
{{/each}}

**Usage:**

```jsx
{{usage}}
```

---

{{/each}}

### Layout Components

{{#each LAYOUT_COMPONENTS}}

#### {{name}}

**Location:** `{{path}}`
**Purpose:** {{purpose}}

**Props:**

{{#each props}}
- `{{name}}` ({{type}}) {{#if required}}**Required**{{/if}} - {{description}}
{{/each}}

**Responsive Behavior:**

{{responsiveBehavior}}

**Usage:**

```jsx
{{usage}}
```

---

{{/each}}

## Design System

### Colors

{{#each COLORS}}
- **{{name}}:** {{value}} - {{usage}}
{{/each}}

### Typography

{{#each TYPOGRAPHY}}
- **{{name}}:** {{styles}} - {{usage}}
{{/each}}

### Spacing

{{#each SPACING}}
- **{{name}}:** {{value}} - {{usage}}
{{/each}}

### Breakpoints

{{#each BREAKPOINTS}}
- **{{name}}:** {{value}} - {{description}}
{{/each}}

## Component Guidelines

### Naming Conventions

{{NAMING_CONVENTIONS}}

### Props Best Practices

{{PROPS_BEST_PRACTICES}}

### Styling Approach

{{STYLING_APPROACH}}

### Testing Components

{{#if HAS_COMPONENT_TESTS}}
Components should include:

- Unit tests for logic
- Snapshot tests for rendering
- Accessibility tests

**Test Location:** `{{TEST_LOCATION}}`
**Test Command:** `{{TEST_COMMAND}}`
{{/if}}

## Storybook

{{#if HAS_STORYBOOK}}
Interactive component documentation is available in Storybook.

**Start Storybook:**

```bash
{{STORYBOOK_COMMAND}}
```

**Storybook URL:** {{STORYBOOK_URL}}
{{/if}}

## Creating New Components

### Component Template

```jsx
{{COMPONENT_TEMPLATE}}
```

### Component Checklist

- [ ] Props are properly typed
- [ ] Component is responsive
- [ ] Accessibility standards met
- [ ] Unit tests written
- [ ] Storybook story created
- [ ] Documentation updated

---

*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*
