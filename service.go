package device

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/conthing/device-sdk-go/sdk/cache"
	"github.com/conthing/device-sdk-go/sdk/clients"
	"github.com/conthing/device-sdk-go/sdk/controller"
	"github.com/conthing/device-sdk-go/sdk/provision"

	"github.com/conthing/device-sdk-go/sdk/config"

	"github.com/conthing/device-sdk-go/pkg/models"
	"github.com/conthing/device-sdk-go/sdk/common"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/types"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

var svc *Service

type Service struct {
	svcInfo      *common.ServiceInfo
	discovery    models.ProtocolDiscovery
	initAttempts int
	initialized  bool
	stopped      bool
	cw           *Watchers
	asyncCh      chan *models.AsyncValues
	startTime    time.Time
}

func (s *Service) Name() string {
	return common.ServiceName
}

func (s *Service) Version() string {
	return common.ServiceVersion
}

func (s *Service) Start(errChan chan error) (err error) {
	err = clients.InitDependencyClients()

	if err != nil {
		return err
	}

	if config.RegistryClient != nil {
		go config.ListenForConfigChange()
	}

	err = selfRegistry()
	if err != nil {
		return fmt.Errorf("Couldn't register to metadata service")
	}

	// initialize devices, objects & profiles
	cache.InitCache()
	err = provision.LoadProfiles(common.CurrentConfig.Device.ProfilesDir)
	if err != nil {
		return fmt.Errorf("Failed to create the pre-defined Device Profiles")
	}

	err = provision.LoadDevices(common.CurrentConfig.DeviceList)
	if err != nil {
		return fmt.Errorf("Failed to create the pre-defined Devices")
	}

	s.cw = newWatchers()

	// initialize driver
	// if common.CurrentConfig.Service.EnableAsyncReadings {
	// 	s.asyncCh = make(chan *dsModels.AsyncValues, common.CurrentConfig.Service.AsyncBufferSize)
	// 	//go processAsyncResults()
	// }
	err = common.Driver.Initialize(common.LoggingClient, s.asyncCh)
	if err != nil {
		return fmt.Errorf("Driver.Initialize failure: %v", err)
	}

	// Setup REST API
	r := controller.InitRestRoutes()

	//autoevent.GetManager().StartAutoEvents()
	http.TimeoutHandler(nil, time.Millisecond*time.Duration(s.svcInfo.Timeout), "Request timed out")

	// TODO: call ListenAndServe in a goroutine

	common.LoggingClient.Info(fmt.Sprintf("*Service Start() called, name=%s, version=%s", common.ServiceName, common.ServiceVersion))

	go func() {
		errChan <- http.ListenAndServe(common.Colon+strconv.Itoa(s.svcInfo.Port), r)
	}()

	common.LoggingClient.Info("Listening on port: " + strconv.Itoa(common.CurrentConfig.Service.Port))
	common.LoggingClient.Info("Service started in: " + time.Since(s.startTime).String())

	common.LoggingClient.Debug("*Service Start() exit")

	return err

}

func NewService(serviceName string, serviceVersion string, confProfile string, confDir string, useRegistry string, proto models.ProtocolDriver) (*Service, error) {
	startTime := time.Now()
	if svc != nil {
		err := fmt.Errorf("NewService: service already exist!\n")
		return nil, err
	}
	if len(serviceName) == 0 {
		err := fmt.Errorf("NewService: empty name specified\n")
		return nil, err
	}
	common.ServiceName = serviceName

	config, err := config.LoadConfig(useRegistry, confProfile, confDir)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error loading config file:%v\n", err)
		os.Exit(1)
	}

	common.CurrentConfig = config

	if len(serviceVersion) == 0 {
		err := fmt.Errorf("NewService: empty version number specified\n")
		return nil, err
	}

	common.ServiceVersion = serviceVersion

	if proto == nil {
		err := fmt.Errorf("NewService: no Driver specified\n")
		return nil, err
	}

	svc = &Service{}
	svc.startTime = startTime
	svc.svcInfo = &config.Service
	common.Driver = proto
	return svc, nil
}

func selfRegistry() error {
	common.LoggingClient.Debug("Trying to find Device Service: " + common.ServiceName)

	ctx := context.WithValue(context.Background(), common.CorrelationHeader, uuid.New().String())
	ds, err := common.DeviceServiceClient.DeviceServiceForName(common.ServiceName, ctx)

	if err != nil {
		if errsc, ok := err.(*types.ErrServiceClient); ok && (errsc.StatusCode == http.StatusNotFound) {
			common.LoggingClient.Info(fmt.Sprintf("Device Service %s doesn't exist ,creating a new one ", ds.Name))
			ds, err = createNewDeviceService()
		} else {
			common.LoggingClient.Error(fmt.Sprintf("DeviceServicForName failed: %v", err))
			return err
		}
	} else {
		common.LoggingClient.Info(fmt.Sprintf("Device Service %s exists", ds.Name))
	}
	common.LoggingClient.Debug(fmt.Sprintf("Device Service in Core MetaData: %s", ds.Name))
	common.CurrentDeviceService = ds
	svc.initialized = true
	return nil
}

func createNewDeviceService() (contract.DeviceService, error) {
	addr, err := makeNewAddressable()
	if err != nil {
		common.LoggingClient.Error(fmt.Sprintf("makeNewAddressable failed: %v", err))
		return contract.DeviceService{}, err
	}
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	ds := contract.DeviceService{
		Name:           common.ServiceName,
		Labels:         svc.svcInfo.Labels,
		OperatingState: "ENABLED",
		Addressable:    *addr,
		AdminState:     "UNLOCKED",
	}
	ds.Origin = millis

	ctx := context.WithValue(context.Background(), common.CorrelationHeader, uuid.New().String())
	id, err := common.DeviceServiceClient.Add(&ds, ctx)
	if err != nil {
		common.LoggingClient.Error(fmt.Sprintf("Add Deviceservice: %s; failed: %v", common.ServiceName, err))
		return contract.DeviceService{}, err
	}
	if err = common.VerifyIdFormat(id, "Device Service"); err != nil {
		return contract.DeviceService{}, err
	}

	// NOTE - this differs from Addressable and Device objects,
	// neither of which require the '.Service'prefix
	ds.Id = id
	common.LoggingClient.Debug("New deviceservice Id: " + ds.Id)

	return ds, nil
}

func makeNewAddressable() (*contract.Addressable, error) {
	// check whether there has been an existing addressable
	ctx := context.WithValue(context.Background(), common.CorrelationHeader, uuid.New().String())
	addr, err := common.AddressableClient.AddressableForName(common.ServiceName, ctx)
	if err != nil {
		if errsc, ok := err.(*types.ErrServiceClient); ok && (errsc.StatusCode == http.StatusNotFound) {
			common.LoggingClient.Info(fmt.Sprintf("Addressable %s doesn't exist, creating a new one", common.ServiceName))
			millis := time.Now().UnixNano() / int64(time.Millisecond)
			addr = contract.Addressable{
				Timestamps: contract.Timestamps{
					Origin: millis,
				},
				Name:       common.ServiceName,
				HTTPMethod: http.MethodPost,
				Protocol:   common.HttpProto,
				Address:    svc.svcInfo.Host,
				Port:       svc.svcInfo.Port,
				Path:       common.APICallbackRoute,
			}
			id, err := common.AddressableClient.Add(&addr, ctx)
			if err != nil {
				common.LoggingClient.Error(fmt.Sprintf("Add addressable failed %s, error: %v", addr.Name, err))
				return nil, err
			}
			if err = common.VerifyIdFormat(id, "Addressable"); err != nil {
				return nil, err
			}
			addr.Id = id
		} else {
			common.LoggingClient.Error(fmt.Sprintf("AddressableForName failed: %v", err))
			return nil, err
		}
	} else {
		common.LoggingClient.Info(fmt.Sprintf("Addressable %s exists", common.ServiceName))
	}

	return &addr, nil
}
