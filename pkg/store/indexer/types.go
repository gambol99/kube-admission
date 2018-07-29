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

package indexer

// Interface is the contract to the index
type Interface interface {
	// Delete remove the id from the index
	Delete(string) error
	// DeleleByQuery removes all the documents which match the query
	DeleteByQuery(string) (int, error)
	// Index is responsible is add a document the index
	Index(string, interface{}) error
	// Search is responsible for searching the index
	Search(string) ([]string, error)
	// Query is responsible for searching the index
	Query(interface{}) ([]string, error)
	// Size returns the size of the index
	Size() (uint64, error)
}
