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

package utils

import (
	"errors"
	"time"

	"github.com/jpillora/backoff"
)

// Retry attempts to perform an operation for x time
func Retry(attempts int, min time.Duration, jitter bool, fn func() error) error {
	b := &backoff.Backoff{Min: min, Factor: 1, Jitter: jitter}

	// @step: give it a go once before jumping in
	for i := 0; i < attempts; i++ {
		if err := fn(); err == nil {
			return nil
		}

		time.Sleep(b.Duration())
	}

	return errors.New("operation failed")
}
