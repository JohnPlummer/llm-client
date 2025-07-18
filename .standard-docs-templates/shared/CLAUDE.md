# Claude Configuration for {{PROJECT_NAME}}

## Project Context

{{PROJECT_NAME}} is a {{PROJECT_TYPE}} project built with {{FRAMEWORKS}}.

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

## AI Assistant Guidelines

### Before Making Any Changes

1. **Read project-overview.md** to understand the project's purpose and architecture
2. **Check development-setup.md** for project-specific setup and commands
3. **Review key-components.md** to understand how components relate to each other
4. **Look at recent-changes.md** to understand recent project evolution

### When Adding Features

1. **Follow established patterns** found in key-components.md
2. **Update relevant documentation** - don't leave docs out of sync
3. **Check deployment-guide.md** if changes affect deployment
4. **Add any new issues to troubleshooting.md** if they might affect others

### When Fixing Bugs

1. **Check troubleshooting.md first** - solution might already be documented
2. **Understand root cause** using project architecture from project-overview.md
3. **Update troubleshooting.md** with new solutions for future reference
4. **Update recent-changes.md** with the fix details

### When Updating Documentation

1. **Keep information in the right place** - don't duplicate across files
2. **Update recent-changes.md** when making significant doc changes
3. **Cross-reference related information** instead of copying it
4. **Maintain consistency** with existing documentation style and structure

## Project-Specific Information

### Quick Reference

- **Primary Language**: {{PRIMARY_LANGUAGE}}
- **Main Framework**: {{FRAMEWORK}}
- **Install Command**: {{INSTALL_COMMAND}}
- **Start Command**: {{START_COMMAND}}
- **Test Command**: {{TEST_COMMAND}}
- **Build Command**: {{BUILD_COMMAND}}

### Key Directories

- **{{MAIN_SOURCE_DIR}}** - Main source code
- **{{TEST_DIR}}** - Test files  
- **{{CONFIG_DIR}}** - Configuration files
- **docs/** - All project documentation

### Important Notes for AI Assistance

- **Always verify** information against current codebase - don't rely on assumptions
- **Refer to documentation** rather than duplicating information
- **Update documentation** when making changes that affect setup, architecture, or usage
- **Follow established patterns** shown in existing code and documented in key-components.md
- **Consider security implications** for all changes, especially in production systems

### Common Patterns to Follow

(Reference key-components.md for detailed patterns)

- **Code Organization**: {{CODE_ORGANIZATION_PATTERN}}
- **Error Handling**: {{ERROR_HANDLING_PATTERN}}
- **Testing Approach**: {{TESTING_APPROACH}}
- **Configuration Management**: {{CONFIG_MANAGEMENT_PATTERN}}

---

## Documentation Update Guidelines

This file serves as a navigation hub for Claude. For detailed project information, refer to the specific documentation files listed above.

---
*Last updated: {{TIMESTAMP}}*

