// Package scorer provides version information and metadata for the llm-client library.
//
// This package exports version constants and utilities for runtime version detection,
// enabling applications to programmatically access library version information for
// logging, compatibility checks, and debugging purposes.
//
// The version follows semantic versioning (semver) principles and is updated with
// each release to reflect changes in functionality, bug fixes, or breaking changes.
package scorer

// Version represents the current semantic version of the llm-client library.
//
// This constant follows semantic versioning format (MAJOR.MINOR.PATCH) and is
// automatically updated during the release process. Applications can use this
// for version logging, compatibility validation, or feature detection.
//
// Current version: 0.11.0 indicates pre-1.0 development phase with potential
// breaking changes between minor versions.
const Version = "0.11.0"

// VersionInfo encapsulates comprehensive version metadata for the llm-client library.
//
// This struct provides structured access to version information, enabling applications
// to perform version comparisons, logging, and runtime compatibility checks without
// parsing version strings manually.
//
// Fields:
//   - Version: Semantic version string (e.g., "0.11.0")
//   - Name: Human-readable library name for identification
type VersionInfo struct {
	// Version contains the semantic version string following semver format
	Version string

	// Name contains the canonical library name for identification purposes
	Name string
}

// GetVersion returns structured version information for the llm-client library.
//
// This function provides the primary interface for applications to access version
// metadata programmatically. The returned VersionInfo struct contains both the
// version string and library name, enabling comprehensive version reporting.
//
// Returns:
//
//	VersionInfo: Struct containing version string and library name
//
// Usage:
//
//	info := GetVersion()
//	log.Printf("Using %s version %s", info.Name, info.Version)
func GetVersion() VersionInfo {
	return VersionInfo{
		Version: Version,
		Name:    "llm-client",
	}
}
