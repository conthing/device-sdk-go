// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2017-2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"net/http"

	"github.com/conthing/device-sdk-go/sdk/common"
	"github.com/conthing/device-sdk-go/sdk/controller/correlation"
	"github.com/gorilla/mux"
)

func InitRestRoutes() *mux.Router {
	r := mux.NewRouter().PathPrefix(common.APIv1Prefix).Subrouter()

	common.LoggingClient.Debug("init status rest controller")
	r.HandleFunc("/ping", statusFunc)

	// common.LoggingClient.Debug("init version rest controller")
	// r.HandleFunc(common.APIVersionRoute, versionFunc).Methods(http.MethodGet)

	common.LoggingClient.Debug("init command rest controller")
	sr := r.PathPrefix("/device").Subrouter()
	sr.HandleFunc("/all/{command}", commandAllFunc).Methods(http.MethodGet, http.MethodPut)
	sr.HandleFunc("/{id}/{command}", commandFunc).Methods(http.MethodGet, http.MethodPut)
	sr.HandleFunc("/name/{name}/{command}", commandFunc).Methods(http.MethodGet, http.MethodPut)

	common.LoggingClient.Debug("init callback rest controller")
	r.HandleFunc("/callback", callbackFunc)

	common.LoggingClient.Debug("init other rest controller")
	r.HandleFunc("/discovery", discoveryFunc).Methods(http.MethodPost)
	r.HandleFunc("debug/transformData/{transformData}", transformFunc).Methods(http.MethodGet)

	common.LoggingClient.Debug("init the metrics and config rest controller each")
	r.HandleFunc("metrics", metricsHandler).Methods(http.MethodGet)
	r.HandleFunc("config", configHandler).Methods(http.MethodGet)

	r.Use(correlation.ManageHeader)
	r.Use(correlation.OnResponseComplete)
	r.Use(correlation.OnRequestBegin)

	return r
}
