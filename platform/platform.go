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
	"path/filepath"

	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/logger"
)

// Platform represents the platform contributions for an application.
type Platform struct {
	// Root is the path to the root directory for the platform contributions.
	Root string

	// Envs is the collection of environment variables contributed by the platform.
	Envs EnvironmentVariables

	// Logger is used to write debug and info to the console.
	Logger logger.Logger
}

// String makes Platform satisfy the Stringer interface.
func (p Platform) String() string {
	return fmt.Sprintf("Platform{ Root: %s, Envs: %s, Logger: %s }", p.Root, p.Envs, p.Logger)
}

// DefaultPlatform creates a new instance of Platform, extracting the Root path from os.Args[1].
func DefaultPlatform(logger logger.Logger) (Platform, error) {
	root, err := internal.OsArgs(1)
	if err != nil {
		return Platform{}, err
	}

	if logger.IsDebugEnabled() {
		contents, err := internal.DirectoryContents(root)
		if err != nil {
			return Platform{}, err
		}
		logger.Debug("Platform contents: %s", contents)
	}

	envs, err := enumerateEnvs(root, logger)
	if err != nil {
		return Platform{}, err
	}

	return Platform{
		Root:   root,
		Envs:   envs,
		Logger: logger,
	}, err
}

func enumerateEnvs(root string, logger logger.Logger) (EnvironmentVariables, error) {
	files, err := filepath.Glob(filepath.Join(root, "env", "*"))
	if err != nil {
		return nil, err
	}

	var e EnvironmentVariables

	for _, file := range files {
		e = append(e, EnvironmentVariable{file, filepath.Base(file), logger})
	}

	logger.Debug("Platform environment variables: %s", e)

	return e, nil
}
