//go:build windows

package appregistry

import (
	"fmt"

	"github.com/gentlemanautomaton/winapp/unpackaged"
	"golang.org/x/sys/windows/registry"
)

// View provides a view of the global or user application registry, either
// 32-bit or 64-bit.
type View struct {
	name   string
	access uint32
	key    registry.Key
}

// Registry path for both HKLM and HKCU.
const root = `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`

// These are the four possible views of installed applications.
var (
	global32 = View{name: "32-bit global", key: registry.LOCAL_MACHINE, access: registry.WOW64_32KEY}
	global64 = View{name: "64-bit global", key: registry.LOCAL_MACHINE, access: registry.WOW64_64KEY}
	user32   = View{name: "32-bit user", key: registry.CURRENT_USER, access: registry.WOW64_32KEY}
	user64   = View{name: "64-bit user", key: registry.CURRENT_USER, access: registry.WOW64_64KEY}
)

var (
	// Global32 accesses the 32-bit global application registry.
	Global32 = global32

	// Global64 accesses the 64-bit global application registry.
	Global64 = global64

	// User32 accesses the 32-bit user application registry.
	User32 = user32

	// User64 accesses the 64-bit user application registry.
	User64 = user64
)

// Views is a slice of all available application views.
var Views = []View{Global32, Global64, User32, User64}

// open opens a registry key in the view.
func (v View) open(k registry.Key, path string, access uint32) (registry.Key, error) {
	return registry.OpenKey(k, path, access|v.access)
}

// root opens the application registry root key in the view.
func (v View) root(access uint32) (registry.Key, error) {
	return v.open(v.key, root, access)
}

// Name returns the name of the view.
func (v View) Name() string {
	return v.name
}

// Add attempts to add the given unpackaged app to the Windows registry.
//
// It returns an error if an application or component with the given ID
// already exists.
func (v View) Add(app unpackaged.App) error {
	root, err := v.root(registry.CREATE_SUB_KEY)
	if err != nil {
		return err
	}
	defer root.Close()

	key, existing, err := registry.CreateKey(root, string(app.ID), registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	if existing {
		return fmt.Errorf("the registry key \"%s\" already exists", app.ID)
	}

	return writeApp(app, key)
}

// Remove attempts to remove the given unpackaged app ID from the Windows
// registry.
//
// It does _not_ run the component's uninstall command. It is the caller's
// responsibility to properly remove the component's files before removing it
// from the registry.
func (v View) Remove(id unpackaged.AppID) error {
	root, err := v.root(registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return err
	}
	defer root.Close()

	return registry.DeleteKey(root, string(id))
}

// Add attempts to add the given application component to the Windows
// registry.
func (v View) Get(id unpackaged.AppID) (unpackaged.App, error) {
	root, err := v.root(registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return unpackaged.App{}, err
	}
	defer root.Close()

	key, err := v.open(root, string(id), registry.QUERY_VALUE)
	if err != nil {
		return unpackaged.App{}, err
	}
	defer key.Close()

	var app unpackaged.App
	if err := readApp(&app, key); err != nil {
		return unpackaged.App{}, err
	}

	return app, nil
}

// Add attempts to add the given application component to the Windows
// registry.
//
// It returns an error if a component with the given ID already exists.
func (v View) List() (unpackaged.AppList, error) {
	root, err := v.root(registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return nil, fmt.Errorf("failed to read the application registry: %w", err)
	}
	defer root.Close()

	ids, err := root.ReadSubKeyNames(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read application registry sub keys: %w", err)
	}

	var out unpackaged.AppList
	for _, id := range ids {
		key, err := v.open(root, id, registry.QUERY_VALUE)
		if err != nil {
			return nil, fmt.Errorf("failed to open application \"%s\": %w", id, err)
		}

		app := unpackaged.App{ID: unpackaged.AppID(id)}
		err = readApp(&app, key)

		key.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read attributes for application \"%s\": %w", id, err)
		}

		out = append(out, app)
	}

	return out, nil
}
