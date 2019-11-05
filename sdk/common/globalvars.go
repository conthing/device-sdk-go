package common

import (
	"github.com/edgexfoundry/go-mod-core-contracts/clients/coredata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/conthing/device-sdk-go/pkg/models"
)

var (
	ServiceName           string
	ServiceVersion        string
	CurrentConfig         *Config
	CurrentDeviceService  contract.DeviceService
	UseRegistry           bool
	ServiceLocked         bool
	Driver                models.ProtocolDriver
	EventClient           coredata.EventClient
	AddressableClient     metadata.AddressableClient
	DeviceClient          metadata.DeviceClient
	DeviceServiceClient   metadata.DeviceServiceClient
	DeviceProfileClient   metadata.DeviceProfileClient
	LoggingClient         logger.LoggingClient
	ValueDescriptorClient coredata.ValueDescriptorClient
)
