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

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/buildpack/libbuildpack/logger"
)

// EnvironmentVariable represents an environment variable provided by the platform.
type EnvironmentVariable struct {
	// File is the location of the environment variable's contents
	File string

	// Name is the name of the environment variable
	Name string

	logger logger.Logger
}

// Set sets the environment variable content in the current process environment.
func (e EnvironmentVariable) Set() error {
	value, err := e.value()
	if err != nil {
		return err
	}

	e.logger.Debug("Setting environment variable: %s <= %s", e.Name, value)
	return os.Setenv(e.Name, value)
}

// String makes EnvironmentVariable satisfy the Stringer interface.
func (e EnvironmentVariable) String() string {
	return fmt.Sprintf("EnvironmentVariable{ Name: %s, File: %s, logger: %s }", e.Name, e.File, e.logger)
}

func (e EnvironmentVariable) value() (string, error) {
	b, err := ioutil.ReadFile(e.File)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
