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

package logger_test

import (
	"bytes"
	"testing"

	"github.com/buildpacks/libbuildpack/v2/internal"
	"github.com/buildpacks/libbuildpack/v2/logger"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLogger(t *testing.T) {
	spec.Run(t, "Logger", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("writes output to debug writer", func() {
			var debug bytes.Buffer

			logger := logger.NewLogger(&debug, nil)
			logger.Debug("%s %s", "test-string-1", "test-string-2")

			g.Expect(debug.String()).To(gomega.Equal("test-string-1 test-string-2\n"))
		})

		it("does not write to debug if not configured", func() {
			logger := logger.NewLogger(nil, nil)
			logger.Debug("%s %s", "test-string-1", "test-string-2")
		})

		it("reports debug enabled when configured", func() {
			var debug bytes.Buffer

			g.Expect(logger.NewLogger(&debug, nil).IsDebugEnabled()).To(gomega.BeTrue())
		})

		it("reports debug disabled when not configured", func() {
			g.Expect(logger.NewLogger(nil, nil).IsDebugEnabled()).To(gomega.BeFalse())
		})

		it("writes output to info writer", func() {
			var info bytes.Buffer

			logger := logger.NewLogger(nil, &info)
			logger.Info("%s %s", "test-string-1", "test-string-2")

			g.Expect(info.String()).To(gomega.Equal("test-string-1 test-string-2\n"))
		})

		it("does not write to info if not configured", func() {
			logger := logger.NewLogger(nil, nil)
			logger.Info("%s %s", "test-string-1", "test-string-2")
		})

		it("reports info enabled when configured", func() {
			var info bytes.Buffer

			g.Expect(logger.NewLogger(nil, &info).IsInfoEnabled()).To(gomega.BeTrue())
		})

		it("reports info disabled when not configured", func() {
			g.Expect(logger.NewLogger(nil, nil).IsInfoEnabled()).To(gomega.BeFalse())
		})

		it("suppresses debug output", func() {
			root := internal.ScratchDir(t, "logger")
			c, d := internal.ReplaceConsole(t)
			defer d()

			logger, err := logger.DefaultLogger(root)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			logger.Debug("test-debug-output")
			logger.Info("test-info-output")

			g.Expect(c.Err(t)).NotTo(gomega.ContainSubstring("test-debug-output"))
			g.Expect(c.Out(t)).To(gomega.ContainSubstring("test-info-output"))
		})

		it("allows debug output if BP_DEBUG is set", func() {
			root := internal.ScratchDir(t, "logger")
			c, d := internal.ReplaceConsole(t)
			defer d()
			defer internal.ReplaceEnv(t, "BP_DEBUG", "")()

			logger, err := logger.DefaultLogger(root)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			logger.Debug("test-debug-output")
			logger.Info("test-info-output")

			g.Expect(c.Err(t)).To(gomega.ContainSubstring("test-debug-output"))
			g.Expect(c.Out(t)).To(gomega.ContainSubstring("test-info-output"))
		})

		it("allows debug output if platform/env/BP_DEBUG is set", func() {
			root := internal.ScratchDir(t, "logger")
			internal.TouchTestFile(t, root, "env", "BP_DEBUG")
			c, d := internal.ReplaceConsole(t)
			defer d()

			logger, err := logger.DefaultLogger(root)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			logger.Debug("test-debug-output")
			logger.Info("test-info-output")

			g.Expect(c.Err(t)).To(gomega.ContainSubstring("test-debug-output"))
			g.Expect(c.Out(t)).To(gomega.ContainSubstring("test-info-output"))
		})
	}, spec.Report(report.Terminal{}))
}
