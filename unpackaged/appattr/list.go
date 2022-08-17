//go:build windows

package appattr

// List is a list of application attributes.
type List []Value

// Get returns the first attribute in the list with the given name, if
// present. It returns false if an attribute with that name is not present.
func (list List) Get(name string) (v Value, ok bool) {
	for i := range list {
		if list[i].Name == name {
			return list[i], true
		}
	}
	return Value{}, false
}

// Get returns the first attribute in the list with the given name as a
// string. If the attribute is not present it returns an empty string.
func (list List) GetString(name string) string {
	attr, ok := list.Get(name)
	if !ok {
		return ""
	}
	return attr.Data
}
