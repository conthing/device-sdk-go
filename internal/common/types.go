package common

import (
	"fmt"
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
