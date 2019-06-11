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

package platform_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/logger"
	"github.com/buildpack/libbuildpack/platform"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestPlatform(t *testing.T) {
	spec.Run(t, "Platform", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("enumerates platform environment variables", func() {
			root := internal.ScratchDir(t, "platform")

			internal.WriteTestFile(t, filepath.Join(root, "env", "TEST_KEY"), "test-value")

			platform, err := platform.DefaultPlatform(root, logger.Logger{})
			g.Expect(err).To(Succeed())

			g.Expect(platform.EnvironmentVariables).To(HaveKey("TEST_KEY"))
		})
	}, spec.Report(report.Terminal{}))
}
