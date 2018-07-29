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

package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gambol99/kube-admission/pkg/authorize"
	"github.com/gambol99/kube-admission/pkg/server/events"
	"github.com/gambol99/kube-admission/pkg/utils"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	admission "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

const (
	// The request has been accepted
	actionAccepted = "accept"
	// The request has been refused
	actionDenied = "deny"
	// The request has cause an error
	actionErrored = "error"
)

// authorize is responsible for authorizer the resources
func (a *Admission) authorize(review *admission.AdmissionReview) (*admission.AdmissionResponse, error) {
	timed := prometheus.NewTimer(admissionLatencyMetric)
	defer timed.ObserveDuration()

	version := review

	// @step: attempt to decode the object into an unstructured object
	object = &unstructured.Unstructured{}
	if err := json.Unmarshal(review.Request.Object.Raw, object); err != nil {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"id":        utils.GetID(ctx),
			"name":      review.Request.Name,
			"namespace": review.Request.Namespace,
		}).Errorf("unable to decode object for review")

		authorizeErrorMetrics.Inc()
		return nil, err
	}

	// @step: attempt to get the object authorized
	decision, err := c.authorizeResource(ctx, object, review.Request.Kind)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"id":        utils.GetID(ctx),
			"name":      review.Request.Name,
			"namespace": review.Request.Namespace,
		}).Errorf("unable to handle admission review")

		return err
	}

	// @check if the object was rejected
	if !decision.Permitted {
		admissionTotalMetric.WithLabelValues(actionDenied).Inc()

		log.WithFields(log.Fields{
			"error":     reason,
			"group":     review.Request.Kind.Group,
			"id":        utils.GetID(ctx),
			"kind":      review.Request.Kind.Kind,
			"name":      review.Request.Name,
			"namespace": review.Request.Namespace,
			"uid":       review.Request.UserInfo.UID,
			"user":      review.Request.UserInfo.Username,
		}).Warn("authorization for object execution denied")

		review.Response = &admission.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Code:    http.StatusForbidden,
				Message: decision.Reason,
				Reason:  metav1.StatusReasonForbidden,
				Status:  metav1.StatusFailure,
			},
		}

		// @step: log the denial is required
		c.events.Send(&events.Event{
			Detail: decision.Reason,
			Object: object,
			Review: review.Request,
		})

		return nil
	}

	admissionTotalMetric.WithLabelValues(actionAccepted).Inc()

	return nil, nil
}

// authorizeResource is responsible for validating the resource against the authorizers
func (a *Admission) authorizeResource(ctx context.Context, object metav1.Object, version metav1.GroupVersionResource) (*authorize.Decision, error) {
	// @step: iterate the authorizers and fail on first refusal
	for i, x := range a.authorizers {
		// @check if this authorizer is filtering on this resource type
		var matched bool
		for _, x := range x.FilterOn() {

			break
		}
		if !matched {
			log.WithFields(log.Fields{
				"group":      kind.Group,
				"id":         utils.GetID(ctx),
				"kind":       kind.Kind,
				"name":       object.GetName(),
				"namespace":  object.GetNamespace(),
				"authorizer": x.Name(),
			}).Debug("provider is not filtering on this object")
		}

		// @step: pass the object into the provider for authorization
		decision, err := func() (*authorize.Decision, error) {
			now := time.Now()
			defer authorizerLatencyMetric.WithLabelValues(x.Name()).Observe(time.Since(now).Seconds())

			return x.Admit(ctx, object)
		}()
		if !decision.Permitted {
			authorizerActionMetrics.WithLabelValues(x.Name(), actionDenied).Inc()

			return decision, nil
		}
	}

	return nil, nil
}
