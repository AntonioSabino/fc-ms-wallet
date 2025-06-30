package createaccount

import (
	"errors"
	"testing"

	"github.com/AntonioSabino/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

type ClientGatewayMock struct {
	mock.Mock
}

func (m *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) Save(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func TestCreateAccountUseCase_Execute(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "john@example.com")

	accountGateway := &AccountGatewayMock{}
	clientGateway := &ClientGatewayMock{}

	clientGateway.On("Get", "123").Return(client, nil)
	accountGateway.On("Save", mock.Anything).Return(nil)

	uc := NewCreateAccountUseCase(accountGateway, clientGateway)

	input := CreateAccountInputDTO{
		ClientID: "123",
	}

	output, err := uc.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)

	accountGateway.AssertExpectations(t)
	clientGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "Save", 1)
	clientGateway.AssertNumberOfCalls(t, "Get", 1)
}

func TestCreateAccountUseCase_ExecuteWithClientNotFound(t *testing.T) {
	accountGateway := &AccountGatewayMock{}
	clientGateway := &ClientGatewayMock{}

	clientGateway.On("Get", "123").Return(nil, errors.New("client not found"))

	uc := NewCreateAccountUseCase(accountGateway, clientGateway)

	input := CreateAccountInputDTO{
		ClientID: "123",
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "client not found", err.Error())

	clientGateway.AssertExpectations(t)
	clientGateway.AssertNumberOfCalls(t, "Get", 1)
	accountGateway.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateAccountUseCase_ExecuteWithAccountGatewayError(t *testing.T) {
	client, _ := entity.NewClient("John Doe", "john@example.com")

	accountGateway := &AccountGatewayMock{}
	clientGateway := &ClientGatewayMock{}

	clientGateway.On("Get", "123").Return(client, nil)
	accountGateway.On("Save", mock.Anything).Return(errors.New("database error"))

	uc := NewCreateAccountUseCase(accountGateway, clientGateway)

	input := CreateAccountInputDTO{
		ClientID: "123",
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "database error", err.Error())

	accountGateway.AssertExpectations(t)
	clientGateway.AssertExpectations(t)
	accountGateway.AssertNumberOfCalls(t, "Save", 1)
	clientGateway.AssertNumberOfCalls(t, "Get", 1)
}

func TestCreateAccountUseCase_ExecuteWithEmptyClientID(t *testing.T) {
	accountGateway := &AccountGatewayMock{}
	clientGateway := &ClientGatewayMock{}

	clientGateway.On("Get", "").Return(nil, errors.New("client id is required"))

	uc := NewCreateAccountUseCase(accountGateway, clientGateway)

	input := CreateAccountInputDTO{
		ClientID: "",
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Equal(t, "client id is required", err.Error())

	clientGateway.AssertExpectations(t)
	clientGateway.AssertNumberOfCalls(t, "Get", 1)
	accountGateway.AssertNumberOfCalls(t, "Save", 0)
}

func TestNewCreateAccountUseCase(t *testing.T) {
	accountGateway := &AccountGatewayMock{}
	clientGateway := &ClientGatewayMock{}

	uc := NewCreateAccountUseCase(accountGateway, clientGateway)

	assert.NotNil(t, uc)
	assert.Equal(t, accountGateway, uc.AccountGateway)
	assert.Equal(t, clientGateway, uc.ClientGateway)
}
