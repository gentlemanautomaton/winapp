//go:build windows

package unpackaged

import "github.com/gentlemanautomaton/winapp/unpackaged/appattr"

// App is an unpackaged Windows application or software component stored
// within the Windows registry.
type App struct {
	ID         AppID
	Attributes appattr.List
}
