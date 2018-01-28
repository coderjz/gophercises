// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/coderjz/gophercises/exercise7/cmd"
	"github.com/coderjz/gophercises/exercise7/task"
)

func main() {
	db, err := bolt.Open("task.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	manager := &task.Manager{}
	manager.SetDB(db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(task.TaskBucket())
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	cmd.Execute(manager)
}
