package autoevent

import (
	"fmt"
	"sync"

	"github.com/conthing/device-sdk-go/sdk/common"

	contract "github.com/edgexfoundry/go-mod-core-contracts/models"

	"github.com/conthing/device-sdk-go/sdk/cache"
)

type Manager interface {
	StartAutoEvents()
	StopAutoEvents()
	RestartForDevice(deviceName string)
	StopForDevice(deviceName string)
}

var (
	createOnce sync.Once
	m          *manager
	mutex      sync.Mutex
)

type manager struct {
	execsMap  map[string][]Executor
	startOnce sync.Once
}

func (m *manager) StartAutoEvents() {
	mutex.Lock()
	m.startOnce.Do(func() {
		for _, d := range cache.Devices().All() {
			execs := triggerExecutors(d.Name, d.AutoEvents)
			m.execsMap[d.Name] = execs
		}
	})
	mutex.Unlock()
}

func (m *manager) StopAutoEvents() {
	mutex.Lock()
	for k, v := range m.execsMap {
		for _, e := range v {
			e.Stop()
		}
		delete(m.execsMap, k)
	}
	mutex.Unlock()
}

func triggerExecutors(deviceName string, autoEvents []contract.AutoEvent) []Executor {
	var execs []Executor
	for _, autoEvent := range autoEvents {
		exec, err := NewExecutor(deviceName, autoEvent)
		if err != nil {
			common.LoggingClient.Error(fmt.Sprintf("AutoEvent for resource %s cannot be created,%v", autoEvent.Resource, err))
			continue
		}

		execs = append(execs, exec)
		go exec.Run()
	}
	return execs
}

func (m *manager) RestartForDevice(deviceName string) {
	m.StopForDevice(deviceName)
	d, ok := cache.Devices().ForName(deviceName)
	if !ok {
		common.LoggingClient.Error(fmt.Sprintf("there is no Device %s in cache to start AutoEvent", deviceName))
	}

	mutex.Lock()
	execs := triggerExecutors(deviceName, d.AutoEvents)
	m.execsMap[deviceName] = execs
	mutex.Unlock()
}

func (m *manager) StopForDevice(deviceName string) {
	mutex.Lock()
	execs, ok := m.execsMap[deviceName]
	if ok {
		for _, e := range execs {
			e.Stop()
		}
		delete(m.execsMap, deviceName)
	}
	mutex.Unlock()
}

func GetManager() Manager {
	createOnce.Do(func() {
		m = &manager{execsMap: make(map[string][]Executor)}
	})
	return m
}
