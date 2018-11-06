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

package platform_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack/internal"
	platformPkg "github.com/buildpack/libbuildpack/platform"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestEnvironmentVariable(t *testing.T) {
	spec.Run(t, "EnvironmentVariable", testEnvironmentVariable, spec.Report(report.Terminal{}))
}

func testEnvironmentVariable(t *testing.T, when spec.G, it spec.S) {

	it("sets a platform environment variable", func() {
		root := internal.ScratchDir(t, "platform")
		defer internal.ProtectEnv(t, "TEST_KEY")()

		if err := internal.WriteToFile(strings.NewReader("test-value"), filepath.Join(root, "platform", "env", "TEST_KEY"), 0644); err != nil {
			t.Fatal(err)
		}

		e := platformPkg.EnvironmentVariable{File: filepath.Join(root, "platform", "env", "TEST_KEY"), Name: "TEST_KEY"}

		if err := e.Set(); err != nil {
			t.Fatal(err)
		}

		if os.Getenv("TEST_KEY") != "test-value" {
			t.Errorf("os.GetEnv(\"TEST_KEY\") = %s, expected test-value", os.Getenv("TEST_KEY"))
		}
	})
}
