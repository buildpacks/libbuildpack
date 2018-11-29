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

func TestEnvironmentVariables(t *testing.T) {
	spec.Run(t, "EnvironmentVariables", testEnvironmentVariables, spec.Report(report.Terminal{}))
}

func testEnvironmentVariables(t *testing.T, when spec.G, it spec.S) {

	it("reports environment variable containment", func() {
		envs := platformPkg.EnvironmentVariables{
			platformPkg.EnvironmentVariable{Name: "TEST_KEY"},
		}

		contains := envs.Contains("TEST_KEY")
		if !contains {
			t.Errorf("Platform.Envs.Contains(\"TEST_KEY\") = %t, expected true", contains)
		}

		contains = envs.Contains("TEST_KEY_2")
		if contains {
			t.Errorf("Platform.Envs.Contains(\"TEST_KEY_2\") = %t, expected false", contains)
		}
	})

	it("sets all platform environment variables", func() {
		root := internal.ScratchDir(t, "platform")
		defer internal.ProtectEnv(t, "TEST_KEY_1", "TEST_KEY_2")()
		if err := internal.WriteToFile(strings.NewReader("test-value-1"), filepath.Join(root, "env", "TEST_KEY_1"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := internal.WriteToFile(strings.NewReader("test-value-2"), filepath.Join(root, "env", "TEST_KEY_2"), 0644); err != nil {
			t.Fatal(err)
		}

		envs := platformPkg.EnvironmentVariables{
			platformPkg.EnvironmentVariable{File: filepath.Join(root, "env", "TEST_KEY_1"), Name: "TEST_KEY_1"},
			platformPkg.EnvironmentVariable{File: filepath.Join(root, "env", "TEST_KEY_2"), Name: "TEST_KEY_2"},
		}

		if err := envs.SetAll(); err != nil {
			t.Fatal(err)
		}

		if os.Getenv("TEST_KEY_1") != "test-value-1" {
			t.Errorf("os.GetEnv(\"TEST_KEY_1\") = %s, expected test-value-1", os.Getenv("TEST_KEY_1"))
		}

		if os.Getenv("TEST_KEY_2") != "test-value-2" {
			t.Errorf("os.GetEnv(\"TEST_KEY_2\") = %s, expected test-value-2", os.Getenv("TEST_KEY_2"))
		}
	})

}
