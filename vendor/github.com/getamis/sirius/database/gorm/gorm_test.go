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

package gorm

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/getamis/sirius/database"
	"github.com/getamis/sirius/database/mysql"
	"github.com/getamis/sirius/test"
)

var mySQLContainer *test.MySQLContainer

var _ = Describe("Test GORM", func() {
	Describe("create a new database connection", func() {
		Context("using MySQL with invalid configuration", func() {
			It("should be failed", func() {
				db, err := New(
					"mysql",
					database.DriverOption(
						mysql.Connector("http", "127.0.0.1", "3306"),
						mysql.UserInfo("username", "password")),
				)
				Expect(err).ShouldNot(BeNil())
				Expect(db).Should(BeNil())
			})
		})

		Context("using MySQL container", func() {
			It("should be ok", func() {
				opt1, opt2, opt3 := mysql.DSNToOptions(mySQLContainer.URL)
				db, err := New(
					"mysql",
					database.Retry(1*time.Second, 3*time.Second),
					database.DriverOption(opt1, opt2, opt3),
				)
				Expect(err).Should(BeNil())
				Expect(db).ShouldNot(BeNil())
			})
		})
	})
})

var _ = BeforeSuite(func() {
	mySQLContainer, _ = test.NewMySQLContainer()
	mySQLContainer.Start()
})

var _ = AfterSuite(func() {
	mySQLContainer.Stop()
})

func TestGORMFactorySuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GORM Factory Suite")
}
