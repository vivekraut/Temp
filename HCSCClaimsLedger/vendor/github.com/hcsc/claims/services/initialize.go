package services

import (
	"errors"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hcsc/claims/data"
	"github.com/hcsc/claims/query"
)

func CreateCustomer(customerJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateCustomer start ")
	var customer data.Customer
	err := json.Unmarshal([]byte(customerJSON), &customer)
	if err != nil {
		fmt.Println("Failed to unmarshal  customer ")
	}
	fmt.Println("Customer ID : ",customer.CustomerID)
	err = stub.PutState(customer.CustomerID, []byte(customerJSON))
	if err != nil {
		fmt.Println("Failed to create customer ")
	}
	fmt.Println("Created Cusotmer Contract with Key : "+ customer.CustomerID)
	fmt.Println("In initialize.CreateCustomer end ")
	return nil,nil
}

func CreateMember(subscriberID string, customerJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateMember start ")

	var customer data.Customer
	err := json.Unmarshal([]byte(customerJSON), &customer)
	if err != nil {
		fmt.Println("Failed to unmarshal  customer ")
	}
	fmt.Println("Customer ID : ",customer.CustomerID)
	key := subscriberID+"_"+customer.CustomerID

	err = stub.PutState(key, []byte(customerJSON))
	if err != nil {
		fmt.Println("Failed to create customer ")
	}
	fmt.Println("Created Cusotmer Contract with Key : "+ customer.CustomerID)
	fmt.Println("In initialize.CreateMember end ")
	return nil,nil
}
func CreateMemberRelation(subscriberID string,memberRelationJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateMemberRelation start ")

	var memberRelation []data.MemberRelation
	err := json.Unmarshal([]byte(memberRelationJSON), &memberRelation)
	if err != nil {
		fmt.Println("Failed to unmarshal memberRelationJSON ")
	}
	//key := cMemberPrefix+memberRelation.SubscriberID+"_"+memberRelation.MemberID
	key := data.CMemberPrefix+"_"+subscriberID

	err = stub.PutState(key, []byte(memberRelationJSON))
	fmt.Println("Created Member Relation with Key :  "+key)
	fmt.Println("In initialize.CreateMemberRelation end ")
	return nil,nil
}

func CreateContractDefinition(contractDefinition data.ContractDefinition, stub shim.ChaincodeStubInterface) ([]byte, error) {
//func CreateContractDefinition(contractDefinitionJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {

	fmt.Println("In initialize.CreateContractDefinition start ")
	contractDefinitionBytes, err := json.Marshal(contractDefinition)
	if err != nil {
		fmt.Println("Failed to Marshal contract definition ")
	}
	key :=data.ContractPrefix+"_"+contractDefinition.ContractID
	err = stub.PutState(key, contractDefinitionBytes)

	if err != nil {
		fmt.Println("Failed to create contract with key  "+key)
	}

	fmt.Println("In initialize.CreateContractDefinition end ")
	return nil,nil
}

func CreateContractList(contractList []string,stub shim.ChaincodeStubInterface)([]byte, error) {
	fmt.Println("In initialize.CreateContractList start ")

	contractListBytes, err := json.Marshal(contractList)
	if err != nil {
		fmt.Println("Failed to Marshal contractList ")
	}
	key := data.AllContractPrefix
	err = stub.PutState(key,contractListBytes)
	fmt.Println("Added ContractsList  with Key :  "+key)
	fmt.Println("In initialize.CreateContractList end ")

	return nil, nil
}


func AddContractToContractList(contractID string,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.AddContractToContractList start ")

 	var contractList []string
	contractList,err := query.GetContractsList(stub)
	if err !=nil {
		fmt.Println("Error retrieving Contract List " )
		//return nil, errors.New("Error retrieving Members of Contract " + contractID)
	}
	if (len(contractList) > 0 ){
		for _,cID := range contractList {
			if (cID == contractID){
				fmt.Printf("Contract %v already added to the contractList   " , contractID)
				return nil, errors.New("Contract "+contractID +" already added to the contractsList "  )
			}
		}
		contractList = append(contractList, contractID)
		CreateContractList(contractList,stub)
		fmt.Printf("Contract %v is added to the contractList " ,contractID )

	}else {
		contractList = append(contractList, contractID)
		CreateContractList(contractList,stub)
		fmt.Printf("1st Contract %v is added to the contractList  \n" , contractID )
	}
	fmt.Println("ContractList :  ",contractList)
	fmt.Println("In initialize.AddContractToContractList end ")
	return nil,nil
}

func AddMembersToContract(contractID string, memberID string,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.AddMembersToContract start ")

 	var membersOfContract data.ContractMembers
	membersOfContract,err := query.GetMembersOfContract(contractID,stub)
	if err !=nil {
		fmt.Println("Error retrieving Members of Contract  " , contractID )
		//return nil, errors.New("Error retrieving Members of Contract " + contractID)
	}

	memberIDs := membersOfContract.MemberIDs
	if (len(memberIDs) > 0 ){
		for _,memID := range memberIDs {
			if (memID == memberID){
				fmt.Printf("Member %v already added to the contract %v  " , memberID,contractID )
				return nil, errors.New("Member "+memberID +" already added to the contract " + contractID )
			}
		}
		memberIDs = append(memberIDs, memberID)
		membersOfContract.MemberIDs = memberIDs
		CreateMembersOfContract(contractID,membersOfContract,stub)
		fmt.Printf("Member %v is added to the contract %v  " , memberID,contractID )

	}else {
		membersOfContract.ContractID=contractID
		var memberIDs []string
		memberIDs = append(memberIDs, memberID)
		membersOfContract.MemberIDs = memberIDs
		CreateMembersOfContract(contractID,membersOfContract,stub)
		fmt.Printf("1st Member %v is added to the contract %v  \n" , memberID,contractID )
	}
	fmt.Println("Members :  ",membersOfContract)
	fmt.Println("In initialize.AddMembersToContract end ")
	return nil,nil
}

func CreateMembersOfContract(contractID string, membersOfContract data.ContractMembers,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateMembersOfContract start ")

	membersOfContractBytes, err := json.Marshal(membersOfContract)
	if err != nil {
		fmt.Println("Failed to Marshal membersOfContract ")
	}
	key := data.ContractMemsPrefix+"_"+contractID

	err = stub.PutState(key,membersOfContractBytes)
	fmt.Println("Added Members to contract with Key :  "+key)

	fmt.Println("In initialize.CreateMembersOfContract end ")
	return nil,nil
}

func CreateCustomerBalance(customerBalance data.CustomerBalance , stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateCustomerBalance start ")
	customerBalanceBytes, err := json.Marshal(customerBalance)
	if err != nil {
		fmt.Println("Failed to Marshal customerBalance ")
	}
	key := data.BalancePrefix+customerBalance.ContractID+"_"+customerBalance.MemberID
	err = stub.PutState(key, customerBalanceBytes)
	fmt.Println("Created balance with Key :   "+key)
	fmt.Println("In initialize.CreateCustomerBalance end  ")
	return nil,nil
}

func CreateCustomerFamilyBalance(customerFBalance data.CustomerFamilyBalance , stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateCustomerBalance start ")
	customerFBalanceBytes, err := json.Marshal(customerFBalance)
	if err != nil {
		fmt.Println("Failed to Marshal customerFBalance ")
	}
	key := data.BalancePrefix+"_"+customerFBalance.ContractID
	err = stub.PutState(key, customerFBalanceBytes)
	fmt.Println("Created balance with Key :   "+key)
	fmt.Println("In initialize.CreateCustomerFamilyBalance end  ")
	return nil,nil
}

func UpdateCustomerAndFamilyBalances(allBalance data.CustomerAllBalance , stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.UpdateCustomerAndFamilyBalances start ")

	var customerBalance data.CustomerBalance
	customerBalance.ContractID = allBalance.ContractID
	customerBalance.MemberID = allBalance.MemberID
	customerBalance.I_Ded_Balance = allBalance.I_Ded_Balance
	customerBalance.I_OOP_Balance = allBalance.I_OOP_Balance
	CreateCustomerBalance(customerBalance,stub)

	var customerFBalance data.CustomerFamilyBalance
	customerFBalance.ContractID = allBalance.ContractID
	customerFBalance.F_Ded_Balance = allBalance.F_Ded_Balance
	customerFBalance.F_OOP_Balance = allBalance.F_OOP_Balance

	fmt.Println("In initialize.UpdateCustomerAndFamilyBalances end ")
	CreateCustomerFamilyBalance(customerFBalance,stub)

	return nil,nil
}

func CreateCustomerLimit(customerLimit data.CustomerLimit, stub shim.ChaincodeStubInterface) ([]byte, error) {
//func CreateCustomerLimit(customerLimitJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateCustomerLimit start ")
	// err := json.Unmarshal([]byte(customerLimitJSON), &customerLimit)
	// if err != nil {
	// 	fmt.Println("Failed to unmarshal customerLimitJSON ")
	// }
	key := data.LimitPrefix+customerLimit.ContractID+"_"+customerLimit.MemberID
	customerLimitBytes, err := json.Marshal(customerLimit)
	if err != nil {
		fmt.Println("Failed to Marshal customerLimit ")
	}
	err = stub.PutState(key, customerLimitBytes)
	if err != nil {
		fmt.Println("Failed to create  customer Limit ")
	}

	fmt.Println("In initialize.CreateCustomerLimit end ")

	return nil,nil
}

func CreateClaims(claim data.Claims, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateClaims start ")
	claimBytes, err := json.Marshal(claim)
	if err != nil {
		fmt.Println("Failed to Marshal claims ")
	}
	key := data.ClaimPrefix+claim.ContractID+"_"+claim.MemberID+"_"+claim.ClaimID
	err = stub.PutState(key, claimBytes)
	fmt.Println("Created Claim  with Key :   "+key)
	fmt.Println("In initialize.CreateClaims end ")
	return nil,nil
}


//Reset methods
func ResetCustomeBalances(args []string,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.ResetCustomeBalances  start ")

	contractID := args[0]

	contractDefinition,err := query.GetContract(contractID,stub)
	if err != nil {
		fmt.Println("Error receiving  the Customer Contract")
		return nil, errors.New("Error receiving  the Customer Contract")
	}
	fmt.Println("Contract :", contractDefinition)


	membersOfContract, err := query.GetMembersOfContract(contractID,stub)
	if (err != nil ){
		fmt.Println("No Members Found. ")
	}
	memberIDs := membersOfContract.MemberIDs
	fmt.Printf("Members of Contract %v is %v \n", contractID,membersOfContract)
	for _,memID := range memberIDs{
		fmt.Println("Member ID :  ",memID)
		var customerBalance data.CustomerBalance
		customerBalance.MemberID=memID
		customerBalance.ContractID=contractID
		customerBalance.I_Ded_Balance=0
		customerBalance.I_OOP_Balance=0

		_, err = CreateCustomerBalance(customerBalance,stub)
		if (err != nil){
			fmt.Println("Error resetting balance for member:  ",memID)
		}
	}
	if (contractDefinition.ContractType == "F"){
		var customerFBalance data.CustomerFamilyBalance
		customerFBalance.ContractID = contractID
		customerFBalance.F_Ded_Balance = 0
		customerFBalance.F_OOP_Balance = 0

		CreateCustomerFamilyBalance(customerFBalance, stub)
	}
	// customerBalanceJSON :=	`{"MemberID":"CU1001", "ContractID":"CC1001","I_Ded_Balance":0,"I_OOP_Balance":0,"F_Ded_Balance":0,"F_OOP_Balance":0}`
	// var customerBalance data.CustomerBalance
	// err := json.Unmarshal([]byte(customerBalanceJSON), &customerBalance)
	// if err !=nil {
	// 	fmt.Println("Error unmarshalling customerBalance ")
	// }
	// _, err = CreateCustomerBalance(customerBalance,stub)
	//
	// if err != nil {
	// 	fmt.Println("Updating Customer Balance " )
	// }
	fmt.Println("In initialize.ResetCustomeBalances end ")

	return nil,nil

}


func addToCustomerContracts(customerID string, contractID string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.addToCustomerContracts start ")

	key :=data.ContractPrefix+"_"+customerID
	customerContracts, err := query.GetContracts(key,stub)
	if  err != nil {
		fmt.Println("Error Getting Customer Contracts")
	}
	customerContracts = append(customerContracts, contractID)
	tempBytes, _ := json.Marshal(&customerContracts)
	err = stub.PutState(key,tempBytes)

	fmt.Println("Added contract to contract list")
	fmt.Println("In initialize.addToCustomerContracts end ")
	return nil,err
}

func AddClaimIDToMemContClaims(customerID string, contractID string, claimID string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.AddClaimIDToMemContClaims start ")

	key := data.ContractClaimsPrefix+"_"+contractID+"_"+customerID
	contractClaims, err := query.GetMemberClaimIDsOfContract(customerID,contractID,stub)
	if  err != nil {
		fmt.Printf("Error Getting Claims of the  Member %v of the  Contract %v \n ", customerID,contractID)
	}
	contractClaims= append(contractClaims, claimID)
	tempBytes, _ := json.Marshal(&contractClaims)
	err = stub.PutState(key,tempBytes)

	fmt.Println("Added claim to contractClaims list ",contractClaims)
	fmt.Println("In initialize.AddClaimIDToMemContClaims end ")
	return nil,err
}

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
					fmt.Println("Txn Found at index %v Remove and update with latest status" ,index )
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

func updateLastIndTxn(lastTxn data.LastTxnOfContract, stub shim.ChaincodeStubInterface)([]byte, error){
	fmt.Println("In services.updateLastIndTxn start ")
	key := data.LastTxnsIndPrefix+"_"+lastTxn.ContractID
	fmt.Println("Last ind txn to be updated : ",lastTxn )
	tempBytes, _ := json.Marshal(&lastTxn)
	err := stub.PutState(key,tempBytes)
	if (err != nil) {
		fmt.Println("Error updating last txn of individual for contract ",err )
	}
	fmt.Println("In services.updateLastIndTxn end ")

	return nil, nil
}

func UpdateLastFamMemTxn(lastTxn data.LastTxnOfContract, stub shim.ChaincodeStubInterface)([]byte, error){
	fmt.Println("In services.updateLastFamMemTxn start ")

	var lastTxnOfFamily []data.LastTxnOfContract
	lastTxnOfFamily, err := query.GetLastFamilyTxns(lastTxn.ContractID,stub)

	if  err != nil {
		fmt.Println("Error Getting Last Txns of Individuals  : " ,err)
	}
	fmt.Println("Last Ind txns retrieved from ledger : ", lastTxnOfFamily)
	memberFound := false
	if len(lastTxnOfFamily) > 0 {
		for index, txn := range lastTxnOfFamily {
			if (txn.MemberID == lastTxn.MemberID){
				fmt.Println("Updating last ind txn of member ",lastTxn.MemberID)
				lastTxnOfFamily[index] = lastTxn
				memberFound = true
				break
			}
		}
	}else {
		memberFound = true
		lastTxnOfFamily = append(lastTxnOfFamily, lastTxn)
	}

	if (!memberFound) {
		lastTxnOfFamily = append(lastTxnOfFamily, lastTxn)
	}
	key := data.LastTxnsFamPrefix+"_"+lastTxn.ContractID
	tempBytes, _ := json.Marshal(&lastTxnOfFamily)
	err = stub.PutState(key,tempBytes)
	fmt.Println("In services.updateLastIndTxn end ")

	return nil, nil
}

func InitializeCustomerContract(args []string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.InitializeCustomerContract start ")
	//1. Create Customer
	customerJSON :=`{"CustomerID":"CU1001", "FirstName":"Lisa","Johnson":"Kenny","Gender":"F","DOB":"01-01-1960"}`
	CreateCustomer(customerJSON, stub )
	//2. Create Member of Customer Family
	customerJSON =`{"CustomerID":"CU1001_01", "FirstName":"Alex","LastName":"Johnson","Gender":"M","DOB":"01-01-1960"}`
	CreateCustomer(customerJSON, stub )
	customerJSON =`{"CustomerID":"CU1001_01", "FirstName":"Sanjay","LastName":"Johnson","Gender":"M","DOB":"01-01-1985"}`
	CreateCustomer(customerJSON, stub )
	customerJSON =`{"CustomerID":"CU1001_01", "FirstName":"Roberto","LastName":"Johnson","Gender":"M","DOB":"01-01-1987"}`
	CreateCustomer(customerJSON, stub )


	//3. Create Member Relation
	memberRelationJSON :=`[{"SubscriberID":"CU1001", "MemberID":"CU1001","Relationship":"Self"},
	{"SubscriberID":"CU1001", "MemberID":"CU1001_01","Relationship":"Spouse"},
	{"SubscriberID":"CU1001", "MemberID":"CU1001_02","Relationship":"Dependent"}
	{"SubscriberID":"CU1001", "MemberID":"CU1001_03","Relationship":"Dependent"}
	]`
	CreateMemberRelation("CU1001",memberRelationJSON,stub )

	//4. Create Contract Definition
	contractDefinitionJSON :=`{"ContractID":"CC1001", "SubscriberID":"CU1001","ContractStartDate":"10-10-2016",
		"ContractEndDate":"11-12-2016","ContractType":"I","F_Ded_Limit":300,"F_OOP_Limit":500,
		"Copay_feed_OOP":true,"Ded_feed_OOP":true
		}`
	//CreateContractDefinition(contractDefinitionJSON,stub )
	var contractDefinition data.ContractDefinition
	err := json.Unmarshal([]byte(contractDefinitionJSON), &contractDefinition)
	if err != nil {
		fmt.Println("Failed to unmarshal contractDefinitionJSON ")
	}
	CreateContractDefinition(contractDefinition,stub )

	AddContractToContractList("CC1001",stub)

	//5. Set Contract Limits

	AddMembersToContract("CC1001","CU1001",stub)
	AddMembersToContract("CC1001","CU1001_01",stub)
	AddMembersToContract("CC1001","CU1001_02",stub)

	customerLimitJSON :=	`{"MemberID":"CU1001", "ContractID":"CC1001","I_Ded_Limit":200,"I_OOP_Limit":300}`
//	CreateCustomerLimit(customerLimitJSON,stub )
	var customerLimit data.CustomerLimit
	err = json.Unmarshal([]byte(customerLimitJSON), &customerLimit)
	if err !=nil {

		fmt.Println("Error creating customer Limits ")
	}
	CreateCustomerLimit(customerLimit,stub )

	customerLimitJSON =	`{"MemberID":"CU1001_01", "ContractID":"CC1001","I_Ded_Limit":200,"I_OOP_Limit":300}`
	//CreateCustomerLimit(customerLimitJSON,stub )
	err = json.Unmarshal([]byte(customerLimitJSON), &customerLimit)
	if err !=nil {

		fmt.Println("Error creating customer Limits ")
	}
	CreateCustomerLimit(customerLimit,stub )

	customerLimitJSON =	`{"MemberID":"CU1001_02", "ContractID":"CC1001","I_Ded_Limit":200,"I_OOP_Limit":300}`
	//CreateCustomerLimit(customerLimitJSON,stub )
	err = json.Unmarshal([]byte(customerLimitJSON), &customerLimit)
	if err !=nil {

		fmt.Println("Error creating customer Limits ")
	}
	CreateCustomerLimit(customerLimit,stub )
	//6. Set Customer Balances
	customerBalanceJSON :=	`{"MemberID":"CU1001", "ContractID":"CC1001","I_Ded_Banalce":0,"I_OOP_Balance":0}`
	var customerBalance data.CustomerBalance
	err = json.Unmarshal([]byte(customerBalanceJSON), &customerBalance)
	if err !=nil {

		fmt.Println("Error creating customerBalance ")
	}
	CreateCustomerBalance(customerBalance,stub )

	customerBalanceJSON =	`{"MemberID":"CU1001_01", "ContractID":"CC1001","I_Ded_Banalce":0,"I_OOP_Balance":0}`
	err = json.Unmarshal([]byte(customerBalanceJSON), &customerBalance)
	if err !=nil {

		fmt.Println("Error creating customerBalance ")
	}
	CreateCustomerBalance(customerBalance,stub )
	customerBalanceJSON =	`{"MemberID":"CU1001_02", "ContractID":"CC1001","I_Ded_Banalce":0,"I_OOP_Balance":0}`
	err = json.Unmarshal([]byte(customerBalanceJSON), &customerBalance)
	if err !=nil {

		fmt.Println("Error creating customerBalance ")
	}
	CreateCustomerBalance(customerBalance,stub )

	//6. Set Customer Balances
	customerFBalanceJSON :=	`{"ContractID":"CC1001","F_Ded_Balance":0,"F_OOP_Balance":0}`
	var customerFBalance data.CustomerFamilyBalance
	err = json.Unmarshal([]byte(customerFBalanceJSON), &customerFBalance)
	if err !=nil {

		fmt.Println("Error creating customerFBalance ")
	}

	CreateCustomerFamilyBalance(customerFBalance, stub)
//-------------

	//1. Create Customer for P001 & P002
	customerJSON =`{"CustomerID":"112200", "FirstName":"Lisa","Johnson":"Kenny","Gender":"F","DOB":"01-01-1960"}`
	CreateCustomer(customerJSON, stub )
	//2. Create Member of Customer Family
	customerJSON =`{"CustomerID":"112200_01", "FirstName":"Alex","LastName":"Johnson","Gender":"M","DOB":"01-01-1960"}`
	CreateCustomer(customerJSON, stub )
	customerJSON =`{"CustomerID":"112200_02", "FirstName":"Sanjay","LastName":"Johnson","Gender":"M","DOB":"01-01-1985"}`
	CreateCustomer(customerJSON, stub )
	customerJSON =`{"CustomerID":"112200_03", "FirstName":"Roberto","LastName":"Johnson","Gender":"M","DOB":"01-01-1987"}`
	CreateCustomer(customerJSON, stub )

	//3. Create Member Relation
	memberRelationJSON =`[
	{"SubscriberID":"112200", "MemberID":"112200","Relationship":"Self"},
	{"SubscriberID":"112200", "MemberID":"112200_01","Relationship":"Spouse"},
	{"SubscriberID":"112200", "MemberID":"112200_02","Relationship":"Dependent"}
	]`
	CreateMemberRelation("112200",memberRelationJSON,stub )

	_,err = InitializeCustomerContract_P001(stub)
	_,err  =InitializeCustomerContract_P002(stub)

	_,err  =InitializeCustomerContract_P003(stub)

	// _,err = InitializeIndCustomerContract("P003","10034",stub)
	//1. Create Customer for P003

	fmt.Println("In initialize.InitializeCustomerContract end.")
	return nil,nil
}

var contractDefinition data.ContractDefinition
var customerLimit data.CustomerLimit
var customerBalance data.CustomerBalance

func InitializeIndCustomerContract(contractID string, subscriberID string,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.InitializeIndCustomerContract  start ")
	InitializeCustomerContract_P003(stub);
	CreateContractDefinition(contractDefinition,stub )
	AddContractToContractList(contractDefinition.ContractID,stub)
	AddMembersToContract(contractDefinition.ContractID,contractDefinition.SubscriberID,stub)
	CreateCustomerLimit(customerLimit,stub )
	CreateCustomerBalance(customerBalance,stub )
	fmt.Println("In initialize.InitializeIndCustomerContract  end ")

	return nil,nil
}

// func InitializeCustomerContract_P003(stub shim.ChaincodeStubInterface) ([]byte, error) {
//
// 	fmt.Println("In initialize.InitializeCustomerContract_P003  start ")
// 	contractDefinitionJSON :=`{"ContractID":"P003", "SubscriberID":"10034","ContractStartDate":"10-10-2016",
// 	  "ContractEndDate":"11-12-2017","ContractType":"I","F_Ded_Limit":0,"F_OOP_Limit":0,
// 	  "Copay_feed_OOP":true,"Ded_feed_OOP":true
// 	  }`
//
// 	err := json.Unmarshal([]byte(contractDefinitionJSON), &contractDefinition)
// 	if err != nil {
// 		fmt.Println("Failed to unmarshal contractDefinitionJSON ")
// 	}
// 	customerJSON :=`{"CustomerID":"10034", "FirstName":"Lisa.I","Johnson":"Kenny","Gender":"F","DOB":"01-01-1960"}`
// 	CreateCustomer(customerJSON, stub )
//
// 	memberRelationJSON :=`[
// 	{"SubscriberID":"10034", "MemberID":"10034","Relationship":"Self"}
// 	]`
// 	CreateMemberRelation("10034",memberRelationJSON,stub )
//
// 	customerLimitJSON :=	`{"MemberID":"10034", "ContractID":"P003","I_Ded_Limit":200,"I_OOP_Limit":300}`
// 	err = json.Unmarshal([]byte(customerLimitJSON), &customerLimit)
// 	if err !=nil {
// 		fmt.Println("Error creating customer Limits ")
// 	}
// 	//6. Set Customer Balances
// 	customerBalanceJSON :=	`{"MemberID":"10034", "ContractID":"P003","I_Ded_Balance":0,"I_OOP_Balance":0}`
// 	err = json.Unmarshal([]byte(customerBalanceJSON), &customerBalance)
// 	if err !=nil {
// 	  fmt.Println("Error creating customerBalance ")
// 	}
// 	fmt.Println("In initialize.InitializeCustomerContract_P003  end ")
//
// 	return nil,nil
// }

func InitializeCustomerContract_P001(stub shim.ChaincodeStubInterface) ([]byte, error) {

	fmt.Println("In initialize.InitializeCustomerContract P001 start ")
	contractDefinitionJSON :=`{"ContractID":"P001", "SubscriberID":"112200","ContractStartDate":"10-10-2016",
	  "ContractEndDate":"11-12-2017","ContractType":"I","F_Ded_Limit":0,"F_OOP_Limit":0,
	  "Copay_feed_OOP":true,"Ded_feed_OOP":true
	  }`

	// Create Contract Definition
	//CreateContractDefinition(contractDefinitionJSON,stub )
	var contractDefinition data.ContractDefinition
	err := json.Unmarshal([]byte(contractDefinitionJSON), &contractDefinition)
	if err != nil {
		fmt.Println("Failed to unmarshal contractDefinitionJSON ")
	}
	CreateContractDefinition(contractDefinition,stub )
	AddContractToContractList(contractDefinition.ContractID,stub)
	//5. Set Contract Limits
	AddMembersToContract(contractDefinition.ContractID,contractDefinition.SubscriberID,stub)

	customerLimitJSON :=	`{"MemberID":"112200", "ContractID":"P001","I_Ded_Limit":200,"I_OOP_Limit":300}`
	var customerLimit data.CustomerLimit
	err = json.Unmarshal([]byte(customerLimitJSON), &customerLimit)
	if err !=nil {

		fmt.Println("Error creating customer Limits ")
	}
	CreateCustomerLimit(customerLimit,stub )

	//6. Set Customer Balances
	customerBalanceJSON :=	`{"MemberID":"112200", "ContractID":"P001","I_Ded_Banalce":0,"I_OOP_Balance":0}`
	var customerBalance data.CustomerBalance
	err = json.Unmarshal([]byte(customerBalanceJSON), &customerBalance)
	if err !=nil {

	  fmt.Println("Error creating customerBalance ")
	}
	CreateCustomerBalance(customerBalance,stub )

	fmt.Println("In initialize.InitializeCustomerContract  P001 end ")

	return nil,nil
}

func InitializeCustomerContract_P003(stub shim.ChaincodeStubInterface) ([]byte, error) {

	fmt.Println("In initialize.InitializeCustomerContract_P003  start ")
	contractDefinitionJSON :=`{"ContractID":"P003", "SubscriberID":"10034","ContractStartDate":"10-10-2016",
	  "ContractEndDate":"11-12-2017","ContractType":"I","F_Ded_Limit":0,"F_OOP_Limit":0,
	  "Copay_feed_OOP":true,"Ded_feed_OOP":true
	  }`


	// Create Contract Definition
	//CreateContractDefinition(contractDefinitionJSON,stub )
	var contractDefinition data.ContractDefinition
	err := json.Unmarshal([]byte(contractDefinitionJSON), &contractDefinition)
	if err != nil {
		fmt.Println("Failed to unmarshal contractDefinitionJSON ")
	}
	CreateContractDefinition(contractDefinition,stub )
	AddContractToContractList(contractDefinition.ContractID,stub)

	customerJSON :=`{"CustomerID":"10034", "FirstName":"LisaI","Johnson":"Kenny","Gender":"F","DOB":"01-01-1960"}`
	CreateCustomer(customerJSON, stub )

	memberRelationJSON :=`[
	{"SubscriberID":"10034", "MemberID":"10034","Relationship":"Self"}
	]`
	CreateMemberRelation("10034",memberRelationJSON,stub )
	//5. Set Contract Limits
	AddMembersToContract("P003","10034",stub)

	customerLimitJSON :=	`{"MemberID":"10034", "ContractID":"P003","I_Ded_Limit":100,"I_OOP_Limit":200}`
	var customerLimit data.CustomerLimit
	err = json.Unmarshal([]byte(customerLimitJSON), &customerLimit)
	if err !=nil {

		fmt.Println("Error creating customer Limits ")
	}
	CreateCustomerLimit(customerLimit,stub )

	//6. Set Customer Balances
	customerBalanceJSON :=	`{"MemberID":"10034", "ContractID":"P003","I_Ded_Banalce":0,"I_OOP_Balance":0}`
	var customerBalance data.CustomerBalance
	err = json.Unmarshal([]byte(customerBalanceJSON), &customerBalance)
	if err !=nil {

	  fmt.Println("Error creating customerBalance ")
	}
	CreateCustomerBalance(customerBalance,stub )

	fmt.Println("In initialize.InitializeCustomerContract_P003   end ")

	return nil,nil
}

func InitializeCustomerContract_P002(stub shim.ChaincodeStubInterface) ([]byte, error) {

	fmt.Println("In initialize.InitializeCustomerContract P002 start ")
	contractDefinitionJSON :=`{"ContractID":"P002", "SubscriberID":"112200","ContractStartDate":"10-10-2016",
	  "ContractEndDate":"11-12-2017","ContractType":"F","F_Ded_Limit":300,"F_OOP_Limit":500,
	  "Copay_feed_OOP":true,"Ded_feed_OOP":true
	  }`

	//4. Create Contract Definition
	//CreateContractDefinition(contractDefinitionJSON,stub )
	var contractDefinition data.ContractDefinition
	err := json.Unmarshal([]byte(contractDefinitionJSON), &contractDefinition)
	if err != nil {
		fmt.Println("Failed to unmarshal contractDefinitionJSON ")
	}
	CreateContractDefinition(contractDefinition,stub )
	AddContractToContractList(contractDefinition.ContractID,stub)

	AddMembersToContract(contractDefinition.ContractID,"112200",stub)
	AddMembersToContract(contractDefinition.ContractID,"112200_01",stub)
	AddMembersToContract(contractDefinition.ContractID,"112200_02",stub)
	//5. Set Contract Limits

	customerLimitJSON :=	`{"MemberID":"112200", "ContractID":"P002","I_Ded_Limit":200,"I_OOP_Limit":300}`
//	CreateCustomerLimit(customerLimitJSON,stub )
	var customerLimit data.CustomerLimit
	err = json.Unmarshal([]byte(customerLimitJSON), &customerLimit)
	if err !=nil {

		fmt.Println("Error creating customer Limits ")
	}
	CreateCustomerLimit(customerLimit,stub )

	customerLimitJSON =	`{"MemberID":"112200_01", "ContractID":"P002","I_Ded_Limit":200,"I_OOP_Limit":300}`
	//CreateCustomerLimit(customerLimitJSON,stub )
	err = json.Unmarshal([]byte(customerLimitJSON), &customerLimit)
	if err !=nil {

		fmt.Println("Error creating customer Limits ")
	}
	CreateCustomerLimit(customerLimit,stub )

	customerLimitJSON =	`{"MemberID":"112200_02", "ContractID":"P002","I_Ded_Limit":200,"I_OOP_Limit":300}`
	//CreateCustomerLimit(customerLimitJSON,stub )
	err = json.Unmarshal([]byte(customerLimitJSON), &customerLimit)
	if err !=nil {

		fmt.Println("Error creating customer Limits ")
	}
	CreateCustomerLimit(customerLimit,stub )

	//6. Set Customer Balances
	customerBalanceJSON :=	`{"MemberID":"112200", "ContractID":"P002","I_Ded_Banalce":0,"I_OOP_Balance":0}`
	var customerBalance data.CustomerBalance
	err = json.Unmarshal([]byte(customerBalanceJSON), &customerBalance)
	if err !=nil {

	  fmt.Println("Error creating customerBalance ")
	}
	CreateCustomerBalance(customerBalance,stub )

	customerBalanceJSON =	`{"MemberID":"112200_01", "ContractID":"P002","I_Ded_Banalce":0,"I_OOP_Balance":0}`
	err = json.Unmarshal([]byte(customerBalanceJSON), &customerBalance)
	if err !=nil {

	  fmt.Println("Error creating customerBalance... ")
	}
	CreateCustomerBalance(customerBalance,stub )
	customerBalanceJSON =	`{"MemberID":"112200_02", "ContractID":"P002","I_Ded_Banalce":0,"I_OOP_Balance":0}`
	err = json.Unmarshal([]byte(customerBalanceJSON), &customerBalance)
	if err !=nil {

	  fmt.Println("Error creating customerBalance ")
	}
	CreateCustomerBalance(customerBalance,stub )

	//6. Set Customer Balances
	customerFBalanceJSON :=	`{"ContractID":"P002","F_Ded_Balance":0,"F_OOP_Balance":0}`
	var customerFBalance data.CustomerFamilyBalance
	err = json.Unmarshal([]byte(customerFBalanceJSON), &customerFBalance)
	if err !=nil {

	  fmt.Println("Error creating customerFBalance ")
	}

	CreateCustomerFamilyBalance(customerFBalance, stub)
	fmt.Println("In initialize.InitializeCustomerContract  P002 end ")

	return nil,nil
}
