package db

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jadejashaktisinh/gh-git-action-cli/config"
	"go.etcd.io/bbolt"
)

type Record struct {
	RunID      int64             `json:"run_id"`
	Timestamp  time.Time         `json:"timestamp"`
	Repository string            `json:"repository"`
	Workflow   string            `json:"workflow"`
	Branch     string            `json:"branch"`
	Inputs     map[string]string `json:"inputs"`
	Conclusion string            `json:"conclusion"`
}

var db *bbolt.DB
var bucketName = []byte("runs")

func InitDB() error {
	configDir := config.GetConfigPath()
	dbPath := filepath.Join(configDir, "history.db")

	var err error
	db, err = bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func SaveRun(record Record) error {
	return db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		
		id, _ := b.NextSequence()
		recordKey := fmt.Sprintf("%010d", id) // Use sequence for sorting or just use timestamp

		data, err := json.Marshal(record)
		if err != nil {
			return err
		}

		return b.Put([]byte(recordKey), data)
	})
}

func GetHistory() ([]Record, error) {
	var records []Record

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.ForEach(func(k, v []byte) error {
			var r Record
			if err := json.Unmarshal(v, &r); err != nil {
				return err
			}
			records = append(records, r)
			return nil
		})
	})

	// Return in reverse chronological order (newest first)
	for i, j := 0, len(records)-1; i < j; i, j = i+1, j-1 {
		records[i], records[j] = records[j], records[i]
	}

	return records, err
}
