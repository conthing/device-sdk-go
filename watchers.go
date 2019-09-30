package device

import (
	"sync"

	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type Watchers struct {
	devices map[string]models.Device
}

var (
	wcOnce   sync.Once
	watchers *Watchers
)

func newWatchers() *Watchers {
	wcOnce.Do(func() {
		watchers = &Watchers{}
	})
	return watchers
}
