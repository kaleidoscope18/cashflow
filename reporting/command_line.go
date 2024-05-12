package reporting

import (
	"cashflow/models"
	"encoding/json"
	"fmt"
)

// TODO use package termtables :D
func PrintCommandLine(transactionService models.TransactionService) {
	js, _ := json.MarshalIndent(transactionService.ListTransactions("2000/01/03"), "", " ")
	fmt.Println(string(js))
}
