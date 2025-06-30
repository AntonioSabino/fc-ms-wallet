package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	client1, _ := NewClient("John", "j@j.com")
	account1 := NewAccount(client1)
	client2, _ := NewClient("Jane", "jane@j.com")
	account2 := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	t.Run("should create a transaction", func(t *testing.T) {
		transaction, err := NewTransaction(account1, account2, 100)
		assert.Nil(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, account1, transaction.AccountFrom)
		assert.Equal(t, account2, transaction.AccountTo)
		assert.Equal(t, 100.0, transaction.Amount)
		assert.NotEmpty(t, transaction.ID)
		assert.NotEmpty(t, transaction.CreatedAt)
		assert.Equal(t, 900.0, account1.Balance)
		assert.Equal(t, 1100.0, account2.Balance)
	})

	t.Run("should return error when account from is nil", func(t *testing.T) {
		transaction, err := NewTransaction(nil, account2, 100)
		assert.NotNil(t, err)
		assert.Nil(t, transaction)
		assert.Equal(t, "account from cannot be nil", err.Error())
	})

	t.Run("should return error when account to is nil", func(t *testing.T) {
		transaction, err := NewTransaction(account1, nil, 100)
		assert.NotNil(t, err)
		assert.Nil(t, transaction)
		assert.Equal(t, "account to cannot be nil", err.Error())
	})

	t.Run("should return error when amount is less than or equal to zero", func(t *testing.T) {
		transaction, err := NewTransaction(account1, account2, 0)
		assert.NotNil(t, err)
		assert.Nil(t, transaction)
		assert.Equal(t, "amount must be greater than zero", err.Error())

		transaction, err = NewTransaction(account1, account2, -10)
		assert.NotNil(t, err)
		assert.Nil(t, transaction)
		assert.Equal(t, "amount must be greater than zero", err.Error())
	})

	t.Run("should return error when account from has insufficient funds", func(t *testing.T) {
		accountWithLowBalance := NewAccount(client1)
		accountWithLowBalance.Credit(50)

		transaction, err := NewTransaction(accountWithLowBalance, account2, 100)
		assert.NotNil(t, err)
		assert.Nil(t, transaction)
		assert.Equal(t, "insufficient funds in account from", err.Error())
	})
}

func TestTransaction_Validate(t *testing.T) {
	client1, _ := NewClient("John", "j@j.com")
	account1 := NewAccount(client1)
	client2, _ := NewClient("Jane", "jane@j.com")
	account2 := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	t.Run("should return nil when transaction is valid", func(t *testing.T) {
		transaction := &Transaction{
			AccountFrom: account1,
			AccountTo:   account2,
			Amount:      100,
		}
		err := transaction.Validate()
		assert.Nil(t, err)
	})

	t.Run("should return error when account from is nil", func(t *testing.T) {
		transaction := &Transaction{
			AccountTo: account2,
			Amount:    100,
		}
		err := transaction.Validate()
		assert.NotNil(t, err)
		assert.Equal(t, "account from cannot be nil", err.Error())
	})

	t.Run("should return error when account to is nil", func(t *testing.T) {
		transaction := &Transaction{
			AccountFrom: account1,
			Amount:      100,
		}
		err := transaction.Validate()
		assert.NotNil(t, err)
		assert.Equal(t, "account to cannot be nil", err.Error())
	})

	t.Run("should return error when amount is less than or equal to zero", func(t *testing.T) {
		transaction := &Transaction{
			AccountFrom: account1,
			AccountTo:   account2,
			Amount:      0,
		}
		err := transaction.Validate()
		assert.NotNil(t, err)
		assert.Equal(t, "amount must be greater than zero", err.Error())

		transaction.Amount = -10
		err = transaction.Validate()
		assert.NotNil(t, err)
		assert.Equal(t, "amount must be greater than zero", err.Error())
	})

	t.Run("should return error when account from has insufficient funds", func(t *testing.T) {
		accountWithLowBalance := NewAccount(client1)
		accountWithLowBalance.Credit(50)

		transaction := &Transaction{
			AccountFrom: accountWithLowBalance,
			AccountTo:   account2,
			Amount:      100,
		}
		err := transaction.Validate()
		assert.NotNil(t, err)
		assert.Equal(t, "insufficient funds in account from", err.Error())
	})
}

func TestTransaction_Commit(t *testing.T) {
	client1, _ := NewClient("John", "j@j.com")
	account1 := NewAccount(client1)
	client2, _ := NewClient("Jane", "jane@j.com")
	account2 := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	transaction := &Transaction{
		AccountFrom: account1,
		AccountTo:   account2,
		Amount:      100,
	}

	transaction.Commit()
	assert.Equal(t, 900.0, account1.Balance)
	assert.Equal(t, 1100.0, account2.Balance)
}
