package models

type WithDate interface {
	GetDate() string
}

func (t Transaction) GetDate() string {
	return t.Date
}

func (t ComputedTransaction) GetDate() string {
	return t.Date
}

func (t Balance) GetDate() string {
	return t.Date
}
