package services

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hcsc/claims/data"
	"github.com/hcsc/claims/query"
)

//NOTE: This method is not being used.
func AddToLastTxnsList(lastTxn data.LastTxnOfContract, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.addToLast10Txns start ")
	var lastTxns []data.LastTxnOfContract
	lastTxns, err := query.GetLastTxnsList(stub)
	fmt.Println("No. of last txns : ", len(lastTxns))
	if  err != nil {
		fmt.Println("Error Getting Last Txns  : " ,lastTxn.ContractID)
		lastTxns = append(lastTxns,lastTxn)
	} else {
		//Find if the txn is already exists in the lastTxns list
		if (len(lastTxns) > 0){
			txnFound := false
			for index,txn := range lastTxns{
				if (txn.ContractID == lastTxn.ContractID &&
					txn.SubscriberID == lastTxn.SubscriberID &&
					txn.MemberID == lastTxn.MemberID &&
					txn.ClaimID == lastTxn.ClaimID &&
					txn.Transaction.TransactionID == lastTxn.Transaction.TransactionID &&
					txn.Transaction.AccumType == lastTxn.Transaction.AccumType ){
					fmt.Println("Txn Found at index %v. Remove and update with latest status" ,index )
					txnFound = true
					lastTxns = append(lastTxns[:index],lastTxns[index+1:]...)
					fmt.Println("Adjusted length of last txns : ", len(lastTxns))
					break;
				}
			}
			if (!txnFound){
				if (len(lastTxns) == data.NoOfLastTxns ) {
					fmt.Println("Reached Limit. Hence removing first txn and adding new")
					lastTxns = lastTxns[1:]
				}
			}
		}
		lastTxns = append(lastTxns,lastTxn)
	}
	fmt.Println("No. of last txns : ", len(lastTxns))
	key := data.LastTxnsPrefix
	tempBytes, _ := json.Marshal(&lastTxns)
	err = stub.PutState(key,tempBytes)
	fmt.Println("In initialize.addToLast10Txns end ")
	return nil, nil
}

// func InitializeCustomers(stub shim.ChaincodeStubInterface) ([]byte, error) {
// 	fmt.Println("In initialize.InitializeCustomers  start ")
// 	//1. Create Customers
// 	customerJSON :=`{"CustomerID":"112200", "FirstName":"Lisa","LastName":"Kenny","Gender":"F","DOB":"01-01-1960"}`
// 	CreateCustomer(customerJSON, stub )
// 	customerJSON =`{"CustomerID":"10034", "FirstName":"Marie","LastName":"Kenny","Gender":"F","DOB":"01-01-1960"}`
// 	CreateCustomer(customerJSON, stub )
// 	customerJSON =`{"CustomerID":"10035", "FirstName":"Alex","LastName":"Johnson","Gender":"M","DOB":"01-01-1960"}`
// 	CreateCustomer(customerJSON, stub )
// 	customerJSON =`{"CustomerID":"10036", "FirstName":"Sanjay","LastName":"Johnson","Gender":"M","DOB":"01-01-1985"}`
// 	CreateCustomer(customerJSON, stub )
// 	customerJSON =`{"CustomerID":"10037", "FirstName":"Roberto","LastName":"Johnson","Gender":"M","DOB":"01-01-1987"}`
// 	CreateCustomer(customerJSON, stub )
//
// 	// Create Subscriber 112200 Relationship with members
// 	memberRelationJSON :=`[
// 	{"SubscriberID":"112200", "MemberID":"112200","Relationship":"Self"}
// 	]`
// 	CreateMemberRelation("112200",memberRelationJSON,stub )
//
// 	// Create Subscriber 10034 Relationship with members
// 	memberRelationJSON =`[
// 	{"SubscriberID":"10034", "MemberID":"10034","Relationship":"Self"},
// 	{"SubscriberID":"10034", "MemberID":"10035","Relationship":"Spouse"},
// 	{"SubscriberID":"10034", "MemberID":"10036","Relationship":"Dependent"}
// 	]`
// 	CreateMemberRelation("10034",memberRelationJSON,stub )
//
// 	fmt.Println("In initialize.InitializeCustomers  end ")
// 	return nil,nil
// }
//
// func InitializeIndCon(stub shim.ChaincodeStubInterface) ([]byte, error) {
// 	fmt.Println("In initialize.InitializeIndCon  start ")
// 	contractDefinitionJSON :=`{"ContractID":"P001", "SubscriberID":"112200","ContractStartDate":"10-10-2016",
// 	  "ContractEndDate":"11-12-2017","ContractType":"I","F_Ded_Limit":0,"F_OOP_Limit":0,
// 	  "Copay_feed_OOP":true,"Ded_feed_OOP":true
// 	  }`
// 	_,err := CreateContract(contractDefinitionJSON,300,500,stub)
// 	if (err != nil ){
// 		fmt.Println("Error creating Family Contract")
// 	}
// 	fmt.Println("In initialize.InitializeIndCon  P001 end ")
// 	return nil,nil
// }
//
// func InitializeFamilyContract(stub shim.ChaincodeStubInterface) ([]byte, error) {
// 	fmt.Println("In initialize.InitializeFamilyContract start ")
// 	contractDefinitionJSON :=`{"ContractID":"P101", "SubscriberID":"10034","ContractStartDate":"10-10-2016",
// 	  "ContractEndDate":"11-12-2017","ContractType":"F","F_Ded_Limit":300,"F_OOP_Limit":500,
// 	  "Copay_feed_OOP":true,"Ded_feed_OOP":true
// 	  }`
// 	_,err := CreateContract(contractDefinitionJSON, 300,500,stub)
// 	if (err != nil ){
// 		fmt.Println("Error creating Family Contract")
// 	}
// 	AddMemberToContract("P101","10035",300,500,stub)
// 	AddMemberToContract("P101","10036",300,500,stub)
// 	fmt.Println("In initialize.InitializeFamilyContract  end ")
// 	return nil,nil
// }
