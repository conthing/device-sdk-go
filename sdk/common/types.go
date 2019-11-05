package common

import (
	"fmt"

	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

type ServiceInfo struct {
	Host string

	Port int

	ConnectRetries int

	Labels []string

	OpenMsg string

	Timeout int

	EnableAsyncReadings bool

	AsyncBufferSize int
}

type Config struct {
	Service ServiceInfo

	Registry RegistryService

	Clients map[string]ClientInfo

	Logging LoggingInfo

	Writable WritableInfo

	Device DeviceInfo

	DeviceList []DeviceConfig `consul:"-"`
}

type DeviceConfig struct {
	Name        string
	Profile     string
	Description string
	Labels      []string
	Protocols   map[string]contract.ProtocolProperties
	AutoEvents  []contract.AutoEvent
}

type DeviceInfo struct {
	DataTransform  bool
	InitCmd        string
	InitCmdArgs    string
	MaxCmdOps      int
	MaxCmdValueLen int
	RemoveCmd      string
	RemoveCmdArgs  string
	ProfilesDir    string
}

type WritableInfo struct {
	LogLevel string
}

type RegistryService struct {
	Host          string
	Port          int
	Type          string
	Timeout       int
	CheckInterval string
	FailLimit     int
	FailWaitTime  int64
}

type ClientInfo struct {
	Name     string
	Host     string
	Port     int
	Protocol string
	Timeout  int
}

type LoggingInfo struct {
	EnableRemote bool
	File         string
}

func (c ClientInfo) Url() string {
	url := fmt.Sprintf("%s://%s:%v", c.Protocol, c.Host, c.Port)
	return url
}

type Telemetry struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects uint64
}
