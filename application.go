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

package libbuildpack

import (
	"fmt"
	"os"
	"path/filepath"
)

// Application represents the application being processed by buildpacks.
type Application struct {
	// Root is the path to the root directory of the application.
	Root string
}

// String makes Application satisfy the Stringer interface.
func (a Application) String() string {
	return fmt.Sprintf("Application{ root: %s }", a.Root)
}

func (a Application) applicationContents() []string {
	var contents []string

	filepath.Walk(a.Root, func(path string, info os.FileInfo, err error) error {
		contents = append(contents, path)
		return nil
	})

	return contents
}

// DefaultApplication creates a new instance of application, extracting the Root path from the working directory.
func DefaultApplication(logger Logger) (Application, error) {
	root, err := os.Getwd()
	if err != nil {
		return Application{}, err
	}

	return NewApplication(root, logger), nil
}

// NewApplication creates a new instance of Application, configuring the Root path.
func NewApplication(root string, logger Logger) Application {
	a := Application{root}

	if logger.IsDebugEnabled() {
		logger.Debug("Application contents: %s", a.applicationContents())
	}

	return a
}
