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

package libbuildpack_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/buildpack/libbuildpack"
	"github.com/buildpack/libbuildpack/internal"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildPlan(t *testing.T) {
	spec.Run(t, "BuildPlan", testBuildPlan, spec.Report(report.Terminal{}))
}

func testBuildPlan(t *testing.T, when spec.G, it spec.S) {

	logger := libbuildpack.NewLogger(nil, nil)

	expected := libbuildpack.BuildPlan{
		"alpha": libbuildpack.BuildPlanDependency{
			"version": "alpha-version",
			"name":    "alpha-name",
		},
		"bravo": libbuildpack.BuildPlanDependency{
			"name": "bravo-name",
		},
	}

	it("unmarshals default from os.Stdin", func() {
		console, d := internal.ReplaceConsole(t)
		defer d()

		console.In(t, `[alpha]
  version = "alpha-version"
  name = "alpha-name"

[bravo]
  name = "bravo-name"
  `)

		buildPlan, err := libbuildpack.DefaultBuildPlan(logger)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(buildPlan, expected) {
			t.Errorf("BuildPlan = %s, wanted %s", buildPlan, expected)
		}
	})

	it("unmarshals from reader", func() {
		in := strings.NewReader(`[alpha]
  version = "alpha-version"
  name = "alpha-name"

[bravo]
  name = "bravo-name"
  `)

		buildPlan, err := libbuildpack.NewBuildPlan(in, logger)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(buildPlan, expected) {
			t.Errorf("BuildPlan = %s, wanted %s", buildPlan, expected)
		}

	})

	it("returns a dependency by name", func() {
		expected := libbuildpack.BuildPlanDependency{}
		p := libbuildpack.BuildPlan{"test-dependency": expected}

		actual := p["test-dependency"]
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("BuildPlan[\"test-dependency\"] = %s, expected %s", actual, expected)
		}
	})

	it("returns nil if a named dependency does not exist", func() {
		p := libbuildpack.BuildPlan{}

		actual := p["test-dependency"]
		if actual != nil {
			t.Errorf("BuildPlan[\"test-dependency\"] = %s, expected nil", actual)
		}
	})

	it("returns the dependency version", func() {
		d := libbuildpack.BuildPlanDependency{"version": "1.*"}

		actual, err := d.Version()
		if err != nil {
			t.Fatal(err)
		}

		expected, err := semver.NewVersion("1.0")
		if err != nil {
			t.Fatal(err)
		}

		if !actual.Check(expected) {
			t.Errorf("BuildPlanDependency.Version = %v, expected %v", actual, expected)
		}
	})

	it("returns error if dependency version does not exist", func() {
		d := libbuildpack.BuildPlanDependency{}

		_, err := d.Version()
		if err == nil {
			t.Errorf("BuildPlanDependency.Version = nil, expected not nil")
		}
	})

}
