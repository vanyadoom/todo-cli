package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"todo-cli/task"
)

type Storage struct {
	filename string
	Tasks    []task.Task
}

func NewStorage(filename string) *Storage {
	return &Storage{
		filename: filename,
		Tasks:    []task.Task{},
	}
}

func (s *Storage) Load() error {

	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		s.Tasks = []task.Task{}
		return nil
	}

	err = json.Unmarshal(data, &s.Tasks)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Save() error {
	data, err := json.MarshalIndent(s.Tasks, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.filename, data, 0644)
	if err != nil {
		return err
	}

	return nil

}

func (s *Storage) Add(t *task.Task) {
	s.Tasks = append(s.Tasks, *t)
}

func (s *Storage) Delete(id string) bool {
	for i, t := range s.Tasks {
		if strings.HasPrefix(t.ID, id) {
			s.Tasks = append(s.Tasks[:i], s.Tasks[i+1:]...)
			return true
		}
	}
	return false
}

func (s *Storage) FindByID(id string) *task.Task {
	for i := range s.Tasks {
		if strings.HasPrefix(s.Tasks[i].ID, id) {
			return &s.Tasks[i]
		}
	}
	return nil
}

func (s *Storage) GetByStatus(status task.Status) []task.Task {
	result := []task.Task{}

	for _, t := range s.Tasks {
		if t.Status == status {
			result = append(result, t)
		}
	}

	return result
}
