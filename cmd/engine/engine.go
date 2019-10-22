package main

import (
	"clouditor"
	aws2 "clouditor/accounts/aws"
	"clouditor/discovery"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/fatih/structs"
)

var wg sync.WaitGroup

var ch1 = make(chan clouditor.Event)
var ch2 = make(chan clouditor.Event)

var db clouditor.Database

func main() {
	fmt.Printf("Hello!\n")

	db = &clouditor.LegacyDatabase{}
	if err := db.Connect(); err != nil {
		fmt.Printf("Error connecting to database: %s", err)
		return
	}

	account := aws2.Account{}
	if err := db.GetAccountById("AWS", &account); err == nil {
		fmt.Printf("Account: %+v", account)
	} else {
		fmt.Printf("Error: %+v", err)
	}

	bus := clouditor.NewEventBus()

	bus.Subscribe("discovered", ch1)
	bus.Subscribe("updated", ch2)

	go LambdaScanner(bus)

	for {
		select {
		case e := <-ch1:
			fmt.Printf("Event: %+v\n", e)
		case e := <-ch2:
			fmt.Printf("Event: %+v\n", e)
		}
	}
}

func LambdaScanner(bus *clouditor.EventBus) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigDisable,
		SharedConfigFiles: []string{},
	}))

	svc := lambda.New(sess, &aws.Config{Credentials: credentials.NewCredentials(aws2.NewDatabaseCredentialsProvider(db)), Region: aws.String("eu-central-1")})

	for {
		result, err := svc.ListFunctions(nil)
		if err != nil {
			fmt.Println("Cannot list functions")
			os.Exit(0)
		}

		for _, f := range result.Functions {
			f2 := MyFunctionConfiguration(*f)
			bus.Publish("discovered", discovery.Asset{Name: f2.GetName(), Description: aws.StringValue(f.Description)})

			s := structs.New(f2)

			for _, field := range s.Fields() {
				fmt.Printf("field: %+v\n", field.Name())

				if field.IsExported() {
					fmt.Printf("value: %+v\n", field.Value())
				}
			}
		}

		time.Sleep(5 * time.Minute)
	}
}

type MyFunctionConfiguration lambda.FunctionConfiguration

func (f *MyFunctionConfiguration) GetName() string {
	return aws.StringValue(f.FunctionName)
}

func (f *MyFunctionConfiguration) GetID() string {
	return aws.StringValue(f.FunctionArn)
}

func (f *MyFunctionConfiguration) GetType() string {
	return "Serverless"
}
