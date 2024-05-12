package domain

import (
	"cashflow/models"
	"cashflow/repository"
	"fmt"
	"reflect"
	"testing"
)

func setup() models.TransactionService {
	storage := repository.New("Mocked")
	return New(storage)
}

func TestWriteBalanceWithoutDate(t *testing.T) {
	service := setup()

	result := service.WriteBalance(23.22, nil)
	if result.Amount != 23.22 {
		t.Errorf(`The balance should have been %f but was %f`, 23.22, result.Amount)
	}
}

func TestWriteBalanceRoundedToTwoDecimalsWithDate(t *testing.T) {
	service := setup()

	dateInput := "2000/01/02"

	result := service.WriteBalance(23.222222, &dateInput)

	if result.Amount != 23.22 {
		t.Errorf(`The balance should have been %f but was %f`, 23.22, result.Amount)
	}

	if result.Date != dateInput {
		t.Errorf(`The date should have been %s but was %s`, dateInput, result.Date)
	}
}

func TestListBalances(t *testing.T) {
	service := setup()
	result := service.ListBalances()

	if len(result) == 0 || reflect.ValueOf(result[0]).Kind() != reflect.Ptr {
		t.Errorf(`There should be an array of pointers to balances, instead got %s`, fmt.Sprint(result))
	}
}
