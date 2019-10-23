package main

import (
	"clouditor"
	"clouditor/discovery"
	aws3 "clouditor/discovery/aws"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/fatih/structs"

	"github.com/Azure/go-autorest/autorest/azure/auth"
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

	account := aws3.Account{}
	if err := db.GetAccountById("AWS", &account); err == nil {
		fmt.Printf("Account: %+v", account)
	} else {
		fmt.Printf("Error: %+v", err)
	}

	bus := clouditor.NewEventBus()

	bus.Subscribe("discovered", ch1)
	bus.Subscribe("updated", ch2)

	go DoScan(bus)
	go AzureStuff()

	for {
		select {
		case e := <-ch1:
			fmt.Printf("Event: %+v\n", e)
		case e := <-ch2:
			fmt.Printf("Event: %+v\n", e)
		}
	}
}

func AzureStuff() {
	// create a VirtualNetworks client
	vnetClient := network.NewVirtualNetworksClient("<subscriptionID>")

	// create an authorizer from env vars or Azure Managed Service Idenity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		vnetClient.Authorizer = authorizer
	}
}

func DoScan(bus *clouditor.EventBus) {
	scanner := aws3.LambdaScanner{}
	scanner.Init(db)

	for {
		functions, _ := scanner.List()

		for _, f := range functions {
			f2 := aws3.LambdaFunction(*f)
			bus.Publish("discovered", discovery.Asset{Name: f2.Name(), Description: aws.StringValue(f.Description)})

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
