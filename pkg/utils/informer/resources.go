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

package informer

import (
	"fmt"

	v1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	v1beta1 "k8s.io/api/admissionregistration/v1beta1"
	v1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	v1beta2 "k8s.io/api/apps/v1beta2"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	v2beta1 "k8s.io/api/autoscaling/v2beta1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v2alpha1 "k8s.io/api/batch/v2alpha1"
	certificatesv1beta1 "k8s.io/api/certificates/v1beta1"
	coordinationv1beta1 "k8s.io/api/coordination/v1beta1"
	corev1 "k8s.io/api/core/v1"
	eventsv1beta1 "k8s.io/api/events/v1beta1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	rbacv1alpha1 "k8s.io/api/rbac/v1alpha1"
	rbacv1beta1 "k8s.io/api/rbac/v1beta1"
	schedulingv1alpha1 "k8s.io/api/scheduling/v1alpha1"
	schedulingv1beta1 "k8s.io/api/scheduling/v1beta1"
	settingsv1alpha1 "k8s.io/api/settings/v1alpha1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1alpha1 "k8s.io/api/storage/v1alpha1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
)

// ResourceNames returns a list of supported resource versions
func ResourceNames() []string {
	var list []string
	for k, _ := range ResourceVersions() {
		list = append(list, k)
	}

	return list
}

// ResourceVersion returns a map of supported resources versions
func ResourceVersions() map[string]schema.GroupVersionResource {
	return map[string]schema.GroupVersionResource{
		niceVersion(appsv1beta1.SchemeGroupVersion.WithResource("controllerrevisions")):                appsv1beta1.SchemeGroupVersion.WithResource("controllerrevisions"),
		niceVersion(appsv1beta1.SchemeGroupVersion.WithResource("deployments")):                        appsv1beta1.SchemeGroupVersion.WithResource("deployments"),
		niceVersion(appsv1beta1.SchemeGroupVersion.WithResource("statefulsets")):                       appsv1beta1.SchemeGroupVersion.WithResource("statefulsets"),
		niceVersion(autoscalingv1.SchemeGroupVersion.WithResource("horizontalpodautoscalers")):         autoscalingv1.SchemeGroupVersion.WithResource("horizontalpodautoscalers"),
		niceVersion(batchv1.SchemeGroupVersion.WithResource("jobs")):                                   batchv1.SchemeGroupVersion.WithResource("jobs"),
		niceVersion(batchv1beta1.SchemeGroupVersion.WithResource("cronjobs")):                          batchv1beta1.SchemeGroupVersion.WithResource("cronjobs"),
		niceVersion(certificatesv1beta1.SchemeGroupVersion.WithResource("certificatesigningrequests")): certificatesv1beta1.SchemeGroupVersion.WithResource("certificatesigningrequests"),
		niceVersion(coordinationv1beta1.SchemeGroupVersion.WithResource("leases")):                     coordinationv1beta1.SchemeGroupVersion.WithResource("leases"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("componentstatuses")):                       corev1.SchemeGroupVersion.WithResource("componentstatuses"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("configmaps")):                              corev1.SchemeGroupVersion.WithResource("configmaps"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("endpoints")):                               corev1.SchemeGroupVersion.WithResource("endpoints"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("events")):                                  corev1.SchemeGroupVersion.WithResource("events"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("limitranges")):                             corev1.SchemeGroupVersion.WithResource("limitranges"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("namespaces")):                              corev1.SchemeGroupVersion.WithResource("namespaces"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("nodes")):                                   corev1.SchemeGroupVersion.WithResource("nodes"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("persistentvolumeclaims")):                  corev1.SchemeGroupVersion.WithResource("persistentvolumeclaims"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("persistentvolumes")):                       corev1.SchemeGroupVersion.WithResource("persistentvolumes"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("pods")):                                    corev1.SchemeGroupVersion.WithResource("pods"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("podtemplates")):                            corev1.SchemeGroupVersion.WithResource("podtemplates"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("replicationcontrollers")):                  corev1.SchemeGroupVersion.WithResource("replicationcontrollers"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("resourcequotas")):                          corev1.SchemeGroupVersion.WithResource("resourcequotas"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("secrets")):                                 corev1.SchemeGroupVersion.WithResource("secrets"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("serviceaccounts")):                         corev1.SchemeGroupVersion.WithResource("serviceaccounts"),
		niceVersion(corev1.SchemeGroupVersion.WithResource("services")):                                corev1.SchemeGroupVersion.WithResource("services"),
		niceVersion(eventsv1beta1.SchemeGroupVersion.WithResource("events")):                           eventsv1beta1.SchemeGroupVersion.WithResource("events"),
		niceVersion(extensionsv1beta1.SchemeGroupVersion.WithResource("daemonsets")):                   extensionsv1beta1.SchemeGroupVersion.WithResource("daemonsets"),
		niceVersion(extensionsv1beta1.SchemeGroupVersion.WithResource("deployments")):                  extensionsv1beta1.SchemeGroupVersion.WithResource("deployments"),
		niceVersion(extensionsv1beta1.SchemeGroupVersion.WithResource("ingresses")):                    extensionsv1beta1.SchemeGroupVersion.WithResource("ingresses"),
		niceVersion(extensionsv1beta1.SchemeGroupVersion.WithResource("podsecuritypolicies")):          extensionsv1beta1.SchemeGroupVersion.WithResource("podsecuritypolicies"),
		niceVersion(networkingv1.SchemeGroupVersion.WithResource("networkpolicies")):                   networkingv1.SchemeGroupVersion.WithResource("networkpolicies"),
		niceVersion(policyv1beta1.SchemeGroupVersion.WithResource("poddisruptionbudgets")):             policyv1beta1.SchemeGroupVersion.WithResource("poddisruptionbudgets"),
		niceVersion(policyv1beta1.SchemeGroupVersion.WithResource("podsecuritypolicies")):              policyv1beta1.SchemeGroupVersion.WithResource("podsecuritypolicies"),
		niceVersion(rbacv1.SchemeGroupVersion.WithResource("clusterrolebindings")):                     rbacv1.SchemeGroupVersion.WithResource("clusterrolebindings"),
		niceVersion(rbacv1.SchemeGroupVersion.WithResource("clusterroles")):                            rbacv1.SchemeGroupVersion.WithResource("clusterroles"),
		niceVersion(rbacv1.SchemeGroupVersion.WithResource("rolebindings")):                            rbacv1.SchemeGroupVersion.WithResource("rolebindings"),
		niceVersion(rbacv1.SchemeGroupVersion.WithResource("roles")):                                   rbacv1.SchemeGroupVersion.WithResource("roles"),
		niceVersion(rbacv1alpha1.SchemeGroupVersion.WithResource("clusterrolebindings")):               rbacv1alpha1.SchemeGroupVersion.WithResource("clusterrolebindings"),
		niceVersion(rbacv1alpha1.SchemeGroupVersion.WithResource("clusterroles")):                      rbacv1alpha1.SchemeGroupVersion.WithResource("clusterroles"),
		niceVersion(rbacv1alpha1.SchemeGroupVersion.WithResource("rolebindings")):                      rbacv1alpha1.SchemeGroupVersion.WithResource("rolebindings"),
		niceVersion(rbacv1alpha1.SchemeGroupVersion.WithResource("roles")):                             rbacv1alpha1.SchemeGroupVersion.WithResource("roles"),
		niceVersion(rbacv1beta1.SchemeGroupVersion.WithResource("clusterrolebindings")):                rbacv1beta1.SchemeGroupVersion.WithResource("clusterrolebindings"),
		niceVersion(rbacv1beta1.SchemeGroupVersion.WithResource("clusterroles")):                       rbacv1beta1.SchemeGroupVersion.WithResource("clusterroles"),
		niceVersion(rbacv1beta1.SchemeGroupVersion.WithResource("rolebindings")):                       rbacv1beta1.SchemeGroupVersion.WithResource("rolebindings"),
		niceVersion(rbacv1beta1.SchemeGroupVersion.WithResource("roles")):                              rbacv1beta1.SchemeGroupVersion.WithResource("roles"),
		niceVersion(schedulingv1alpha1.SchemeGroupVersion.WithResource("priorityclasses")):             schedulingv1alpha1.SchemeGroupVersion.WithResource("priorityclasses"),
		niceVersion(schedulingv1beta1.SchemeGroupVersion.WithResource("priorityclasses")):              schedulingv1beta1.SchemeGroupVersion.WithResource("priorityclasses"),
		niceVersion(settingsv1alpha1.SchemeGroupVersion.WithResource("podpresets")):                    settingsv1alpha1.SchemeGroupVersion.WithResource("podpresets"),
		niceVersion(storagev1.SchemeGroupVersion.WithResource("storageclasses")):                       storagev1.SchemeGroupVersion.WithResource("storageclasses"),
		niceVersion(storagev1alpha1.SchemeGroupVersion.WithResource("volumeattachments")):              storagev1alpha1.SchemeGroupVersion.WithResource("volumeattachments"),
		niceVersion(storagev1beta1.SchemeGroupVersion.WithResource("storageclasses")):                  storagev1beta1.SchemeGroupVersion.WithResource("storageclasses"),
		niceVersion(storagev1beta1.SchemeGroupVersion.WithResource("volumeattachments")):               storagev1beta1.SchemeGroupVersion.WithResource("volumeattachments"),
		niceVersion(v1.SchemeGroupVersion.WithResource("controllerrevisions")):                         v1.SchemeGroupVersion.WithResource("controllerrevisions"),
		niceVersion(v1.SchemeGroupVersion.WithResource("daemonsets")):                                  v1.SchemeGroupVersion.WithResource("daemonsets"),
		niceVersion(v1.SchemeGroupVersion.WithResource("deployments")):                                 v1.SchemeGroupVersion.WithResource("deployments"),
		niceVersion(v1.SchemeGroupVersion.WithResource("replicasets")):                                 v1.SchemeGroupVersion.WithResource("replicasets"),
		niceVersion(v1.SchemeGroupVersion.WithResource("statefulsets")):                                v1.SchemeGroupVersion.WithResource("statefulsets"),
		niceVersion(v1alpha1.SchemeGroupVersion.WithResource("initializerconfigurations")):             v1alpha1.SchemeGroupVersion.WithResource("initializerconfigurations"),
		niceVersion(v1alpha1.SchemeGroupVersion.WithResource("mutatingwebhookconfigurations")):         v1alpha1.SchemeGroupVersion.WithResource("mutatingwebhookconfigurations"),
		niceVersion(v1beta1.SchemeGroupVersion.WithResource("validatingwebhookconfigurations")):        v1beta1.SchemeGroupVersion.WithResource("validatingwebhookconfigurations"),
		niceVersion(v1beta2.SchemeGroupVersion.WithResource("controllerrevisions")):                    v1beta2.SchemeGroupVersion.WithResource("controllerrevisions"),
		niceVersion(v1beta2.SchemeGroupVersion.WithResource("daemonsets")):                             v1beta2.SchemeGroupVersion.WithResource("daemonsets"),
		niceVersion(v1beta2.SchemeGroupVersion.WithResource("deployments")):                            v1beta2.SchemeGroupVersion.WithResource("deployments"),
		niceVersion(v1beta2.SchemeGroupVersion.WithResource("statefulsets")):                           v1beta2.SchemeGroupVersion.WithResource("statefulsets"),
		niceVersion(v2alpha1.SchemeGroupVersion.WithResource("cronjobs")):                              v2alpha1.SchemeGroupVersion.WithResource("cronjobs"),
		niceVersion(v2beta1.SchemeGroupVersion.WithResource("horizontalpodautoscalers")):               v2beta1.SchemeGroupVersion.WithResource("horizontalpodautoscalers"),
	}
}

// niceVersion returns a more friendly resource string
func niceVersion(version schema.GroupVersionResource) string {
	if version.Group != "" {
		return fmt.Sprintf("%s/%s/%s", version.Group, version.Version, version.Resource)
	}

	return fmt.Sprintf("%s/%s", version.Version, version.Resource)
}
