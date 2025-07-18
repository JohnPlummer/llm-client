# Claude Configuration for {{PROJECT_NAME}} (Backend API)

## Project Context

{{PROJECT_NAME}} is a backend API project built with {{FRAMEWORKS}}.

{{PROJECT_DESCRIPTION}}

## Documentation Navigation

This project uses a structured documentation system. When helping with this project, **always refer to the relevant documentation files** rather than making assumptions.

### Core Documentation Files

- **[docs/project-overview.md](docs/project-overview.md)** - **Start here** for project architecture, goals, and technical overview
- **[docs/development-setup.md](docs/development-setup.md)** - **Essential** for environment setup, dependencies, and getting started
- **[docs/key-components.md](docs/key-components.md)** - **Reference** for main modules, their responsibilities, and relationships
- **[docs/deployment-guide.md](docs/deployment-guide.md)** - **Use** for deployment processes, environments, and procedures
- **[docs/troubleshooting.md](docs/troubleshooting.md)** - **Check first** for common issues and debugging approaches
- **[docs/recent-changes.md](docs/recent-changes.md)** - **Review** for recent updates and changelog

### API-Specific Documentation

- **[docs/api-reference.md](docs/api-reference.md)** - **Complete API documentation** with endpoints, parameters, and examples
- **[docs/database-schema.md](docs/database-schema.md)** - **Database structure** and relationships
- **OpenAPI Specification** - Look for `openapi.yaml`, `openapi.json`, or `swagger.json` in the project root or `docs/` directory

## AI Assistant Guidelines for API Development

### Before Making API Changes

1. **Check for OpenAPI spec** - Look for `openapi.yaml`, `openapi.json`, or `swagger.json` files
2. **Read api-reference.md** to understand existing endpoints and patterns
3. **Check database-schema.md** for data relationships and constraints
4. **Review troubleshooting.md** for common API issues and solutions
5. **Understand security patterns** from existing endpoint implementations

### When Adding New Endpoints

1. **Follow RESTful principles** and existing URL patterns shown in api-reference.md
2. **Implement proper authentication** - check existing patterns in codebase
3. **Add comprehensive validation** for all inputs
4. **Update API documentation**:
   - **If using OpenAPI**: Update the OpenAPI spec file first, then validate it
   - **If using manual docs**: Update api-reference.md with new endpoint documentation
5. **Add appropriate error handling** and return consistent error formats
6. **Consider rate limiting** and security implications

### When Modifying Database Schema

1. **Create migrations** - never modify schema directly
2. **Update database-schema.md** with schema changes
3. **Consider backwards compatibility** for existing API endpoints
4. **Test data migration** thoroughly
5. **Update API documentation**:
   - **If using OpenAPI**: Update response schemas in the spec
   - **If using manual docs**: Update api-reference.md with new response formats

### When Debugging API Issues

1. **Check troubleshooting.md first** - solution might already be documented
2. **Verify authentication** and authorization logic
3. **Check database connections** and query performance
4. **Review error logs** for specific error patterns
5. **Test with different input scenarios**

### API Development Best Practices

- **RESTful Design**: Use appropriate HTTP methods and status codes
- **Consistent Responses**: Follow established response format patterns
- **Input Validation**: Validate all inputs before processing
- **Error Handling**: Return meaningful error messages with appropriate HTTP status codes
- **Security**: Implement proper authentication, authorization, and input sanitization
- **Documentation**: Keep API documentation updated with all endpoint changes
  - **OpenAPI projects**: Update spec first, then validate and regenerate docs
  - **Manual projects**: Update api-reference.md directly
- **Testing**: Write comprehensive tests for all endpoints

### Common API Patterns to Follow

(Reference api-reference.md and existing code for examples)

- **Authentication**: {{AUTH_PATTERN}}
- **Error Responses**: {{ERROR_RESPONSE_PATTERN}}
- **Pagination**: {{PAGINATION_PATTERN}}
- **Input Validation**: {{VALIDATION_PATTERN}}
- **Database Queries**: {{DB_QUERY_PATTERN}}

### Security Considerations for APIs

- **Authentication**: Verify user identity on all protected endpoints
- **Authorization**: Check permissions before allowing data access
- **Input Sanitization**: Sanitize all user inputs to prevent injection attacks
- **Rate Limiting**: Implement appropriate rate limits to prevent abuse
- **HTTPS**: Ensure all API communication uses HTTPS in production
- **Logging**: Log security events without exposing sensitive data

### Database Best Practices

- **Migrations**: Use migrations for all schema changes
- **Transactions**: Use database transactions for complex operations
- **Indexing**: Consider query performance when designing database queries
- **Connection Pooling**: Use connection pooling for database connections
- **Backup Strategy**: Ensure database backup procedures are in place

## Project-Specific Information

### Quick Reference

- **Primary Language**: {{PRIMARY_LANGUAGE}}
- **API Framework**: {{API_FRAMEWORK}}
- **Database**: {{DATABASE_TYPE}}
- **Authentication**: {{AUTH_TYPE}}
- **Install Command**: {{INSTALL_COMMAND}}
- **Start Command**: {{START_COMMAND}}
- **Test Command**: {{TEST_COMMAND}}
- **Migration Command**: {{MIGRATION_COMMAND}}

### Key Directories

- **{{API_DIR}}** - API route handlers
- **{{MODELS_DIR}}** - Data models
- **{{MIDDLEWARE_DIR}}** - Middleware functions
- **{{TESTS_DIR}}** - Test files
- **{{CONFIG_DIR}}** - Configuration files
- **{{MIGRATIONS_DIR}}** - Database migrations

### Important Notes for AI Assistance

- **Check for OpenAPI spec first** - Look for `openapi.yaml`, `openapi.json`, or `swagger.json`
- **Always check api-reference.md** for endpoint documentation before making changes
- **Refer to database-schema.md** for data relationships and constraints
- **Follow existing authentication patterns** shown in the codebase
- **Update documentation** when making API changes (OpenAPI spec or manual docs)
- **Consider security implications** for all API modifications
- **Test thoroughly** including edge cases and error scenarios

## OpenAPI Workflow Support

### If Project Uses OpenAPI

**Common file locations**: `openapi.yaml`, `openapi.json`, `swagger.json`, `docs/openapi.yaml`

**When making API changes**:

1. **Update OpenAPI spec first** - This is the source of truth
2. **Validate the spec**: Use tools like `swagger-parser` or `redocly lint`
3. **Regenerate documentation** if using automated tools (Swagger UI, ReDoc)
4. **Update client SDKs** if using code generation
5. **Test with updated spec** using tools like Postman or Prism mock server

**Common OpenAPI tools to be aware of**:

- **Swagger Editor**: For editing and validating specs
- **Swagger UI**: For interactive documentation
- **ReDoc**: Alternative documentation generator
- **OpenAPI Generator**: For client SDK generation
- **Prism**: For mock servers and contract testing

### If Project Uses Manual Documentation

**When making API changes**:

1. **Update api-reference.md** with new endpoint details
2. **Include request/response examples** for clarity
3. **Update error response documentation** as needed
4. **Consider migrating to OpenAPI** for better tooling support

### Migration Considerations

- **OpenAPI offers better tooling**: Code generation, validation, mock servers
- **Manual docs offer more flexibility**: Custom formatting, complex explanations
- **Hybrid approach**: OpenAPI spec + additional manual documentation for complex topics

---

*This file serves as a navigation hub for Claude working on API development. For detailed API information, refer to api-reference.md, database-schema.md, and any OpenAPI specifications.*

*Last updated: {{TIMESTAMP}}*
