package services

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"smthtozip/internal/models"
	"strconv"
	"strings"
	"sync"
)

type Archiver struct {
	storage fileStorage
	mu      sync.Mutex
}

func NewArchiver(storage fileStorage) *Archiver {
	return &Archiver{storage: storage}
}

func (a *Archiver) Process(task *models.Task, baseURL string) error {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	log.Println("ya tut")
	log.Println(task.URLs)
	for _, url := range task.URLs {
		if isAllowedType(url) {
			if err := addFileToZip(zipWriter, url); err != nil {
				log.Println("ya oshibka")
				continue
			}
		} else {
			log.Println("ya oshibka2")
		}
	}

	if err := zipWriter.Close(); err != nil {
		return err
	}

	filename := strconv.Itoa(task.ID) + ".zip"
	if err := a.storage.Save(filename, buf); err != nil {
		return err
	}

	a.mu.Lock()
	task.ArchivePath = baseURL + "/archives/" + filename + "?download=1"
	task.Status = "completed"
	a.mu.Unlock()

	return nil
}

func isAllowedType(url string) bool {
	ext := strings.ToLower(filepath.Ext(url))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".pdf"
}

func addFileToZip(zipWriter *zip.Writer, url string) error {
	log.Println("url is: " + url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileName := filepath.Base(url)
	writer, err := zipWriter.Create(fileName)
	if err != nil {
		log.Println("failed to create file in archive")
		return err
	}

	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		log.Println("failed to write to archive")
	}
	return err
}
