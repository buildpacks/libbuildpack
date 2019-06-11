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

package stack_test

import (
	"os"
	"testing"

	"github.com/buildpack/libbuildpack/internal"
	"github.com/buildpack/libbuildpack/logger"
	"github.com/buildpack/libbuildpack/stack"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestStack(t *testing.T) {
	spec.Run(t, "Stack", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("extracts value from CNB_STACK_ID", func() {
			defer internal.ReplaceEnv(t, "CNB_STACK_ID", "test-stack-name")()

			g.Expect(stack.DefaultStack(logger.Logger{})).To(Equal(stack.Stack("test-stack-name")))
		})

		it("returns error when CNB_STACK_ID not set", func() {
			defer internal.ProtectEnv(t, "CNB_STACK_ID")()
			g.Expect(os.Unsetenv("CNB_STACK_ID")).Should(Succeed())

			_, err := stack.DefaultStack(logger.Logger{})
			g.Expect(err).To(MatchError("CNB_STACK_ID not set"))
		})
	}, spec.Report(report.Terminal{}))
}
