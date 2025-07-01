package gateway

import "github.com/AntonioSabino/fc-ms-wallet/internal/entity"

type TransactionGateway interface {
	Save(transaction *entity.Transaction) error
}
