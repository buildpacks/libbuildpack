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

package logger_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack/internal"
	loggerPkg "github.com/buildpack/libbuildpack/logger"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestLogger(t *testing.T) {
	spec.Run(t, "Logger", testLogger, spec.Random(), spec.Report(report.Terminal{}))
}

func testLogger(t *testing.T, when spec.G, it spec.S) {

	it("writes output to debug writer", func() {
		var debug bytes.Buffer

		logger := loggerPkg.NewLogger(&debug, nil)
		logger.Debug("%s %s", "test-string-1", "test-string-2")

		if debug.String() != "test-string-1 test-string-2\n" {
			t.Errorf("debug = %s, wanted test-string-1 test-string-2\\n", debug.String())
		}
	})

	it("does not write to debug if not configured", func() {
		var debug io.Writer

		logger := loggerPkg.NewLogger(debug, nil)
		logger.Debug("%s %s", "test-string-1", "test-string-2")
	})

	it("reports debug enabled when configured", func() {
		var debug bytes.Buffer

		if !loggerPkg.NewLogger(&debug, nil).IsDebugEnabled() {
			t.Errorf("IsDebugEnabled = false, expected true")
		}
	})

	it("reports debug disabled when not configured", func() {
		var debug io.Writer

		if loggerPkg.NewLogger(debug, nil).IsDebugEnabled() {
			t.Errorf("IsDebugEnabled = true, expected false")
		}
	})

	it("writes output to info writer", func() {
		var info bytes.Buffer

		logger := loggerPkg.NewLogger(nil, &info)
		logger.Info("%s %s", "test-string-1", "test-string-2")

		if info.String() != "test-string-1 test-string-2\n" {
			t.Errorf("info = %s, wanted test-string-1 test-string-2\\n", info.String())
		}
	})

	it("does not write to info if not configured", func() {
		var info io.Writer

		logger := loggerPkg.NewLogger(nil, info)
		logger.Info("%s %s", "test-string-1", "test-string-2")
	})

	it("reports info enabled when configured", func() {
		var info bytes.Buffer

		if !loggerPkg.NewLogger(nil, &info).IsInfoEnabled() {
			t.Errorf("IsInfoEnabled = false, expected true")
		}
	})

	it("reports info disabled when not configured", func() {
		var info io.Writer

		if loggerPkg.NewLogger(nil, info).IsInfoEnabled() {
			t.Errorf("IsInfoEnabled = true, expected false")
		}
	})

	it("suppresses debug output", func() {
		c, d := internal.ReplaceConsole(t)
		defer d()

		logger := loggerPkg.DefaultLogger()

		logger.Debug("test-debug-output")
		logger.Info("test-info-output")

		if strings.Contains(c.Err(t), "test-debug-output") {
			t.Errorf("stderr contained test-debug-output, expected not")
		}

		if !strings.Contains(c.Out(t), "test-info-output") {
			t.Errorf("stdout did not contain test-info-output, expected to")
		}
	})

	it("allows debug output if BP_DEBUG is set", func() {
		c, d := internal.ReplaceConsole(t)
		defer d()

		defer internal.ReplaceEnv(t, "BP_DEBUG", "")()

		logger := loggerPkg.DefaultLogger()

		logger.Debug("test-debug-output")
		logger.Info("test-info-output")

		if !strings.Contains(c.Err(t), "test-debug-output") {
			t.Errorf("stderr did not contain test-debug-output, expected to")
		}

		if !strings.Contains(c.Out(t), "test-info-output") {
			t.Errorf("stdout did not contain test-info-output, expected to")
		}
	})
}
