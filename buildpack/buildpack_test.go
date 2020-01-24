/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package buildpack_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpacks/libbuildpack/v2/buildpack"
	"github.com/buildpacks/libbuildpack/v2/internal"
	"github.com/buildpacks/libbuildpack/v2/logger"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildpack(t *testing.T) {
	spec.Run(t, "Buildpack", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("unmarshals default from buildpack.toml", func() {
			root := internal.ScratchDir(t, "buildpack")
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"))()

			internal.WriteTestFile(t, filepath.Join(root, "buildpack.toml"), `[buildpack]
id = "buildpack-id"
name = "buildpack-name"
version = "buildpack-version"

[[stacks]]
id = 'stack-id'
build-images = ["build-image-tag"]
run-images = ["run-image-tag"]

[metadata]
test-key = "test-value"
`)

			g.Expect(buildpack.DefaultBuildpack(logger.Logger{})).To(gomega.Equal(buildpack.Buildpack{
				Info: buildpack.Info{
					ID:      "buildpack-id",
					Name:    "buildpack-name",
					Version: "buildpack-version",
				},
				Stacks: []buildpack.Stack{
					{
						ID:          "stack-id",
						BuildImages: buildpack.BuildImages{"build-image-tag"},
						RunImages:   buildpack.RunImages{"run-image-tag"},
					},
				},
				Metadata: buildpack.Metadata{"test-key": "test-value"},
				Root:     root,
			}))
		})
	}, spec.Report(report.Terminal{}))
}
