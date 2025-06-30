package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("John Doe", "john.doe@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, "John Doe", client.Name)
	assert.Equal(t, "john.doe@example.com", client.Email)
	assert.False(t, client.CreatedAt.IsZero())
	assert.False(t, client.UpdatedAt.IsZero())
}

func TestCreateNewClientWhenArgsAreInvalid(t *testing.T) {
	client, err := NewClient("", "")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestUpdateClient(t *testing.T) {
	client, _ := NewClient("Jane Doe", "jane.doe@example.com")
	err := client.Update("Jane Smith", "jane.smith@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "Jane Smith", client.Name)
	assert.Equal(t, "jane.smith@example.com", client.Email)
}

func TestUpdateClientWithInvalidArgs(t *testing.T) {
	client, _ := NewClient("Jane Doe", "jane.doe@example.com")
	err := client.Update("", "")
	assert.Error(t, err, "name is required")
}

func TestAddAccount(t *testing.T) {
	client, _ := NewClient("Alice", "alice@example.com")
	account := NewAccount(client)
	err := client.AddAccount(account)
	assert.NoError(t, err)
	assert.Contains(t, client.Accounts, account)
}
