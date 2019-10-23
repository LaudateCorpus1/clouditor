package clouditor

// ServiceAccount is anything which contains credentials to connect to a particular Cloud or Cluster
type ServiceAccount interface {
	AccountID() string
	User() string
	IsAutoDiscovered() bool
	/*Validate()*/
}

type BaseAccount struct {
	AccountID        string `bson:"accountId"`
	User             string `bson:"user"`
	IsAutoDiscovered bool   `bson:"autoDiscovered"`
}
