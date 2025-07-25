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
