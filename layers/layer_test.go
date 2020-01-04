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

package layers_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/layers"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLayer(t *testing.T) {
	spec.Run(t, "Layer", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		when("layer content metadata", func() {

			type metadata struct {
				Alpha string
				Bravo int
			}

			var (
				root  string
				layer layers.Layer
			)

			it.Before(func() {
				root = internal.ScratchDir(t, "layer")
				layer = layers.Layers{Root: root}.Layer("test-layer")
			})

			it("reads layer content metadata", func() {
				internal.WriteTestFile(t, filepath.Join(root, "test-layer.toml"), `[metadata]
Alpha = "test-value"
Bravo = 1
`)

				var actual metadata
				g.Expect(layer.ReadMetadata(&actual)).To(gomega.Succeed())

				g.Expect(actual).To(gomega.Equal(metadata{"test-value", 1}))
			})

			it("does not read layer content metadata if it does not exist", func() {
				var actual metadata
				g.Expect(layer.ReadMetadata(&actual)).To(gomega.Succeed())

				g.Expect(actual).To(gomega.Equal(metadata{}))
			})

			it("remove layer content metadata", func() {
				internal.WriteTestFile(t, filepath.Join(root, "test-layer.toml"), `[metadata]
Alpha = "test-value"
Bravo = 1
`)

				g.Expect(layer.RemoveMetadata()).To(gomega.Succeed())
				g.Expect("").NotTo(gomega.BeAnExistingFile())
			})

			it("writes layer content metadata", func() {
				g.Expect(layer.WriteMetadata(metadata{"test-value", 1},
					layers.Build, layers.Cache, layers.Launch)).To(gomega.Succeed())

				g.Expect(filepath.Join(root, "test-layer.toml")).To(internal.HaveContent(`build = true
cache = true
launch = true

[metadata]
  Alpha = "test-value"
  Bravo = 1
`))
			})

			it("writes a profile file", func() {
				g.Expect(layer.WriteProfile("test-name", "%s-%d", "test-string", 1)).To(gomega.Succeed())

				g.Expect(filepath.Join(root, "test-layer", "profile.d", "test-name")).To(internal.HaveContent("test-string-1"))
			})
		})

		when("environment files", func() {

			var (
				root  string
				layer layers.Layer
			)

			it.Before(func() {
				root = internal.ScratchDir(t, "layer")
				layer = layers.Layer{Root: root}
			})

			when("build", func() {

				it("writes an append environment file", func() {
					g.Expect(layer.AppendBuildEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.build", "TEST_NAME.append")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a default environment file", func() {
					g.Expect(layer.DefaultBuildEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.build", "TEST_NAME.default")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a delimiter environment file", func() {
					g.Expect(layer.DelimiterBuildEnv("TEST_NAME", "test-delimiter")).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.build", "TEST_NAME.delim")).To(internal.HaveContent("test-delimiter"))
				})

				it("writes an override environment file", func() {
					g.Expect(layer.OverrideBuildEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.build", "TEST_NAME.override")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a prepend environment file", func() {
					g.Expect(layer.PrependBuildEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.build", "TEST_NAME.prepend")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a prepend path environment file", func() {
					g.Expect(layer.PrependPathBuildEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.build", "TEST_NAME")).To(internal.HaveContent("test-string-1"))
				})
			})

			when("launch", func() {

				it("writes an append environment file", func() {
					g.Expect(layer.AppendLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.launch", "TEST_NAME.append")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a default environment file", func() {
					g.Expect(layer.DefaultLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.launch", "TEST_NAME.default")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a delimiter environment file", func() {
					g.Expect(layer.DelimiterLaunchEnv("TEST_NAME", "test-delimiter")).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.launch", "TEST_NAME.delim")).To(internal.HaveContent("test-delimiter"))
				})

				it("writes an override environment file", func() {
					g.Expect(layer.OverrideLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.launch", "TEST_NAME.override")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a prepend environment file", func() {
					g.Expect(layer.PrependLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.launch", "TEST_NAME.prepend")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a prepend path environment file", func() {
					g.Expect(layer.PrependPathLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env.launch", "TEST_NAME")).To(internal.HaveContent("test-string-1"))
				})
			})

			when("shared", func() {

				it("writes an append environment file", func() {
					g.Expect(layer.AppendSharedEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env", "TEST_NAME.append")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a default environment file", func() {
					g.Expect(layer.DefaultSharedEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env", "TEST_NAME.default")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a delimiter environment file", func() {
					g.Expect(layer.DelimiterSharedEnv("TEST_NAME", "test-delimiter")).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env", "TEST_NAME.delim")).To(internal.HaveContent("test-delimiter"))
				})

				it("writes an override environment file", func() {
					g.Expect(layer.OverrideSharedEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env", "TEST_NAME.override")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a prepend environment file", func() {
					g.Expect(layer.PrependSharedEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env", "TEST_NAME.prepend")).To(internal.HaveContent("test-string-1"))
				})

				it("writes a prepend path environment file", func() {
					g.Expect(layer.PrependPathSharedEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(gomega.Succeed())

					g.Expect(filepath.Join(root, "env", "TEST_NAME")).To(internal.HaveContent("test-string-1"))
				})
			})
		})
	}, spec.Report(report.Terminal{}))
}
