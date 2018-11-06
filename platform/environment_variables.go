/*
 * Copyright 2018 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package platform

// EnvironmentVariables is a collection of EnvironmentVariable instances.
type EnvironmentVariables []EnvironmentVariable

// Contains returns whether an EnvironmentVariables contains a given variable by name.
func (e EnvironmentVariables) Contains(name string) bool {
	for _, ev := range e {
		if ev.Name == name {
			return true
		}
	}

	return false
}

// SetAll sets all of the environment variable content in the current process environment.
func (e EnvironmentVariables) SetAll() error {
	for _, ev := range e {
		if err := ev.Set(); err != nil {
			return err
		}
	}

	return nil
}
