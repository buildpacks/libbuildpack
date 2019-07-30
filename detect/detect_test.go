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

package detect_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/buildpack/libbuildpack/detect"
	"github.com/buildpack/libbuildpack/internal"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var root string

		it.Before(func() {
			root = internal.ScratchDir(t, "detect")
		})

		it("contains default values", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer internal.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			console, e := internal.ReplaceConsole(t)
			defer e()

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

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			console.In(t, `[alpha]
  version = "alpha-version"
  name = "alpha-name"

[bravo]
  name = "bravo-name"
`)

			g.Expect(d.BuildPlan.Init()).To(gomega.Succeed())

			g.Expect(d.Application).NotTo(gomega.BeZero())
			g.Expect(d.Buildpack).NotTo(gomega.BeZero())
			g.Expect(d.BuildPlan).NotTo(gomega.BeZero())
			g.Expect(d.BuildPlanWriter).NotTo(gomega.BeZero())
			g.Expect(d.Logger).NotTo(gomega.BeZero())
			g.Expect(d.Platform).NotTo(gomega.BeZero())
			g.Expect(d.Services).NotTo(gomega.BeZero())
			g.Expect(d.Stack).NotTo(gomega.BeZero())
		})

		it("returns code when erroring", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer internal.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			internal.TouchTestFile(t, root, "buildpack.toml")

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(d.Error(42)).To(gomega.Equal(42))
		})

		it("returns 100 when failing", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer internal.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			internal.TouchTestFile(t, root, "buildpack.toml")

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(d.Fail()).To(gomega.Equal(detect.FailStatusCode))
		})

		it("returns 0 and BuildPlan when passing", func() {
			defer internal.ReplaceWorkingDirectory(t, root)()
			defer internal.ReplaceEnv(t, "CNB_STACK_ID", "test-stack")()
			defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan.toml"))()

			internal.TouchTestFile(t, root, "buildpack.toml")

			d, err := detect.DefaultDetect()
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(d.Pass(buildplan.BuildPlan{
				"alpha": buildplan.Dependency{Version: "test-version"},
			})).To(gomega.Equal(detect.PassStatusCode))

			g.Expect(filepath.Join(root, "plan.toml")).To(internal.HaveContent(`[alpha]
  version = "test-version"
`))
		})
	}, spec.Report(report.Terminal{}))
}
