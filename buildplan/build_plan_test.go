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

package buildplan_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/buildpack/libbuildpack/internal"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildPlan(t *testing.T) {
	spec.Run(t, "BuildPlan", testBuildPlan, spec.Random(), spec.Report(report.Terminal{}))
}

func testBuildPlan(t *testing.T, when spec.G, it spec.S) {

	it("unmarshals from os.Stdin", func() {
		console, d := internal.ReplaceConsole(t)
		defer d()

		console.In(t, `[alpha]
  version = "alpha-version"

  [alpha.metadata]
    test-key = "test-value"

[bravo]
`)

		buildPlan := buildplan.BuildPlan{}
		if err := buildPlan.Init(); err != nil {
			t.Fatal(err)
		}

		expected := buildplan.BuildPlan{
			"alpha": buildplan.Dependency{
				Version:  "alpha-version",
				Metadata: buildplan.Metadata{"test-key": "test-value"},
			},
			"bravo": buildplan.Dependency{
			},
		}

		if !reflect.DeepEqual(buildPlan, expected) {
			t.Errorf("BuildPlan = %s, wanted %s", buildPlan, expected)
		}
	})

	it("marshals to os.Args[2]", func() {
		root := internal.ScratchDir(t, "buildPlan")
		defer internal.ReplaceArgs(t, filepath.Join(root, "bin", "test"), filepath.Join(root, "platform"), filepath.Join(root, "plan"))()

		buildPlan := buildplan.BuildPlan{
			"alpha": buildplan.Dependency{
				Version:  "alpha-version",
				Metadata: buildplan.Metadata{"test-key": "test-value"},
			},
			"bravo": buildplan.Dependency{
			},
		}

		if err := buildPlan.Write(); err != nil {
			t.Fatal(err)
		}

		internal.BeFileLike(t, filepath.Join(root, "plan", "alpha"), 0644, `version = "alpha-version"

[metadata]
  test-key = "test-value"
`)

		internal.BeFileLike(t, filepath.Join(root, "plan", "bravo"), 0644, `version = ""
`)
	})

	it("returns a dependency by name", func() {
		expected := buildplan.Dependency{}
		p := buildplan.BuildPlan{"test-dependency": expected}

		actual := p["test-dependency"]
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("BuildPlan[\"test-dependency\"] = %s, expected %s", actual, expected)
		}
	})

	it("returns nil if a named dependency does not exist", func() {
		p := buildplan.BuildPlan{}

		actual, ok := p["test-dependency"]
		if ok {
			t.Errorf("BuildPlan[\"test-dependency\"] = %s, expected nil", actual)
		}
	})
}
