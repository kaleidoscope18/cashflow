package reporting

import (
	"cashflow/models"
	"encoding/json"
	"fmt"
)

// TODO use package termtables :D
func PrintCommandLine(transactionService models.TransactionService) {
	w := "2000/01/03"
	js, _ := json.MarshalIndent(transactionService.ListTransactions(&w), "", " ")
	fmt.Println(string(js))
}
