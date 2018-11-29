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
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/logger"
	platformPkg "github.com/buildpack/libbuildpack/platform"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestPlatform(t *testing.T) {
	spec.Run(t, "Platform", testPlatform, spec.Random(), spec.Report(report.Terminal{}))
}

func testPlatform(t *testing.T, when spec.G, it spec.S) {

	it("extracts root from os.Args[1]", func() {
		root := internal.ScratchDir(t, "platform")
		defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"))()

		platform, err := platformPkg.DefaultPlatform(logger.Logger{})
		if err != nil {
			t.Fatal(err)
		}

		if platform.Root != filepath.Join(root, "platform") {
			t.Errorf("Platform.Root = %s, wanted %s", platform.Root, root)
		}
	})

	it("enumerates platform environment variables", func() {
		root := internal.ScratchDir(t, "platform")
		defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"))()
		if err := internal.WriteToFile(strings.NewReader("test-value"), filepath.Join(root, "platform", "env", "TEST_KEY"), 0644); err != nil {
			t.Fatal(err)
		}

		platform, err := platformPkg.DefaultPlatform(logger.Logger{})
		if err != nil {
			t.Fatal(err)
		}

		if platform.Envs[0].Name != "TEST_KEY" {
			t.Errorf("Platform.Envs[0].Name = %s, expected TEST_KEY", platform.Envs[0])
		}
	})
}
