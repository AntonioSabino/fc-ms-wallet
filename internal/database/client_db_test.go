package database

import (
	"database/sql"
	"testing"

	"github.com/AntonioSabino/fc-ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
	_ "modernc.org/sqlite"
)

type ClientDBTestSuite struct {
	suite.Suite
	db       *sql.DB
	clientDB *ClientDB
}

func (s *ClientDBTestSuite) SetupTest() {
	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	s.Nil(err)
	s.db = db
	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255))")
	s.clientDB = NewClientDB(db)
}

func (s *ClientDBTestSuite) TearDownTest() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE clients")
}

func TestClientDBTestSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTestSuite))
}

func (s *ClientDBTestSuite) TestSaveClient() {
	client, _ := entity.NewClient("Jane Doe", "jane.doe@example.com")
	err := s.clientDB.Save(client)
	s.Nil(err)

	retrievedClient, err := s.clientDB.Get(client.ID)
	s.Nil(err)
	s.Equal(client.ID, retrievedClient.ID)
}

func (s *ClientDBTestSuite) TestGetClient() {
	client, _ := entity.NewClient("John Doe", "john.doe@example.com")
	s.clientDB.Save(client)

	retrievedClient, err := s.clientDB.Get(client.ID)
	s.Nil(err)
	s.Equal(client.ID, retrievedClient.ID)
}
