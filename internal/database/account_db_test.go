package database

import (
	"database/sql"
	"testing"

	"github.com/AntonioSabino/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
	client    *entity.Client
}

func (s *AccountDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance decimal, created_at date, FOREIGN KEY(client_id) REFERENCES clients(id))")

	s.accountDB = NewAccountDB(db)
	s.client, _ = entity.NewClient("Jane Doe", "jane.doe@example.com")
}

func (s *AccountDBTestSuite) TearDownTest() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE clients")
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (s *AccountDBTestSuite) TestSaveAccount() {
	account := entity.NewAccount(s.client)
	err := s.accountDB.Save(account)
	s.Nil(err)
}

func (s *AccountDBTestSuite) TestFindByID() {
	s.db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)", s.client.ID, s.client.Name, s.client.Email, s.client.CreatedAt)
	account := entity.NewAccount(s.client)
	err := s.accountDB.Save(account)
	s.Nil(err)
	retrievedAccount, err := s.accountDB.FindByID(account.ID)
	s.Nil(err)
	s.NotNil(retrievedAccount)
	s.Equal(account.ID, retrievedAccount.ID)
	s.Equal(s.client.ID, retrievedAccount.Client.ID)
	s.Equal(s.client.Name, retrievedAccount.Client.Name)
	s.Equal(s.client.Email, retrievedAccount.Client.Email)
}
