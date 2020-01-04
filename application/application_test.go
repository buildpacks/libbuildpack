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

package application_test

import (
	"testing"

	"github.com/buildpack/libbuildpack/application"
	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/logger"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestApplication(t *testing.T) {
	spec.Run(t, "Application", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("extracts root from working directory", func() {
			root := internal.ScratchDir(t, "application")
			defer internal.ReplaceWorkingDirectory(t, root)()

			application, err := application.DefaultApplication(logger.Logger{})
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(application.Root).To(gomega.Equal(root))
		})
	}, spec.Report(report.Terminal{}))
}
