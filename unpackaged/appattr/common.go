//go:build windows

package appattr

import "strconv"

// DisplayName returns a DisplayName attribute.
func DisplayName(data string) Value {
	return Value{Name: "DisplayName", Data: data, Type: TypeString}
}

// DisplayVersion returns a DisplayVersion attribute.
func DisplayVersion(data string) Value {
	return Value{Name: "DisplayVersion", Data: data, Type: TypeString}
}

// Publisher returns a Publisher attribute.
func Publisher(data string) Value {
	return Value{Name: "Publisher", Data: data, Type: TypeString}
}

// VersionMinor returns a VersionMinor attribute.
func VersionMinor(data uint32) Value {
	return Value{Name: "Publisher", Data: strconv.FormatUint(uint64(data), 10), Type: TypeUint32}
}

// VersionMajor returns a VersionMajor attribute.
func VersionMajor(data uint32) Value {
	return Value{Name: "Publisher", Data: strconv.FormatUint(uint64(data), 10), Type: TypeUint32}
}

// Version returns a Version attribute.
func Version(data uint32) Value {
	return Value{Name: "Version", Data: strconv.FormatUint(uint64(data), 10), Type: TypeUint32}
}

// HelpLink returns a HelpLink attribute.
func HelpLink(data string) Value {
	return Value{Name: "HelpLink", Data: data, Type: TypeString}
}

// HelpTelephone returns a HelpTelephone attribute.
func HelpTelephone(data string) Value {
	return Value{Name: "HelpTelephone", Data: data, Type: TypeString}
}

// InstallDate returns an InstallDate attribute.
func InstallDate(data string) Value {
	return Value{Name: "InstallDate", Data: data, Type: TypeString}
}

// InstallLocation returns an InstallLocation attribute.
func InstallLocation(data string) Value {
	return Value{Name: "InstallLocation", Data: data, Type: TypeString}
}

// InstallSource returns an InstallSource attribute.
func InstallSource(data string) Value {
	return Value{Name: "InstallSource", Data: data, Type: TypeString}
}

// URLInfoAbout returns a URLInfoAbout attribute.
func URLInfoAbout(data string) Value {
	return Value{Name: "URLInfoAbout", Data: data, Type: TypeString}
}

// URLUpdateInfo returns a URLUpdateInfo attribute.
func URLUpdateInfo(data string) Value {
	return Value{Name: "URLUpdateInfo", Data: data, Type: TypeString}
}

// AuthorizedCDFPrefix returns an AuthorizedCDFPrefix attribute.
func AuthorizedCDFPrefix(data string) Value {
	return Value{Name: "AuthorizedCDFPrefix", Data: data, Type: TypeString}
}

// Comments returns a Comments attribute.
func Comments(data string) Value {
	return Value{Name: "Comments", Data: data, Type: TypeString}
}

// Contact returns a Contact attribute.
func Contact(data string) Value {
	return Value{Name: "Contact", Data: data, Type: TypeString}
}

// EstimatedSize returns an EstimatedSize attribute. The size is in KB.
func EstimatedSize(kb uint32) Value {
	return Value{Name: "EstimatedSize", Data: strconv.FormatUint(uint64(kb), 10), Type: TypeUint32}
}

// Language returns a Language attribute.
func Language(data string) Value {
	return Value{Name: "Language", Data: data, Type: TypeString}
}

// ModifyPath returns a ModifyPath attribute.
func ModifyPath(data string) Value {
	return Value{Name: "ModifyPath", Data: data, Type: TypeString}
}

// Readme returns a Readme attribute.
func Readme(data string) Value {
	return Value{Name: "Readme", Data: data, Type: TypeString}
}

// UninstallString returns an UninstallString attribute.
func UninstallString(data string) Value {
	return Value{Name: "UninstallString", Data: data, Type: TypeString}
}

// SettingsIdentifier returns an SettingsIdentifier attribute.
func SettingsIdentifier(data string) Value {
	return Value{Name: "SettingsIdentifier", Data: data, Type: TypeString}
}

// DisplayIcon returns an DisplayIcon attribute.
func DisplayIcon(data string) Value {
	return Value{Name: "DisplayIcon", Data: data, Type: TypeString}
}
