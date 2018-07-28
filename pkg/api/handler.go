/*
Copyright 2018 Rohith Jayawardene <gambol99@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"net/http"
	"time"

	"github.com/gambol99/kube-admission/pkg/utils"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	admission "k8s.io/api/admission/v1beta1"
)

// admissionHandler is responsible for handling inbound authorization requests
func (a *Server) admissionHandler(ctx echo.Context) error {
	review := &admission.AdmissionReview{}

	// @step: we need to unmarshal the review
	err := ctx.Bind(review)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("unable to decode the request")

		return ctx.NoContent(http.StatusBadRequest)
	}

	// @step: inject a request correlation id into the context
	id := utils.SetID(ctx.Request().Context(), uuid.NewV1().String())

	// @step: apply the policy against the review
	now := time.Now()
	review.Response, err = a.config.AdmissionHandlerFunc(id, review)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("unable to validation request against policy")

		return ctx.NoContent(http.StatusInternalServerError)
	}
	requestLatencyMetric.Observe(time.Since(now).Seconds())

	return ctx.JSON(http.StatusOK, review)
}

// healthHandler is responsible for return the service health
func (a *Server) healthHandler(ctx echo.Context) error {
	status, err := a.config.HealthHandlerFunc(ctx.Request().Context())
	if err != nil {
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.String(http.StatusOK, status)
}
