package query

import (
	"encoding/json"
  "errors"
	"fmt"
	"strings"
	"strconv"
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

	if function == "searchTxns" {
		fmt.Println("Invoking SearchTxns " + function)
		var allTxns []data.LastTxnOfContract
		allTxns,err := SearchTxns(args,stub)
		if err != nil {
			fmt.Println("Error In Search operation : ",err)
			return nil, errors.New("Error In Search operation.")
		}
		fmt.Println("All success, returning search details")
		return json.Marshal(allTxns)
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

	if function == "getCustomerRelationships" {
		fmt.Println("Invoking GetSubscriberMemberRelation " + function)
		memberRelations,err := GetSubscriberMemberRelation(args[0], stub)
		if err != nil {
			fmt.Println("Error receiving  the Customer Relationships")
			return nil, errors.New("Error receiving  the Customer Relationships")
		}
		fmt.Println("All success, returning the Customer Relationships")
		return json.Marshal(memberRelations)
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

	if function == "getContract" {
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
	}
	//GetAllClaimsOfContract
	if function == "getAllClaimsOfContract" {
		fmt.Println("Invoking GetClaimsOfContract " + function)
		contractID := args[0]
		var claimList []data.Claims
		claimList,err := GetClaimsOfContract(contractID,stub)
		if err != nil {
			fmt.Println("Error receiving  the Claim list ")
			return nil, errors.New("Error receiving  the Claim list")
		}
		fmt.Println("All success, returning the Claim List ")
		return json.Marshal(claimList)
	}

	if function == "getTxnsFromLedger" {
		fmt.Println("Invoking GetTxnsFromLedger  " + function)
		if (len(args) < 1 ){
			return nil, errors.New("Not enough input params. Minimum params expected :1 Max : 2")
		}
		noOftxns,err := strconv.Atoi(args[0])
		if (err != nil){
			return nil, errors.New("Enter numeric value for no of txns to be returned")
		}
		source := ""
		if (len(args) >1 ){
			source = args[1]
		}
		var lastTxns []data.LastTxnOfContract
		lastTxns,err = GetTxnsFromLedger(noOftxns,source,stub)
		if err != nil {
			fmt.Println("Error receiving  last 50 the txns  ")
			return nil, errors.New("Error receiving  last 50 the txns ")
		}
		fmt.Println("All success, returning the last 50 txns ")
		return json.Marshal(lastTxns)
	}

	if function == "getTxnsOfMemCont" {
		fmt.Println("Invoking GetTxnsOfMemCont  " + function)
		if (len(args) < 4 ){
			return nil, errors.New("Not enough input params. Minimum params expected :4 Max : 5")
		}
		subscriberID := args[0]
		memberID := args[1]
		contractID := args[2]
		noOftxns,err := strconv.Atoi(args[3])
		source := ""
		if (len(args) >4 ){
			source = args[4]
		}
		var lastTxns []data.LastTxnOfContract
		lastTxns,err = GetTxnsOfMemCont(contractID,subscriberID,memberID,source,noOftxns,stub)
		if err != nil {
			fmt.Println("Error receiving the Txns of Member & Contract  ")
			return nil, errors.New("Error receiving the Txns of Member & Contract ")
		}
		fmt.Println("All success, returning the txns of member & contract ")
		return json.Marshal(lastTxns)
	}

	if function == "getTxnsOfMem" {
		fmt.Println("Invoking GetTxnsOfMem  " + function)
		if (len(args) < 2 ){
			return nil, errors.New("Not enough input params. Minimum params expected :2 Max : 3")
		}
		memberID := args[0]
		noOftxns,err := strconv.Atoi(args[1])
		source := ""
		if (len(args) >2 ){
			source = args[2]
		}
		var lastTxns []data.LastTxnOfContract
		lastTxns,err = GetTxnsOfMem(memberID,noOftxns,source,stub)
		if err != nil {
			fmt.Println("Error receiving the Txns of Member")
			return nil, errors.New("Error receiving the Txns of Member")
		}
		fmt.Println("All success, returning the txns of member ")
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

	if function == "getContractsList" {
		fmt.Println("Invoking getContractsList  " + function)
		var contractsList []string
		contractsList,err := GetContractsList(stub)
		if err != nil {
			fmt.Println("Error getting  contracts from ledger  ")
			return nil, errors.New("Error getting  contracts from ledger ")
		}
		fmt.Println("All success, returning the last Ind txns List of contract ")
		return json.Marshal(contractsList)
	}

	return nil, errors.New("Received unknown query function name")
}

func GetCustomerDetails(customerID string, stub shim.ChaincodeStubInterface) (data.Customer, error) {
	fmt.Println("In query.GetCustomerDetails start ")
	var customer data.Customer
	customerBytes, err := stub.GetState(customerID)
	if err != nil {
		fmt.Println("Error retrieving Customer Details " + customerID)
		return customer, errors.New("No Customer Found with ID " + customerID)
	}
	err = json.Unmarshal(customerBytes, &customer)
	if err != nil {
		fmt.Println("Error Unmarshalling Customer Details " + customerID)
		return customer, errors.New("No Customer Found with ID" + customerID)
	}
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

func GetSubscriber(contractID string, stub shim.ChaincodeStubInterface)(string, error) {
	fmt.Println("In query.GetSubscriber start ")
	contractDefinition,err := GetContract(contractID,stub)
	if err != nil {
		fmt.Println("No Contract Found with contracID" ,  contractID)
		return "", errors.New("No Contract Found with contracID " + contractID)
	}
	fmt.Println("In query.GetSubscriber end ")
	return contractDefinition.SubscriberID, nil
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

func GetCarrierBalances(memberID string,contractID string, stub shim.ChaincodeStubInterface)(data.ContMemCarrierBalance, error) {
	fmt.Println("In query.GetCarrierBalances start ")
	var contMemCarrierBalance data.ContMemCarrierBalance
	key := data.ContCarrierBalPrefix+"_"+contractID+"_"+memberID
	fmt.Println("Get Carrier Balances with Key :  ",key)
	contMemCarrierBalanceBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving Carrier Balances " + contractID)
		return contMemCarrierBalance, errors.New("Error retrieving Carrier  Balances  " + contractID)
	}
	err = json.Unmarshal(contMemCarrierBalanceBytes, &contMemCarrierBalance)
	fmt.Println("Member Carrier Balances : " , contMemCarrierBalance)
	fmt.Println("In query.GetCarrierBalances end ")
	return contMemCarrierBalance, err
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
	//7. Get Carrier Balances
	var contMemCarrierBalance data.ContMemCarrierBalance
	contMemCarrierBalance,err = GetCarrierBalances(memberID,contractID,stub)
	fmt.Println("Customer Carrier Balances received : " ,contMemCarrierBalance)
	if (err != nil ){
		fmt.Println("No Customer Carrier Balances Found")
	}else {
		carrierBalance := contMemCarrierBalance.CarrierBalances[strings.ToUpper(data.MedicalCarrier)]
		if (carrierBalance.CarrierName != ""){
			memNContDetails.MemContCarrierBal =append(memNContDetails.MemContCarrierBal,carrierBalance)
		}
		carrierBalance = contMemCarrierBalance.CarrierBalances[strings.ToUpper(data.PharmacyCarrier)]
		if (carrierBalance.CarrierName != ""){
			memNContDetails.MemContCarrierBal =append(memNContDetails.MemContCarrierBal,carrierBalance)
		}
		carrierBalance = contMemCarrierBalance.CarrierBalances[strings.ToUpper(data.DentalCarrier)]
		if (carrierBalance.CarrierName != ""){
			memNContDetails.MemContCarrierBal =append(memNContDetails.MemContCarrierBal,carrierBalance)
		}
	}

	fmt.Println("Carrier Balances:  ", memNContDetails.MemContCarrierBal)
	fmt.Println("In services.GetMemberDetailsNHistory end ")
	return memNContDetails,nil
}

func SearchTxns(args []string,stub shim.ChaincodeStubInterface)([]data.LastTxnOfContract, error) {
  fmt.Println("In query.SearchTxns start ")
  // var args []string
	// searchString :=`{"ContractID":"123","SubscriberID":"112200","MemberID":"112200"}`
	// args = append(args,searchString)
	fmt.Println("Input params : ", args)
	b := []byte(args[0])
	noOftxns,err := strconv.Atoi(args[1])

	var f interface{}
	err = json.Unmarshal(b, &f)
	if (err != nil ){
			fmt.Println("Error unmarshalling input params : ",err)
			return nil, errors.New("Provide Inputs in correct format")
	}
	searchMap := f.(map[string]interface{})
	fmt.Println("searchMap : ",searchMap)

	contractID := ""
	if (searchMap["ContractID"] != nil){
		contractID = searchMap["ContractID"].(string)
	}
	subscriberID := ""
	if (searchMap["SubscriberID"] != nil){
		subscriberID = searchMap["SubscriberID"].(string)
	}
	memberID := ""
	if (searchMap["MemberID"] != nil){
		memberID = searchMap["MemberID"].(string)
	}
	source := ""
	if (searchMap["Source"] != nil){
		source = searchMap["Source"].(string)
	}

	// txnStartDt := searchMap["TxnStartDt"]
	// txnEndDt := searchMap["TxnEndDt"]

	var allTxns []data.LastTxnOfContract
	if ((contractID !="" && subscriberID !="" && memberID !=""  ) ||
			(contractID !="" && subscriberID =="" && memberID !="" )){
		fmt.Printf("Search for contract %v,subscriber %v,member %v, source %v \n",contractID,subscriberID,memberID,source)
		allTxns,err = GetTxnsOfMemCont(contractID, subscriberID,memberID,source,noOftxns,stub)
		if (err != nil){
				fmt.Printf("Error getting  txns for contract %v ,subscriber %v, member %v ,source %v \n", contractID,subscriberID,memberID,source)
				return allTxns,errors.New("Error getting txns of contract "+contractID +" of subscriber "+subscriberID+"for Member " +memberID+" from source "+source)
		}
		 return allTxns,nil
	}else if ((contractID !="" && subscriberID !="" && memberID =="") ||
						(contractID !="" && subscriberID ==""  && memberID =="" )||
						(contractID =="" && subscriberID !=""  && memberID =="") ){
		fmt.Printf("Search for contractID %v,subscriberID %v, source %v \n",contractID,subscriberID,source)
		allTxns,err = GetTxnsOfContract(contractID,noOftxns,source,stub)
		if (err != nil){
				fmt.Printf("Error getting  txns of contract %v  of subscriberID %v from source &v   \n", contractID,subscriberID,source)
				return allTxns,errors.New("Error getting txns of contract "+contractID +" of Subscriber " +subscriberID +" from source "+source)
		}
		return allTxns,nil
	}else if (contractID =="" && subscriberID ==""  && memberID !="") {
			fmt.Printf("Search for memberID %v, source %v \n",memberID,source)
			allTxns,err = GetTxnsOfMem(memberID, noOftxns,source,stub)
		if (err != nil){
				fmt.Printf("Error getting  txns of memberID %v from source %v \n", memberID,source)
				return allTxns,errors.New("Error getting txns of Member "+memberID +" from Source "+source)
		}
	}else if (contractID =="" && subscriberID ==""  && memberID =="") {
		fmt.Printf("Search for source %v . Get from Ledger\n",source)
		allTxns,err = GetTxnsFromLedger(noOftxns,source,stub)
		if (err != nil){
				fmt.Println("Error getting  txns from Ledger")
				return allTxns,errors.New("Error getting  txns from Ledger")
		}
	}
  fmt.Println("In query.SearchTxns end ")
	return allTxns, nil
}

func GetIndexOfTxnFromClaim(claim data.Claims,txnID int)(int){
	fmt.Println("In services.GetIndexOfTxnFromClaim start ")
	for index,txn := range claim.Transactions{
		if (txn.TransactionID == txnID){
			fmt.Println("Txn match found. Index is : ", index)
			return index
		}
	}
	fmt.Println("No matching txn found. Retrun 0")
	fmt.Println("In services.GetIndexOfTxnFromClaim end ")
	return 0
}
