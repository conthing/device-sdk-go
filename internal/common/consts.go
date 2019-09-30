package common

import (
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
)

const (
	ClientData     = "Data"
	ClientMetadata = "Metadata"
	ClientLogging  = "Logging"

	Colon      = ":"
	HttpScheme = "http://"
	HttpProto  = "HTTP"

	ConfigDirectory    = "./res"
	ConfigFileName     = "configuration.toml"
	ConfigRegistryStem = "edgex/devices/1.0/"
	WritableKey        = "/Writable"
	RegistryFailLimit  = 3

	APIPingRoute            = clients.ApiPingRoute
	APICallbackRoute        = clients.ApiCallbackRoute
	APIValueDescriptorRoute = clients.ApiValueDescriptorRoute
	APIVersionRoute         = clients.ApiVersionRoute
	APIMetricsRoute         = clients.ApiMetricsRoute
	APIConfigRoute          = clients.ApiConfigRoute
	APIAllCommandRoute      = clients.ApiDeviceRoute + "/all/{command}"
	APIIdCommandRoute       = clients.ApiDeviceRoute + "/{id}/{command}"
	APINameCommandRoute     = clients.ApiDeviceRoute + "/name/{name}/{command}"
	APIDiscoveryRoute       = clients.ApiBase + "discovery"
	APITransformRoute       = clients.ApiBase + "debug/transformData/{transformData}"

	CorrelationHeader = clients.CorrelationHeader
	URLRawQuery       = "urlRawQuery"
)
