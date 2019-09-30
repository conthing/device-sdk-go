package endpoint

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/types"
	"github.com/edgexfoundry/go-mod-registry/registry"
)

type Endpoint struct {
	RegistryClient registry.Client
	passFirstRun   bool
	WG             *sync.WaitGroup
}

func (endpoint Endpoint) Monitor(params types.EndpointParams, ch chan string) {
	for {
		data, err := endpoint.RegistryClient.GetServiceEndpoint(params.ServiceKey)
		if err != nil {
			fmt.Fprintln(os.Stdout, err.Error())
		}
		url := fmt.Sprintf("http://%s:%v%s", data.Host, data.Port, params.Path)
		ch <- url

		if !endpoint.passFirstRun {
			endpoint.WG.Done()
			endpoint.passFirstRun = true
		}
		time.Sleep(time.Second * time.Duration(params.Interval))
	}
}
