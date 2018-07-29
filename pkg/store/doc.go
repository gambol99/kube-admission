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

// Package store is the code used to maintain a document store of objects
// of interest. The store pulls in resources from the kubernetes api / informers
// and populates the store with resources. All resources are referenced by the
// api group, kind and name i.e 'core/v1/namespace/test'. For core/v1 this can be
// shortened to 'namespace/test'
package store

/*
  # list the namespaces in the store
	items = store.list("namespaces");
	# list all the services in a namespace
	store.namespace("default").kind("pods").get("pod1")
	store.namespace("all").kind("services").list
	store.kind("nodes").list()
	store.kind("namespaces).get("default")
	store.kind("namespace").set("
*/
