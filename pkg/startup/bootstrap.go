package startup

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/conthing/device-sdk-go"
	"github.com/conthing/device-sdk-go/pkg/models"
)

var (
	confProfile string
	confDir     string
	useRegistry string
)

func BootStrap(serviceName string, serviceVersion string, driver models.ProtocolDriver) {
	flag.StringVar(&useRegistry, "registry", "", "Indicates")
	flag.StringVar(&useRegistry, "r", "", "Indicates")
	flag.StringVar(&confProfile, "profile", "", "Specify")
	flag.StringVar(&confProfile, "p", "", "Specify")
	flag.StringVar(&confDir, "confdir", "", "Specify")
	flag.StringVar(&confDir, "c", "", "Specify")
	flag.Parse()

	if err := startService(serviceName, serviceVersion, driver); err != nil {
		fmt.Fprintf(os.Stderr, "error:%v\n", err)
		os.Exit(1)
	}
}

func startService(serviceName string, serviceVersion string, driver models.ProtocolDriver) error {
	s, err := device.NewService(serviceName, serviceVersion, confProfile, confDir, useRegistry, driver)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Calling service Start.\n")
	errChan := make(chan error, 2)
	listenForInterrupt(errChan)
	go s.Start(errChan)

	err = <-errChan
	fmt.Fprintf(os.Stdout, "Terminating: v%.\n", err)
	return err
}

func listenForInterrupt(errChan chan error) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
}
