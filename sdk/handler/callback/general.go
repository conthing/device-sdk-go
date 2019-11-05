package callback

import (
	"fmt"

	"github.com/conthing/device-sdk-go/sdk/common"
	contract "github.com/edgexfoundry/go-mod-core-contracts/models"
)

func CallbackHandler(cbAlert contract.CallbackAlert, method string) common.AppError {
	if (cbAlert.Id == "") || (cbAlert.ActionType == "") {
		appErr := common.NewBadRequestError("Missing parameters", nil)
		common.LoggingClient.Error(fmt.Sprintf("Missing callback parameters"))
		return appErr
	}

	if cbAlert.ActionType == contract.DEVICE {
		return handlerDevice(method, cbAlert.Id)
	} else if cbAlert.ActionType == contract.PROFILE {
		return handlerProfile(method, cbAlert.Id)
	}

	common.LoggingClient.Error(fmt.Sprintf("Invalid callback action type : %s", cbAlert.ActionType))
	appErr := common.NewBadRequestError("Invalid callback action type", nil)
	return appErr
}
