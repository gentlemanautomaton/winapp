//go:build windows

package appattr

// Value is an application attribute that can be recorded in the Windows
// registry.
type Value struct {
	Name string
	Data string
	Type Type
}

// String returns a string representation of the attribute.
func (v Value) String() string {
	return v.Name + ": " + v.Data
}
