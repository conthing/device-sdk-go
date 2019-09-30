package models

type ProtocolDiscovery interface {
	Discover() (devices *interface{}, err error)
}
