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

package test

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/getamis/sirius/database/mysql"
)

type MySQLContainer struct {
	dockerContainer *Container
	URL             string
}

func (container *MySQLContainer) Start() error {
	return container.dockerContainer.Start()
}

func (container *MySQLContainer) Suspend() error {
	return container.dockerContainer.Suspend()
}

func (container *MySQLContainer) Stop() error {
	return container.dockerContainer.Stop()
}

func NewMySQLContainer() (*MySQLContainer, error) {
	port := 3306
	password := "my-secret-pw"
	database := "db0"
	connectionString, _ := mysql.ToConnectionString(
		mysql.Connector(mysql.DefaultProtocol, "127.0.0.1", fmt.Sprintf("%d", port)),
		mysql.Database(database),
		mysql.UserInfo("root", password),
	)
	checker := func(c *Container) error {
		return retry(10, 5*time.Second, func() error {
			db, err := sql.Open("mysql", connectionString)
			if err != nil {
				return err
			}
			defer db.Close()
			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database))
			if err != nil {
				return err
			}
			return nil
		})
	}
	container := &MySQLContainer{
		dockerContainer: NewDockerContainer(
			ImageRepository("mysql"),
			ImageTag("5.7"),
			Port(port),
			DockerEnv(
				[]string{
					fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", password),
					fmt.Sprintf("MYSQL_DATABASE=%s", database),
				},
			),
			HealthChecker(checker),
		),
	}

	container.URL = connectionString

	return container, nil
}
