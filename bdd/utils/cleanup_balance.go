package utils

import "log"

func CleanupBalance(date string, account string, token string) error {
	query := `
		mutation { 
			deleteBalance(input: {
				date: "` + date + `"
				account: "` + account + `"
			})
		}
		`
	err := PostGraphQL(query, "deleteBalance", token, nil)
	if err != nil {
		log.Println("error on deleting balance:", err)
		return err
	}
	return nil
}
