package repository

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"database/sql"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetTransactionById_Success(t *testing.T) {
	ctx := context.Background()
	id := "123"
	expectedTransaction := models.Transaction{
		Id:          "123",
		Description: "Test Transaction",
		Amount:      100.0,
		Date:        "2023-10-01",
		Recurrency:  "monthly",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "description", "amount", "date", "recurrency"}).
		AddRow(expectedTransaction.Id, expectedTransaction.Description, expectedTransaction.Amount, expectedTransaction.Date, expectedTransaction.Recurrency)

	mock.ExpectQuery("SELECT id, description, amount, date, recurrency FROM transactions WHERE id = ?").
		WithArgs(id).
		WillReturnRows(rows)

	repo := &mysqlRepository{
		db:    db,
		mutex: sync.RWMutex{},
	}

	transaction, err := repo.GetTransactionById(ctx, id)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if transaction.Id != expectedTransaction.Id || transaction.Description != expectedTransaction.Description || transaction.Amount != expectedTransaction.Amount || transaction.Date != utils.ParseDate(expectedTransaction.Date) || transaction.Recurrency != expectedTransaction.Recurrency {
		t.Errorf("expected %v, got %v", expectedTransaction, transaction)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetTransactionById_NotFound(t *testing.T) {
	ctx := context.Background()
	id := "nonexistent"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, description, amount, date, recurrency FROM transactions WHERE id = ?").
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	repo := &mysqlRepository{
		db:    db,
		mutex: sync.RWMutex{},
	}

	_, err = repo.GetTransactionById(ctx, id)
	if err == nil || err.Error() != "transaction not found" {
		t.Errorf("expected 'transaction not found' error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListTransactions_Success(t *testing.T) {
	ctx := context.Background()
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "description", "amount", "date", "recurrency"}).
		AddRow(1, "Test Transaction", 100.0, "2023-06-15", "monthly")

	mock.ExpectQuery("SELECT id, description, amount, date, recurrency FROM transactions WHERE date BETWEEN \\? AND \\? ORDER BY date ASC").
		WithArgs(from, to).
		WillReturnRows(rows)

	repo := &mysqlRepository{
		db:    db,
		mutex: sync.RWMutex{},
	}

	transactions, err := repo.ListTransactions(ctx, from, to)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if len(transactions) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(transactions))
	}

	if transactions[0].Description != "Test Transaction" {
		t.Errorf("expected 'Test Transaction', got %s", transactions[0].Description)
	}
}

func TestListTransactions_InvalidDateFormat(t *testing.T) {
	ctx := context.Background()
	from := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, description, amount, date, recurrency FROM transactions WHERE date BETWEEN \\? AND \\? ORDER BY date ASC").
		WithArgs(from, "invalid")

	repo := &mysqlRepository{
		db:    db,
		mutex: sync.RWMutex{},
	}

	transactions, err := repo.ListTransactions(ctx, from, to)
	if err == nil {
		t.Errorf("expected an error due to invalid date format")
	}

	if len(transactions) != 0 {
		t.Errorf("expected 0 transactions due to invalid date format, got %d", len(transactions))
	}
}

func TestInsertTransaction_Success(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := &mysqlRepository{
		db:    mockDB,
		mutex: sync.RWMutex{},
	}

	transaction := models.Transaction{
		Amount:      100.0,
		Date:        "2023-10-01",
		Description: "Test Transaction",
		Recurrency:  "monthly",
	}

	mock.ExpectPrepare("INSERT INTO transactions").
		ExpectExec().
		WithArgs(transaction.Amount, transaction.Date, transaction.Description, transaction.Recurrency).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.InsertTransaction(ctx, transaction)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if id != "1" {
		t.Errorf("expected id to be '1', got %s", id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsertTransaction_SQLPreparationError(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := &mysqlRepository{
		db:    mockDB,
		mutex: sync.RWMutex{},
	}

	transaction := models.Transaction{
		Amount:      100.0,
		Date:        "2023-10-01",
		Description: "Test Transaction",
		Recurrency:  "monthly",
	}

	mock.ExpectPrepare("INSERT INTO transactions").
		WillReturnError(fmt.Errorf("preparation error"))

	_, err = repo.InsertTransaction(ctx, transaction)
	if err == nil {
		t.Error("expected an error but got none")
	}

	if err.Error() != "preparation error" {
		t.Errorf("expected 'preparation error', got %s", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEditTransaction_Success(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := &mysqlRepository{
		db:    mockDB,
		mutex: sync.RWMutex{},
	}

	transaction := models.Transaction{
		Id:          "1",
		Description: "Updated Description",
		Amount:      100.0,
		Recurrency:  "monthly",
		Date:        "2023-10-01",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE transactions").
		WithArgs(transaction.Description, transaction.Amount, transaction.Recurrency, transaction.Date, transaction.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	id, err := repo.EditTransaction(ctx, transaction)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if id != transaction.Id {
		t.Errorf("expected id %s, got %s", transaction.Id, id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteTransaction_Success(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := &mysqlRepository{
		db:    mockDB,
		mutex: sync.RWMutex{},
	}

	id := "1"
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM transactions WHERE id = ?").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	result, err := repo.DeleteTransaction(ctx, id)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if result != id {
		t.Errorf("expected %s, got %s", id, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDeleteTransaction_NotFound(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := &mysqlRepository{
		db:    mockDB,
		mutex: sync.RWMutex{},
	}

	id := "1"
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM transactions WHERE id = ?").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 0))
	mock.ExpectRollback()

	result, err := repo.DeleteTransaction(ctx, id)
	if err == nil {
		t.Errorf("expected error, got none")
	}
	if result != id {
		t.Errorf("expected %s, got %s", id, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestEditTransaction_InvalidID(t *testing.T) {
	ctx := context.Background()
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	repo := &mysqlRepository{
		db:    mockDB,
		mutex: sync.RWMutex{},
	}

	transaction := models.Transaction{
		Id:          "999",
		Description: "Updated Description",
		Amount:      100.0,
		Recurrency:  "monthly",
		Date:        "2023-10-01",
	}

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE transactions").
		WithArgs(transaction.Description, transaction.Amount, transaction.Recurrency, transaction.Date, transaction.Id).
		WillReturnError(sql.ErrNoRows)
	mock.ExpectRollback()

	id, err := repo.EditTransaction(ctx, transaction)
	if err == nil {
		t.Errorf("expected error but got none")
	}

	if id != transaction.Id {
		t.Errorf("expected id %s, got %s", transaction.Id, id)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
