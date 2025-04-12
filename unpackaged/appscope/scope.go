package appscope

// Scope identifies the scope of an application's installation.
type Scope string

// Recognized application scopes.
const (
	Machine Scope = "machine"
	User    Scope = "user"
)
