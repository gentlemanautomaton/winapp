//go:build windows

package unpackaged_test

import (
	"flag"
	"fmt"
	"time"

	"github.com/gentlemanautomaton/winapp/unpackaged"
	"github.com/gentlemanautomaton/winapp/unpackaged/appattr"
	"github.com/gentlemanautomaton/winapp/unpackaged/appregistry"
)

func Example() {
	var install, uninstall bool
	flag.BoolVar(&install, "i", false, "install app")
	flag.BoolVar(&uninstall, "u", false, "uninstall app")
	flag.Parse()

	view := appregistry.Global64
	app := unpackaged.App{
		ID: "winapp-test-app",
		Attributes: appattr.List{
			appattr.DisplayName("Test Application (winapp)"),
			appattr.DisplayVersion("0.0.0"),
			appattr.Publisher("Test Publisher"),
			appattr.InstallDate(time.Now().Format(appattr.DateLayout)),
			appattr.EstimatedSize(1024),
		},
	}

	if install {
		fmt.Printf("Adding unpackaged application \"%s\" to the %s application registry...", app.ID, view.Name())
		if err := view.Add(app); err != nil {
			fmt.Printf(" failed: %v\n", err)
		} else {
			fmt.Printf(" done.\n")
		}
	}

	if uninstall {
		fmt.Printf("Removing unpackaged application \"%s\" from the %s application registry...", app.ID, view.Name())
		if err := view.Remove(app.ID); err != nil {
			fmt.Printf(" failed: %v\n", err)
		} else {
			fmt.Printf(" done.\n")
		}
	}
}
