package clouditor

type Database interface {
	Connect() error
	GetAccountById(id string, account ServiceAccount) error
}
