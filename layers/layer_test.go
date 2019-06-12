/*
 * Copyright 2018-2019 the original author or authors.
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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/layers"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLayer(t *testing.T) {
	spec.Run(t, "Layer", func(t *testing.T, when spec.G, it spec.S) {

		g := NewGomegaWithT(t)

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
				var err error
				layer, err = layers.Layers{Root: root}.Layer("test-layer")
				g.Expect(err).NotTo(HaveOccurred())
			})

			it("reads layer content metadata", func() {
				internal.WriteTestFile(t, filepath.Join(root, "test-layer.toml"), `[metadata]
Alpha = "test-value"
Bravo = 1
`)

				var actual metadata
				g.Expect(layer.ReadMetadata(&actual)).To(Succeed())

				g.Expect(actual).To(Equal(metadata{"test-value", 1}))
			})

			it("does not read layer content metadata if it does not exist", func() {
				var actual metadata
				g.Expect(layer.ReadMetadata(&actual)).To(Succeed())

				g.Expect(actual).To(Equal(metadata{}))
			})

			it("remove layer content metadata", func() {
				internal.WriteTestFile(t, filepath.Join(root, "test-layer.toml"), `[metadata]
Alpha = "test-value"
Bravo = 1
`)

				g.Expect(layer.RemoveMetadata()).To(Succeed())
				g.Expect("").NotTo(BeAnExistingFile())
			})

			it("writes layer content metadata", func() {
				g.Expect(layer.WriteMetadata(metadata{"test-value", 1},
					layers.Build, layers.Cache, layers.Launch)).To(Succeed())

				g.Expect(filepath.Join(root, "test-layer.toml")).To(internal.HaveContent(`build = true
cache = true
launch = true

[metadata]
  Alpha = "test-value"
  Bravo = 1
`))
			})

			it("writes a profile file", func() {
				g.Expect(layer.WriteProfile("test-name", "%s-%d", "test-string", 1)).To(Succeed())

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
					g.Expect(layer.AppendBuildEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(Succeed())

					g.Expect(filepath.Join(root, "env.build", "TEST_NAME.append")).To(internal.HaveContent("test-string-1"))
				})

				it("writes an append path environment file", func() {
					g.Expect(layer.AppendPathBuildEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(Succeed())

					g.Expect(filepath.Join(root, "env.build", "TEST_NAME")).To(internal.HaveContent("test-string-1"))
				})

				it("writes an override environment file", func() {
					g.Expect(layer.OverrideBuildEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(Succeed())

					g.Expect(filepath.Join(root, "env.build", "TEST_NAME.override")).To(internal.HaveContent("test-string-1"))
				})
			})

			when("launch", func() {

				it("writes an append environment file", func() {
					g.Expect(layer.AppendLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(Succeed())

					g.Expect(filepath.Join(root, "env.launch", "TEST_NAME.append")).To(internal.HaveContent("test-string-1"))
				})

				it("writes an append path environment file", func() {
					g.Expect(layer.AppendPathLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(Succeed())

					g.Expect(filepath.Join(root, "env.launch", "TEST_NAME")).To(internal.HaveContent("test-string-1"))
				})

				it("writes an override environment file", func() {
					g.Expect(layer.OverrideLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(Succeed())

					g.Expect(filepath.Join(root, "env.launch", "TEST_NAME.override")).To(internal.HaveContent("test-string-1"))
				})

			})

			when("shared", func() {

				it("writes an append environment file", func() {
					g.Expect(layer.AppendSharedEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(Succeed())

					g.Expect(filepath.Join(root, "env", "TEST_NAME.append")).To(internal.HaveContent("test-string-1"))
				})

				it("writes an append path environment file", func() {
					g.Expect(layer.AppendPathSharedEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(Succeed())

					g.Expect(filepath.Join(root, "env", "TEST_NAME")).To(internal.HaveContent("test-string-1"))
				})

				it("writes an override environment file", func() {
					g.Expect(layer.OverrideSharedEnv("TEST_NAME", "%s-%d", "test-string", 1)).To(Succeed())

					g.Expect(filepath.Join(root, "env", "TEST_NAME.override")).To(internal.HaveContent("test-string-1"))
				})

			})
		})

		it("layer root extracted", func() {
			layerName := "test-layer"
			root := internal.ScratchDir(t, "layer")
			l, err := layers.Layers{Root: root}.Layer(layerName)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(l.Root).To(Equal(filepath.Join(root, layerName)))
		})

		when("file existence checks", func() {

			var (
				l layers.Layer
			)

			it.Before(func() {
				root := internal.ScratchDir(t, "layer")
				var err error
				l, err = layers.Layers{Root: root}.Layer("test-layer")
				g.Expect(err).NotTo(HaveOccurred())
			})

			it("exists in root", func() {
				file := "exists.txt"
				err := ioutil.WriteFile(filepath.Join(l.Root, file), []byte("content"), 0600)
				g.Expect(err).NotTo(HaveOccurred())

				exists, err := l.FileExists(file)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(exists).Should(BeTrue())
			})

			it("exists in subdir", func() {
				file := "subdir/subdir2/exists.txt"
				err := os.MkdirAll(filepath.Dir(filepath.Join(l.Root, file)), 0700)
				g.Expect(err).NotTo(HaveOccurred())
				err = ioutil.WriteFile(filepath.Join(l.Root, file), []byte("content"), 0600)
				g.Expect(err).NotTo(HaveOccurred())

				exists, err := l.FileExists(file)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(exists).Should(BeTrue())
			})

			it("does not exist", func() {
				exists, err := l.FileExists("doesnotexist.txt")
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(exists).Should(BeFalse())
			})
		})

	}, spec.Report(report.Terminal{}))
}
