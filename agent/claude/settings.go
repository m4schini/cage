package claude

import "os"

// config.go — example values struct

type PermissionRules struct {
	Allow []string
	Deny  []string
}

type AttributionCfg struct {
	Commits      bool
	PullRequests bool
}

type SandboxCfg struct {
	Enabled             bool
	UnsandboxedCommands []string
	WeakerNested        bool // true only in unprivileged Docker
}

type ClaudeSettings struct {
	// Auth
	ApiKey       string // ANTHROPIC_API_KEY
	BaseURL      string // optional gateway / proxy
	ApiTimeoutMs string // e.g. "300000"

	// Models (empty = use Claude Code defaults)
	Model       string // top-level model override
	SonnetModel string // ANTHROPIC_DEFAULT_SONNET_MODEL
	HaikuModel  string // ANTHROPIC_DEFAULT_HAIKU_MODEL
	OpusModel   string // ANTHROPIC_DEFAULT_OPUS_MODEL

	// Behaviour
	AutoUpdatesChannel    string // "stable" | "beta" | "off"
	AutoCompactPct        int    // 0 = disable, 80 = default
	AlwaysThinkingEnabled *bool
	ShowTurnDuration      *bool
	SpinnerTipsEnabled    *bool

	Permissions PermissionRules
	Attribution AttributionCfg
	Sandbox     SandboxCfg
}

// Example values for a container / CI context:
var ContainerDefaults = ClaudeSettings{
	ApiKey:             os.Getenv("ANTHROPIC_API_KEY"),
	Model:              "claude-sonnet-4-6",
	HaikuModel:         "claude-haiku-4-5-20251001",
	AutoUpdatesChannel: "stable",
	AutoCompactPct:     80,
	Permissions: PermissionRules{
		Allow: []string{"Bash(git *)", "Bash(go *)", "Bash(make *)"},
		Deny:  []string{"Read(.env)", "Read(.env.*)", "Read(~/.ssh/**)"},
	},
	Attribution: AttributionCfg{Commits: false, PullRequests: false},
	Sandbox: SandboxCfg{
		Enabled:             true,
		WeakerNested:        true, // needed inside Docker without --privileged
		UnsandboxedCommands: []string{"git", "docker"},
	},
}
