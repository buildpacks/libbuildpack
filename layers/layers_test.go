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
	"testing"

	"github.com/buildpack/libbuildpack/internal"
	layersPkg "github.com/buildpack/libbuildpack/layers"
	"github.com/buildpack/libbuildpack/logger"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLayers(t *testing.T) {
	spec.Run(t, "Layers", testLayers, spec.Random(), spec.Report(report.Terminal{}))
}

func testLayers(t *testing.T, when spec.G, it spec.S) {

	it("extracts root from os.Args[3]", func() {
		root := internal.ScratchDir(t, "layers")
		defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plans"), filepath.Join(root, "layers"))()

		layers, err := layersPkg.DefaultLayers(logger.Logger{})
		if err != nil {
			t.Fatal(err)
		}

		if layers.Root != filepath.Join(root, "layers") {
			t.Errorf("Laters.Root = %s, expected = launch-root", layers.Root)
		}
	})

	it("creates a layer with root based on its name", func() {
		root := internal.ScratchDir(t, "layers")

		layers := layersPkg.Layers{Root: root}
		layer := layers.Layer("test-layer")

		expected := filepath.Join(root, "test-layer")

		if layer.Root != expected {
			t.Errorf("Layer.Root = %s, expected %s", layer.Root, root)
		}
	})

	it("writes launch metadata", func() {
		root := internal.ScratchDir(t, "layers")
		layers := layersPkg.Layers{Root: root}

		lm := layersPkg.Metadata{
			Processes: layersPkg.Processes{
				layersPkg.Process{Type: "web", Command: "command-1"},
				layersPkg.Process{Type: "task", Command: "command-2"},
			},
		}

		if err := layers.WriteMetadata(lm); err != nil {
			t.Fatal(err)
		}

		internal.BeFileLike(t, filepath.Join(root, "launch.toml"), 0644, `[[processes]]
  type = "web"
  command = "command-1"

[[processes]]
  type = "task"
  command = "command-2"
`)
	})
}
