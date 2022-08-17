//go:build windows

package appregistry

import (
	"errors"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/gentlemanautomaton/winapp/unpackaged"
	"github.com/gentlemanautomaton/winapp/unpackaged/appattr"
	"golang.org/x/sys/windows/registry"
)

// writeApp writes an app to the given registry key.
func writeApp(app unpackaged.App, key registry.Key) error {
	for _, attr := range app.Attributes {
		switch attr.Type {
		case appattr.TypeString:
			if err := key.SetStringValue(attr.Name, attr.Data); err != nil {
				return err
			}
		case appattr.TypeExpand:
			if err := key.SetExpandStringValue(attr.Name, attr.Data); err != nil {
				return err
			}
		case appattr.TypeUint32:
			data, err := strconv.ParseUint(attr.Data, 10, 32)
			if err != nil {
				return err
			}
			if err := key.SetDWordValue(attr.Name, uint32(data)); err != nil {
				return err
			}
		default:
			return errors.New("unrecongized or unsupported attribute type")
		}
	}

	return nil
}

// readApp reads an app from the given registry key.
func readApp(app *unpackaged.App, key registry.Key) error {
	var attrs appattr.List

	names, err := key.ReadValueNames(0)
	if err != nil {
		return err
	}

	buf := make([]byte, 4096)
	for _, name := range names {
		n, typ, err := key.GetValue(name, buf)
		if err == syscall.ERROR_MORE_DATA && n < 1048576 {
			buf = make([]byte, n)
			n, typ, err = key.GetValue(name, buf)
		}
		if err != nil {
			return err
		}

		data := buf[:n]

		switch typ {
		case registry.SZ, registry.EXPAND_SZ:
			val := ""
			if len(data) > 0 {
				u := (*[1 << 29]uint16)(unsafe.Pointer(&data[0]))[: len(data)/2 : len(data)/2]
				val = syscall.UTF16ToString(u)
			}
			attrs = append(attrs, appattr.Value{Name: name, Data: val, Type: appattr.Type(typ)})
		case registry.DWORD:
			if len(data) != 4 {
				return errors.New("DWORD value is not 4 bytes long")
			}
			var val32 uint32
			copy((*[4]byte)(unsafe.Pointer(&val32))[:], data)
			attrs = append(attrs, appattr.Value{Name: name, Data: strconv.Itoa(int(val32)), Type: appattr.Type(typ)})
		}
	}

	app.Attributes = attrs

	return nil
}
