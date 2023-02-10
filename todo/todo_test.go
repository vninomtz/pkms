package todo_test

import (
	"os"
	"testing"

	"github.com/vninomtz/swe-notes/todo"
)

func TestAdd(t *testing.T) {
	tasks := todo.Tasks{}
	name := "New Task"
	tasks.Add(name)
	if tasks[0].Name != name {
		t.Errorf("Expected %q, got %q instead.", name, tasks[0].Name)
	}
}

func TestComplete(t *testing.T) {
	tasks := todo.Tasks{}
	name := "New Task"
	tasks.Add(name)
	if tasks[0].Name != name {
		t.Errorf("Expected %q, got %q instead.", name, tasks[0].Name)
	}

	if tasks[0].Done {
		t.Errorf("New task should no be completed.")
	}

	tasks.Complete(1)

	if !tasks[0].Done {
		t.Errorf("New task should be completed.")
	}
}

func TestDelete(t *testing.T) {
	tasks := todo.Tasks{}
	names := []string{
		"New Task 1",
		"New Task 2",
		"New Task 3",
	}
	for _, v := range names {
		tasks.Add(v)
	}

	if tasks[0].Name != names[0] {
		t.Errorf("Expected %q, got %q instead.", names[0], tasks[0].Name)
	}

	tasks.Delete(2)

	if len(tasks) != 2 {
		t.Errorf("Expected list length %d, got %q instead.", 2, len(tasks))
	}

	if tasks[1].Name != names[2] {
		t.Errorf("Expected %q, got %q instead.", names[2], tasks[1].Name)
	}

}

func TestSaveGet(t *testing.T) {
	l1 := todo.Tasks{}
	l2 := todo.Tasks{}

	name := "New Task"
	l1.Add(name)

	tf, err := os.CreateTemp("", "")

	if err != nil {
		t.Fatalf("Error creating temp file %s", err)
	}

	defer os.Remove(tf.Name())

	if err := l1.Save(tf.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}

	if err := l2.Get(tf.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}

	if l1[0].Name != l2[0].Name {
		t.Errorf("Task %q should match %q task", l1[0].Name, l2[0].Name)
	}

}
