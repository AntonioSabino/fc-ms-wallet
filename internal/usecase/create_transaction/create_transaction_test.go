package createtransaction

import (
	"errors"
	"testing"

	"github.com/AntonioSabino/fc-ms-wallet/internal/entity"
	"github.com/AntonioSabino/fc-ms-wallet/internal/event"
	"github.com/AntonioSabino/fc-ms-wallet/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Save(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Account), args.Error(1)
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	clientFrom, _ := entity.NewClient("John Doe", "john@example.com")
	clientTo, _ := entity.NewClient("Jane Doe", "jane@example.com")

	accountFrom := entity.NewAccount(clientFrom)
	accountFrom.Credit(100.0) // Adding balance to account

	accountTo := entity.NewAccount(clientTo)

	transactionGateway := &TransactionGatewayMock{}
	accountGateway := &AccountGatewayMock{}

	accountGateway.On("FindByID", "account-from-id").Return(accountFrom, nil)
	accountGateway.On("FindByID", "account-to-id").Return(accountTo, nil)
	transactionGateway.On("Save", mock.Anything).Return(nil)

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	uc := NewCreateTransactionUseCase(transactionGateway, accountGateway, dispatcher, event)

	input := CreateTransactionInputDTO{
		AccountIDFrom: "account-from-id",
		AccountIDTo:   "account-to-id",
		Amount:        50.0,
	}

	output, err := uc.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)

	// Verify that balances were updated
	assert.Equal(t, 50.0, accountFrom.Balance)
	assert.Equal(t, 50.0, accountTo.Balance)

	transactionGateway.AssertExpectations(t)
	accountGateway.AssertExpectations(t)
	transactionGateway.AssertNumberOfCalls(t, "Save", 1)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 2)
}

func TestCreateTransactionUseCase_ExecuteWithAccountFromNotFound(t *testing.T) {
	transactionGateway := &TransactionGatewayMock{}
	accountGateway := &AccountGatewayMock{}

	accountGateway.On("FindByID", "account-from-id").Return(nil, errors.New("account not found"))

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	uc := NewCreateTransactionUseCase(transactionGateway, accountGateway, dispatcher, event)

	input := CreateTransactionInputDTO{
		AccountIDFrom: "account-from-id",
		AccountIDTo:   "account-to-id",
		Amount:        50.0,
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "account not found", err.Error())

	accountGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 1)
	transactionGateway.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateTransactionUseCase_ExecuteWithAccountToNotFound(t *testing.T) {
	clientFrom, _ := entity.NewClient("John Doe", "john@example.com")
	accountFrom := entity.NewAccount(clientFrom)
	accountFrom.Credit(100.0)

	transactionGateway := &TransactionGatewayMock{}
	accountGateway := &AccountGatewayMock{}

	accountGateway.On("FindByID", "account-from-id").Return(accountFrom, nil)
	accountGateway.On("FindByID", "account-to-id").Return(nil, errors.New("account not found"))

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	uc := NewCreateTransactionUseCase(transactionGateway, accountGateway, dispatcher, event)

	input := CreateTransactionInputDTO{
		AccountIDFrom: "account-from-id",
		AccountIDTo:   "account-to-id",
		Amount:        50.0,
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "account not found", err.Error())

	accountGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 2)
	transactionGateway.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateTransactionUseCase_ExecuteWithInsufficientFunds(t *testing.T) {
	clientFrom, _ := entity.NewClient("John Doe", "john@example.com")
	clientTo, _ := entity.NewClient("Jane Doe", "jane@example.com")

	accountFrom := entity.NewAccount(clientFrom)
	accountFrom.Credit(30.0) // Less than the transaction amount

	accountTo := entity.NewAccount(clientTo)

	transactionGateway := &TransactionGatewayMock{}
	accountGateway := &AccountGatewayMock{}

	accountGateway.On("FindByID", "account-from-id").Return(accountFrom, nil)
	accountGateway.On("FindByID", "account-to-id").Return(accountTo, nil)

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	uc := NewCreateTransactionUseCase(transactionGateway, accountGateway, dispatcher, event)

	input := CreateTransactionInputDTO{
		AccountIDFrom: "account-from-id",
		AccountIDTo:   "account-to-id",
		Amount:        50.0,
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "insufficient funds in account from", err.Error())

	accountGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 2)
	transactionGateway.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateTransactionUseCase_ExecuteWithZeroAmount(t *testing.T) {
	clientFrom, _ := entity.NewClient("John Doe", "john@example.com")
	clientTo, _ := entity.NewClient("Jane Doe", "jane@example.com")

	accountFrom := entity.NewAccount(clientFrom)
	accountFrom.Credit(100.0)

	accountTo := entity.NewAccount(clientTo)

	transactionGateway := &TransactionGatewayMock{}
	accountGateway := &AccountGatewayMock{}

	accountGateway.On("FindByID", "account-from-id").Return(accountFrom, nil)
	accountGateway.On("FindByID", "account-to-id").Return(accountTo, nil)

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	uc := NewCreateTransactionUseCase(transactionGateway, accountGateway, dispatcher, event)

	input := CreateTransactionInputDTO{
		AccountIDFrom: "account-from-id",
		AccountIDTo:   "account-to-id",
		Amount:        0.0,
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "amount must be greater than zero", err.Error())

	accountGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 2)
	transactionGateway.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateTransactionUseCase_ExecuteWithNegativeAmount(t *testing.T) {
	clientFrom, _ := entity.NewClient("John Doe", "john@example.com")
	clientTo, _ := entity.NewClient("Jane Doe", "jane@example.com")

	accountFrom := entity.NewAccount(clientFrom)
	accountFrom.Credit(100.0)

	accountTo := entity.NewAccount(clientTo)

	transactionGateway := &TransactionGatewayMock{}
	accountGateway := &AccountGatewayMock{}

	accountGateway.On("FindByID", "account-from-id").Return(accountFrom, nil)
	accountGateway.On("FindByID", "account-to-id").Return(accountTo, nil)

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	uc := NewCreateTransactionUseCase(transactionGateway, accountGateway, dispatcher, event)

	input := CreateTransactionInputDTO{
		AccountIDFrom: "account-from-id",
		AccountIDTo:   "account-to-id",
		Amount:        -10.0,
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "amount must be greater than zero", err.Error())

	accountGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 2)
	transactionGateway.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateTransactionUseCase_ExecuteWithTransactionGatewayError(t *testing.T) {
	clientFrom, _ := entity.NewClient("John Doe", "john@example.com")
	clientTo, _ := entity.NewClient("Jane Doe", "jane@example.com")

	accountFrom := entity.NewAccount(clientFrom)
	accountFrom.Credit(100.0)

	accountTo := entity.NewAccount(clientTo)

	transactionGateway := &TransactionGatewayMock{}
	accountGateway := &AccountGatewayMock{}

	accountGateway.On("FindByID", "account-from-id").Return(accountFrom, nil)
	accountGateway.On("FindByID", "account-to-id").Return(accountTo, nil)
	transactionGateway.On("Save", mock.Anything).Return(errors.New("database error"))

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	uc := NewCreateTransactionUseCase(transactionGateway, accountGateway, dispatcher, event)

	input := CreateTransactionInputDTO{
		AccountIDFrom: "account-from-id",
		AccountIDTo:   "account-to-id",
		Amount:        50.0,
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "database error", err.Error())

	// Note: The transaction was committed (balances changed) but save failed
	assert.Equal(t, 50.0, accountFrom.Balance)
	assert.Equal(t, 50.0, accountTo.Balance)

	transactionGateway.AssertExpectations(t)
	accountGateway.AssertExpectations(t)
	transactionGateway.AssertNumberOfCalls(t, "Save", 1)
	accountGateway.AssertNumberOfCalls(t, "FindByID", 2)
}

func TestNewCreateTransactionUseCase(t *testing.T) {
	transactionGateway := &TransactionGatewayMock{}
	accountGateway := &AccountGatewayMock{}

	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreated()

	uc := NewCreateTransactionUseCase(transactionGateway, accountGateway, dispatcher, event)

	assert.NotNil(t, uc)
	assert.Equal(t, transactionGateway, uc.TransactionGateway)
	assert.Equal(t, accountGateway, uc.AccountGateway)
	assert.Equal(t, dispatcher, uc.EventDispatcher)
	assert.Equal(t, event, uc.TransactionCreated)
}
