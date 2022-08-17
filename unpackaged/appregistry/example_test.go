//go:build windows

package appregistry_test

import (
	"fmt"

	"github.com/gentlemanautomaton/winapp/unpackaged/appregistry"
)

func Example() {
	for _, view := range appregistry.Views {
		fmt.Printf("---- %s ----\n", view.Name())

		apps, err := view.List()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		for _, app := range apps {
			line := string(app.ID)
			attrs := app.Attributes
			if name := attrs.GetString("DisplayName"); name != "" {
				line += ": " + name
			}
			if version := attrs.GetString("DisplayVersion"); version != "" {
				line += " [" + version + "]"
			}
			fmt.Printf("  %s\n", line)
		}
	}
}
