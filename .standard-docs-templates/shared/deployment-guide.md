# {{PROJECT_NAME}} - Deployment Guide

## Deployment Overview
{{PROJECT_TYPE}} deployment guide for {{PROJECT_NAME}}.

## Environments
{{#each ENVIRONMENTS}}
### {{name}}
- **URL:** {{url}}
- **Branch:** {{branch}}
- **Auto-deploy:** {{autoDeploy}}
- **Description:** {{description}}
{{/each}}

## Prerequisites
{{#each DEPLOYMENT_PREREQUISITES}}
- {{this}}
{{/each}}

## Build Process
{{#if BUILD_STEPS}}
{{#each BUILD_STEPS}}
### {{step}}. {{title}}
```bash
{{command}}
```
{{#if description}}
{{description}}
{{/if}}
{{/each}}
{{/if}}

## Deployment Steps

### Manual Deployment
{{#each MANUAL_DEPLOYMENT_STEPS}}
{{step}}. {{description}}
{{#if command}}
```bash
{{command}}
```
{{/if}}
{{/each}}

### Automated Deployment (CI/CD)
{{#if HAS_CI_CD}}
This project uses {{CI_CD_PLATFORM}} for automated deployments.

**Pipeline Configuration:** `{{CI_CD_CONFIG_FILE}}`

**Deployment Triggers:**
{{#each DEPLOYMENT_TRIGGERS}}
- {{trigger}}: {{description}}
{{/each}}
{{/if}}

## Environment Variables
{{#each DEPLOYMENT_ENV_VARIABLES}}
- `{{name}}`: {{description}} {{#if required}}(Required){{/if}}
{{/each}}

## Health Checks
{{#if HAS_HEALTH_CHECKS}}
- **Health Check URL:** {{HEALTH_CHECK_URL}}
- **Expected Response:** {{HEALTH_CHECK_RESPONSE}}
{{/if}}

## Monitoring & Logging
{{#if HAS_MONITORING}}
- **Monitoring Dashboard:** {{MONITORING_URL}}
- **Log Aggregation:** {{LOG_PLATFORM}}
{{/if}}

## Rollback Procedures
{{#each ROLLBACK_STEPS}}
{{step}}. {{description}}
{{#if command}}
```bash
{{command}}
```
{{/if}}
{{/each}}

## Troubleshooting Deployment Issues
Common deployment problems and solutions can be found in [troubleshooting.md](./troubleshooting.md#deployment-issues).

---
*Last updated: {{TIMESTAMP}}*
*Generated with [standard-docs](https://github.com/johnplummer/standard-docs)*
