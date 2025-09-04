package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type JSONOutput struct {
	UserID string `json:"user_id"`
	Count  int    `json:"count"`
}

func WriteJSONOutput(consumer *ConsumerStr, outputDir string) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	eventCounts := consumer.eventCounts

	createdUsers := []JSONOutput{}
	updatedUsers := []JSONOutput{}
	deletedUsers := []JSONOutput{}

	for userID, events := range eventCounts {
		if count := events["created"]; count > 0 {
			createdUsers = append(createdUsers, JSONOutput{UserID: userID, Count: count})
		}

		if count := events["updated"]; count > 0 {
			updatedUsers = append(updatedUsers, JSONOutput{UserID: userID, Count: count})
		}

		if count := events["deleted"]; count > 0 {
			deletedUsers = append(deletedUsers, JSONOutput{UserID: userID, Count: count})
		}
	}

	createdData, err := json.Marshal(createdUsers)
	if err != nil {
		return err
	}
	createdDataFilepath := filepath.Join(outputDir, "created.json")
	os.WriteFile(createdDataFilepath, createdData, 0644)

	updatedData, err := json.Marshal(updatedUsers)
	if err != nil {
		return err
	}
	updatedDataFilepath := filepath.Join(outputDir, "updated.json")
	os.WriteFile(updatedDataFilepath, updatedData, 0644)

	deletedData, err := json.Marshal(deletedUsers)
	if err != nil {
		return err
	}
	deletedDataFilepath := filepath.Join(outputDir, "deleted.json")
	if err := os.WriteFile(deletedDataFilepath, deletedData, 0644); err != nil {
		return err
	}

	return nil
}
