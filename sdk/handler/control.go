package handler

import (
	"fmt"

	"github.com/conthing/device-sdk-go/internal/common"
)

func DiscoveryHandler(requestMap map[string]string) {
	common.LoggingClient.Info(fmt.Sprintf("service:discovery request"))
}

func TransformHandler(requestMap map[string]string) (map[string]string, common.AppError) {
	common.LoggingClient.Info(fmt.Sprintf("service : transform request: transformData: %s", requestMap["transformData"]))
	return requestMap, nil
}
