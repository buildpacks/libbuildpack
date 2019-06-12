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

package application_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/application"
	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/logger"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestApplication(t *testing.T) {
	spec.Run(t, "Application", func(t *testing.T, when spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var (
			root      string
			app       application.Application
			wdCleanUp func()
		)

		it.Before(func() {
			root = internal.ScratchDir(t, "application")
			wdCleanUp = internal.ReplaceWorkingDirectory(t, root)
			var err error
			app, err = application.DefaultApplication(logger.Logger{})
			g.Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			wdCleanUp()
		})

		it("extracts root from working directory", func() {
			g.Expect(app.Root).To(Equal(root))
		})

		when("file existence checks", func() {

			it("exists in root", func() {
				file := "exists.txt"
				err := ioutil.WriteFile(filepath.Join(app.Root, file), []byte("content"), 0600)
				g.Expect(err).NotTo(HaveOccurred())

				exists, err := app.FileExists(file)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(exists).Should(BeTrue())
			})

			it("exists in subdir", func() {
				file := "subdir/subdir2/exists.txt"
				err := os.MkdirAll(filepath.Dir(filepath.Join(app.Root, file)), 0700)
				g.Expect(err).NotTo(HaveOccurred())
				err = ioutil.WriteFile(filepath.Join(app.Root, file), []byte("content"), 0600)
				g.Expect(err).NotTo(HaveOccurred())

				exists, err := app.FileExists(file)
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(exists).Should(BeTrue())
			})

			it("does not exist", func() {
				exists, err := app.FileExists("doesnotexist.txt")
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(exists).Should(BeFalse())
			})
		})
	}, spec.Report(report.Terminal{}))
}
