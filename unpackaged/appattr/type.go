//go:build windows

package appattr

import "golang.org/x/sys/windows/registry"

// Type identifies the type of an attribute.
//
// When attributes are written to the registry, their type determines the type
// of registry value that will be written.
type Type uint32

// Supported attribute types.
const (
	TypeNone   Type = registry.NONE
	TypeString Type = registry.SZ
	TypeExpand Type = registry.EXPAND_SZ
	TypeUint32 Type = registry.DWORD
)
