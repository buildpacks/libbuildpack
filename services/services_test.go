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

package services_test

import (
	"testing"

	"github.com/buildpacks/libbuildpack/v2/internal"
	"github.com/buildpacks/libbuildpack/v2/logger"
	"github.com/buildpacks/libbuildpack/v2/platform"
	"github.com/buildpacks/libbuildpack/v2/services"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestServices(t *testing.T) {
	spec.Run(t, "Services", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		it("is empty with no CNB_SERVICES", func() {
			s, err := services.DefaultServices(platform.Platform{}, logger.Logger{})
			g.Expect(err).To(gomega.Succeed())

			g.Expect(s).To(gomega.BeEmpty())
		})

		it("parses CNB_SERVICES", func() {
			defer internal.ReplaceEnv(t, "CNB_SERVICES", `{
  "elephantsql": [
    {
      "name": "elephantsql-binding-c6c60",
      "binding_name": "elephantsql-binding-c6c60",
      "instance_name": "elephantsql-c6c60",
      "label": "elephantsql",
      "tags": [
        "postgres",
        "postgresql",
        "relational"
      ],
      "plan": "turtle",
      "credentials": {
        "uri": "postgres://exampleuser:examplepass@babar.elephantsql.com:5432/exampleuser"
      }
    }
  ],
  "sendgrid": [
    {
      "name": "mysendgrid",
      "binding_name": null,
      "instance_name": "mysendgrid",
      "label": "sendgrid",
      "tags": [
        "smtp"
      ],
      "plan": "free",
      "credentials": {
        "hostname": "smtp.sendgrid.net",
        "username": "QvsXMbJ3rK",
        "password": "HCHMOYluTv"
      }
    }
  ]
}`)()

			s, err := services.DefaultServices(platform.Platform{}, logger.Logger{})
			g.Expect(err).To(gomega.Succeed())

			g.Expect(s).To(gomega.ConsistOf(services.Services{
				{
					BindingName:  "elephantsql-binding-c6c60",
					Credentials:  services.Credentials{"uri": "postgres://exampleuser:examplepass@babar.elephantsql.com:5432/exampleuser"},
					InstanceName: "elephantsql-c6c60",
					Label:        "elephantsql",
					Plan:         "turtle",
					Tags:         []string{"postgres", "postgresql", "relational"},
				},
				{
					Credentials:  services.Credentials{"hostname": "smtp.sendgrid.net", "password": "HCHMOYluTv", "username": "QvsXMbJ3rK"},
					InstanceName: "mysendgrid",
					Label:        "sendgrid",
					Plan:         "free",
					Tags:         []string{"smtp"},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
