package utils

func CleanupTransactions(transactionIds []string, token string) error {
	for _, id := range transactionIds {
		query := `
		mutation { 
			deleteTransaction(input: {
				id: "` + id + `"
				account: "chequing"
			})
		}
		`
		err := PostGraphQL(query, "deleteTransaction", token, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
