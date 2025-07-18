# {{PROJECT_NAME}} - Database Schema

## Database Overview

- **Database Type:** {{DATABASE_TYPE}}
- **Version:** {{DATABASE_VERSION}}
- **Connection:** {{DATABASE_CONNECTION_INFO}}

## Schema Diagram

<!-- Include database schema diagram here -->

## Tables

{{#each TABLES}}

### {{name}}

{{#if description}}
{{description}}
{{/if}}

**Columns:**

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
{{#each columns}}
| {{name}} | {{type}} | {{constraints}} | {{description}} |
{{/each}}

**Indexes:**

{{#each indexes}}
- `{{name}}` on {{columns}} {{#if unique}}(UNIQUE){{/if}}
{{/each}}

**Relationships:**

{{#each relationships}}
- {{type}} to `{{table}}` via {{columns}}
{{/each}}

---

{{/each}}

## Migrations

### Migration History

{{#each MIGRATIONS}}
- **{{version}}** ({{date}}) - {{description}}
{{/each}}

### Running Migrations

```bash
{{MIGRATION_COMMAND}}
```

### Creating New Migrations

```bash
{{CREATE_MIGRATION_COMMAND}}
```

## Seed Data

{{#if HAS_SEED_DATA}}

```bash
{{SEED_COMMAND}}
```

{{/if}}

## Backup & Restore

{{#if HAS_BACKUP_PROCEDURES}}

### Backup

```bash
{{BACKUP_COMMAND}}
```

### Restore

```bash
{{RESTORE_COMMAND}}
```

{{/if}}

## Performance Considerations

{{#each PERFORMANCE_NOTES}}
- **{{title}}:** {{description}}
{{/each}}

## Data Access Patterns

{{#each ACCESS_PATTERNS}}

### {{name}}

{{description}}

**Queries:**

{{#each queries}}

```sql
{{query}}
```

{{/each}}
{{/each}}

---

*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*
