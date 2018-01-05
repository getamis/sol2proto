// Copyright 2017 AMIS Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SQL Options", func() {
	Describe("Test options", func() {
		Context("Retry options", func() {
			It("should match", func() {
				delay := 1 * time.Second
				timeout := 10 * time.Minute
				expectedOptions := &Options{
					RetryDelay:   delay,
					RetryTimeout: timeout,
				}

				options := &Options{}
				fn := Retry(delay, timeout)
				fn(options)

				Expect(options.RetryDelay).Should(Equal(expectedOptions.RetryDelay))
				Expect(options.RetryTimeout).Should(Equal(expectedOptions.RetryTimeout))
			})
		})

		Context("Table", func() {
			It("should match", func() {
				tableName := "This is a table name"
				expectedOptions := &Options{
					TableName: tableName,
				}

				options := &Options{}
				fn := Table(tableName)
				fn(options)

				Expect(options.TableName).Should(Equal(expectedOptions.TableName))
			})
		})

		Context("Logging", func() {
			It("should match", func() {
				logging := true
				expectedOptions := &Options{
					Logging: logging,
				}

				options := &Options{}
				fn := Logging(logging)
				fn(options)

				Expect(options.Logging).Should(Equal(expectedOptions.Logging))
			})
		})

		Context("Driver", func() {
			It("should match", func() {
				driver := "This is a driver name"
				expectedOptions := &Options{
					Driver: driver,
				}

				options := &Options{}
				fn := Driver(driver)
				fn(options)

				Expect(options.Driver).Should(Equal(expectedOptions.Driver))
			})
		})

		Context("Driver options", func() {
			It("should match", func() {
				arbitraryOptions := []interface{}{
					"stringOption",
					1,
					"foo",
				}
				expectedOptions := &Options{
					DriverOptions: arbitraryOptions,
				}

				options := &Options{}
				fn := DriverOption(arbitraryOptions...)
				fn(options)

				Expect(options.DriverOptions).Should(Equal(expectedOptions.DriverOptions))
			})
		})
	})
})

func TestSQLOptionSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SQL Option Suite")
}
