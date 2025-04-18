//go:build windows

package appregistry

import (
	"errors"
	"fmt"
	"syscall"

	"github.com/gentlemanautomaton/winapp/appcode"
	"github.com/gentlemanautomaton/winapp/unpackaged"
	"github.com/gentlemanautomaton/winapp/unpackaged/appscope"
	"golang.org/x/sys/windows/registry"
)

// View provides a view of the machine or user application registry, either
// 32-bit or 64-bit.
type View struct {
	name   string
	arch   appcode.Architecture
	scope  appscope.Scope
	key    registry.Key
	access uint32
}

// Registry path for both HKLM and HKCU.
const root = `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`

// These are the four possible views of installed applications.
var (
	machine32 = View{name: "32-bit machine", arch: appcode.X86, scope: appscope.Machine, key: registry.LOCAL_MACHINE, access: registry.WOW64_32KEY}
	machine64 = View{name: "64-bit machine", arch: appcode.X64, scope: appscope.Machine, key: registry.LOCAL_MACHINE, access: registry.WOW64_64KEY}
	user32    = View{name: "32-bit user", arch: appcode.X86, scope: appscope.User, key: registry.CURRENT_USER, access: registry.WOW64_32KEY}
	user64    = View{name: "64-bit user", arch: appcode.X64, scope: appscope.User, key: registry.CURRENT_USER, access: registry.WOW64_64KEY}
)

var (
	// Machine32 accesses the 32-bit machine application registry.
	Machine32 = machine32

	// Machine64 accesses the 64-bit machine application registry.
	Machine64 = machine64

	// User32 accesses the 32-bit user application registry.
	User32 = user32

	// User64 accesses the 64-bit user application registry.
	User64 = user64
)

// Views is a slice of all available application views.
var Views = []View{Machine32, Machine64, User32, User64}

// ViewFor returns a view for the given architecture and scope. It returns an
// error if either are not recognized.
func ViewFor(arch appcode.Architecture, scope appscope.Scope) (View, error) {
	switch arch {
	case appcode.X64:
		switch scope {
		case appscope.Machine:
			return machine64, nil
		case appscope.User:
			return user64, nil
		case "":
			return View{}, errors.New("missing application scope")
		default:
			return View{}, fmt.Errorf("unrecognized application scope: %s", scope)
		}
	case appcode.X86:
		switch scope {
		case appscope.Machine:
			return machine32, nil
		case appscope.User:
			return user32, nil
		case "":
			return View{}, errors.New("missing application scope")
		default:
			return View{}, fmt.Errorf("unrecognized application scope: %s", scope)
		}
	case "":
		return View{}, errors.New("missing application architecture")
	default:
		return View{}, fmt.Errorf("unrecognized application architecture: %s", arch)
	}
}

// open opens a registry key in the view.
func (v View) open(k registry.Key, path string, access uint32) (registry.Key, error) {
	if v.name == "" {
		return 0, errors.New("use of unprepared application registry view")
	}
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

// Architecture returns the architecture of the view.
func (v View) Architecture() appcode.Architecture {
	return v.arch
}

// Scope returns the scope of the view.
func (v View) Scope() appscope.Scope {
	return v.scope
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

// Contains returns true if the given unpackaged app ID is present in the
// Windows registry.
func (v View) Contains(id unpackaged.AppID) (bool, error) {
	root, err := v.root(registry.ENUMERATE_SUB_KEYS)
	if err != nil {
		return false, err
	}
	defer root.Close()

	key, err := v.open(root, string(id), registry.QUERY_VALUE)
	if err != nil {
		if err == syscall.ERROR_FILE_NOT_FOUND {
			return false, nil
		}
		return false, err
	}
	defer key.Close()

	return true, nil
}

// Get attempts to retrieve the requested unpackaged app from the Windows
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
