package aws

import (
	"clouditor"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

var sess *session.Session
var cfg aws.Config

func init() {
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigDisable,
		SharedConfigFiles: []string{},
	}))
}

type LambdaScanner struct {
	client *lambda.Lambda
}

type MyInterface interface {
	NewRequest(*request.Operation, interface{}, interface{}) *request.Request
}

type MyFunc func(client.ConfigProvider, ...*aws.Config) *lambda.Lambda

func (l *LambdaScanner) Init(db clouditor.Database) {
	l.client = createClient(db, lambda.New).(*lambda.Lambda)
}

func (l *LambdaScanner) List() (LambdaFunctions, error) {
	result, err := l.client.ListFunctions(nil)
	if err != nil {
		return nil, err
	}

	return result.Functions, nil
}

type LambdaFunction lambda.FunctionConfiguration
type LambdaFunctions []*lambda.FunctionConfiguration

func (f *LambdaFunction) Name() string {
	return aws.StringValue(f.FunctionName)
}

func (f *LambdaFunction) ID() string {
	return aws.StringValue(f.FunctionArn)
}

func (f *LambdaFunction) Type() string {
	return "Serverless"
}
