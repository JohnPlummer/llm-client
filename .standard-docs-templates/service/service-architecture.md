# {{PROJECT_NAME}} - Service Architecture

## Service Overview

**Service Type:** {{SERVICE_TYPE}}
**Runtime:** {{RUNTIME}}
**Framework:** {{FRAMEWORK}}

## Architecture Diagram

<!-- Include service architecture diagram here -->

## Core Components

### Service Entry Point

**Location:** `{{ENTRY_POINT}}`
**Purpose:** {{ENTRY_POINT_PURPOSE}}

### Business Logic

{{#each BUSINESS_LOGIC_COMPONENTS}}
#### {{name}}

**Location:** `{{path}}`
**Purpose:** {{purpose}}
**Dependencies:** {{dependencies}}
{{/each}}

### Data Access Layer

{{#each DATA_ACCESS_COMPONENTS}}
#### {{name}}

**Location:** `{{path}}`
**Purpose:** {{purpose}}
**Data Source:** {{dataSource}}
{{/each}}

### External Integrations

{{#each EXTERNAL_INTEGRATIONS}}
#### {{name}}

**Type:** {{type}}
**Endpoint:** {{endpoint}}
**Authentication:** {{authentication}}
**Purpose:** {{purpose}}
{{/each}}

## Service Patterns

### Message Processing

{{#if HAS_MESSAGE_PROCESSING}}
**Message Broker:** {{MESSAGE_BROKER}}
**Queue/Topic:** {{QUEUE_TOPIC}}

**Message Handlers:**

{{#each MESSAGE_HANDLERS}}
- **{{name}}:** {{description}}
{{/each}}

{{/if}}

### Scheduled Tasks

{{#if HAS_SCHEDULED_TASKS}}
**Scheduler:** {{SCHEDULER}}

**Tasks:**

{{#each SCHEDULED_TASKS}}
- **{{name}}:** {{schedule}} - {{description}}
{{/each}}

{{/if}}

### Background Processing

{{#if HAS_BACKGROUND_PROCESSING}}
**Worker Pool:** {{WORKER_POOL}}
**Job Queue:** {{JOB_QUEUE}}

**Job Types:**

{{#each JOB_TYPES}}
- **{{name}}:** {{description}}
{{/each}}

{{/if}}

## Configuration

### Environment Variables

{{#each ENVIRONMENT_VARIABLES}}
- `{{name}}` - {{description}} {{#if required}}(Required){{/if}}
{{/each}}

### Configuration Files

{{#each CONFIG_FILES}}
- **{{name}}:** `{{path}}` - {{description}}
{{/each}}

## Error Handling

### Error Categories

{{#each ERROR_CATEGORIES}}
#### {{name}}

**Handling Strategy:** {{strategy}}
**Retry Policy:** {{retryPolicy}}
**Notification:** {{notification}}
{{/each}}

### Circuit Breaker

{{#if HAS_CIRCUIT_BREAKER}}
**Pattern:** {{CIRCUIT_BREAKER_PATTERN}}
**Failure Threshold:** {{FAILURE_THRESHOLD}}
**Recovery Time:** {{RECOVERY_TIME}}

{{/if}}

## Monitoring & Observability

### Metrics

{{#each METRICS}}
- **{{name}}:** {{description}} ({{type}})
{{/each}}

### Health Checks

{{#if HAS_HEALTH_CHECKS}}
**Endpoint:** {{HEALTH_CHECK_ENDPOINT}}
**Checks:**

{{#each HEALTH_CHECKS}}
- {{name}}: {{description}}
{{/each}}

{{/if}}

### Logging

**Log Level:** {{LOG_LEVEL}}
**Log Format:** {{LOG_FORMAT}}
**Log Destination:** {{LOG_DESTINATION}}

**Structured Logging Fields:**

{{#each LOG_FIELDS}}
- {{name}}: {{description}}
{{/each}}

### Distributed Tracing

{{#if HAS_DISTRIBUTED_TRACING}}
**Tracing System:** {{TRACING_SYSTEM}}
**Trace ID Header:** {{TRACE_ID_HEADER}}

{{/if}}

## Performance Characteristics

### Resource Requirements

- **CPU:** {{CPU_REQUIREMENTS}}
- **Memory:** {{MEMORY_REQUIREMENTS}}
- **Storage:** {{STORAGE_REQUIREMENTS}}

### Throughput

- **Expected Load:** {{EXPECTED_LOAD}}
- **Peak Load:** {{PEAK_LOAD}}
- **Response Time:** {{RESPONSE_TIME}}

### Scaling

{{#if HAS_AUTO_SCALING}}
**Auto-scaling:** {{AUTO_SCALING_CONFIG}}
**Scale Triggers:** {{SCALE_TRIGGERS}}

{{/if}}

## Security

### Authentication

{{#if HAS_AUTHENTICATION}}
**Method:** {{AUTH_METHOD}}
**Token Validation:** {{TOKEN_VALIDATION}}

{{/if}}

### Authorization

{{#if HAS_AUTHORIZATION}}
**Model:** {{AUTHORIZATION_MODEL}}
**Policies:** {{AUTHORIZATION_POLICIES}}

{{/if}}

### Data Protection

{{#each DATA_PROTECTION}}
- **{{type}}:** {{description}}
{{/each}}

## Dependencies

### Internal Dependencies

{{#each INTERNAL_DEPENDENCIES}}
- **{{name}}:** {{description}}
{{/each}}

### External Dependencies

{{#each EXTERNAL_DEPENDENCIES}}
- **{{name}}:** {{description}} ({{criticality}})
{{/each}}

---
*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*

