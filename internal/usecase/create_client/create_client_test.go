package createclient

import (
	"testing"

	"github.com/AntonioSabino/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientGatewayMock struct {
	mock.Mock
}

func (m *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
}

func (m *ClientGatewayMock) Save(client *entity.Client) error {
	args := m.Called(client)
	return args.Error(0)
}

func TestCreateClientUseCase_Execute(t *testing.T) {
	m := &ClientGatewayMock{}
	m.On("Save", mock.Anything).Return(nil)

	uc := NewCreateClientUseCase(m)

	input := CreateClientInputDTO{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	output, err := uc.Execute(input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "John Doe", output.Name)
	assert.Equal(t, "john@example.com", output.Email)
	assert.NotEmpty(t, output.ID)
	assert.NotEmpty(t, output.CreatedAt)
	assert.NotEmpty(t, output.UpdatedAt)

	m.AssertExpectations(t)
	m.AssertNumberOfCalls(t, "Save", 1)
}

func TestCreateClientUseCase_ExecuteWithInvalidName(t *testing.T) {
	m := &ClientGatewayMock{}

	uc := NewCreateClientUseCase(m)

	input := CreateClientInputDTO{
		Name:  "",
		Email: "john@example.com",
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Error(t, err, "name is required")

	m.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateClientUseCase_ExecuteWithInvalidEmail(t *testing.T) {
	m := &ClientGatewayMock{}

	uc := NewCreateClientUseCase(m)

	input := CreateClientInputDTO{
		Name:  "John Doe",
		Email: "",
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Error(t, err, "email is required")

	m.AssertNumberOfCalls(t, "Save", 0)
}

func TestCreateClientUseCase_ExecuteWithGatewayError(t *testing.T) {
	m := &ClientGatewayMock{}
	m.On("Save", mock.Anything).Return(assert.AnError)

	uc := NewCreateClientUseCase(m)

	input := CreateClientInputDTO{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	output, err := uc.Execute(input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
	assert.Error(t, err)

	m.AssertExpectations(t)
	m.AssertNumberOfCalls(t, "Save", 1)
}
