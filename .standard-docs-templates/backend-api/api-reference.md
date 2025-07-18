# {{PROJECT_NAME}} - API Reference

{{#if HAS_OPENAPI_SPEC}}

## OpenAPI Specification

This API is documented using OpenAPI 3.0. The complete specification is available at:

- **Spec File**: `{{OPENAPI_SPEC_PATH}}`
- **Interactive Documentation**: {{#if SWAGGER_UI_URL}}[Swagger UI]({{SWAGGER_UI_URL}}){{/if}}{{#if REDOC_URL}}[ReDoc]({{REDOC_URL}}){{/if}}

### Quick Reference

{{#if SWAGGER_UI_URL}}
**🔗 [View Interactive API Documentation]({{SWAGGER_UI_URL}})**
{{/if}}

The OpenAPI specification includes:

- Complete endpoint documentation with request/response schemas
- Authentication and authorization details
- Error response formats
- Request/response examples
- Model definitions and validation rules

### Using the OpenAPI Spec

```bash
# View the raw specification
cat {{OPENAPI_SPEC_PATH}}

# Generate client SDKs
npx @openapitools/openapi-generator-cli generate -i {{OPENAPI_SPEC_PATH}} -g typescript-fetch -o ./client

# Validate the specification
npx swagger-parser validate {{OPENAPI_SPEC_PATH}}
```

### API Overview

{{else}}

## API Overview

This API documentation provides manual documentation for all endpoints. Consider adopting OpenAPI 3.0 for better tooling support.

### Base URL

{{#each ENVIRONMENTS}}

- **{{name}}:** {{baseUrl}}
{{/each}}

### Authentication

{{#if HAS_AUTH}}
{{AUTH_TYPE}} authentication is required for most endpoints.

{{#if AUTH_TYPE === 'Bearer'}}
Include the token in the Authorization header:

```text
Authorization: Bearer <your-token>
```

{{/if}}

{{#if AUTH_TYPE === 'API Key'}}
Include the API key in the header:

```text
{{API_KEY_HEADER}}: <your-api-key>
```

{{/if}}
{{/if}}

### API Endpoints

{{#each ENDPOINTS}}

#### {{method}} {{path}}

{{#if description}}
{{description}}
{{/if}}

{{#if auth}}
**Authentication:** {{auth}}
{{/if}}

**Parameters:**

{{#each parameters}}

- `{{name}}` ({{type}}) {{#if required}}**Required**{{/if}} - {{description}}
{{/each}}

**Request Body:**

{{#if requestBody}}

```json
{{requestBody}}
```

{{/if}}

**Response:**

{{#each responses}}
**{{status}}** - {{description}}

```json
{{example}}
```

{{/each}}

**Example Request:**

```bash
curl -X {{method}} "{{../BASE_URL}}{{path}}" \
{{#if ../HAS_AUTH}}  -H "Authorization: Bearer <token>" \{{/if}}
{{#if requestBody}}  -H "Content-Type: application/json" \
  -d '{{requestBody}}'{{/if}}
```

---

{{/each}}

### Data Models

{{#each MODELS}}

#### {{name}}

{{#if description}}
{{description}}
{{/if}}

```json
{{schema}}
```

**Properties:**

{{#each properties}}

- `{{name}}` ({{type}}) {{#if required}}**Required**{{/if}} - {{description}}
{{/each}}

---

{{/each}}
{{/if}}

## Common Response Formats

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data
  },
  "message": "Success message",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {}
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Common Error Codes

{{#each ERROR_CODES}}

- `{{code}}` ({{status}}) - {{description}}
{{/each}}

## Rate Limiting

{{#if HAS_RATE_LIMITING}}

- **Rate limit:** {{RATE_LIMIT}}
- **Rate limit headers:**
  - `X-RateLimit-Limit`: Total requests allowed
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Reset time (Unix timestamp)

{{/if}}

## Pagination

{{#if HAS_PAGINATION}}
Large result sets are paginated using {{PAGINATION_TYPE}}.

**Query Parameters:**

{{#each PAGINATION_PARAMS}}

- `{{name}}` - {{description}}
{{/each}}

**Response Headers:**

{{#each PAGINATION_HEADERS}}

- `{{name}}` - {{description}}
{{/each}}

{{/if}}

## Webhooks

{{#if HAS_WEBHOOKS}}

{{#each WEBHOOKS}}

### {{name}}

**Event:** {{event}}
**URL:** Your configured webhook URL
**Method:** POST

**Payload:**

```json
{{payload}}
```

{{/each}}

{{/if}}

## SDK & Libraries

{{#each SDKS}}

- **{{language}}:** {{#if url}}[{{name}}]({{url}}){{else}}{{name}}{{/if}}
{{/each}}

{{#if HAS_OPENAPI_SPEC}}

## Updating the OpenAPI Spec

When making API changes:

1. **Update the spec first** (`{{OPENAPI_SPEC_PATH}}`)
2. **Validate the spec**: `npx swagger-parser validate {{OPENAPI_SPEC_PATH}}`
3. **Regenerate documentation** if using automated tools
4. **Update client SDKs** if using code generation
5. **Test the changes** with the updated spec

### OpenAPI Tools

- **Swagger Editor**: Edit and validate specs
- **Swagger UI**: Generate interactive documentation
- **ReDoc**: Alternative documentation generator
- **OpenAPI Generator**: Generate client SDKs and server stubs
- **Prism**: Mock server from OpenAPI specs

{{/if}}

---

*Last updated: {{TIMESTAMP}}*
{{#if HAS_OPENAPI_SPEC}}
*API specification: {{OPENAPI_SPEC_PATH}}*
{{/if}}
