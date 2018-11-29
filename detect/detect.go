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

package detect

import (
	"fmt"

	applicationPkg "github.com/buildpack/libbuildpack/application"
	buildpackPkg "github.com/buildpack/libbuildpack/buildpack"
	buildplanPkg "github.com/buildpack/libbuildpack/buildplan"
	loggerPkg "github.com/buildpack/libbuildpack/logger"
	platformPkg "github.com/buildpack/libbuildpack/platform"
	stackPkg "github.com/buildpack/libbuildpack/stack"
)

const (
	// FailStatusCode is the status code returned for fail.
	FailStatusCode = 100

	// PassStatusCode is the status code returned for pass.
	PassStatusCode = 0
)

// Detect represents all of the components available to a buildpack at detect time.
type Detect struct {
	// Application is the application being processed by the buildpack.
	Application applicationPkg.Application

	// Buildpack represents the metadata associated with a buildpack.
	Buildpack buildpackPkg.Buildpack

	// BuildPlan represents dependencies contributed by previous builds.
	BuildPlan buildplanPkg.BuildPlan

	// BuildPlanWriter is the writer used to write the BuildPlan in Pass().
	BuildPlanWriter buildplanPkg.Writer

	// Logger is used to write debug and info to the console.
	Logger loggerPkg.Logger

	// Platform represents components contributed by the platform to the buildpack.
	Platform platformPkg.Platform

	// Stack is the stack currently available to the application.
	Stack string
}

// Error signals an error during detection by exiting with a specified non-zero, non-100 status code.
func (d Detect) Error(code int) int {
	d.Logger.Debug("Detection produced an error. Exiting with %d.", code)
	return code
}

// Fail signals an unsuccessful detection by exiting with a 100 status code.
func (d Detect) Fail() int {
	d.Logger.Debug("Detection failed. Exiting with %d.", FailStatusCode)
	return FailStatusCode
}

// Pass signals a successful detection by exiting with a 0 status code.
func (d Detect) Pass(buildPlan buildplanPkg.BuildPlan) (int, error) {
	d.Logger.Debug("Detection passed. Exiting with %d.", PassStatusCode)

	if err := buildPlan.Write(d.BuildPlanWriter); err != nil {
		return -1, err
	}

	return PassStatusCode, nil
}

// String makes Detect satisfy the Stringer interface.
func (d Detect) String() string {
	return fmt.Sprintf("Detect{ Application: %s, Buildpack: %s, BuildPlan: %s, Logger: %s, Platform: %s, Stack: %s }",
		d.Application, d.Buildpack, d.BuildPlan, d.Logger, d.Platform, d.Stack)
}

// DefaultDetect creates a new instance of Detect using default values.
func DefaultDetect() (Detect, error) {
	logger := loggerPkg.DefaultLogger()

	application, err := applicationPkg.DefaultApplication(logger)
	if err != nil {
		return Detect{}, err
	}

	buildpack, err := buildpackPkg.DefaultBuildpack(logger)
	if err != nil {
		return Detect{}, err
	}

	buildPlan := buildplanPkg.BuildPlan{}

	buildPlanWriter := buildplanPkg.DefaultWriter

	platform, err := platformPkg.DefaultPlatform(logger)
	if err != nil {
		return Detect{}, err
	}

	stack, err := stackPkg.DefaultStack(logger)
	if err != nil {
		return Detect{}, err
	}

	return Detect{
		application,
		buildpack,
		buildPlan,
		buildPlanWriter,
		logger,
		platform,
		stack,
	}, nil
}
