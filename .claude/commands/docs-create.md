# Create Documentation System

You are tasked with setting up a comprehensive documentation system for this project using the standard-docs template system.

## Your Task

1. **Analyze the project structure** to understand the codebase
2. **Detect project type** (backend-api, frontend-webapp, service, mobile-app, library, or monorepo)
3. **Set up documentation** using appropriate templates
4. **Create a CLAUDE.md** configuration file for future reference

## Step-by-Step Process

### 1. Project Analysis

- Scan the project structure to understand the codebase
- Identify the primary programming language and frameworks
- Determine the project type based on files and dependencies
- Look for existing documentation to understand current state
- **For Backend APIs**: Check for OpenAPI specs (`openapi.yaml`, `openapi.json`, `swagger.json`)

### 2. Documentation Setup

- Create a `docs/` directory if it doesn't exist
- Generate appropriate documentation files based on project type:
  - **All Projects**: project-overview.md, development-setup.md, key-components.md, deployment-guide.md, troubleshooting.md, recent-changes.md
  - **Backend API**: api-reference.md, database-schema.md
  - **Frontend Web App**: component-library.md, routing-guide.md
  - **Service**: service-architecture.md
  - **Mobile App**: platform-guides.md
  - **Library**: package-usage.md
  - **Monorepo**: workspace-structure.md, cross-package-dependencies.md

### 3. Content Generation

For each document:

- **Fill in project-specific information** based on your analysis
- **Extract relevant details** from the codebase (API endpoints, components, etc.)
- **For Backend APIs with OpenAPI**: Reference existing spec and generate complementary documentation
- **For Backend APIs without OpenAPI**: Create comprehensive manual API documentation
- **Provide meaningful content** rather than just placeholder text
- **Include examples** and practical information

### 4. CLAUDE.md Configuration and Deduplication

**IMPORTANT: Before creating or modifying CLAUDE.md, check for existing CLAUDE.md files:**

1. **Check for existing CLAUDE.md** in the project root
2. **If CLAUDE.md exists:**
   - Analyze the existing content for information that should be moved to documentation
   - Identify duplicated information that exists in both CLAUDE.md and the new docs
   - Display proposed changes to the user, explaining:
     - What information will be moved from CLAUDE.md to specific documentation files
     - What information will be replaced with references to the new documentation
     - Why these changes improve the documentation structure (avoiding duplication)
   - Ask the user for approval before proceeding with CLAUDE.md modifications
   - Wait for user confirmation before continuing

3. **Create or update CLAUDE.md** to:
   - Reference all generated documentation instead of duplicating information
   - Include only project-specific context that doesn't belong in standard docs
   - Provide development commands and workflows
   - Serve as a central reference that points to the appropriate documentation files
   - Avoid duplicating information that now exists in the docs/ directory

## Key Information to Extract

### For All Projects

- Project name and description
- Main dependencies and frameworks
- Development setup requirements
- Build and deployment processes
- Key directories and file structure

### For Backend APIs

- API endpoints and routes (check OpenAPI spec if available)
- Database schema and models
- Authentication methods
- External integrations
- OpenAPI specification analysis (if `openapi.yaml`, `openapi.json`, or `swagger.json` exists)

### For Frontend Apps

- Component structure and patterns
- Routing configuration
- State management approach
- Build and deployment process

### For Services

- Service architecture and patterns
- Message processing capabilities
- Configuration and environment variables
- Health checks and monitoring

### For Mobile Apps

- Platform-specific configurations
- Navigation structure
- Push notification setup
- Platform deployment processes

### For Libraries

- Installation methods
- API surface and usage examples
- TypeScript definitions
- Browser/Node.js compatibility

### For Monorepos

- Workspace structure and organization
- Package dependencies and relationships
- Build and deployment strategies
- Cross-package development workflows

## Quality Guidelines

1. **Be Comprehensive** - Include all relevant information found in the codebase
2. **Be Accurate** - Only include information you can verify from the code
3. **Be Practical** - Focus on information developers will actually use
4. **Be Current** - Base information on the current state of the codebase
5. **Be Consistent** - Use consistent formatting and terminology

## Output Format

After completing the documentation setup, provide:

1. A summary of what was generated
2. Key insights about the project structure
3. Recommendations for maintaining the documentation
4. Next steps for the development team

Remember: The goal is to create documentation that will be genuinely useful for developers working on this project, not just template files with placeholder content.
