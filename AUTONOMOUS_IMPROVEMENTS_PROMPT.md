# Autonomous Improvement Execution Prompt

Execute the improvements listed in IMPROVEMENTS.md autonomously, following these rules:

1. Work through each step sequentially from Step 1 to Step 12
2. DO NOT ask for permission to run commands or make changes
3. DO NOT ask for confirmation before proceeding to the next step
4. DO run `make check` after each implementation
5. DO fix any errors or warnings that arise before proceeding
6. DO commit changes after each successful step
7. DO update the IMPROVEMENTS.md file to mark steps as complete
8. DO continue until all 12 steps are complete

## Execution Instructions

For each step in IMPROVEMENTS.md:

1. Read the step requirements carefully
2. Implement the changes as described
3. Run `make check` immediately after implementation
4. If there are errors or warnings:
   - Fix them without asking for permission
   - Run `make check` again until it passes
5. Update IMPROVEMENTS.md to mark the step as complete (add ✅ to the title)
6. Commit the changes with the specified commit message
7. Proceed immediately to the next step

## Error Handling

- If `make check` fails, analyze the error and fix it autonomously
- If tests fail, update or add tests as needed to make them pass
- If compilation fails, fix the code to resolve the issue
- Do not stop or ask for help unless encountering an unrecoverable error

## Expected Behavior

Start with Step 1 and continue through Step 12 without interruption. The final result should be:
- All 12 improvements implemented
- All tests passing
- All steps marked as complete in IMPROVEMENTS.md
- 12+ commits in git history (one per step minimum)
- A fully improved codebase

Begin immediately with Step 1: Fix Init Error Handling.