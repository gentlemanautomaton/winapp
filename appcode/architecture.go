package appcode

// Architecture identifies the processor architecture targeted by application
// code.
type Architecture string

// Architectures recognized by windows.
//
// https://learn.microsoft.com/en-us/windows/msix/package/device-architecture
// https://learn.microsoft.com/en-us/windows/win32/sysinfo/image-file-machine-constants
const (
	X64 Architecture = "x64"
	X86 Architecture = "x86"

	// TODO: Consider adding these, but think through the impact on
	// unpackaged/appregistry:
	//ARM32 Type = "ARM"
	//ARM64 Type = "ARM64"
)
