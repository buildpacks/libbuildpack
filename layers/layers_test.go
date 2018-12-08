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
	"github.com/buildpack/libbuildpack/layers"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLayers(t *testing.T) {
	spec.Run(t, "Layers", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var root string

		it.Before(func() {
			root = internal.ScratchDir(t, "layers")
		})

		it("creates a layer with root based on its name", func() {
			layer := layers.Layers{Root: root}.Layer("test-layer")

			g.Expect(layer.Root).To(Equal(filepath.Join(root, "test-layer")))
		})

		it("writes launch metadata", func() {
			g.Expect(layers.Layers{Root: root}.WriteMetadata(layers.Metadata{
				Processes: layers.Processes{
					layers.Process{Type: "web", Command: "command-1"},
					layers.Process{Type: "task", Command: "command-2"},
				},
			})).To(Succeed())

			g.Expect(filepath.Join(root, "launch.toml")).To(internal.HaveContent(`[[processes]]
  type = "web"
  command = "command-1"

[[processes]]
  type = "task"
  command = "command-2"
`))
		})
	}, spec.Report(report.Terminal{}))
}
