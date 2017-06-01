package services

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


func InitializeCustomerContract(args []string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.InitializeCustomerContract start ")

	//1. Create Customers in System and their relationships
	_,err := InitializeCustomers(stub)
	if (err != nil){
		fmt.Println("Error creating Customers ")
	}
	//2. Create an Individual Contract
	_,err = InitializeIndCon(stub)
	if (err != nil){
		fmt.Println("Error creating IND contract")
	}
	//3. Create an Family Contract
	_,err  =InitializeFamilyContract(stub)
	if (err != nil){
		fmt.Println("Error creating FAMILY contract")
	}
	fmt.Println("In initialize.InitializeCustomerContract end.")
	return nil,nil
}

func InitializeCustomers(stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.InitializeCustomers  start ")
	//1. Create Customers
	customerJSON :=`{"CustomerID":"114545", "FirstName":"Jimmy","LastName":"Grant","Gender":"M","DOB":"07-20-1988"}`
	CreateCustomer(customerJSON, stub )
	// Create Subscriber 112200 Relationship with members
	memberRelationJSON :=`[
	{"SubscriberID":"114545", "MemberID":"114545","Relationship":"Self"}
	]`
	CreateMemberRelation("114545",memberRelationJSON,stub )

	customerJSON =`{"CustomerID":"110067", "FirstName":"Woody","LastName":"Binns","Gender":"M","DOB":"01-23-1980"}`
	CreateCustomer(customerJSON, stub )
	memberRelationJSON =`[
	{"SubscriberID":"110067", "MemberID":"110067","Relationship":"Self"}
	]`
	CreateMemberRelation("110067",memberRelationJSON,stub )

	customerJSON =`{"CustomerID":"114366", "FirstName":"Carla","LastName":"Simon","Gender":"F","DOB":"08-12-1974"}`
	CreateCustomer(customerJSON, stub )
	memberRelationJSON =`[
	{"SubscriberID":"114366", "MemberID":"114366","Relationship":"Self"}
	]`
	CreateMemberRelation("114366",memberRelationJSON,stub )

	customerJSON =`{"CustomerID":"115033", "FirstName":"Frank","LastName":"Debonair","Gender":"M","DOB":"10-15-1985"}`
	CreateCustomer(customerJSON, stub )
	memberRelationJSON =`[
	{"SubscriberID":"115033", "MemberID":"115033","Relationship":"Self"}
	]`
	CreateMemberRelation("115033",memberRelationJSON,stub )

	customerJSON =`{"CustomerID":"112222", "FirstName":"Debbie","LastName":"Binns","Gender":"F","DOB":"05-15-1983"}`
	CreateCustomer(customerJSON, stub )
	customerJSON =`{"CustomerID":"112233", "FirstName":"Wayne","LastName":"Binns","Gender":"M","DOB":"11-13-2016"}`
	CreateCustomer(customerJSON, stub )
	// Create Subscriber 112222 Relationship with members
	memberRelationJSON =`[
	{"SubscriberID":"112222", "MemberID":"112222","Relationship":"Self"},
	{"SubscriberID":"112222", "MemberID":"112233","Relationship":"Dependent"}
	]`
	CreateMemberRelation("112222",memberRelationJSON,stub )

	fmt.Println("In initialize.InitializeCustomers  end ")
	return nil,nil
}

func InitializeIndCon(stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.InitializeIndCon  start ")
	contractDefinitionJSON :=`{"ContractID":"P004", "SubscriberID":"114545","ContractStartDate":"10-10-2016",
	  "ContractEndDate":"12-30-2017","ContractType":"I","F_Ded_Limit":0,"F_OOP_Limit":0,
	  "Copay_feed_OOP":true,"Ded_feed_OOP":true
	  }`
	_,err := CreateContract(contractDefinitionJSON,300,500,stub)
	if (err != nil ){
		fmt.Println("Error creating Individual Contract")
	}

	contractDefinitionJSON =`{"ContractID":"P005", "SubscriberID":"114366","ContractStartDate":"10-10-2016",
		"ContractEndDate":"12-30-2017","ContractType":"I","F_Ded_Limit":0,"F_OOP_Limit":0,
		"Copay_feed_OOP":true,"Ded_feed_OOP":true
		}`
	_,err = CreateContract(contractDefinitionJSON,300,500,stub)
	if (err != nil ){
		fmt.Println("Error creating Individual Contract")
	}

	contractDefinitionJSON =`{"ContractID":"P006", "SubscriberID":"115033","ContractStartDate":"10-10-2016",
		"ContractEndDate":"12-30-2017","ContractType":"I","F_Ded_Limit":0,"F_OOP_Limit":0,
		"Copay_feed_OOP":true,"Ded_feed_OOP":true
		}`
	_,err = CreateContract(contractDefinitionJSON,300,500,stub)
	if (err != nil ){
		fmt.Println("Error creating Individual Contract")
	}

	contractDefinitionJSON =`{"ContractID":"P007", "SubscriberID":"110067","ContractStartDate":"10-10-2016",
	  "ContractEndDate":"12-30-2017","ContractType":"I","F_Ded_Limit":0,"F_OOP_Limit":0,
	  "Copay_feed_OOP":true,"Ded_feed_OOP":true
	  }`
	_,err = CreateContract(contractDefinitionJSON,300,500,stub)
	if (err != nil ){
		fmt.Println("Error creating Individual Contract")
	}
	fmt.Println("In initialize.InitializeIndCon  P001 end ")
	return nil,nil
}

func InitializeFamilyContract(stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.InitializeFamilyContract start ")
	contractDefinitionJSON :=`{"ContractID":"P101", "SubscriberID":"112222","ContractStartDate":"10-10-2016",
	  "ContractEndDate":"12-30-2017","ContractType":"F","F_Ded_Limit":500,"F_OOP_Limit":1000,
	  "Copay_feed_OOP":true,"Ded_feed_OOP":true
	  }`
	_,err := CreateContract(contractDefinitionJSON, 300,500,stub)
	if (err != nil ){
		fmt.Println("Error creating Family Contract")
	}
	AddMemberToContract("P101","112233",300,500,stub)

	fmt.Println("In initialize.InitializeFamilyContract  end ")
	return nil,nil
}

//
//
//
