package aws

import (
	"clouditor"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

type Account struct {
	clouditor.BaseAccount `bson:",inline"`
	AccessKeyID           string `bson:"accessKeyId"`
	SecretAccessKey       string `bson:"secretAccessKey"`
}

func (a *Account) AccountID() string {
	return a.BaseAccount.AccountID
}

func (a *Account) User() string {
	return a.BaseAccount.User
}

func (a *Account) IsAutoDiscovered() bool {
	return a.BaseAccount.IsAutoDiscovered
}

func (a *Account) ResolveCredentials() credentials.Value {
	return credentials.Value{
		AccessKeyID:     a.AccessKeyID,
		SecretAccessKey: a.SecretAccessKey,
		ProviderName:    "DatabaseCredentialsProvider",
	}
}

type DatabaseCredentialsProvider struct {
	db          clouditor.Database
	credentials *credentials.Credentials
	retrieved   bool
}

func NewDatabaseCredentialsProvider(db clouditor.Database) *DatabaseCredentialsProvider {
	return &DatabaseCredentialsProvider{db, credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
	}), false}
}

func (p *DatabaseCredentialsProvider) Retrieve() (credentials.Value, error) {
	var account Account
	if err := p.db.GetAccountById("AWS", &account); err != nil {
		return credentials.Value{}, err
	}

	// forward request to regular credential chain
	if account.IsAutoDiscovered() {
		return p.credentials.Get()
	}

	p.retrieved = true

	return account.ResolveCredentials(), nil
}

func (p *DatabaseCredentialsProvider) IsExpired() bool {
	var account Account
	if err := p.db.GetAccountById("AWS", &account); err != nil {
		return true
	}

	// forward request to regular credential chain
	if account.IsAutoDiscovered() {
		return p.credentials.IsExpired()
	}

	return !p.retrieved
}
