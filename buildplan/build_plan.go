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

package buildplan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/buildpack/libbuildpack/internal"
)

// BuildPlan represents the dependencies contributed by a build.  Note that you may need to call Init() to load contents
// from os.Stdin.
type BuildPlan map[string]Dependency

// Init initializes the BuildPlan by reading os.Stdin.  Will block until os.Stdin is closed.
func (b BuildPlan) Init() error {
	if _, err := toml.DecodeReader(os.Stdin, &b); err != nil {
		return err
	}

	return nil
}

// Write writes the build plan to a collection files rooted at os.Args[2].
func (b BuildPlan) Write() error {
	root, err := internal.OsArgs(2)
	if err != nil {
		return err
	}

	for name, dep := range b {
		s, err := internal.ToTomlString(dep)
		if err != nil {
			return err
		}

		if err := internal.WriteToFile(strings.NewReader(s), filepath.Join(root, name), 0644); err != nil {
			return err
		}
	}

	return nil
}

// String makes BuildPlan satisfy the Stringer interface.
func (b BuildPlan) String() string {
	var entries []string

	for k, v := range b {
		entries = append(entries, fmt.Sprintf("%s: %s", k, v))
	}

	return fmt.Sprintf("BuildPlan{ %s }", strings.Join(entries, ", "))
}
