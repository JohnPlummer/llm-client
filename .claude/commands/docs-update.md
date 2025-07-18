# Update Documentation from Git History

You are tasked with updating project documentation based on recent git history and changes since the last documentation update.

## Your Task

1. **Analyze git history** to identify significant changes
2. **Update relevant documentation** based on changes found
3. **Refresh the recent-changes.md** document
4. **Ensure documentation accuracy** with current codebase

## Step-by-Step Process

### 1. Git History Analysis

- Check git log since last documentation update (look for timestamp in docs)
- Identify significant commits, PRs, and changes
- Categorize changes by type (features, bug fixes, refactoring, etc.)
- Note any breaking changes or major architectural updates

### 2. Documentation Impact Assessment

For each significant change, determine which documentation needs updates:

- **New features** → Update project-overview.md, key-components.md, relevant guides
- **API changes** → Update api-reference.md (or OpenAPI spec if present), examples
- **Configuration changes** → Update development-setup.md, deployment-guide.md
- **New dependencies** → Update development-setup.md, project-overview.md
- **Bug fixes** → Update troubleshooting.md if relevant
- **Performance improvements** → Update relevant architecture documents

### 3. Update Documentation Files

- **project-overview.md**: Update features, architecture, tech stack
- **development-setup.md**: Update dependencies, setup steps, commands
- **key-components.md**: Update component descriptions, new components
- **api-reference.md**: Update endpoints, parameters, examples (or reference OpenAPI spec if present)
- **deployment-guide.md**: Update deployment processes, environments
- **troubleshooting.md**: Add new known issues and solutions
- **Project-specific docs**: Update based on relevant changes

### 4. Update recent-changes.md

Create comprehensive changelog with:

- **Recent Updates**: Last 30 days of significant changes
- **New Features**: Added functionality with descriptions
- **Bug Fixes**: Resolved issues and their solutions
- **Breaking Changes**: Changes that affect existing functionality
- **Dependencies**: Updated packages and versions
- **Performance**: Improvements and optimizations

### 5. CLAUDE.md Deduplication Check

**IMPORTANT: Before updating documentation, check for CLAUDE.md deduplication:**

1. **Check for existing CLAUDE.md** in the project root
2. **If CLAUDE.md exists:**
   - Analyze the existing content for information that duplicates what's in the updated documentation
   - Identify outdated information in CLAUDE.md that should be removed or updated
   - Display proposed changes to the user, explaining:
     - What information will be removed from CLAUDE.md because it's now in documentation
     - What information will be replaced with references to the updated documentation
     - Why these changes improve the documentation structure (avoiding duplication)
   - Ask the user for approval before proceeding with CLAUDE.md modifications
   - Wait for user confirmation before continuing

3. **Update CLAUDE.md** to:
   - Reference updated documentation instead of duplicating information
   - Remove outdated information that's now covered in the docs/ directory
   - Update references to point to the correct documentation files
   - Maintain only project-specific context that doesn't belong in standard docs

## Information to Extract from Git

### Commit Analysis

- Commit messages and descriptions
- File changes and modifications
- Author and timestamp information
- Associated pull request data

### Change Categories

- **Features**: New functionality added
- **Fixes**: Bug fixes and corrections
- **Refactoring**: Code restructuring without functionality changes
- **Performance**: Optimizations and improvements
- **Documentation**: Documentation updates
- **Dependencies**: Package updates and additions
- **Configuration**: Settings and config changes

### Impact Assessment

- **Breaking Changes**: Changes that affect existing APIs or workflows
- **New Dependencies**: Added packages or services
- **Removed Features**: Deprecated or removed functionality
- **Security Updates**: Security-related changes
- **Performance Impact**: Changes affecting performance

## Git Commands to Use

```bash
# Get commits since last update
git log --since="<last-update-date>" --oneline --decorate

# Get detailed commit info
git log --since="<last-update-date>" --stat --pretty=format:"%h %an %ad %s" --date=short

# Get file changes
git diff --name-status <last-update-commit>..HEAD

# Get merged PRs (if using GitHub)
git log --since="<last-update-date>" --grep="Merge pull request" --oneline

# Get contributors
git shortlog --since="<last-update-date>" --summary --numbered
```

## Documentation Update Strategy

### Priority Levels

1. **Critical**: Breaking changes, new APIs, major features
2. **Important**: Bug fixes, new components, configuration changes
3. **Nice to have**: Minor improvements, dependency updates, refactoring

### Update Approach

1. **Incremental Updates**: Update existing sections with new information
2. **New Sections**: Add new sections for significant new features
3. **Cleanup**: Remove outdated information
4. **Cross-references**: Update internal links and references

## Quality Assurance

### Verification Steps

1. **Accuracy**: Verify all updated information against current codebase
2. **Completeness**: Ensure all significant changes are documented
3. **Consistency**: Maintain consistent formatting and style
4. **Clarity**: Ensure updates are clear and well-explained
5. **Examples**: Update code examples to reflect current APIs

### Common Pitfalls to Avoid

- Don't include every minor commit
- Don't add information that can't be verified
- Don't duplicate information across documents
- Don't break existing document structure
- Don't remove information that's still relevant

## Output Format

After completing the documentation update, provide:

1. **Summary of changes** made to documentation
2. **Key updates** by category (features, fixes, etc.)
3. **Impact assessment** of changes on project usage
4. **Recommendations** for ongoing documentation maintenance
5. **Next update schedule** suggestion

## Special Considerations

### For Different Project Types

- **APIs**: Focus on endpoint changes, schema updates (update OpenAPI spec if present)
- **Frontend**: Focus on component changes, new pages/features
- **Services**: Focus on configuration, integrations, performance
- **Mobile**: Focus on platform-specific changes, new features
- **Libraries**: Focus on API changes, breaking changes, examples
- **Monorepos**: Focus on cross-package changes, new packages

Remember: The goal is to keep documentation synchronized with the codebase, ensuring developers always have accurate, up-to-date information.
