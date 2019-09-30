package main

import (
	"fmt"

	"github.com/conthing/device-sdk-go/example/driver"
	"github.com/conthing/device-sdk-go/pkg/startup"
)

const (
	serviceName = "device-simple"
	version     = "0.0.1"
)

func main() {
	fmt.Println("1")
	sd := driver.SimpleDriver{}
	startup.BootStrap(serviceName, version, &sd)
}
