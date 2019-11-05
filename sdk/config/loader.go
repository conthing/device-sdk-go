package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/pelletier/go-toml"

	"github.com/edgexfoundry/go-mod-registry/pkg/types"

	"github.com/edgexfoundry/go-mod-registry/registry"
	"github.com/conthing/device-sdk-go/sdk/common"
)

var (
	RegistryClient registry.Client
)

func LoadConfig(useRegistry string, profile string, confDir string) (configuration *common.Config, err error) {
	fmt.Fprintf(os.Stdout, "Init: useRegistry: %v profile: %s confDir: %s\n", useRegistry, profile, confDir)

	var registryMsg string

	e := NewEnvironment()
	if useRegistry != "" {
		configuration = &common.Config{}
		useRegistry = e.OverrideUseRegistryFromEnvironment(useRegistry)
		err = parseRegistryPath(useRegistry, configuration)
		if err != nil {
			return
		}

		registryMsg = "Registry in registry ..."
		registryConfig := types.Config{
			Host:       configuration.Registry.Host,
			Port:       configuration.Registry.Port,
			Type:       configuration.Registry.Type,
			Stem:       common.ConfigRegistryStem,
			CheckRoute: common.APIPingRoute,
			ServiceKey: common.ServiceName,
		}

		RegistryClient, err = registry.NewRegistryClient(registryConfig)
		if err != nil {
			return nil, fmt.Errorf("connection to Registry could not be made : v%", err.Error())
		}

		if err := checkRegistryUp(configuration); err != nil {
			return nil, err
		}

		hasConfiguration, err := RegistryClient.HasConfiguration()
		if err != nil {
			return nil, fmt.Errorf("could not verify that Registry already has configuration:v%", err.Error())
		}

		if hasConfiguration {
			rawConfig, err := RegistryClient.GetConfiguration(configuration)
			if err != nil {
				return nil, fmt.Errorf("could not get configuration from Registry: v%", err.Error())
			}

			actual, ok := rawConfig.(*common.Config)

			if !ok {
				return nil, fmt.Errorf("configuration from Registry failed type check")
			}

			configuration = actual
		} else {
			fmt.Fprintf(os.Stdout, "Pushing configuration into Registry...")
			_, configTree, err := loadConfigFromFile(profile, confDir)
			if err != nil {
				return nil, err
			}

			err = RegistryClient.PutConfigurationToml(e.OverrideFromEnvironment(configTree), true)
			if err != nil {
				return nil, fmt.Errorf("could not push configuration from Registry: v%", err.Error())
			}

			err = configTree.Unmarshal(configuration)
			if err != nil {
				return nil, fmt.Errorf("could not marshal configTree to configuration: v%", err.Error())
			}
		}

		registryConfig = types.Config{
			Host:          configuration.Registry.Host,
			Port:          configuration.Registry.Port,
			Type:          configuration.Registry.Type,
			Stem:          common.ConfigRegistryStem,
			CheckInterval: configuration.Registry.CheckInterval,
			CheckRoute:    common.APIPingRoute,
			ServiceKey:    common.ServiceName,
			ServiceHost:   configuration.Service.Host,
			ServicePort:   configuration.Service.Port,
		}

		RegistryClient, err = registry.NewRegistryClient(registryConfig)
		if err != nil {
			return nil, fmt.Errorf("connection to Registry could not be made : v%", err.Error())
		}

		err = RegistryClient.Register()
		if err != nil {
			return nil, fmt.Errorf("could not registry service with Registry : v%", err.Error())
		}

	} else {
		registryMsg = "Bypassing registration in registry..."
		configuration, _, err = loadConfigFromFile(profile, confDir)
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(registryMsg)
	return configuration, nil

}

func parseRegistryPath(registryUrl string, config *common.Config) error {
	u, err := url.Parse(registryUrl)
	if err != nil {
		fmt.Fprintf(os.Stdout, "The format of Registry path from argument is wrong: ", err.Error())
		return err
	}

	port, err := strconv.Atoi(u.Port())
	if err != nil {
		fmt.Fprintf(os.Stdout, "The port format of Registry path from argument is wrong: ", err.Error())
		return err
	}

	config.Registry.Type = u.Scheme
	config.Registry.Host = u.Host
	config.Registry.Port = port
	return nil
}

func checkRegistryUp(config *common.Config) error {
	registryUrl := common.BuildAddr(config.Registry.Host, strconv.Itoa(config.Registry.Port))
	fmt.Println("Check registry is up...", registryUrl)
	config.Registry.FailLimit = common.RegistryFailLimit
	fails := 0
	for fails < config.Registry.FailLimit {
		if RegistryClient.IsAlive() {
			break
		}

		time.Sleep(time.Second * time.Duration(config.Registry.FailWaitTime))
		fails++
	}

	if fails >= config.Registry.FailLimit {
		return errors.New("can't get connection to Registry")
	}
	return nil
}

func loadConfigFromFile(profile string, confDir string) (config *common.Config, tree *toml.Tree, err error) {
	if len(confDir) == 0 {
		confDir = common.ConfigDirectory
	}

	if len(profile) > 0 {
		confDir = confDir + "/" + profile
	}

	path := path.Join(confDir, common.ConfigFileName)
	absPath, err := filepath.Abs(path)
	if err != nil {
		err = fmt.Errorf("Could not create absolute path to load configuration: %s,%v", path, err.Error())
		return nil, nil, err
	}

	fmt.Fprintf(os.Stdout, fmt.Sprintf("Loading configuration from: %s\n", absPath))

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("could not load configuration file; invalid TOML (%s)", path)
		}
	}()

	config = &common.Config{}
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("Could not load configuration file (%s):%v\nBe sure to change to program folder or set working directory.", path, err.Error())
	}

	err = toml.Unmarshal(contents, config)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal configuration struct (%s):%v", path, err.Error)
	}

	tree, err = toml.LoadBytes(contents)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to unmarshal configuration tree (%s):%v", path, err.Error)
	}

	return config, tree, nil
}

func ListenForConfigChange() {
	if RegistryClient == nil {
		common.LoggingClient.Error("listenForConfigChanges() registry client not set")
		return
	}

	common.LoggingClient.Info("listen for config changes from Registry")

	errChannel := make(chan error)
	updateChannel := make(chan interface{})

	RegistryClient.WatchForChanges(updateChannel, errChannel, &common.WritableInfo{}, common.WritableKey)

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-signalChan:
			return
		case ex := <-errChannel:
			common.LoggingClient.Error(ex.Error())
		case raw, ok := <-updateChannel:
			if ok {
				actual, ok := raw.(*common.WritableInfo)
				if !ok {
					common.LoggingClient.Error("listenForConfigChanges() type check failed")
				}
				common.CurrentConfig.Writable = *actual
				common.LoggingClient.Info("Writable configuration has been updated,Setting log level to " + common.CurrentConfig.Writable.LogLevel)
				common.LoggingClient.SetLogLevel(common.CurrentConfig.Writable.LogLevel)
			} else {
				return
			}
		}
	}
}
