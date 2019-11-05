package callback

import (
	"context"
	"fmt"
	"net/http"

	"github.com/conthing/device-sdk-go/sdk/autoevent"

	"github.com/conthing/device-sdk-go/sdk/provision"

	"github.com/conthing/device-sdk-go/sdk/cache"

	"github.com/conthing/device-sdk-go/sdk/common"
	"github.com/google/uuid"
)

func handlerDevice(method string, id string) common.AppError {
	ctx := context.WithValue(context.Background(), common.CorrelationHeader, uuid.New().String())
	if method == http.MethodPost {
		device, err := common.DeviceClient.Device(id, ctx)
		if err != nil {
			appErr := common.NewBadRequestError(err.Error(), err)
			common.LoggingClient.Error(fmt.Sprintf("Cannot find the device %s from Core Metadata: %v", id, err))
			return appErr
		}

		_, exist := cache.Profiles().ForName(device.Profile.Name)
		if exist == false {
			err = cache.Profiles().Add(device.Profile)
			if err == nil {
				provision.CreateDescriptorsFromProfile(&device.Profile)
				common.LoggingClient.Info(fmt.Sprintf("Added device profile %s", device.Profile.Id))
			} else {
				appErr := common.NewServerError(err.Error(), err)
				common.LoggingClient.Error(fmt.Sprintf("Couldn't add device profile %s: %v", device.Profile.Name, err.Error()))
				return appErr
			}
		}

		err = cache.Devices().Add(device)
		if err == nil {
			common.LoggingClient.Info(fmt.Sprintf("Added device %s", id))
		} else {
			appErr := common.NewServerError(err.Error(), err)
			common.LoggingClient.Error(fmt.Sprintf("Couldn't add device %s: %v", id, err.Error()))
			return appErr
		}

		common.LoggingClient.Debug(fmt.Sprintf("Handler - starting AutoEvents for device %s", device.Name))
		autoevent.GetManager().RestartForDevice(device.Name)
	} else if method == http.MethodPut {
		device, err := common.DeviceClient.Device(id, ctx)
		if err != nil {
			appErr := common.NewBadRequestError(err.Error(), err)
			common.LoggingClient.Error(fmt.Sprintf("Cannot find the device %s from Core Metadata: %v", id, err.Error()))
			return appErr
		}

		err = cache.Devices().Update(device)
		if err == nil {
			common.LoggingClient.Info(fmt.Sprintf("Updated device %s", id))
		} else {
			appErr := common.NewServerError(err.Error(), err)
			common.LoggingClient.Error(fmt.Sprintf("Couldn't update device %s: %v", id, err.Error()))
			return appErr
		}

		common.LoggingClient.Debug(fmt.Sprintf("Handler - restarting AutoEvents for device %s", device.Name))
		autoevent.GetManager().RestartForDevice(device.Name)
	} else if method == http.MethodDelete {
		if device, ok := cache.Devices().ForId(id); ok {
			common.LoggingClient.Debug(fmt.Sprintf("Handler - stopping AutoEvents for device %s", device.Name))
			autoevent.GetManager().StopForDevice(device.Name)
		}

		err := cache.Devices().Remove(id)
		if err == nil {
			common.LoggingClient.Info(fmt.Sprintf("Remove device %s", id))
		} else {
			appErr := common.NewServerError(err.Error(), err)
			common.LoggingClient.Error(fmt.Sprintf("Couldn't remove device %s: %v", id, err.Error()))
			return appErr
		}
	} else {
		common.LoggingClient.Error(fmt.Sprintf("Invalid device method type :%s", method))
		appErr := common.NewServerError("Invalid device method", nil)
		return appErr
	}

	return nil
}
