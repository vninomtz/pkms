package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type task struct {
	Name        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Tasks []task

func (t *Tasks) Add(name string) {
	new_task := task{
		Name:        name,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, new_task)
}

func (t *Tasks) Complete(i int) error {
	ls := *t
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}
	// Adjust index for 0 based index
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

func (t *Tasks) Delete(i int) error {
	ls := *t
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}
	// Adjust index for 0 based index
	*t = append(ls[:i-1], ls[i:]...)

	return nil
}

func (t *Tasks) Save(filename string) error {
	js, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, js, 0644)
}

func (t *Tasks) Get(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	if len(file) == 0 {
		return nil
	}

	return json.Unmarshal(file, t)
}
