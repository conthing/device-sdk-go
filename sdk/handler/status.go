// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package handler

import (
	"github.com/conthing/device-sdk-go/sdk/common"
)

func StatusHandler() string {
	return common.ServiceVersion
}
