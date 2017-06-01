package services

import (
	"errors"
	"encoding/json"
	"fmt"
	"time"
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

	fmt.Println("Previous Customer DOB  : ",customer.DOB)
	 t, _ := time.Parse("01-02-2006", customer.DOB)
	 customer.DOB = t.Format("01-02-2006")
	 fmt.Println("Customer DOB  : ",customer.DOB)

	customerBytes, err := json.Marshal(customer)
	if err != nil {
		fmt.Println("Failed to Marshal customer  ")
	}

	//err = stub.PutState(customer.CustomerID, []byte(customerJSON))
	err = stub.PutState(customer.CustomerID, customerBytes)
	if err != nil {
		fmt.Println("Failed to create customer ")
	}
	fmt.Println("Created Cusotmer Contract with Key : "+ customer.CustomerID)
	fmt.Println("In initialize.CreateCustomer end ")
	return nil,nil
}

func CreateMemberRelation(subscriberID string,memberRelationJSON string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateMemberRelation start ")

	var memberRelation []data.MemberRelation
	err := json.Unmarshal([]byte(memberRelationJSON), &memberRelation)
	if err != nil {
		fmt.Println("Failed to unmarshal memberRelationJSON ")
	}
	_,err = CreateMemRelation(subscriberID,memberRelation,stub)
	if (err != nil){
		//fmt.Println("Error creating Member Relation with Key :  "+key)
		return nil, errors.New("Error creating Member Relation for subscriber: "+subscriberID)
	}

	fmt.Println("In initialize.CreateMemberRelation end ")
	return nil,nil
}

func CreateMemRelation(subscriberID string, memberRelation []data.MemberRelation, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateMemRelation start ")

	memberRelationBytes, err := json.Marshal(memberRelation)
	if err != nil {
		fmt.Println("Failed to Marshal Member Relation ")
	}
	key := data.CMemberPrefix+"_"+subscriberID
	err = stub.PutState(key, memberRelationBytes)
	if (err != nil){
		fmt.Println("Error creating Member Relation with Key :  "+key)
		return nil, errors.New("Error creating Member Relation for subscriber: "+subscriberID)
	}
	fmt.Println("Created Member Relation with Key :  "+key)

	fmt.Println("In initialize.CreateMemRelation end ")
	return nil,nil
}

func CreateContractDefinition(contractDefinition data.ContractDefinition, stub shim.ChaincodeStubInterface) ([]byte, error) {
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

func CreateContract(contractJSON string,iDedLimit float64,iOOPLimit float64,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateContract start ")
	var contract data.ContractDefinition
	err := json.Unmarshal([]byte(contractJSON), &contract)
	if err != nil {
		fmt.Println("Failed to unmarshal Contract ")
		return nil,errors.New("Failed to unmarshal Contract ")
	}
	_,err = CreateContractDefinition(contract,stub)
	if err != nil {
		fmt.Println("Failed to unmarshal  customer ")
	}
	AddContractToContractList(contract.ContractID,stub)
	AddMemberToContract(contract.ContractID,contract.SubscriberID,iDedLimit,iOOPLimit,stub)

	if (contract.ContractType == "F"){
		var customerFBalance data.CustomerFamilyBalance
		customerFBalance.ContractID = contract.ContractID
		customerFBalance.F_Ded_Balance = 0
		customerFBalance.F_OOP_Balance = 0
		_,err = CreateCustomerFamilyBalance(customerFBalance, stub)
		if err !=nil {
			fmt.Println("Error creating customerFBalance ")
		}
	}
	fmt.Println("In initialize.CreateContract end ")
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
	if (err != nil){
		fmt.Printf("Failed create/update Contract List ")
		return nil,errors.New("Failed to create/update Contract List ")
	}
	fmt.Println("Added ContractsList  with Key :  "+key)
	fmt.Println("In initialize.CreateContractList end ")

	return nil, nil
}

func AddContractToContractList(contractID string,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.AddContractToContractList start ")
 	var contractList []string
	contractList,err := query.GetContractsList(stub)
	if err !=nil {
		fmt.Println("No Contracts found. " )
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

func AddMemberToContract(contractID string, memberID string,iDedLimit float64, iOOPLimit float64,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.AddMemberToContract start ")
	contractDefinition,err := query.GetContract(contractID,stub)
	if err != nil {
		fmt.Println("Error receiving  the Customer Contract")
		return nil, errors.New("Error receiving  the Customer Contract")
	}
	var membersOfContract data.ContractMembers
	var memberIDs []string
	membersOfContract,err = query.GetMembersOfContract(contractID,stub)
	if err !=nil {
		fmt.Println("Error retrieving Members of Contract  " , contractID )
	}
 	memberAdded := false
	memberIDs = membersOfContract.MemberIDs
	if (len(memberIDs)  == 0){
		membersOfContract.ContractID=contractID
		memberIDs = append(memberIDs, memberID)
		memberAdded = true
		fmt.Printf("1st Member %v is added to the contract %v  \n" , memberID,contractID )
	}else if (len(memberIDs) > 0 ) {
		if (contractDefinition.ContractType == "I"){
			fmt.Printf("Individual policy %v already has member %v \n", contractID,memberID )
			 return nil, errors.New("Individual policy "+contractID+" already has member."+memberID)
		}else {
			for _,memID := range memberIDs {
			 if (memID == memberID){
				 fmt.Printf("Member %v already added to the contract %v  " , memberID,contractID )
				 return nil, errors.New("Member "+memberID +" already added to the contract " + contractID )
			 }
		 }
		 memberIDs = append(memberIDs, memberID)
		 memberAdded = true
		 fmt.Printf("Member %v is added to the contract %v  " , memberID,contractID )
		}

	}

	if (memberAdded){
		membersOfContract.MemberIDs = memberIDs
		CreateMembersOfContract(contractID,membersOfContract,stub)
		var customerLimit data.CustomerLimit
		customerLimit.MemberID = memberID
		customerLimit.ContractID = contractID
		customerLimit.I_Ded_Limit = iDedLimit
		customerLimit.I_OOP_Limit = iOOPLimit
		_,err = CreateCustomerLimit(customerLimit,stub )
		if err !=nil {
			fmt.Println("Error creating customer Limits ")
		}

		var customerBalance data.CustomerBalance
		customerBalance.MemberID = memberID
		customerBalance.ContractID = contractID
		customerBalance.I_Ded_Balance = 0
		customerBalance.I_OOP_Balance = 0
		_,err =CreateCustomerBalance(customerBalance,stub )
		if err !=nil {
			fmt.Println("Error creating customer Balances ")
		}
	}
	fmt.Println("Members :  ",membersOfContract)
	fmt.Println("In initialize.AddMemberToContract end ")
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

func CreateCarrierBalances(contMemCarrierBalance data.ContMemCarrierBalance , stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.CreateCarrierBalances start ")
	contMemCarrierBalanceBytes, err := json.Marshal(contMemCarrierBalance)
	if err != nil {
		fmt.Println("Failed to Marshal contMemCarrierBalance ")
	}
	key := data.ContCarrierBalPrefix+"_"+contMemCarrierBalance.ContractID+"_"+contMemCarrierBalance.MemberID
	err = stub.PutState(key, contMemCarrierBalanceBytes)
	fmt.Println("Created Carrier balances for memebr with Key :   "+key)
	fmt.Println("In initialize.CreateCarrierBalances end  ")
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
	fmt.Println("In initialize.CreateCustomerLimit start ")
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

func RemoveMemberFromContract(contractID string, memberID string,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In initialize.RemoveMemberFromContract start ")
	var membersOfContract data.ContractMembers
	var memberIDs []string
	membersOfContract,err := query.GetMembersOfContract(contractID,stub)
	if err !=nil {
		fmt.Println("Error retrieving Members of Contract  " , contractID )
	}
	memberIDs = membersOfContract.MemberIDs
	if (len(memberIDs) > 0){
		fmt.Printf("Members of the contract %v are : %v \n " , contractID ,memberIDs)
		for index,memID := range memberIDs {
		 if (memID == memberID){
			 fmt.Printf("Member %v found in the contract %v \n " , memberID,contractID )
			 newmemberIDs:= append(memberIDs[:index], memberIDs[index+1:]...)
			 membersOfContract.MemberIDs = newmemberIDs
			 fmt.Printf("Member %v remove from the contract %v  \n" , memberID,contractID )
			 _,err = CreateMembersOfContract(contractID,membersOfContract,stub)
			 if (err !=nil ){
				 fmt.Printf("Error updating the members of the contract %v \n " , contractID)
				 return nil, errors.New("Error updating the members of the contract  " + contractID )
			 }
			 fmt.Printf("Members of the contract %v are : %v \n " , contractID ,memberIDs)
			 return nil,nil
		 }
	 }
	}
	fmt.Println("In initialize.RemoveMemberFromContract end ")
	return nil, errors.New("Member "+memberID +" Removed from the contract " + contractID )
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

	fmt.Println("In initialize.ResetCustomeBalances end ")
	return nil,nil
}
///
