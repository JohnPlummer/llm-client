# Analyze Project Documentation

You are tasked with analyzing the current project documentation to assess its quality, completeness, and accuracy.

## Your Task

1. **Review existing documentation** for completeness and accuracy
2. **Identify gaps and inconsistencies** in the current documentation
3. **Assess documentation quality** against best practices
4. **Provide recommendations** for improvements

## Step-by-Step Process

### 1. Documentation Inventory

- List all existing documentation files
- Categorize documents by purpose (setup, reference, guides, etc.)
- Note last update dates and freshness
- Identify missing standard documentation
- **Check for CLAUDE.md file** and analyze for duplicated information with docs/

### 2. Content Analysis

For each document:

- **Accuracy**: Verify information against current codebase
- **Completeness**: Check if all relevant information is covered
- **Clarity**: Assess if content is clear and well-structured
- **Usefulness**: Evaluate practical value for developers
- **Currency**: Check if information is up-to-date

### 3. Gap Analysis

Identify missing documentation:

- **Essential docs**: PROJECT_OVERVIEW, DEVELOPMENT_SETUP, etc.
- **Project-specific docs**: API_REFERENCE, COMPONENT_LIBRARY, etc.
- **Process docs**: Deployment, troubleshooting, contributing
- **Technical docs**: Architecture, database schema, etc.

### 4. CLAUDE.md Deduplication Analysis

**If CLAUDE.md exists:**

- **Identify duplicated information** between CLAUDE.md and docs/ directory
- **Analyze information placement** - what belongs in CLAUDE.md vs docs/
- **Assess reference structure** - does CLAUDE.md properly reference docs instead of duplicating?
- **Recommend consolidation** - suggest moving duplicated content to appropriate docs and updating CLAUDE.md to reference them

### 5. Quality Assessment

Rate documentation on:

- **Completeness** (1-10): How much of the project is documented
- **Accuracy** (1-10): How well docs match current codebase
- **Usability** (1-10): How helpful docs are for developers
- **Maintenance** (1-10): How well docs are maintained

## Analysis Framework

### Document Categories to Review

#### Essential Documentation

- **project-overview.md**: Project description, goals, architecture
- **development-setup.md**: Setup instructions, dependencies
- **key-components.md**: Main modules and their purposes
- **deployment-guide.md**: Deployment processes and environments
- **troubleshooting.md**: Common issues and solutions
- **recent-changes.md**: Recent updates and changelog

#### Project-Specific Documentation

- **Backend API**: api-reference.md, database-schema.md (also check for OpenAPI specs)
- **Frontend**: component-library.md, routing-guide.md
- **Service**: service-architecture.md
- **Mobile**: platform-guides.md
- **Library**: package-usage.md
- **Monorepo**: workspace-structure.md, cross-package-dependencies.md

#### Process Documentation

- **README.md**: Project introduction and quick start
- **CONTRIBUTING.md**: Contribution guidelines
- **CODE_OF_CONDUCT.md**: Community standards
- **CHANGELOG.md**: Version history and changes
- **SECURITY.md**: Security policies and reporting
- **CLAUDE.md**: AI assistant context and project reference (check for duplication with docs/)

### Quality Metrics

#### Content Quality

- **Accuracy**: Information matches current codebase
- **Completeness**: All relevant topics are covered
- **Clarity**: Content is clear and well-organized
- **Examples**: Practical examples and code snippets
- **Consistency**: Consistent style and formatting

#### Technical Quality

- **Code Examples**: Working, tested examples
- **Links**: Internal and external links work
- **Formatting**: Proper markdown formatting
- **Structure**: Logical organization and flow
- **Searchability**: Easy to find information

#### Maintenance Quality

- **Freshness**: Recently updated content
- **Ownership**: Clear maintainers/owners
- **Process**: Documentation update process
- **Automation**: Automated checks and updates
- **Version Control**: Proper git history

## Specific Checks to Perform

### For All Documentation

- Last updated timestamps
- Broken internal/external links
- Outdated information (versions, URLs, etc.)
- Missing sections or incomplete content
- Inconsistent formatting or style

### For Setup Documentation

- Installation instructions work on clean environment
- Dependencies are current and available
- Environment setup steps are complete
- Common setup issues are addressed

### For API Documentation

- All endpoints are documented (in api-reference.md or OpenAPI spec)
- Request/response examples are current
- Authentication methods are explained
- Error codes and responses are covered
- OpenAPI spec is valid and up-to-date (if present)

### For Component Documentation

- All major components are documented
- Usage examples are provided
- Props/parameters are explained
- Component relationships are clear

### For Deployment Documentation

- All deployment targets are covered
- Environment variables are documented
- Deployment steps are complete and current
- Rollback procedures are explained

## Analysis Output

### Documentation Report

1. **Executive Summary**: Overall state of documentation
2. **Quality Scores**: Ratings for each category
3. **Gap Analysis**: Missing or incomplete documentation
4. **Issues Found**: Specific problems identified
5. **Recommendations**: Prioritized improvement suggestions

### Detailed Findings

For each document reviewed:

- **Status**: Complete, Partial, Missing, Outdated
- **Issues**: Specific problems found
- **Recommendations**: Suggested improvements
- **Priority**: High, Medium, Low

### Action Plan

1. **Immediate Actions**: Critical fixes needed
2. **Short-term Goals**: Improvements for next sprint
3. **Long-term Strategy**: Sustainable documentation process
4. **Resource Requirements**: Time and people needed

## Common Documentation Issues

### Content Issues

- Outdated information
- Missing examples
- Incomplete procedures
- Unclear explanations
- Inconsistent terminology

### Technical Issues

- Broken links
- Incorrect code examples
- Missing dependencies
- Wrong file paths
- Outdated screenshots

### Process Issues

- No update process
- Unclear ownership
- No review process
- No automation
- Poor discoverability

## Best Practices Assessment

### Structure

- Clear navigation and organization
- Logical information hierarchy
- Consistent formatting
- Table of contents where needed
- Cross-references between documents

### Content

- Clear, concise writing
- Practical examples
- Step-by-step instructions
- Troubleshooting sections
- Regular updates

### Maintenance

- Clear ownership
- Regular review cycles
- Automated testing
- Version control
- Feedback mechanisms

## Recommendations Framework

### Priority Levels

1. **Critical**: Blocking developer productivity
2. **Important**: Significant impact on efficiency
3. **Nice to have**: Quality of life improvements

### Improvement Types

- **Create**: New documentation needed
- **Update**: Existing docs need refresh
- **Fix**: Specific issues to resolve
- **Enhance**: Quality improvements
- **Automate**: Process improvements
- **Deduplicate**: Remove duplicated information between CLAUDE.md and docs/

Remember: The goal is to provide actionable insights that will help improve the documentation's value to the development team.
