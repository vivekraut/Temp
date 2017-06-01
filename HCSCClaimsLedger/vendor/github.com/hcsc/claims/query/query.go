package query

import (
	"encoding/json"
  "errors"
	"fmt"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hcsc/claims/data"

)

func Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Printf("In query.Query  function %v with args %v  \n", function, args)
	if function == "getMemberDetailsNHistory" {
		fmt.Println("Invoking GetMemberDetailsNHistory " + function)
		subscriberID := args[0]
		memberID := args[1]
		contractID := args[2]
		claimID := args[3]

		var memNContDetails data.MemberNContractDetails
		memNContDetails,err := GetMemberDetailsNHistory(contractID, subscriberID,memberID,
			claimID,stub)

		if err != nil {
			fmt.Println("Error receiving  the Member History")
			return nil, errors.New("Error receiving  the Member History")
		}
		fmt.Println("All success, returning the Member History and details")
		return json.Marshal(memNContDetails)
	}

	if function == "getMembersOfContract" {
		fmt.Println("Invoking GetMembersOfContract " + function)
		var contractMembers data.ContractMembers
		contractMembers,err := GetMembersOfContract(args[0], stub)
		if err != nil {
			fmt.Println("Error receiving  the Members of Contract")
			return nil, errors.New("Error receiving  Members of Contract")
		}
		fmt.Println("All success, returning Members of Contract")
		return json.Marshal(contractMembers)
	}

	if function == "getCustomerDetails" {
		fmt.Println("Invoking GetCustomerDetails " + function)
		var customer data.Customer
		customer,err := GetCustomerDetails(args[0], stub)
		if err != nil {
			fmt.Println("Error receiving  the Customer")
			return nil, errors.New("Error receiving  the Customer")
		}
		fmt.Println("All success, returning the Customer details")
		return json.Marshal(customer)
	}

	if function == "getCustomerLimits" {
		fmt.Println("Invoking getCustomerLimits " + function)
		customerID := args[0]
		contractID := args[1]
		var customerLimit data.CustomerLimit
		customerLimit,err := GetCustomerLimits(customerID, contractID,stub)
		if err != nil {
			fmt.Println("Error receiving  the Customer Limits")
			return nil, errors.New("Error receiving  the Customer Limits")
		}
		fmt.Println("All success, returning the Customer Limits")
		return json.Marshal(customerLimit)
	}

	if function == "getCustomerContractBalance" {
		fmt.Println("Invoking GetCustomerContractBalance " + function)
		customerID := args[0]
		contractID := args[1]
		var customerAllBalance data.CustomerAllBalance
		customerAllBalance,err := GetAllBalancesOfCustomer(customerID, contractID,stub)
		if err != nil {
			fmt.Println("Error receiving  all balances of Customer ")
			return nil, errors.New("Error receiving  all balances of Customer")
		}
		fmt.Println("All success, returning all balances of Customer ")
		return json.Marshal(customerAllBalance)
	}
	if function == "getCustomerBalance" {
		fmt.Println("Invoking GetCustomerBalance " + function)
		customerID := args[0]
		contractID := args[1]
		var customerBalance data.CustomerBalance
		customerBalance,err := GetCustomerBalance(customerID, contractID,stub)
		if err != nil {
			fmt.Println("Error receiving  the Customer Balances")
			return nil, errors.New("Error receiving  the Customer Balances")
		}
		fmt.Println("All success, returning the Customer Balance")
		return json.Marshal(customerBalance)
	}
	if function == "getCustomerContract" {
		fmt.Println("Invoking GetContract " + function)
		contractID := args[0]
		var contractDefinition data.ContractDefinition
		contractDefinition,err := GetContract(contractID,stub)
		if err != nil {
			fmt.Println("Query: Error receiving  the Customer Contract 5")
			return nil, errors.New("Error receiving  the Customer Contract 6")
		}
		fmt.Println("All success, returning the Customer Contract")
		return json.Marshal(contractDefinition)
	}

	if function == "getClaim" {
		fmt.Println("Invoking getClaim " + function)
		customerID := args[0]
		contractID := args[1]
		claimID := args[2]
		var claim data.Claims
		claim,err := GetClaim(customerID, contractID,claimID,stub)
		if err != nil {
			fmt.Println("Error receiving  the Claim ")
			return nil, errors.New("Error receiving  the Claim")
		}
		fmt.Println("All success, returning the Claim ")
		return json.Marshal(claim)
	}

	if function == "getMemberClaimIDsOfContract" {
		fmt.Println("Invoking GetMemberClaimIDsOfContract " + function)
		customerID := args[0]
		contractID := args[1]

		var claimList []string
		claimList,err := GetMemberClaimIDsOfContract(customerID,contractID,stub)
		if err != nil {
			fmt.Println("Error receiving  the Claim list ")
			return nil, errors.New("Error receiving  the Claim list")
		}
		fmt.Println("All success, returning the Claim List ")
		return json.Marshal(claimList)
	}//GetAllClaimsOfContract
	if function == "getAllClaimsOfContract" {
		fmt.Println("Invoking GetAllClaimsOfContract " + function)
		customerID := args[0]
		contractID := args[1]

		var claimList []data.Claims
		claimList,err := GetAllClaimsOfContract(customerID, contractID,stub)
		if err != nil {
			fmt.Println("Error receiving  the Claim list ")
			return nil, errors.New("Error receiving  the Claim list")
		}
		fmt.Println("All success, returning the Claim List ")
		return json.Marshal(claimList)
	}

	if function == "getLastTxnsList" {
		fmt.Println("Invoking GetLastTxnsList  " + function)

		var lastTxns []data.LastTxnOfContract

		lastTxns,err := GetLastTxnsList(stub)
		if err != nil {
			fmt.Println("Error receiving  the Last Txns of Contract  ")
			return nil, errors.New("Error receiving  the Last Txns of Contract ")
		}
		fmt.Println("All success, returning the last txns List of contract ")
		return json.Marshal(lastTxns)
	}

	if function == "getAllLastTxns" {
		fmt.Println("Invoking GetAllLastTxns  " + function)
		
		var lastTxns []data.LastTxnOfContract
		lastTxns,err := GetAllLastTxns(stub)
		if err != nil {
			fmt.Println("Error receiving  last 50 the txns  ")
			return nil, errors.New("Error receiving  last 50 the txns ")
		}
		fmt.Println("All success, returning the last 50 txns ")
		return json.Marshal(lastTxns)
	}

	if function == "getLastTxnsOfMem" {
		fmt.Println("Invoking GetLastTxnsOfMem  " + function)
		subscriberID := args[0]
		memberID := args[1]
		contractID := args[2]
		var lastTxns []data.LastTxnOfContract

		lastTxns,err := GetLastTxnsOfMem(contractID,subscriberID,memberID,stub)
		if err != nil {
			fmt.Println("Error receiving  the Last Txns of Contract  ")
			return nil, errors.New("Error receiving  the Last Txns of Contract ")
		}
		fmt.Println("All success, returning the last txns List of contract ")
		return json.Marshal(lastTxns)
	}

	if function == "getLastIndTxns" {
		fmt.Println("Invoking GetLastIndTxns  " + function)
		contractID := args[0]
		var lastTxn data.LastTxnOfContract
		lastTxn,err := GetLastIndTxn(contractID,stub)
		if err != nil {
			fmt.Println("Error receiving  the Last Ind Txns of Contract  ")
			return nil, errors.New("Error receiving  the Last Ind Txns of Contract ")
		}
		fmt.Println("All success, returning the last Ind txns List of contract ")
		return json.Marshal(lastTxn)
	}

	//
	return nil, errors.New("Received unknown query function name")
}

func GetCustomerDetails(customerID string, stub shim.ChaincodeStubInterface) (data.Customer, error) {
	fmt.Println("In query.GetCustomerDetails start ")
	var customer data.Customer
	customerBytes, err := stub.GetState(customerID)
	if err != nil {
		fmt.Println("Error retrieving Customer Details " + customerID)
		return customer, errors.New("Error retrieving Customer Details " + customerID)
	}
	err = json.Unmarshal(customerBytes, &customer)
	fmt.Println("Customer   : " , customer);
	fmt.Println("In query.GetCustomerDetails end ")
	return customer,nil
}

func GetCustomerLimits(customerID string,contractID string, stub shim.ChaincodeStubInterface)(data.CustomerLimit, error) {
	fmt.Println("In query.GetCustomerLimits start ")
	var customerLimits data.CustomerLimit
	customerLimitBytes, err := stub.GetState(data.LimitPrefix+contractID+"_"+customerID)
	if err != nil {
		fmt.Println("Error retrieving Customer Limits " + customerID)
		return customerLimits, errors.New("Error retrieving Customer Limits " + customerID)
	}
	err = json.Unmarshal(customerLimitBytes, &customerLimits)
	fmt.Println("CustomerLimits  : " , customerLimits);
	fmt.Println("In query.GetCustomerLimits end ")
	if err != nil {
		fmt.Println("Error unmarshalling Customer Limits " + customerID)
		return customerLimits, errors.New("Error unmarshalling Customer Limits " + customerID)
	}
	return customerLimits, nil
}



func GetAllBalancesOfCustomer(customerID string,contractID string, stub shim.ChaincodeStubInterface)(data.CustomerAllBalance, error) {
	fmt.Println("In query.GetAllBalancesOfCustomer start ")
	var allBalances data.CustomerAllBalance
	var customerBalance data.CustomerBalance
	customerBalanceBytes, err := stub.GetState(data.BalancePrefix+contractID+"_"+customerID)
	if err != nil {
	fmt.Println("Error retrieving Customer Balance " + customerID)
	return allBalances, errors.New("Error retrieving Customer Balnace " + customerID)
	}
	err = json.Unmarshal(customerBalanceBytes, &customerBalance)
	fmt.Println("Customer Balance  : " , customerBalance);

	var customerFBalance data.CustomerFamilyBalance
	customerFBalanceBytes, err := stub.GetState(data.BalancePrefix+"_"+contractID)
	if err != nil {
	fmt.Println("Error retrieving Customer Balance " + customerID)
	return allBalances, errors.New("Error retrieving Customer Balnace " + customerID)
	}
	err = json.Unmarshal(customerFBalanceBytes, &customerFBalance)
	fmt.Println("Customer Family Balance  : " , customerFBalance);


	allBalances.ContractID = contractID
	allBalances.MemberID = customerID
	allBalances.I_Ded_Balance = customerBalance.I_Ded_Balance
	allBalances.I_OOP_Balance = customerBalance.I_OOP_Balance
	allBalances.F_Ded_Balance = customerFBalance.F_Ded_Balance
	allBalances.F_OOP_Balance = customerFBalance.F_OOP_Balance

	fmt.Println("In query.GetAllBalancesOfCustomer end ")

	return allBalances, err
}

func GetCustomerBalance(customerID string,contractId string, stub shim.ChaincodeStubInterface)(data.CustomerBalance, error) {
	fmt.Println("In query.GetCustomerBalance start ")
	var customerBalance data.CustomerBalance
	customerBalanceBytes, err := stub.GetState(data.BalancePrefix+contractId+"_"+customerID)
	if err != nil {
	fmt.Println("Error retrieving Customer Balance " + customerID)
	return customerBalance, errors.New("Error retrieving Customer Balnace " + customerID)
	}
	err = json.Unmarshal(customerBalanceBytes, &customerBalance)
	fmt.Println("Customer Balance  : " , customerBalance);
	fmt.Println("In query.GetCustomerBalance end ")
	return customerBalance, err
}

func GetCustomerFBalance(customerID string,contractId string, stub shim.ChaincodeStubInterface)(data.CustomerFamilyBalance, error) {
	fmt.Println("In query.GetCustomerFBalance start ")
	var customerFBalance data.CustomerFamilyBalance
	customerFBalanceBytes, err := stub.GetState(data.BalancePrefix+"_"+contractId)
	if err != nil {
	fmt.Println("Error retrieving Customer Family Balance " + contractId)
	return customerFBalance, errors.New("Error retrieving Customer Family Balnace " + contractId)
	}
	err = json.Unmarshal(customerFBalanceBytes, &customerFBalance)
	fmt.Println("Customer Balance  : " , customerFBalance);
	fmt.Println("In query.GetCustomerFBalance end ")
	return customerFBalance, err
}



func GetContract(contractId string, stub shim.ChaincodeStubInterface)(data.ContractDefinition, error) {
	fmt.Println("In query.GetContract start ")
	var customerContract data.ContractDefinition
	customerContractBytes, err := stub.GetState(data.ContractPrefix+"_"+contractId)

	if err != nil {
		fmt.Println("Error retrieving Customer Contract 1 " + contractId)
		return customerContract, errors.New("Error retrieving  Contract " + contractId)
	}
	err = json.Unmarshal(customerContractBytes, &customerContract)
	fmt.Println("Customer Contract  : " , customerContract);
	if (err != nil) {
			fmt.Println("Error receiving the customer contract 2 ", err)
			return customerContract, err
	}
		fmt.Println("In query.GetContract end ")
	return customerContract, nil
}

func GetContractsList(stub shim.ChaincodeStubInterface)([]string, error) {
	fmt.Println("In query.GetContractsList start ")
	var contractsList []string

	key := data.AllContractPrefix
	contractListBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving ContractsList")
		return contractsList, errors.New("Error retrieving ContractsList")
	}
	err = json.Unmarshal(contractListBytes, &contractsList)
	fmt.Println("ContractList : " , contractsList);
	fmt.Println("In query.GetContractsList end ")

	return contractsList, nil
}

func GetSubscriberMemberRelation(subscriberID string, stub shim.ChaincodeStubInterface)([]data.MemberRelation, error) {
	fmt.Println("In query.GetSubscriberMemberRelation start ")
	key := data.CMemberPrefix+"_"+subscriberID
	var memberRelations []data.MemberRelation
	memberRelationBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving Customer Members" , subscriberID)
		return memberRelations, errors.New("Error retrieving Customer Members " + subscriberID)
	}
	err = json.Unmarshal(memberRelationBytes, &memberRelations)
	fmt.Println("Members   : " , memberRelations);
	fmt.Println("In query.GetSubscriberMemberRelation end ")

	return memberRelations, err
}

func GetMembersOfContract(contractID string, stub shim.ChaincodeStubInterface)(data.ContractMembers, error) {
	fmt.Println("In query.GetContractMembers start ")

	key := data.ContractMemsPrefix+"_"+contractID
	var membersOfContract data.ContractMembers
	membersOfContractBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving Contract Members" , contractID)
		return membersOfContract, errors.New("Error retrieving Contract Members " + contractID)
	}

	err = json.Unmarshal(membersOfContractBytes, &membersOfContract)
	fmt.Println("Contract Members   : " , membersOfContract);
	fmt.Println("In query.GetContractMembers end ")

	return membersOfContract, nil
}
func GetClaim(customerID string,contractID string,claimID string, stub shim.ChaincodeStubInterface)(data.Claims, error) {
	fmt.Println("In query.GetClaim start ")
	var claim data.Claims
	claimBytes, err := stub.GetState(data.ClaimPrefix+contractID+"_"+customerID+"_"+claimID)
	if err != nil {
		fmt.Println("Error retrieving GetClaim  " + claimID)
		return claim, errors.New("Error retrieving Claim  " + claimID)
	}
	err = json.Unmarshal(claimBytes, &claim)
	fmt.Println("Claim   : " , claim);
	fmt.Println("In query.GetClaim end ")
	return claim, err
}

func GetContracts(customerID string, stub shim.ChaincodeStubInterface)([]string, error) {
	fmt.Println("In query.GetContracts start ")

	var customerContracts []string
	key :=data.ContractPrefix+"_"+customerID
	tempBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving Contracts for customer  ",customerID)
		return customerContracts, errors.New("Error retrieving for customer  "+customerID )
	}
	err = json.Unmarshal(tempBytes, &customerContracts)

	fmt.Println("In query.GetContracts end ")
	return customerContracts, err
}

func GetMemberClaimIDsOfContract(memberID string, contractID string,stub shim.ChaincodeStubInterface)([]string, error) {
	fmt.Println("In query.GetMemberClaimIDsOfContract start ")

	var contractClaims []string

	key := data.ContractClaimsPrefix+"_"+contractID+"_"+memberID
	tempBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving member claims of contract  ",contractID)
		return contractClaims, errors.New("Error retrieving member claims of contract  "+contractID )

	}
	err = json.Unmarshal(tempBytes, &contractClaims)

	if err != nil {
		fmt.Println("Error unmarshalling member claims of contract " ,err)
		return contractClaims,err
	}

	fmt.Printf("Claims of Member %v from contract %v  \n ", memberID,contractID,contractClaims)
	fmt.Println("In query.GetMemberClaimIDsOfContract end ")
	return contractClaims, nil
}

func GetAllClaimsOfContract(customerID string, contractID string,stub shim.ChaincodeStubInterface)([]data.Claims, error) {
	fmt.Println("In query.GetAllClaimsOfContract start ")
	var claimsList []data.Claims
	var contractClaims []string

	contractClaims, err := GetMemberClaimIDsOfContract(customerID, contractID,stub)
	if err != nil {
		fmt.Println("Error retrieving claims of contract  ",contractID)
		return claimsList, errors.New("Error retrieving claims of contract  "+contractID )
	}

	if (len(contractClaims) > 0 ){

		for _, claimID := range contractClaims {
			var claim data.Claims
			claim, err = GetClaim(customerID ,contractID ,claimID,stub)
			if err != nil {
				return claimsList,err
			}
			fmt.Println("Claim : " , claim)
 			claimsList = append(claimsList,claim)

		}
	}else {
		fmt.Println("No Claims found for contract : "+contractID )

	}

	fmt.Println("In query.GetAllClaimsOfContract end ")
	return claimsList, err
}

func GetAllClaimsOfContract_New(subscriberID string, contractID string,stub shim.ChaincodeStubInterface)([]data.Claims, error) {
	fmt.Println("In query.GetAllClaimsOfContract start ")
	var claimsList []data.Claims
	var contractClaims []string

	//1. Get all the memberIDs of the contract
	memberID := "0"
	contractClaims, err := GetMemberClaimIDsOfContract(memberID, contractID,stub)
	if err != nil {
		fmt.Println("Error retrieving claims of contract  ",contractID)
		return claimsList, errors.New("Error retrieving claims of contract  "+contractID )
	}

	if (len(contractClaims) > 0 ){

		for _, claimID := range contractClaims {
			var claim data.Claims
			claim, err = GetClaim(memberID ,contractID ,claimID,stub)
			if err != nil {
				return claimsList,err
			}
			fmt.Println("Claim : " , claim)
 			claimsList = append(claimsList,claim)

		}
	}else {
		fmt.Println("No Claims found for contract : "+contractID )

	}

	fmt.Println("In query.GetAllClaimsOfContract end ")
	return claimsList, err
}

func GetLastTxnsList(stub shim.ChaincodeStubInterface)([]data.LastTxnOfContract, error) {
	fmt.Println("In query.GetLastTxnsList start ")

	var lastTxns []data.LastTxnOfContract

	key := data.LastTxnsPrefix
	fmt.Println("Key : ",key)
	tempBytes, err := stub.GetState(key)
	//fmt.Println("Got object : ",tempBytes)
	if err != nil {
		fmt.Println("Error retrieving last txns of contract  ")
		return lastTxns, errors.New("Error retrieving last txns of contract ")

	}

	err = json.Unmarshal(tempBytes, &lastTxns)

	if err != nil {
		fmt.Println("Error unmarshalling last transactions " ,err)
		if strings.Contains(err.Error(), "unexpected end") {

			return lastTxns, nil
		}

		return lastTxns,err
	}

	fmt.Println("In query.GetLastTxnsList end ")
	return lastTxns, nil
}


func GetLastIndTxn(contractID string, stub shim.ChaincodeStubInterface)(data.LastTxnOfContract, error) {
	fmt.Println("In query.GetLastIndTxn start ")

	var lastTxn data.LastTxnOfContract
	key := data.LastTxnsIndPrefix+"_"+contractID
	fmt.Println("Key : ",key)
	tempBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving last txns for   ",key)
		return lastTxn, errors.New("Error retrieving last txns for " +key)
	}
	err = json.Unmarshal(tempBytes, &lastTxn)

	if err != nil {
		fmt.Println("Error unmarshalling last transactions " ,err)
		if strings.Contains(err.Error(), "unexpected end") {
			return lastTxn, errors.New("Empty")
		}
		return lastTxn, err
	}

	fmt.Println("Last Ind txn found is :  ",lastTxn )
	fmt.Println("In query.GetLastIndTxn end ")
	return lastTxn, nil
}


func GetLastFamilyTxns(contractID string, stub shim.ChaincodeStubInterface)([]data.LastTxnOfContract, error) {
	fmt.Println("In query.GetLastFamilyTxns start ")

	var lastTxns []data.LastTxnOfContract
	key := data.LastTxnsFamPrefix+"_"+contractID
	fmt.Println("Key : ",key)
	tempBytes, err := stub.GetState(key)

	if err != nil {
		fmt.Println("Error retrieving last txns for   ",key)
		return lastTxns, errors.New("Error retrieving last txns for " +key)

	}
	err = json.Unmarshal(tempBytes, &lastTxns)

	if err != nil {
		fmt.Println("Error unmarshalling last transactions " ,err)
		if strings.Contains(err.Error(), "unexpected end") {
			return lastTxns, errors.New("Empty")
		}
		return lastTxns, err
	}

	fmt.Println("In query.GetLastFamilyTxns end ")
	return lastTxns, nil
}


func GetMemberDetailsNHistory(contractID string, subscriberID string,memberID string,
	claimID string,stub shim.ChaincodeStubInterface) (data.MemberNContractDetails, error) {
	fmt.Println("In services.GetMemberDetailsNHistory start ")
 	var memNContDetails data.MemberNContractDetails

	//1. Get Contract Definition
	contractDefinition,err := GetContract(contractID,stub)
	if err != nil {
		fmt.Println("Error receiving  the Customer Contract")
		return memNContDetails, errors.New("Error receiving  the Customer Contract")
	}

	fmt.Println("Contract :", contractDefinition)
	memNContDetails.Contract = contractDefinition

	//2. Get Member Limits
	customerLimit,err := GetCustomerLimits(memberID,contractID,stub)
	if err !=nil {
		fmt.Println("Error retrieving Customer Limits " , memberID )
		return memNContDetails, errors.New("Error retrieving Customer Limits " + memberID)
	}
	fmt.Println("Customer Limits received : " ,customerLimit)

	memNContDetails.I_Ded_Limit = customerLimit.I_Ded_Limit
	memNContDetails.I_OOP_Limit = customerLimit.I_OOP_Limit

	//3. Get MemberDetails Demographics

	customerDetails,err := GetCustomerDetails(memberID,stub)
	if (err != nil){
		fmt.Println("No details found for customer :  ",memberID)
	}
	memNContDetails.MemberInfo = customerDetails

	//4. Get Contract Members and Relationships
	var memberRelations []data.MemberRelation
	memberRelations, err =	GetSubscriberMemberRelation(subscriberID, stub)
	if (err != nil){
		fmt.Println("No Relationships found for subscriberID :  ",subscriberID)
	}
	memNContDetails.MemberInfo = customerDetails

	membersOfContract, err := GetMembersOfContract(contractID,stub)
	if (err != nil ){
		fmt.Println("No Members Found. ")
	}
	fmt.Printf("Members of Contract %v is %v \n", contractID,membersOfContract)

	memberIDs := membersOfContract.MemberIDs
	for _,memID := range memberIDs{
		fmt.Println("Member ID :  ",memID)
		customerDetails,err = GetCustomerDetails(memID,stub)
		if (err != nil){
			fmt.Println("No details found for customer :  ",memID)
		}
		var subMemRelation data.SubMemRelation

		subMemRelation.SubscriberID = subscriberID
		subMemRelation.MemberID = memID
		subMemRelation.FirstName = customerDetails.FirstName
		subMemRelation.LastName = customerDetails.LastName

		for _,memRel := range memberRelations{
			if (memRel.MemberID == memID ){
					subMemRelation.Relationship = memRel.Relationship
					break
			}

		}
		memNContDetails.SubMemRelation = append(memNContDetails.SubMemRelation,subMemRelation)
		fmt.Printf("Member %v details :  \n",memID,customerDetails)

	}

	//5. Get Customer Balances
	if (contractDefinition.ContractType == "I" ){
		var customerBalance data.CustomerBalance
		customerBalance,err := GetCustomerBalance(memberID,contractID,stub)

		if err != nil {
			fmt.Println("Error retrieving Customer Balance " , memberID)
			//return nil, errors.New("Error retrieving Customer Balance " + memberID)
		}
		fmt.Println("Customer Balance received : " ,customerBalance)
		memNContDetails.I_Ded_Balance = customerBalance.I_Ded_Balance
		memNContDetails.I_OOP_Balance = customerBalance.I_OOP_Balance

	}
	if (contractDefinition.ContractType == "F" ){

		var customerAllBalance data.CustomerAllBalance
		customerAllBalance,err := GetAllBalancesOfCustomer(memberID,contractID,stub)

 	 if err != nil {
 		 fmt.Println("Error retrieving Customer Balance " , memberID )
 		// return nil, errors.New("Error retrieving Customer Balance " + memberID)
 	 }

 	 fmt.Println("Customer All Balance received : " ,customerAllBalance)
 	 memNContDetails.I_Ded_Balance = customerAllBalance.I_Ded_Balance
 	 memNContDetails.I_OOP_Balance = customerAllBalance.I_OOP_Balance
 	 memNContDetails.F_Ded_Balance = customerAllBalance.I_Ded_Balance
 	 memNContDetails.F_OOP_Balance = customerAllBalance.I_Ded_Balance

  }
	//6. Get Txn details of Member



	fmt.Println("In services.GetMemberDetailsNHistory end ")
	return memNContDetails,nil
}
