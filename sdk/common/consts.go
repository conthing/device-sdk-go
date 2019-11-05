package common

import (
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
)

const (
	ClientData     = "Data"
	ClientMetadata = "Metadata"
	ClientLogging  = "Logging"

	APIv1Prefix = "/api/v1"
	Colon       = ":"
	HttpScheme  = "http://"
	HttpProto   = "HTTP"

	ConfigDirectory    = "./res"
	ConfigFileName     = "configuration.toml"
	ConfigRegistryStem = "edgex/devices/1.0/"
	WritableKey        = "/Writable"
	RegistryFailLimit  = 3

	APIPingRoute            = APIv1Prefix + "/ping"
	APICallbackRoute        = APIv1Prefix + "/callback"
	APIValueDescriptorRoute = APIv1Prefix + "/valuedescriptor"
	APIDiscoveryRoute       = APIv1Prefix + "/discovery"

	IdVar        string = "id"
	NameVar      string = "name"
	CommandVar   string = "command"
	GetCmdMethod string = "get"
	SetCmdMethod string = "set"

	CorrelationHeader = clients.CorrelationHeader
)
