package models

import "sync"

type Task struct {
	ID          int      `json:"id"`
	URLs        []string `json:"url"`
	Status      string   `json:"status"`
	ArchivePath string   `json:"archivepath"`
	mu          sync.RWMutex
}

type AddURLsRequest struct {
	URLs []string `json:"urls" binding:"required,min=1,max=3"`
}
