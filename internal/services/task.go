package services

import (
	"errors"
	"log"
	"smthtozip/internal/models"
	"sync"
)

type TaskService struct {
	tasks       map[int]*models.Task
	nextID      int
	archiver    *Archiver
	mu          sync.Mutex
	activeTasks chan struct{}
}

func NewTaskService(archiver *Archiver) *TaskService {
	return &TaskService{
		tasks:       make(map[int]*models.Task),
		nextID:      1,
		activeTasks: make(chan struct{}, 3),
		archiver:    archiver,
	}
}

func (ts *TaskService) CreateTask() (*models.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	task := &models.Task{
		ID:     ts.nextID,
		Status: "pending",
		URLs:   []string{},
	}

	ts.tasks[task.ID] = task
	ts.nextID++

	return task, nil
}

func (ts *TaskService) Delete(id int) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if _, err := ts.tasks[id]; !err {
		return errors.New("task not found")
	}
	delete(ts.tasks, id)
	return nil
}

func (ts *TaskService) GetByID(id int) (*models.Task, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	task, exists := ts.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}
	// Возвращаем указатель на существующую задачу
	return task, nil
}

func (ts *TaskService) processTask(id int) {
	select {
	case ts.activeTasks <- struct{}{}:
		defer func() { <-ts.activeTasks }()
	default:
		return
	}
	task, exists := ts.GetByID(id)
	if exists != nil {
		return
	}

	ts.mu.Lock()
	task.Status = "proccessing"
	ts.mu.Unlock()

	if err := ts.archiver.Process(task, "http://localhost:8080"); err != nil {
		ts.mu.Lock()
		task.Status = "failed"
		ts.mu.Unlock()
	}
}

func (ts *TaskService) AddURL(taskID int, url string) error {
	task, exists := ts.GetByID(taskID)
	if exists != nil {
		return errors.New("task not found")
	}

	if len(task.URLs) >= 3 {
		return errors.New("maximum URLs per task reached")
	}

	task.URLs = append(task.URLs, url)
	log.Println(task.URLs)
	return nil
}
