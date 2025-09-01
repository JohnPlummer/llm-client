package scorer

// Version represents the current version of the llm-client library
const Version = "0.11.0"

// VersionInfo provides version information for the library
type VersionInfo struct {
	Version string
	Name    string
}

// GetVersion returns version information for the library
func GetVersion() VersionInfo {
	return VersionInfo{
		Version: Version,
		Name:    "llm-client",
	}
}