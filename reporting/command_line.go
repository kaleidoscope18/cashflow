package reporting

import (
	"cashflow/models"
	"encoding/json"
	"fmt"
)

// TODO use package termtables :D
func PrintCommandLine(transactionService models.TransactionService) {
	transactions, err := (transactionService).ListTransactions(nil)
	if err != nil {
		panic(err.Error())
	}
	js, _ := json.MarshalIndent(transactions, "", " ")
	fmt.Println(string(js))
}
