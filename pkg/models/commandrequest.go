package models

type CommandRequest struct {

	DeviceResourceName string

	Attributes map[string]string

	Type ValueType
}
