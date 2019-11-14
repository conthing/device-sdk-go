package cache

import (
	"fmt"
	"sync"

	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

var (
	dc *deviceCache
)

type DeviceCache interface {
	All() []contract.Device
	Update(device contract.Device) error
	Remove(id string) error
	Add(device contract.Device) error
	ForId(id string) (contract.Device, bool)
	ForName(name string) (contract.Device, bool)
	UpdateAdminState(id string, state contract.AdminState) error
	RemoveByName(name string) error
}

type deviceCache struct {
	dMap    *sync.Map
	nameMap sync.Map
}

func (d *deviceCache) All() []contract.Device {
	var devices []contract.Device
	f := func(k, v interface{}) bool {
		devices = append(devices, v.(contract.Device))
		return true
	}
	d.dMap.Range(f)
	return devices
}

func (d *deviceCache) Add(device contract.Device) error {
	if _, ok := d.dMap.Load(device.Id); ok {
		return fmt.Errorf("device %s has already existed in cache", device.Name)
	}
	d.dMap.Store(device.Name, &device)
	d.nameMap.Store(device.Id, device.Name)
	return nil
}

func (d *deviceCache) ForId(id string) (contract.Device, bool) {
	name, ok := d.nameMap.Load(id)

	if !ok {
		return contract.Device{}, ok
	}

	if device, ok := d.dMap.Load(name); ok {
		dev := device.(*contract.Device)
		return *dev, ok
	} else {
		return contract.Device{}, ok
	}
}

func (d *deviceCache) ForName(name string) (contract.Device, bool) {
	if device, ok := d.dMap.Load(name); ok {
		dev := device.(*contract.Device)
		return *dev, ok
	} else {
		return contract.Device{}, ok
	}
}

func (d *deviceCache) RemoveByName(name string) error {
	device, ok := d.dMap.Load(name)
	if !ok {
		return fmt.Errorf("device %s does not exist in cache", name)
	}
	dev := device.(contract.Device)
	d.nameMap.Delete(dev.Id)
	d.dMap.Delete(name)
	return nil
}

func (d *deviceCache) Remove(id string) error {
	name, ok := d.nameMap.Load(id)
	if !ok {
		return fmt.Errorf("device %s does not exist in cache", name.(string))
	}
	return d.RemoveByName(name.(string))
}

func (d *deviceCache) Update(device contract.Device) error {
	if err := d.Remove(device.Id); err != nil {
		return err
	}
	return d.Add(device)
}

func (d *deviceCache) UpdateAdminState(id string, state contract.AdminState) error {
	name, ok := d.nameMap.Load(id)
	if !ok {
		return fmt.Errorf("device %s cannot be found in cache", id)
	}

	device, ok := d.dMap.Load(name)
	dev := device.(contract.Device)
	dev.AdminState = state
	return nil

}

func newDeviceCache(devices []contract.Device) DeviceCache {
	var devicesMap sync.Map
	var nameMap sync.Map
	for _, device := range devices {
		devicesMap.Store(device.Name, device)
		nameMap.Store(device.Id, device.Name)
	}
	dc = &deviceCache{dMap: &devicesMap, nameMap: nameMap}
	return dc
}

func Devices() DeviceCache {
	if dc == nil {
		InitCache()
	}
	return dc
}
