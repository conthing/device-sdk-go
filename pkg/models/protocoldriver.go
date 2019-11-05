package models

import (
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

type ProtocolDriver interface {
	Initialize(lc logger.LoggingClient, asyncCh chan<- *AsyncValues) error

	HandleReadCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []CommandRequest) ([]*CommandValue, error)

	HandleWriteCommands(deviceName string, protocols map[string]contract.ProtocolProperties, reqs []CommandRequest, params []*CommandValue) error

	Stop(force bool) error
}
