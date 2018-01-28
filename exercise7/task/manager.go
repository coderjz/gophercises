package task

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

//TaskBucket returns the byte slice for the task bucket
func TaskBucket() []byte {
	return []byte("TaskBucket")
}

//Manager provides our access to the list of tasks
type Manager struct {
	tasks []*task
	db    *bolt.DB
}

//SetDB passes in the boltDB instance
func (m *Manager) SetDB(db *bolt.DB) {
	m.db = db
}

//Add will add a new task with the given name
func (m *Manager) Add(taskName string) {
	err := m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(TaskBucket())
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		task := task{
			ID:     id,
			Name:   taskName,
			IsDone: false,
		}

		buf, err := json.Marshal(&task)
		if err != nil {
			return err
		}

		return b.Put(itob(task.ID), buf)
	})
	if err != nil {
		fmt.Println(err)
	}
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

//List output all tasks that are not yet done
func (m *Manager) List() (string, error) {
	result := ""
	// Assume bucket exists and has keys
	err := m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(TaskBucket())

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task task
			err := json.Unmarshal(v, &task)
			if err != nil {
				return fmt.Errorf("error getting task: %v", err)
			}
			if task.IsDone {
				continue
			}
			result += fmt.Sprintf("#%d: %s\n", task.ID, task.Name)
		}

		return nil
	})
	if err != nil {
		return "", err
	}
	return result, nil
}

//Do marks a task as done
func (m *Manager) Do(taskIndex int) error {

	return m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(TaskBucket())
		taskBytes := b.Get(itob(uint64(taskIndex)))
		if taskBytes == nil {
			return fmt.Errorf("Could not find any task at position %d", taskIndex)
		}
		var task task
		err := json.Unmarshal(taskBytes, &task)
		if err != nil {
			return fmt.Errorf("error getting task: %v", err)
		}

		if task.IsDone {
			return fmt.Errorf("Task at position %d is already done", taskIndex)
		}
		task.IsDone = true

		buf, err := json.Marshal(&task)
		if err != nil {
			return fmt.Errorf("error updating task: %v", err)
		}

		return b.Put(itob(task.ID), buf)
	})
}

type task struct {
	ID     uint64
	Name   string
	IsDone bool
}
