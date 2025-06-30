package entity

import (
	"testing"
	"time"
)

func TestNewAccount(t *testing.T) {
	client := &Client{ID: "1", Name: "Test", Email: "test@example.com"}

	t.Run("should create a new account with valid client", func(t *testing.T) {
		account := NewAccount(client)

		if account == nil {
			t.Fatal("expected account to be created")
			return
		}

		if account.ID == "" {
			t.Error("expected ID to be generated")
		}

		if account.Client != client {
			t.Error("expected client to be set")
		}

		if account.Balance != 0 {
			t.Errorf("expected initial balance to be 0, got %f", account.Balance)
		}

		if account.CreatedAt.IsZero() {
			t.Error("expected CreatedAt to be set")
		}

		if account.UpdatedAt.IsZero() {
			t.Error("expected UpdatedAt to be set")
		}
	})

	t.Run("should return nil when client is nil", func(t *testing.T) {
		account := NewAccount(nil)

		if account != nil {
			t.Error("expected account to be nil")
		}
	})
}

func TestAccount_Credit(t *testing.T) {
	client := &Client{ID: "1", Name: "Test", Email: "test@example.com"}
	account := NewAccount(client)

	t.Run("should credit positive amount", func(t *testing.T) {
		previousUpdatedAt := account.UpdatedAt
		time.Sleep(1 * time.Millisecond) // ensure time difference

		account.Credit(100)

		if account.Balance != 100 {
			t.Errorf("expected balance to be 100, got %f", account.Balance)
		}

		if !account.UpdatedAt.After(previousUpdatedAt) {
			t.Error("expected UpdatedAt to be updated")
		}
	})

	t.Run("should not credit zero amount", func(t *testing.T) {
		previousBalance := account.Balance
		previousUpdatedAt := account.UpdatedAt
		time.Sleep(1 * time.Millisecond)

		account.Credit(0)

		if account.Balance != previousBalance {
			t.Errorf("expected balance to remain %f, got %f", previousBalance, account.Balance)
		}

		if account.UpdatedAt != previousUpdatedAt {
			t.Error("expected UpdatedAt to remain unchanged")
		}
	})

	t.Run("should not credit negative amount", func(t *testing.T) {
		previousBalance := account.Balance
		previousUpdatedAt := account.UpdatedAt
		time.Sleep(1 * time.Millisecond)

		account.Credit(-50)

		if account.Balance != previousBalance {
			t.Errorf("expected balance to remain %f, got %f", previousBalance, account.Balance)
		}

		if account.UpdatedAt != previousUpdatedAt {
			t.Error("expected UpdatedAt to remain unchanged")
		}
	})
}

func TestAccount_Debit(t *testing.T) {
	client := &Client{ID: "1", Name: "Test", Email: "test@example.com"}
	account := NewAccount(client)
	account.Credit(100) // Set initial balance

	t.Run("should debit valid amount", func(t *testing.T) {
		previousUpdatedAt := account.UpdatedAt
		time.Sleep(1 * time.Millisecond)

		account.Debit(50)

		if account.Balance != 50 {
			t.Errorf("expected balance to be 50, got %f", account.Balance)
		}

		if !account.UpdatedAt.After(previousUpdatedAt) {
			t.Error("expected UpdatedAt to be updated")
		}
	})

	t.Run("should not debit amount greater than balance", func(t *testing.T) {
		previousBalance := account.Balance
		previousUpdatedAt := account.UpdatedAt
		time.Sleep(1 * time.Millisecond)

		account.Debit(100)

		if account.Balance != previousBalance {
			t.Errorf("expected balance to remain %f, got %f", previousBalance, account.Balance)
		}

		if account.UpdatedAt != previousUpdatedAt {
			t.Error("expected UpdatedAt to remain unchanged")
		}
	})

	t.Run("should not debit zero amount", func(t *testing.T) {
		previousBalance := account.Balance
		previousUpdatedAt := account.UpdatedAt
		time.Sleep(1 * time.Millisecond)

		account.Debit(0)

		if account.Balance != previousBalance {
			t.Errorf("expected balance to remain %f, got %f", previousBalance, account.Balance)
		}

		if account.UpdatedAt != previousUpdatedAt {
			t.Error("expected UpdatedAt to remain unchanged")
		}
	})

	t.Run("should not debit negative amount", func(t *testing.T) {
		previousBalance := account.Balance
		previousUpdatedAt := account.UpdatedAt
		time.Sleep(1 * time.Millisecond)

		account.Debit(-10)

		if account.Balance != previousBalance {
			t.Errorf("expected balance to remain %f, got %f", previousBalance, account.Balance)
		}

		if account.UpdatedAt != previousUpdatedAt {
			t.Error("expected UpdatedAt to remain unchanged")
		}
	})
}
