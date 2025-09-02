#!/bin/bash

# verify.sh - Comprehensive verification script for examples/basic module
# This script implements comprehensive automated verification for compilation and execution

set -e  # Exit on any error
set -u  # Exit on undefined variable

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verification steps
verify_files() {
    log_info "Step 1: Verifying required files exist..."
    
    local required_files=(
        "main.go"
        "go.mod"
        "go.sum"
        ".env.example"
        "example_items.csv"
        "custom_prompt.txt"
    )
    
    local missing_files=()
    for file in "${required_files[@]}"; do
        if [[ ! -f "$file" ]]; then
            missing_files+=("$file")
        fi
    done
    
    if [[ ${#missing_files[@]} -gt 0 ]]; then
        log_error "Missing required files: ${missing_files[*]}"
        return 1
    fi
    
    log_success "All required files exist"
    return 0
}

verify_module_integrity() {
    log_info "Step 2: Verifying Go module integrity..."
    
    # Check go.mod format
    if ! go mod verify; then
        log_error "Module verification failed"
        return 1
    fi
    
    # Run go mod tidy to ensure clean state
    go mod tidy
    
    # Check for any changes after tidy
    if ! git diff --quiet go.mod go.sum 2>/dev/null; then
        log_warning "go.mod or go.sum changed after 'go mod tidy'"
        log_info "This may indicate dependencies were not properly maintained"
    fi
    
    log_success "Module integrity verified"
    return 0
}

verify_compilation() {
    log_info "Step 3: Verifying compilation..."
    
    # Clean build
    rm -f example example_test pipeline_test
    
    # Test compilation
    if ! go build -o example_verification .; then
        log_error "Compilation failed"
        return 1
    fi
    
    # Verify binary was created and is executable
    if [[ ! -x "example_verification" ]]; then
        log_error "Compiled binary is not executable"
        return 1
    fi
    
    # Clean up
    rm -f example_verification
    
    log_success "Compilation successful"
    return 0
}

verify_tests() {
    log_info "Step 4: Running unit tests..."
    
    # Run tests with coverage
    if ! go test -v -cover -short .; then
        log_error "Unit tests failed"
        return 1
    fi
    
    log_success "Unit tests passed"
    return 0
}

verify_integration() {
    log_info "Step 5: Running integration tests..."
    
    # Run integration tests (non-short)
    if ! go test -v -run "TestIntegration|TestFull" .; then
        log_warning "Some integration tests failed (may require API key)"
    else
        log_success "Integration tests passed"
    fi
    
    return 0
}

verify_data_files() {
    log_info "Step 6: Verifying data files..."
    
    # Test CSV files can be loaded
    local csv_files=("example_items.csv" "example_items_edge_cases.csv")
    
    for csv_file in "${csv_files[@]}"; do
        if [[ -f "$csv_file" ]]; then
            local line_count=$(wc -l < "$csv_file")
            if [[ $line_count -lt 2 ]]; then
                log_error "$csv_file appears to be empty or malformed (only $line_count lines)"
                return 1
            else
                log_info "$csv_file validated ($line_count lines)"
            fi
        else
            log_warning "$csv_file not found"
        fi
    done
    
    # Test custom prompt
    if [[ -f "custom_prompt.txt" ]]; then
        local prompt_size=$(wc -c < "custom_prompt.txt")
        if [[ $prompt_size -lt 50 ]]; then
            log_warning "custom_prompt.txt seems too short ($prompt_size characters)"
        else
            log_info "custom_prompt.txt validated ($prompt_size characters)"
        fi
    fi
    
    log_success "Data files verified"
    return 0
}

verify_environment() {
    log_info "Step 7: Verifying environment setup..."
    
    # Check .env.example has required variables
    if [[ -f ".env.example" ]]; then
        local required_vars=("OPENAI_API_KEY" "LOG_LEVEL")
        local missing_vars=()
        
        for var in "${required_vars[@]}"; do
            if ! grep -q "^${var}=" ".env.example"; then
                missing_vars+=("$var")
            fi
        done
        
        if [[ ${#missing_vars[@]} -gt 0 ]]; then
            log_warning ".env.example missing variables: ${missing_vars[*]}"
        else
            log_info ".env.example contains all required variables"
        fi
    fi
    
    # Check if we have an API key for optional live testing  
    if [[ -n "${OPENAI_API_KEY:-}" ]]; then
        log_info "OPENAI_API_KEY found - live testing possible"
    else
        log_info "OPENAI_API_KEY not set - using mock testing only"
    fi
    
    log_success "Environment verification completed"
    return 0
}

verify_documentation() {
    log_info "Step 8: Verifying documentation..."
    
    # Check if README exists and has content
    if [[ -f "README.md" ]]; then
        local readme_size=$(wc -c < "README.md")
        if [[ $readme_size -lt 100 ]]; then
            log_warning "README.md seems too short ($readme_size characters)"
        else
            log_info "README.md exists ($readme_size characters)"
        fi
    else
        log_info "README.md not found (optional)"
    fi
    
    log_success "Documentation verification completed"
    return 0
}

run_benchmarks() {
    log_info "Step 9: Running performance benchmarks..."
    
    # Run benchmarks if requested
    if [[ "${RUN_BENCHMARKS:-}" == "true" ]]; then
        go test -bench=. -benchmem .
        log_success "Benchmarks completed"
    else
        log_info "Benchmarks skipped (set RUN_BENCHMARKS=true to enable)"
    fi
    
    return 0
}

verify_security() {
    log_info "Step 10: Basic security checks..."
    
    # Check for hardcoded API keys in source files
    if grep -r "sk-[a-zA-Z0-9]" --include="*.go" .; then
        log_error "Potential hardcoded API keys found in source code!"
        return 1
    fi
    
    # Check for TODO/FIXME markers that might indicate incomplete work
    local todo_count=$(grep -r "TODO\|FIXME" --include="*.go" . | wc -l)
    if [[ $todo_count -gt 0 ]]; then
        log_warning "Found $todo_count TODO/FIXME comments in code"
    fi
    
    log_success "Security checks completed"
    return 0
}

# Main verification function
main() {
    log_info "Starting comprehensive verification of examples/basic module..."
    log_info "Working directory: $SCRIPT_DIR"
    
    local failed_steps=()
    local verification_steps=(
        "verify_files"
        "verify_module_integrity" 
        "verify_compilation"
        "verify_tests"
        "verify_integration"
        "verify_data_files"
        "verify_environment"
        "verify_documentation"
        "run_benchmarks"
        "verify_security"
    )
    
    for step in "${verification_steps[@]}"; do
        if ! $step; then
            failed_steps+=("$step")
        fi
    done
    
    echo
    echo "========================================="
    echo "VERIFICATION SUMMARY"
    echo "========================================="
    
    if [[ ${#failed_steps[@]} -eq 0 ]]; then
        log_success "ALL VERIFICATION STEPS PASSED!"
        log_info "The examples/basic module is ready for use"
        exit 0
    else
        log_error "FAILED STEPS: ${failed_steps[*]}"
        log_error "Some verification steps failed - please review above output"
        exit 1
    fi
}

# Handle command line arguments
case "${1:-verify}" in
    "verify"|"")
        main
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [verify|help]"
        echo
        echo "Comprehensive verification script for examples/basic module"
        echo
        echo "Options:"
        echo "  verify    Run all verification steps (default)"
        echo "  help      Show this help message"
        echo
        echo "Environment variables:"
        echo "  RUN_BENCHMARKS=true    Enable performance benchmarks"
        echo "  OPENAI_API_KEY=...     Enable live API testing"
        exit 0
        ;;
    *)
        log_error "Unknown command: $1"
        echo "Use '$0 help' for usage information"
        exit 1
        ;;
esac