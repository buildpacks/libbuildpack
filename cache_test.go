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
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack"
	"github.com/buildpack/libbuildpack/internal"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestCache(t *testing.T) {
	spec.Run(t, "Cache", testCache, spec.Report(report.Terminal{}))
}

func testCache(t *testing.T, when spec.G, it spec.S) {

	logger := libbuildpack.NewLogger(nil, nil)

	it("extracts roots from os.Args[1]", func() {
		defer internal.ReplaceArgs(t, "", "cache-root")()

		cache, err := libbuildpack.DefaultCache(logger)
		if err != nil {
			t.Fatal(err)
		}

		if cache.Root != "cache-root" {
			t.Errorf("Cache.Root = %s, expected = cache-root", cache.Root)
		}
	})

	it("creates a cache layer with root based on its name", func() {
		cache := libbuildpack.Cache{Root: "test-root"}
		layer := cache.Layer("test-layer")

		if layer.Root != "test-root/test-layer" {
			t.Errorf("CacheLayer.Root = %s, expected test-root/test-layer", layer.Root)
		}
	})

	it("writes an append environment file", func() {
		root := internal.ScratchDir(t, "cache")
		layer := libbuildpack.CacheLayer{Root: root}

		if err := layer.AppendEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
			t.Fatal(err)
		}

		f := filepath.Join(root, "env", "TEST_NAME.append")
		exists, err := internal.FileExists(f)
		if err != nil {
			t.Fatal(err)
		}

		if !exists {
			t.Errorf("Expected %s to exist, did not", f)
		}

		content, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != "test-string-1" {
			t.Errorf("%s = %s, expected test-string-1", f, string(content))
		}
	})

	it("writes an append path environment file", func() {
		root := internal.ScratchDir(t, "cache")
		layer := libbuildpack.CacheLayer{Root: root}

		if err := layer.AppendPathEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
			t.Fatal(err)
		}

		f := filepath.Join(root, "env", "TEST_NAME")
		exists, err := internal.FileExists(f)
		if err != nil {
			t.Fatal(err)
		}

		if !exists {
			t.Errorf("Expected %s to exist, did not", f)
		}

		content, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != "test-string-1" {
			t.Errorf("%s = %s, expected test-string-1", f, string(content))
		}
	})

	it("writes an override environment file", func() {
		root := internal.ScratchDir(t, "cache")
		layer := libbuildpack.CacheLayer{Root: root}

		if err := layer.OverrideEnv("TEST_NAME", "%s-%d", "test-string", 1); err != nil {
			t.Fatal(err)
		}

		f := filepath.Join(root, "env", "TEST_NAME.override")
		exists, err := internal.FileExists(f)
		if err != nil {
			t.Fatal(err)
		}

		if !exists {
			t.Errorf("Expected %s to exist, did not", f)
		}

		content, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}

		if string(content) != "test-string-1" {
			t.Errorf("%s = %s, expected test-string-1", f, string(content))
		}
	})
}
