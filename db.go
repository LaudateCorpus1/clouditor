package clouditor

import "clouditor/accounts"

type Database interface {
	Connect() error
	GetAccountById(id string, account accounts.Account) error
}
