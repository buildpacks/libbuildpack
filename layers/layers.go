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

package layers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/logger"
)

// Layers represents the layers for an application.
type Layers struct {
	// Root is the path to the root directory for the layers.
	Root string

	// Logger is used to write debug and info to the console
	Logger logger.Logger
}

// Layer creates a Layer with a specified name.
func (l Layers) Layer(name string) Layer {
	metadata := filepath.Join(l.Root, fmt.Sprintf("%s.toml", name))
	return Layer{filepath.Join(l.Root, name), l.Logger, metadata}
}

// String makes Layers satisfy the Stringer interface.
func (l Layers) String() string {
	return fmt.Sprintf("Layers{ Root: %s, Logger: %s }", l.Root, l.Logger)
}

// WriteMetadata writes launch metadata to the filesystem.
func (l Layers) WriteMetadata(metadata Metadata) error {
	m, err := internal.ToTomlString(metadata)
	if err != nil {
		return err
	}

	f := filepath.Join(l.Root, "launch.toml")

	l.Logger.Debug("Writing launch metadata: %s <= %s", f, m)
	return internal.WriteToFile(strings.NewReader(m), f, 0644)
}
