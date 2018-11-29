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

package layers_test

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack/internal"
	layersPkg "github.com/buildpack/libbuildpack/layers"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLayer(t *testing.T) {
	spec.Run(t, "Layer", testLayer, spec.Report(report.Terminal{}))
}

func testLayer(t *testing.T, when spec.G, it spec.S) {

	when("layer content metadata", func() {

		type metadata struct {
			Alpha string
			Bravo int
		}

		it("reads layer content metadata", func() {
			root := internal.ScratchDir(t, "layer")
			layers := layersPkg.Layers{Root: root}
			layer := layers.Layer("test-layer")

			if err := internal.WriteToFile(strings.NewReader(`[metadata]
Alpha = "test-value"
Bravo = 1
`), filepath.Join(root, "test-layer.toml"), 0644); err != nil {
				t.Fatal(err)
			}

			var actual metadata
			if err := layer.ReadMetadata(&actual); err != nil {
				t.Fatal(err)
			}

			expected := metadata{"test-value", 1}

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("metadata = %v, wanted %v", actual, expected)
			}
		})

		it("does not read layer content metadata if it does not exist", func() {
			root := internal.ScratchDir(t, "layer")
			layers := layersPkg.Layers{Root: root}
			layer := layers.Layer("test-layer")

			var actual metadata
			if err := layer.ReadMetadata(&actual); err != nil {
				t.Fatal(err)
			}

			expected := metadata{}

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("metadata = %v, wanted %v", actual, expected)
			}
		})

		it("remove layer content metadata", func() {
			root := internal.ScratchDir(t, "layer")
			layers := layersPkg.Layers{Root: root}
			layer := layers.Layer("test-layer")

			metadata := filepath.Join(root, "test-layer.toml")
			if err := internal.WriteToFile(strings.NewReader(`Alpha = "test-value"
Bravo = 1
`), metadata, 0644); err != nil {
				t.Fatal(err)
			}

			if err := layer.RemoveMetadata(); err != nil {
				t.Fatal(err)
			}

			exists, err := internal.FileExists(metadata)
			if err != nil {
				t.Fatal(err)
			}

			if exists {
				t.Errorf("%s exists, expected not to", metadata)
			}
		})

		it("writes layer content metadata", func() {
			root := internal.ScratchDir(t, "layer")
			layers := layersPkg.Layers{Root: root}
			layer := layers.Layer("test-layer")

			if err := layer.WriteMetadata(metadata{"test-value", 1}, layersPkg.Build, layersPkg.Cache, layersPkg.Launch); err != nil {
				t.Fatal(err)
			}

			internal.BeFileLike(t, filepath.Join(root, "test-layer.toml"), 0644, `build = true
cache = true
launch = true

[metadata]
  Alpha = "test-value"
  Bravo = 1
`)
		})

		it("writes a profile file", func() {
			root := internal.ScratchDir(t, "layer")
			layer := layersPkg.Layer{Root: root}

			if err := layer.WriteProfile("test-name", "%s-%d", "test-string", 1); err != nil {
				t.Fatal(err)
			}

			internal.BeFileLike(t, filepath.Join(root, "profile.d", "test-name"), 0644, "test-string-1")
		})
	})

	when("environment files", func() {

		when("build", func() {

			it("writes an append environment file", func() {
				root := internal.ScratchDir(t, "cache")
				layer := layersPkg.Layer{Root: root}

				if err := layer.AppendBuildEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
					t.Fatal(err)
				}

				internal.BeFileLike(t, filepath.Join(root, "env.build", "TEST_NAME.append"), 0644, "test-string-1")
			})

			it("writes an append path environment file", func() {
				root := internal.ScratchDir(t, "cache")
				layer := layersPkg.Layer{Root: root}

				if err := layer.AppendPathBuildEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
					t.Fatal(err)
				}

				internal.BeFileLike(t, filepath.Join(root, "env.build", "TEST_NAME"), 0644, "test-string-1")
			})

			it("writes an override environment file", func() {
				root := internal.ScratchDir(t, "cache")
				layer := layersPkg.Layer{Root: root}

				if err := layer.OverrideBuildEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
					t.Fatal(err)
				}

				internal.BeFileLike(t, filepath.Join(root, "env.build", "TEST_NAME.override"), 0644, "test-string-1")
			})
		})

		when("launch", func() {

			it("writes an append environment file", func() {
				root := internal.ScratchDir(t, "cache")
				layer := layersPkg.Layer{Root: root}

				if err := layer.AppendLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
					t.Fatal(err)
				}

				internal.BeFileLike(t, filepath.Join(root, "env.launch", "TEST_NAME.append"), 0644, "test-string-1")
			})

			it("writes an append path environment file", func() {
				root := internal.ScratchDir(t, "cache")
				layer := layersPkg.Layer{Root: root}

				if err := layer.AppendPathLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
					t.Fatal(err)
				}

				internal.BeFileLike(t, filepath.Join(root, "env.launch", "TEST_NAME"), 0644, "test-string-1")
			})

			it("writes an override environment file", func() {
				root := internal.ScratchDir(t, "cache")
				layer := layersPkg.Layer{Root: root}

				if err := layer.OverrideLaunchEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
					t.Fatal(err)
				}

				internal.BeFileLike(t, filepath.Join(root, "env.launch", "TEST_NAME.override"), 0644, "test-string-1")
			})

		})

		when("shared", func() {

			it("writes an append environment file", func() {
				root := internal.ScratchDir(t, "cache")
				layer := layersPkg.Layer{Root: root}

				if err := layer.AppendSharedEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
					t.Fatal(err)
				}

				internal.BeFileLike(t, filepath.Join(root, "env", "TEST_NAME.append"), 0644, "test-string-1")
			})

			it("writes an append path environment file", func() {
				root := internal.ScratchDir(t, "cache")
				layer := layersPkg.Layer{Root: root}

				if err := layer.AppendPathSharedEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
					t.Fatal(err)
				}

				internal.BeFileLike(t, filepath.Join(root, "env", "TEST_NAME"), 0644, "test-string-1")
			})

			it("writes an override environment file", func() {
				root := internal.ScratchDir(t, "cache")
				layer := layersPkg.Layer{Root: root}

				if err := layer.OverrideSharedEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
					t.Fatal(err)
				}

				internal.BeFileLike(t, filepath.Join(root, "env", "TEST_NAME.override"), 0644, "test-string-1")
			})

		})
	})
}
