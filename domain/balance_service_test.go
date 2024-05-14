package domain

import (
	"cashflow/models"
	"cashflow/repository"
	"fmt"
	"reflect"
	"testing"
)

func setupBalanceService() models.BalanceService {
	_, br := repository.Init(models.InMemory)
	defer repository.Close()

	bs := NewBalanceService(br)
	return bs
}

func TestWriteBalanceWithoutDate(t *testing.T) {
	service := setupBalanceService()

	result, _ := service.WriteBalance(23.22, nil)
	if result.Amount != 23.22 {
		t.Errorf(`The balance should have been %f but was %f`, 23.22, result.Amount)
	}
}

func TestWriteBalanceRoundedToTwoDecimalsWithDate(t *testing.T) {
	service := setupBalanceService()

	dateInput := "2000/01/02"

	result, _ := service.WriteBalance(23.222222, &dateInput)

	if result.Amount != 23.22 {
		t.Errorf(`The balance should have been %f but was %f`, 23.22, result.Amount)
	}

	if result.Date != dateInput {
		t.Errorf(`The date should have been %s but was %s`, dateInput, result.Date)
	}
}

func TestListBalances(t *testing.T) {
	service := setupBalanceService()
	result, _ := service.ListBalances()
	expected := []models.Balance{
		{Date: "2000/01/01", Amount: 50},
		{Date: "2000/01/05", Amount: 100},
	}

	if len(result) == 0 || reflect.ValueOf(result[0]).Kind() == reflect.Ptr || !reflect.DeepEqual(expected, result) {
		t.Errorf(`Wrong data, expected %s, instead got %s`, fmt.Sprint(expected), fmt.Sprint(result))
	}
}
