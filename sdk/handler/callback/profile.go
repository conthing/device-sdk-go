package callback

import (
	"context"
	"fmt"
	"net/http"

	"github.com/conthing/device-sdk-go/internal/cache"
	"github.com/conthing/device-sdk-go/internal/provision"

	"github.com/conthing/device-sdk-go/internal/common"
	"github.com/google/uuid"
)

func handlerProfile(method string, id string) common.AppError {
	if method == http.MethodPut {
		ctx := context.WithValue(context.Background(), common.CorrelationHeader, uuid.New().String())
		profile, err := common.DeviceProfileClient.DeviceProfile(id, ctx)
		if err != nil {
			appErr := common.NewBadRequestError(err.Error(), err)
			common.LoggingClient.Error(fmt.Sprintf("Cannot find the device profile %s from Core Metadata: %v", id, err))
			return appErr
		}

		err = cache.Profiles().Update(profile)
		if err == nil {
			provision.CreateDescriptorsFromProfile(&profile)
			common.LoggingClient.Info(fmt.Sprintf("Updated device profile %s", id))
		} else {
			appErr := common.NewServerError(err.Error(), err)
			common.LoggingClient.Error(fmt.Sprintf("Couldn't update device profile %s: %v", id, err.Error()))
			return appErr
		}
	} else {
		common.LoggingClient.Error(fmt.Sprintf("Invalid device profile method: %s", method))
		appErr := common.NewBadRequestError("Invalid device profile method", nil)
		return appErr
	}
	return nil
}
