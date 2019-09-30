package models

import (
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type Event struct {
	models.Event
	EncodedEvent []byte
}

func (e Event) HasBinaryValue() bool {
	if len(e.Readings) > 0 {
		for r := range e.Readings {
			if len(e.Readings[r].BinaryValue) > 0 {
				return true
			}
		}
	}
	return false
}
