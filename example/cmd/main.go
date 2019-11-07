package main

import (
	"fmt"

	"github.com/conthing/device-sdk-go/example/driver"
	"github.com/conthing/device-sdk-go/pkg/startup"
)

const (
	serviceName = "hvac-iracc"
	version     = "1.0.0"
)

func main() {
	fmt.Println("1")
	sd := driver.SimpleDriver{}
	startup.BootStrap(serviceName, version, &sd)
}
