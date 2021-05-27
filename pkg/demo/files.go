package demo

import (
	"os"
	"path/filepath"
	"time"

	"github.com/Cludch/csgo-tools/internal/domain/entity"
	log "github.com/sirupsen/logrus"
)

// Demo holds meta information about a demo file.
type Demo struct {
	ID        entity.ID
	MatchTime time.Time
	Filename  string
}

// ScanDemosDir scans the demos dir and returns all match ids.
func ScanDemosDir(path string) []*Demo {
	var demos []*Demo

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			// Ignore non .dem-files
			if filepath.Ext(path) != ".dem" {
				return nil
			}

			fileName := info.Name()

			// Get file creation date.
			modTime := time.Now()

			stats, err := os.Stat(path)
			if err != nil {
				log.Errorf("Unable to read file stats for %v", fileName)
			} else {
				modTime = stats.ModTime()
			}

			// Add demo
			demos = append(demos, &Demo{ID: entity.NewID(), MatchTime: modTime, Filename: fileName})

			return nil
		})
	if err != nil {
		log.Error(err)
	}

	return demos
}
