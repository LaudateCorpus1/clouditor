package aws

import (
	"clouditor"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func createClient(db clouditor.Database, f MyFunc) MyInterface {
	cfg = aws.Config{Credentials: credentials.NewCredentials(NewDatabaseCredentialsProvider(db)), Region: aws.String("eu-central-1")}

	return f(sess, &cfg)
}
