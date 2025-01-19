package utils

import (
	"log"
)

func CleanupPlan(id string, token string) error {
	query := `
		mutation {
			deletePlan(id: "` + id + `")
		}
	`

	err := PostGraphQL(query, "deletePlan", token, nil)
	if err != nil {
		log.Printf("error on deleting plan %s: %v", id, err)
		return err
	}

	return nil
}
