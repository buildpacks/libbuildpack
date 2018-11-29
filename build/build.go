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

package build

import (
	"fmt"

	applicationPkg "github.com/buildpack/libbuildpack/application"
	buildpackPkg "github.com/buildpack/libbuildpack/buildpack"
	buildplanPkg "github.com/buildpack/libbuildpack/buildplan"
	layersPkg "github.com/buildpack/libbuildpack/layers"
	loggerPkg "github.com/buildpack/libbuildpack/logger"
	platformPkg "github.com/buildpack/libbuildpack/platform"
	stackPkg "github.com/buildpack/libbuildpack/stack"
)

// Build represents all of the components available to a buildpack at build time.
type Build struct {
	// Application is the application being processed by the buildpack.
	Application applicationPkg.Application

	// Buildpack represents the metadata associated with a buildpack.
	Buildpack buildpackPkg.Buildpack

	// BuildPlan represents dependencies contributed by previous builds.
	BuildPlan buildplanPkg.BuildPlan

	// BuildPlanWriter is the writer used to write the BuildPlan in Pass().
	BuildPlanWriter buildplanPkg.Writer

	// Layers represents the launch layers contributed by a buildpack.
	Layers layersPkg.Layers

	// Logger is used to write debug and info to the console.
	Logger loggerPkg.Logger

	// Platform represents components contributed by the platform to the buildpack.
	Platform platformPkg.Platform

	// Stack is the stack currently available to the application.
	Stack string
}

// Failure signals an unsuccessful build by exiting with a specified positive status code.
func (b Build) Failure(code int) int {
	b.Logger.Debug("Build failed. Exiting with %d.", code)
	b.Logger.Info("")
	return code
}

// String makes Build satisfy the Stringer interface.
func (b Build) String() string {
	return fmt.Sprintf("Build{ Application: %s, Buildpack: %s, BuildPlan: %s, Layers: %s, Logger: %s, Platform: %s, Stack: %s }",
		b.Application, b.Buildpack, b.BuildPlan, b.Layers, b.Logger, b.Platform, b.Stack)
}

// Success signals a successful build by exiting with a zero status code.
func (b Build) Success(buildPlan buildplanPkg.BuildPlan) (int, error) {
	b.Logger.Debug("Build success. Exiting with %d.", 0)

	if err := buildPlan.Write(b.BuildPlanWriter); err != nil {
		return -1, err
	}

	return 0, nil
}

// DefaultBuild creates a new instance of Build using default values.
func DefaultBuild() (Build, error) {
	logger := loggerPkg.DefaultLogger()

	application, err := applicationPkg.DefaultApplication(logger)
	if err != nil {
		return Build{}, err
	}

	buildpack, err := buildpackPkg.DefaultBuildpack(logger)
	if err != nil {
		return Build{}, err
	}

	buildPlan := buildplanPkg.BuildPlan{}
	if err := buildPlan.Init(); err != nil {
		return Build{}, err
	}

	buildPlanWriter := buildplanPkg.DefaultWriter

	layers, err := layersPkg.DefaultLayers(logger)
	if err != nil {
		return Build{}, err
	}

	platform, err := platformPkg.DefaultPlatform(logger)
	if err != nil {
		return Build{}, err
	}

	s, err := stackPkg.DefaultStack(logger)
	if err != nil {
		return Build{}, err
	}

	return Build{
		application,
		buildpack,
		buildPlan,
		buildPlanWriter,
		layers,
		logger,
		platform,
		s,
	}, nil
}
