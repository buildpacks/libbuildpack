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
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/logger"
	"github.com/buildpack/libbuildpack/platform"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestEnvironmentVariables(t *testing.T) {
	spec.Run(t, "EnvironmentVariables", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("sets all platform environment variables", func() {
			root := internal.ScratchDir(t, "platform")
			defer internal.ProtectEnv(t, "TEST_KEY_1", "TEST_KEY_2")()

			internal.WriteTestFile(t, filepath.Join(root, "env", "TEST_KEY_1"), "test-value-1")
			internal.WriteTestFile(t, filepath.Join(root, "env", "TEST_KEY_2"), "test-value-2")

			platform, err := platform.DefaultPlatform(root, logger.Logger{})
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(platform.EnvironmentVariables.SetAll()).To(gomega.Succeed())
			g.Expect(os.Getenv("TEST_KEY_1")).To(gomega.Equal("test-value-1"))
			g.Expect(os.Getenv("TEST_KEY_2")).To(gomega.Equal("test-value-2"))
		})

	}, spec.Report(report.Terminal{}))
}
