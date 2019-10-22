package main

import (
	"clouditor"
	"clouditor/discovery"
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/fatih/structs"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

var ch1 = make(chan clouditor.Event)
var ch2 = make(chan clouditor.Event)

func main() {
	fmt.Printf("Hello!\n")

	conn, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(conn))

	txn := dgraphClient.NewTxn()
	defer txn.Discard(context.Background())

	q := `query all($a: string) {
		all(func:allofterms(name, $a)) {
		  name
		}
	  }`

	res, err := txn.QueryWithVars(context.Background(), q, map[string]string{"$a": "Star Wars"})
	fmt.Printf("%s\n", res.Json)

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
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := lambda.New(sess, &aws.Config{Region: aws.String("eu-central-1")})

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
