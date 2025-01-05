package bdd

type contextKey string

func (c contextKey) String() string {
	return "bdd context key - " + string(c)
}

var (
	url          = contextKey("url")
	balances     = contextKey("balances")
	transactions = contextKey("transactions")
	userToken    = contextKey("userToken")
	token        = contextKey("token")
	response     = contextKey("response")
)
