package filesystem

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type FilesystemProvider struct {
	directory string
	files     map[string][]byte // map[endpoint]content
	mutex     sync.RWMutex
	debug     bool
}

func NewFilesystemProvider(directory string, debug bool) (*FilesystemProvider, error) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory %s does not exist", directory)
	}

	return &FilesystemProvider{
		directory: directory,
		files:     make(map[string][]byte),
		debug:     debug,
	}, nil
}

func (lp *FilesystemProvider) LoadExistingFiles() error {
	return filepath.Walk(lp.directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file %s: %v", path, err)
				return nil
			}

			endpoint := lp.fileNameToEndpoint(filepath.Base(path))
			lp.updateFile(endpoint, content)

			if lp.debug {
				log.Printf("Loaded file: %s -> endpoint: %s", path, endpoint)
			}
		}

		return nil
	})
}

func (lp *FilesystemProvider) Watch(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	if err := watcher.Add(lp.directory); err != nil {
		return fmt.Errorf("failed to add directory to watcher: %w", err)
	}

	if lp.debug {
		log.Printf("Watching directory: %s", lp.directory)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			info, err := os.Stat(event.Name)
			if err == nil && info.IsDir() {
				continue
			}

			filename := filepath.Base(event.Name)
			endpoint := lp.fileNameToEndpoint(filename)

			switch {
			case event.Op&fsnotify.Write == fsnotify.Write:
				content, err := os.ReadFile(event.Name)
				if err != nil {
					log.Printf("Error reading file %s: %v", event.Name, err)
					continue
				}
				lp.updateFile(endpoint, content)
				if lp.debug {
					log.Printf("Modified file: %s -> endpoint: %s", filename, endpoint)
				}

			case event.Op&fsnotify.Create == fsnotify.Create:
				content, err := os.ReadFile(event.Name)
				if err != nil {
					log.Printf("Error reading file %s: %v", event.Name, err)
					continue
				}
				lp.updateFile(endpoint, content)
				if lp.debug {
					log.Printf("Created file: %s -> endpoint: %s", filename, endpoint)
				}

			case event.Op&fsnotify.Remove == fsnotify.Remove:
				lp.removeFile(endpoint)
				if lp.debug {
					log.Printf("Removed file: %s -> endpoint: %s", filename, endpoint)
				}

			case event.Op&fsnotify.Rename == fsnotify.Rename:
				lp.removeFile(endpoint)
				if lp.debug {
					log.Printf("Renamed/moved file: %s -> endpoint: %s", filename, endpoint)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func (lp *FilesystemProvider) GetFile(endpoint string) ([]byte, bool) {
	lp.mutex.RLock()
	defer lp.mutex.RUnlock()

	content, exists := lp.files[endpoint]
	return content, exists
}

func (lp *FilesystemProvider) ListEndpoints() []string {
	lp.mutex.RLock()
	defer lp.mutex.RUnlock()

	endpoints := make([]string, 0, len(lp.files))
	for endpoint := range lp.files {
		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}

func (lp *FilesystemProvider) fileNameToEndpoint(filename string) string {
	ext := filepath.Ext(filename)
	if ext != "" {
		return strings.TrimSuffix(filename, ext)
	}
	return filename
}

func (lp *FilesystemProvider) updateFile(endpoint string, content []byte) {
	lp.mutex.Lock()
	defer lp.mutex.Unlock()
	lp.files[endpoint] = content
}

func (lp *FilesystemProvider) removeFile(endpoint string) {
	lp.mutex.Lock()
	defer lp.mutex.Unlock()
	delete(lp.files, endpoint)
}
