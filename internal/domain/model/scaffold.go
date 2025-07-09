package model

// ScaffoldOptions represents the options for generating a scaffold
type ScaffoldOptions struct {
	// Basic options
	AppType      string `json:"appType"`      // "api" only
	DatabaseType string `json:"databaseType"` // "none", "postgresql", "mysql", "sqlite"
	RouterType   string `json:"routerType"`   // "standard", "chi", "echo", etc.
	ConfigType   string `json:"configType"`   // "env", "flags"
	LogFormat    string `json:"logFormat"`    // "json", "text"
	ModulePath   string `json:"modulePath"`   // e.g. "github.com/username/project"

	// Additional features
	Features []string `json:"features"` // List of feature names to include

	// Premium features
	PremiumFeatures []string `json:"premiumFeatures"` // Premium features
}

// Template represents a template that can be used for code generation
type Template struct {
	ID          string `json:"id"`          // Unique identifier
	Name        string `json:"name"`        // Display name
	Description string `json:"description"` // Short description
	Path        string `json:"path"`        // Filesystem path to template
	Type        string `json:"type"`        // "api" only
}

// GeneratedScaffold represents a generated scaffold
type GeneratedScaffold struct {
	ID        string          `json:"id"`        // Unique identifier
	Options   ScaffoldOptions `json:"options"`   // Options used to generate the scaffold
	CreatedAt string          `json:"createdAt"` // Creation timestamp
	FilePath  string          `json:"filePath"`  // Path to the generated ZIP file
	Size      int64           `json:"size"`      // Size of the generated ZIP file in bytes
}

// Feature represents a feature that can be included in a scaffold
type Feature struct {
	ID          string `json:"id"`          // Unique identifier
	Name        string `json:"name"`        // Display name
	Description string `json:"description"` // Short description
	IsPremium   bool   `json:"isPremium"`   // Whether this is a premium feature
}
