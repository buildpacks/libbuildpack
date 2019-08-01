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

package buildpackplan_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildpackplan"
	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/logger"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildpackPlan(t *testing.T) {
	spec.Run(t, "Plan", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("unmarshals from plan", func() {
			root := internal.ScratchDir(t, "buildpack-plan")

			internal.WriteTestFile(t, filepath.Join(root, "plan.toml"), `[[entries]]
  name = "test-entry-1a"
  version = "test-version-1a"
  [entries.metadata]
    test-key-1a = "test-value-1a"

[[entries]]
  name = "test-entry-1b"
  version = "test-version-1b"
  [entries.metadata]
    test-key-1b = "test-value-1b"
`)

			p, err := buildpackplan.DefaultPlans(filepath.Join(root, "plan.toml"), logger.Logger{})
			g.Expect(err).To(gomega.Succeed())

			g.Expect(p).To(gomega.Equal(buildpackplan.Plans{
				Entries: []buildpackplan.Plan{
					{"test-entry-1a", "test-version-1a", buildpackplan.Metadata{"test-key-1a": "test-value-1a"}},
					{"test-entry-1b", "test-version-1b", buildpackplan.Metadata{"test-key-1b": "test-value-1b"}},
				},
			}))
		})

	}, spec.Report(report.Terminal{}))
}
