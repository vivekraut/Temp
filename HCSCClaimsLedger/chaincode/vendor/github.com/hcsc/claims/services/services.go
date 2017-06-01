package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hcsc/claims/data"
	"github.com/hcsc/claims/query"
	"encoding/json"
)

// Proces the claims of the Contract
func ProcessClaim(args []string,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.ProcessClaim start ")
	subscriberID :=args[0]
	memberID :=args[1]
	contractID  :=args[2]
	var accumTypes []data.Accum
	err := json.Unmarshal([]byte(args[3]), &accumTypes)
	fmt.Println("accumTypes are ", accumTypes)
	if err!= nil {
		fmt.Println("Error unmarshal accumtypes  ", err)
	}
	claimID :=args[4]
	adjustTxnID,err := strconv.Atoi(args[5])
	source := args[6]
	errorDes := ""
	fmt.Println("subscriberID  ", subscriberID)
	fmt.Println("memberID  ", memberID)
	var membersOfContract data.ContractMembers
	membersOfContract,err = query.GetMembersOfContract(contractID,stub)
	fmt.Println("Members of Contract: ",membersOfContract)
	if err !=nil {
		fmt.Println("Error retrieving Members of Contract  " , contractID )
	}
 	memberFound := false
	memberIDs := membersOfContract.MemberIDs
	fmt.Printf("No of Members in contract %v is %v \n", contractID, len(memberIDs))
	if (len(memberIDs) > 0 ){
		for _,memID := range memberIDs {
			fmt.Printf("memID : %v , input memberID : %v \n", memID,memberID)
			if (memID == memberID){
				memberFound = true
				fmt.Printf("Found Member %v in Contract %v \n",memberID,contractID )
				 break;
			}
		}
		if (!memberFound){
		 errorDes = "{\"Error\":\"Member " + memberID + " Not Found in Contract "+contractID +"\"}"
		 fmt.Println(errorDes)
		 return nil, errors.New(errorDes)
	 }
	}
	fmt.Println(memberFound)
	var customerLimit data.CustomerLimit
	customerLimit,err = query.GetCustomerLimits(memberID,contractID,stub)
	if err !=nil {
		fmt.Println("Error retrieving Customer Limits " , memberID )
		return nil, errors.New("Error retrieving Customer Limits " + memberID)
	}
	fmt.Println("Customer Limits received : " ,customerLimit)
	var newClaim bool = false
	claim, err := query.GetClaim(memberID,contractID,claimID,stub)
	if err != nil {
		claim.ClaimID =claimID
		claim.ContractID=contractID
		claim.MemberID=memberID
		claim.SubscriberID=subscriberID
		claim.Participant =strings.Title(strings.ToLower(source))
		claim.Status = ""
		claim.TotalClaimAmount = 0
		claim.CreateDTTM = time.Now()
		newClaim = true
		fmt.Println("New Claim : ",claim)
	}
	//TODO - update later
	var txnID int = 0
	contractDefinition,err := query.GetContract(contractID,stub)
	if err != nil {
		fmt.Println("Error receiving  the Customer Contract")
		return nil, errors.New("Error receiving  the Customer Contract")
	}
	fmt.Println("Contract :", contractDefinition)
	if (contractDefinition.ContractType == "I"){
		fmt.Println("Recevied Individual Claim for Contract ", contractID)
		lastTxn, err := query.GetLastIndTxn(contractID,stub)
		if (err != nil ){
			fmt.Println("No Last txn found. Set txnID as : 1")
			txnID = 1
		}else {
			txnID = lastTxn.Transaction.TransactionID + 1
		}
		fmt.Println("Last ind txn for contract  ",lastTxn )
		if accumTypes != nil {
			for _, accum := range accumTypes {
				if (accum.Type == "IIDED"){
					_,err :=ProcessIIDEDClaim(newClaim,accumTypes,claim,customerLimit,txnID,adjustTxnID,stub)
					if (err != nil){
						fmt.Println("Error processing IIDED claims")
					}
					break
				}
				if (accum.Type == "IIOOP"){
					_,err :=ProcessIIOOPClaim(newClaim,accum.Amount,claim,customerLimit,txnID,adjustTxnID,stub)
					if (err != nil){
						fmt.Println("Error processing IIOOP claims")
					}
					break
				}
			}
		}
	}
	if (contractDefinition.ContractType == "F"){
		fmt.Println("Recevied Family Claim for Contract ", contractID)
		lastTxns, err := query.GetLastFamilyTxns(contractID,stub)
		if (err != nil ){
			fmt.Println("No Last txns found . ")
		}
		var memberTxnFound bool = false
		if len(lastTxns) > 0 {
			for _, lastTxn := range lastTxns {
				fmt.Println("Found last txn is from claim    : ",lastTxn.ClaimID)
				if (memberID == lastTxn.MemberID ){
					fmt.Printf("Last txn ID is \n",lastTxn.Transaction.TransactionID )
					txnID = lastTxn.Transaction.TransactionID + 1
					memberTxnFound = true
					break
				}
			}
		}
		if (!memberTxnFound){
			txnID = 1
		}
		if accumTypes != nil {
			for _, accum := range accumTypes {
				if (accum.Type == "IIDED"){
					_,err :=ProcessIFDEDClaim(newClaim,accumTypes,claim,contractDefinition,
						customerLimit,txnID,adjustTxnID,stub)
					if (err != nil){
						fmt.Println("Error processing IFDED claims")
					}
					break
				}
				if (accum.Type == "IIOOP"){
					_,err :=ProcessIFOOPClaim(newClaim,accum.Amount,claim,contractDefinition,
						customerLimit,txnID,adjustTxnID,stub)
					if (err != nil){
						fmt.Println("Error processing IFDED claims")
					}
					break
				}
			}
		}
	}
	fmt.Println("In services.ProcessClaim end ")
	return nil,nil
}

// Proces the IFDED  & IFOOP claim from Family contract
func ProcessIFDEDClaim(newClaim bool, accums []data.Accum,claim data.Claims,
	contractDefinition data.ContractDefinition,
	customerLimit data.CustomerLimit,txnID int, adjustTxnID int,stub shim.ChaincodeStubInterface)(bool, error){

	fmt.Println("In services.ProcessIFDEDClaim start ")
	var customerAllBalance data.CustomerAllBalance
	customerAllBalance,err := query.GetAllBalancesOfCustomer(claim.MemberID,claim.ContractID,stub)
	if err != nil {
		fmt.Println("Error retrieving Customer Balance " , claim.MemberID )
		return false, errors.New("Error retrieving Customer Balance " + claim.MemberID)
	}
	fmt.Println("Customer Balance received : " ,customerAllBalance)

	fDedAccumTrans := data.Transaction{AccumType:"IFDED",TransactionDate:time.Now(),
		TxnUpdatedDate:time.Now(),TransactionID:txnID}
	fOopAccumTrans := data.Transaction{AccumType:"IFOOP",
		TransactionDate:time.Now(),TxnUpdatedDate:time.Now(),TransactionID:txnID}
	iDedAccumTrans := data.Transaction{AccumType:"IIDED",
		TransactionDate:time.Now(),TxnUpdatedDate:time.Now(),TransactionID:txnID}
	iOopAccumTrans := data.Transaction{AccumType:"IIOOP",
		TransactionDate:time.Now(),TxnUpdatedDate:time.Now(),TransactionID:txnID}

	for _, accum := range accums {
		if (accum.Type == "IIDED"){
			fDedAccumTrans.AccumAmount = accum.Amount
			iDedAccumTrans.AccumAmount = accum.Amount
			fOopAccumTrans.AccumAmount = accum.Amount
			iOopAccumTrans.AccumAmount = accum.Amount
		}
		if (accum.Type == "IIOOP"){
			fOopAccumTrans.AccumAmount = accum.Amount
			iOopAccumTrans.AccumAmount = accum.Amount
		}
	}
	fDedBalance := customerAllBalance.F_Ded_Balance
	fOopBalance := customerAllBalance.F_OOP_Balance
	iDedBalance := customerAllBalance.I_Ded_Balance
	iOopBalance := customerAllBalance.I_OOP_Balance

	lastTxn :=data.LastTxnOfContract{ClaimID:claim.ClaimID,ContractID:claim.ContractID,
		MemberID:claim.MemberID,SubscriberID:claim.SubscriberID, Source:claim.Participant}

	if (fDedAccumTrans.AccumAmount < 0 && adjustTxnID > 0){
		fmt.Printf("Adjustment claim has come in with amout %v.  Processing .... \n",fDedAccumTrans.AccumAmount)
		if (adjustTxnID > 0){
			//NOTE Added based on review feedback for -ve testing
			var txns_index_map = make(map[string]int)
			for index,tempTxn := range claim.Transactions{
				if (tempTxn.TransactionID == adjustTxnID){
					fmt.Println("Previous Txn to be adjusted is found. Get indexes of Txns")
					txns_index_map[tempTxn.AccumType] = index
				}
			}
			//Verify if the claim submitted with correct Accums as previous claim
			//in case of IFOOP claim submitted, only 2 txns should be there in with same txnid
			if (txns_index_map != nil ){
				 if (len(txns_index_map) < 4){
					 fmt.Println("This -ve / adjustment claim is submitted wrongly ")
		 			claim.Status = data.ClaimError
				 }else {
					customerAllBalance.F_Ded_Balance = customerAllBalance.F_Ded_Balance + fDedAccumTrans.AccumAmount
					customerAllBalance.F_OOP_Balance = customerAllBalance.F_OOP_Balance + fOopAccumTrans.AccumAmount
					customerAllBalance.I_Ded_Balance = customerAllBalance.I_Ded_Balance + iDedAccumTrans.AccumAmount
					customerAllBalance.I_OOP_Balance = customerAllBalance.I_OOP_Balance + iOopAccumTrans.AccumAmount

					claim.Status = data.ClaimProcessed
					_,err := UpdateCustomerAndFamilyBalances(customerAllBalance, stub)
					if (err !=nil ){
						fmt.Println("Error updating balances ", err)
					}
					updateCarrierBalances(claim.MemberID,claim.ContractID,iDedAccumTrans.AccumAmount,iOopAccumTrans.AccumAmount,claim.Participant,stub)
					//NOTE: basaed on 15/02 discussion add - find old txn which is adjusted with -ve value and mark the status resolved
					for _,txnIndex := range txns_index_map{
						claim.Transactions[txnIndex].Status =data.ClaimAdjusted
					}
					_,err = MarkLastFamTxnReview(false,claim.ContractID,claim.MemberID,claim.ClaimID,contractDefinition,customerLimit,stub )
					if (err != nil ){
							fmt.Println("Error Marking Last txn for review  ", err)
					}
				}
			}
		}else {
			//if No adjusted txn sent - report as error
			fmt.Println("No adjustment  txnID is provided")
		 claim.Status = data.ClaimError
		}
	}else {
		fDedUpdateFlag := false
		fOopUpdateFlag := false
		iDedUpdateFlag := false
		iOopUpdateFlag := false

		fDedUpdateFlag,fDedBalance = InvokeRule("IFDED", fDedAccumTrans.AccumAmount, contractDefinition.F_Ded_Limit,
			fDedBalance, &fDedAccumTrans )
		claim.Status = fDedAccumTrans.Status
		if (fDedUpdateFlag){
			fOopUpdateFlag,fOopBalance = InvokeRule("IFOOP", iDedAccumTrans.AccumAmount, contractDefinition.F_OOP_Limit,
				fOopBalance	, &fOopAccumTrans )
			if (fOopUpdateFlag) {
				iDedUpdateFlag,iDedBalance = InvokeRule("IIDED", fOopAccumTrans.AccumAmount, customerLimit.I_Ded_Limit,
					iDedBalance, &iDedAccumTrans )
				if (iDedUpdateFlag){
					iOopUpdateFlag,iOopBalance = InvokeRule("IIOOP", iOopAccumTrans.AccumAmount, customerLimit.I_OOP_Limit,
						iOopBalance , &iOopAccumTrans )
					if (iOopUpdateFlag){
						customerAllBalance.F_Ded_Balance = fDedBalance
						customerAllBalance.F_OOP_Balance = fOopBalance
						customerAllBalance.I_Ded_Balance = iDedBalance
						customerAllBalance.I_OOP_Balance = iOopBalance
						fmt.Println("Customer balances : ",customerAllBalance)
						_,err := UpdateCustomerAndFamilyBalances(customerAllBalance, stub)
						if (err !=nil ){
							fmt.Println("Error updating balances ", err)
						}
						updateCarrierBalances(claim.MemberID,claim.ContractID,iDedAccumTrans.AccumAmount,iOopAccumTrans.AccumAmount,claim.Participant,stub)
						claim.Status =fDedAccumTrans.Status
						claim.TotalClaimAmount = claim.TotalClaimAmount+fDedAccumTrans.AccumAmount
					}else {
						claim.Status = iOopAccumTrans.Status
					}
				}else {
					claim.Status = iDedAccumTrans.Status
				}
			}else {
				claim.Status = fOopAccumTrans.Status
			}
		}
	}
	fmt.Println("Claim Status is : ", claim.Status)
	fDedAccumTrans.Status = claim.Status
	fOopAccumTrans.Status = claim.Status
	iDedAccumTrans.Status = claim.Status
	iOopAccumTrans.Status = claim.Status

	fDedAccumTrans.AccumBalance =	customerAllBalance.F_Ded_Balance
	fOopAccumTrans.AccumBalance =	customerAllBalance.F_OOP_Balance
	iDedAccumTrans.AccumBalance =	customerAllBalance.I_Ded_Balance
	iOopAccumTrans.AccumBalance =	customerAllBalance.I_OOP_Balance

	claim.Transactions = append(claim.Transactions,fDedAccumTrans)
	claim.Transactions = append(claim.Transactions,fOopAccumTrans)
	claim.Transactions = append(claim.Transactions,iDedAccumTrans)
	claim.Transactions = append(claim.Transactions,iOopAccumTrans)

	fmt.Println("FDED txn : ",fDedAccumTrans)
	_,err = createUpdateClaim(claim,newClaim,stub)
	if err != nil {
		fmt.Println("createUpdateClaim failed")
		return false, err
	}
	lastTxn.Transaction = fDedAccumTrans
	//if (claim.Status != data.ClaimError){
	_,err = UpdateLastFamMemTxn(lastTxn,stub)

	if (err != nil ){
		fmt.Println("Error adding to Last Family Txns")
	}
	//	}
	 fmt.Println("In services.ProcessIFDEDClaim end ")
	 return true,nil
}

// Proces the  IFOOP claim from Family contract
func ProcessIFOOPClaim(newClaim bool, claimAmout float64,claim data.Claims,contractDefinition data.ContractDefinition,
		customerLimit data.CustomerLimit,txnID int,adjustTxnID int, stub shim.ChaincodeStubInterface)(bool, error){
	fmt.Println("In services.ProcessIFOOPClaim start ")
	var customerAllBalance data.CustomerAllBalance
	customerAllBalance,err := query.GetAllBalancesOfCustomer(claim.MemberID,claim.ContractID,stub)
	if err != nil {
		fmt.Println("Error retrieving Customer Balance " , claim.MemberID )
		return false, errors.New("Error retrieving Customer Balance " + claim.MemberID)
	}
	fmt.Println("Customer Balance received : " ,customerAllBalance)

	fOopAccumTrans := data.Transaction{AccumType:"IFOOP", AccumAmount:claimAmout,
		TransactionDate:time.Now(),TxnUpdatedDate:time.Now(),TransactionID:txnID}
	iOopAccumTrans := data.Transaction{AccumType:"IIOOP", AccumAmount:claimAmout,
		TransactionDate:time.Now(),TxnUpdatedDate:time.Now(),TransactionID:txnID}

	fOopBalance := customerAllBalance.F_OOP_Balance
	iOopBalance := customerAllBalance.I_OOP_Balance
	lastTxn :=data.LastTxnOfContract{ClaimID:claim.ClaimID,ContractID:claim.ContractID,
			MemberID:claim.MemberID,SubscriberID:claim.SubscriberID,
				Source:claim.Participant}
	if (claimAmout < 0 ){
		fmt.Printf("Adjustment claim has come in with amout %v.  Processing .... \n",claimAmout)
		if (adjustTxnID > 0){
			//NOTE Added based on review feedback for -ve testing
			var txns_index_map = make(map[string]int)
			for index,tempTxn := range claim.Transactions{
				if (tempTxn.TransactionID == adjustTxnID){
					fmt.Println("Previous Txn to be adjusted is found. Get indexes of Txns")
					txns_index_map[tempTxn.AccumType] = index
				}
			}
			//Verify if the claim submitted with correct Accums as previous claim
			//in case of IFOOP claim submitted, only 2 txns should be there in with same txnid
			if (txns_index_map != nil  && len(txns_index_map) > 2 ){
				fmt.Println("This -ve / adjustment claim is submitted wrongly ")
				claim.Status =data.ClaimError
			}else {
				customerAllBalance.F_OOP_Balance = customerAllBalance.F_OOP_Balance + claimAmout
				customerAllBalance.I_OOP_Balance = customerAllBalance.I_OOP_Balance + claimAmout
				claim.Status =data.ClaimProcessed
				_,err := UpdateCustomerAndFamilyBalances(customerAllBalance, stub)
				if (err !=nil ){
					fmt.Println("Error updating balances ", err)
				}
				updateCarrierBalances(claim.MemberID,claim.ContractID,0,claimAmout,claim.Participant,stub)
				//NOTE: basaed on 15/02 discussion add - find old txn which is adjusted with -ve value and mark the status resolved
				for _,txnIndex := range txns_index_map{
					claim.Transactions[txnIndex].Status =data.ClaimAdjusted
				}

				_,err = MarkLastFamTxnReview(true,claim.ContractID,claim.MemberID,claim.ClaimID,contractDefinition,customerLimit,stub )
				if (err != nil ){
						fmt.Println("Error Marking Last txn for review  ", err)
				}
			}
			}else {
				//if No adjusted txn sent - report as error
				fmt.Println("No adjustment txnID is provided")
			 	claim.Status =data.ClaimError
		}
	} else {
		fOopUpdateFlag := false
		iOopUpdateFlag := false
		fOopUpdateFlag,fOopBalance = InvokeRule("IFOOP", claimAmout, contractDefinition.F_OOP_Limit,
			fOopBalance	, &fOopAccumTrans )
		claim.Status = fOopAccumTrans.Status
		if (fOopUpdateFlag) {
			iOopUpdateFlag,iOopBalance = InvokeRule("IIOOP", claimAmout, customerLimit.I_OOP_Limit,
				iOopBalance , &iOopAccumTrans )
			if (iOopUpdateFlag) {
				customerAllBalance.F_OOP_Balance = fOopBalance
				customerAllBalance.I_OOP_Balance = iOopBalance
				_,err := UpdateCustomerAndFamilyBalances(customerAllBalance, stub)
				if (err !=nil ){
					fmt.Println("Error updating balances ", err)
				}
				updateCarrierBalances(claim.MemberID,claim.ContractID,0,claimAmout,claim.Participant,stub)
			}else {
				claim.Status =iOopAccumTrans.Status
			}
		}
	}
	fOopAccumTrans.Status = claim.Status
	iOopAccumTrans.Status = claim.Status
	fOopAccumTrans.AccumBalance =	customerAllBalance.F_OOP_Balance
	iOopAccumTrans.AccumBalance =	customerAllBalance.I_OOP_Balance
	claim.Transactions = append(claim.Transactions,fOopAccumTrans)
	claim.Transactions = append(claim.Transactions,iOopAccumTrans)
	_,err = createUpdateClaim(claim,newClaim,stub)
	if err != nil {
		fmt.Println("createUpdateClaim failed")
		return false, err
	}
	lastTxn.Transaction =fOopAccumTrans
	//if (claim.Status != data.ClaimError){
	_,err = UpdateLastFamMemTxn(lastTxn,stub)
	if (err != nil ){
		fmt.Println("Error adding to Last Family Txns")
	}
	//}
	fmt.Println("In services.ProcessIFOOPClaim end ")
	return true,nil
}

// Proces the IIDED & IIOOP claim from Individual contract
func ProcessIIDEDClaim(newClaim bool,accums []data.Accum,claim data.Claims,
		customerLimit data.CustomerLimit, txnID int, adjustTxnID int, stub shim.ChaincodeStubInterface)(bool, error){
	fmt.Println("In services.ProcessIIDEDClaim start ")
	var customerBalance data.CustomerBalance
	customerBalance,err := query.GetCustomerBalance(claim.MemberID,claim.ContractID,stub)
	if err != nil {
		fmt.Println("Error retrieving Customer Balance " , claim.MemberID )
		return false, errors.New("Error retrieving Customer Balance " + claim.MemberID)
	}
	fmt.Println("Customer Balance received : " ,customerBalance)
	iDedAccumTrans := data.Transaction{AccumType:"IIDED",
		TransactionDate:time.Now(),TxnUpdatedDate:time.Now(),TransactionID:txnID}
	iOopAccumTrans := data.Transaction{AccumType:"IIOOP",
		TransactionDate:time.Now(),TxnUpdatedDate:time.Now(),TransactionID:txnID}
	for _, accum := range accums {
		if (accum.Type == "IIDED"){
			iDedAccumTrans.AccumAmount = accum.Amount
			iOopAccumTrans.AccumAmount = accum.Amount
		}
		if (accum.Type == "IIOOP"){
			iOopAccumTrans.AccumAmount = accum.Amount
		}
	}
	iDedBalance := customerBalance.I_Ded_Balance
	iOopBalance := customerBalance.I_OOP_Balance

	lastTxn :=data.LastTxnOfContract{ClaimID:claim.ClaimID,ContractID:claim.ContractID,
		MemberID:claim.MemberID,SubscriberID:claim.SubscriberID, Source:claim.Participant}

	if ((iDedAccumTrans.AccumAmount < 0 || iOopAccumTrans.AccumAmount < 0) ){
		if (adjustTxnID > 0) {
			var txns_index_map = make(map[string]int)
			for index,tempTxn := range claim.Transactions{
				if (tempTxn.TransactionID == adjustTxnID){
					fmt.Println("Previous Txn to be adjusted is found. Get indexes of Txns")
					txns_index_map[tempTxn.AccumType] = index
				}
			}
			//Verify if the claim submitted with correct Accums as previous claim
			//in case of IFOOP claim submitted, only 2 txns should be there in with same txnid
			if (txns_index_map != nil ){
				 if (len(txns_index_map) < 2){
					 fmt.Println("This -ve / adjustment claim is submitted wrongly ")
		 			claim.Status =data.ClaimError
				 }else {
					 fmt.Printf("Adjustment claim has come in with amout %v.  Processing .... \n",iDedAccumTrans)
			 		customerBalance.I_Ded_Balance = customerBalance.I_Ded_Balance  + iDedAccumTrans.AccumAmount
			 		customerBalance.I_OOP_Balance = customerBalance.I_OOP_Balance  + iOopAccumTrans.AccumAmount

			 		claim.Status = data.ClaimProcessed
			 		_,err := updateBalances(customerBalance,stub)
			 		if err != nil {
			 			fmt.Println("updateBalances failed")
			 			return false, err
			 		}
			 		updateCarrierBalances(claim.MemberID,claim.ContractID,iDedAccumTrans.AccumAmount,iOopAccumTrans.AccumAmount,claim.Participant,stub)
			 		//NOTE: basaed on 15/02 discussion add - find old txn which is adjusted with -ve value and mark the status resolved
					for _,txnIndex := range txns_index_map{
						claim.Transactions[txnIndex].Status =data.ClaimAdjusted
					}

			 		_,err = MarkLastIndTxnReview(false,claim.ContractID,claim.MemberID,claim.ClaimID,customerLimit,stub )
			 		if (err != nil ){
			 				fmt.Println("Error Marking Last txn for review  ", err)
			 		}
				 }
				}
			}else {
				//if No adjusted txn sent - report as error
				fmt.Println("No adjustment  txnID is provided")
			 	claim.Status = data.ClaimError
		}
	} else {
		// Follow processing business rules
		iDedUpdateFlag := false
		iOopUpdateFlag := false
		iDedUpdateFlag,iDedBalance = InvokeRule("IIDED", iDedAccumTrans.AccumAmount, customerLimit.I_Ded_Limit,
			iDedBalance, &iDedAccumTrans )
		claim.Status = iDedAccumTrans.Status
		if (iDedUpdateFlag){
			iOopUpdateFlag,iOopBalance = InvokeRule("IIOOP", iOopAccumTrans.AccumAmount, customerLimit.I_OOP_Limit,
				iOopBalance , &iOopAccumTrans )

			if (iOopUpdateFlag){
				fmt.Println("Both IIDED and IIOP are true")
				customerBalance.I_Ded_Balance = iDedBalance
				customerBalance.I_OOP_Balance = iOopBalance
				_,err := updateBalances(customerBalance,stub)
				if err != nil {
					fmt.Println("updateBalances failed")
					return false, err
				}
				fmt.Println("claim.MemberID,claim.ContractID " ,claim.MemberID,claim.ContractID)
				updateCarrierBalances(claim.MemberID,claim.ContractID,iDedAccumTrans.AccumAmount,iOopAccumTrans.AccumAmount,claim.Participant,stub)
			}else {
				fmt.Println("IIDED - true and  IIOOP - false. Tx status : ",iOopAccumTrans.Status)
				claim.Status = iOopAccumTrans.Status
			}
		}
	}
	fmt.Println("In services.ProcessIIDEDClaim Claim status : ",claim.Status)
	iDedAccumTrans.Status = claim.Status
	iOopAccumTrans.Status = claim.Status
	iDedAccumTrans.AccumBalance =	customerBalance.I_Ded_Balance
	iOopAccumTrans.AccumBalance =	customerBalance.I_OOP_Balance
	claim.Transactions = append(claim.Transactions,iDedAccumTrans)
	claim.Transactions = append(claim.Transactions,iOopAccumTrans)
	_,err = createUpdateClaim(claim,newClaim,stub)
	if err != nil {
		fmt.Println("createUpdateClaim failed")
		return false, err
	}
	lastTxn.Transaction = iDedAccumTrans

	//if (claim.Status != data.ClaimError) {
	_,err =	updateLastIndTxn(lastTxn,stub)
	if (err != nil ){
		fmt.Println("Error adding Last Ind Txns")
	}
	//}
	fmt.Println("In services.ProcessIIDEDClaim end ")
	return true,nil
}

// Proces the IIOOP claim from Individual contract
func ProcessIIOOPClaim(newClaim bool,claimAmout float64,claim data.Claims,
		customerLimit data.CustomerLimit,txnID int,adjustTxnID int, stub shim.ChaincodeStubInterface)(bool, error){
	fmt.Println("In services.ProcessIIOOPClaim start ")
	var customerBalance data.CustomerBalance
	customerBalance,err := query.GetCustomerBalance(claim.MemberID,claim.ContractID,stub)
	if err != nil {
		fmt.Println("Error retrieving Customer Balance " , claim.MemberID )
		return false, errors.New("Error retrieving Customer Balance " + claim.MemberID)
	}
	fmt.Println("Customer Balance received : " ,customerBalance)
	iOopAccumTrans := data.Transaction{AccumType:"IIOOP", AccumAmount:claimAmout,
		TxnUpdatedDate:time.Now(),TransactionDate:time.Now(),TransactionID:txnID}
	iOopBalance := customerBalance.I_OOP_Balance
	lastTxn :=data.LastTxnOfContract{ClaimID:claim.ClaimID,ContractID:claim.ContractID,
		MemberID:claim.MemberID,SubscriberID:claim.SubscriberID, Source:claim.Participant}
 	if (claimAmout <  0 ) {
		if (adjustTxnID > 0){
			//Added based on review feedback for -ve testing
			var txns_index_map = make(map[string]int)
			for index,tempTxn := range claim.Transactions{
				if (tempTxn.TransactionID == adjustTxnID){
					fmt.Println("Previous Txn to be adjusted is found. Get indexes of Txns")
					txns_index_map[tempTxn.AccumType] = index
				}
			}
			fmt.Println("Verify if the claim submitted with correct Accums as previous claim")
			if (txns_index_map != nil  && len(txns_index_map) > 1 ){
				fmt.Println("This -ve / adjustment claim is submitted wrongly ")
				claim.Status =data.ClaimError
			} else{
				//Submitted -ve claim correctly hence proceed further to adjust
				fmt.Println("Adjustment claim has come in with amout %v.  Processing .... \n",claimAmout)
				customerBalance.I_OOP_Balance = customerBalance.I_OOP_Balance + claimAmout
				claim.Status =data.ClaimProcessed
				_,err := updateBalances(customerBalance, stub)
				if (err !=nil ){
					fmt.Println("Error updating balances ", err)
				}
				fmt.Println("claim.MemberID,claim.ContractID ",claim.MemberID,claim.ContractID)
				updateCarrierBalances(claim.MemberID,claim.ContractID,0,claimAmout,claim.Participant,stub)
				//NOTE: basaed on 15/02 discussion add - find old txn which is adjusted with -ve value and mark the status resolved
				for _,txnIndex := range txns_index_map{
					claim.Transactions[txnIndex].Status =data.ClaimAdjusted
				}
				
				_,err = MarkLastIndTxnReview(true,claim.ContractID,claim.MemberID,claim.ClaimID,customerLimit,stub )
				if (err != nil ){
						fmt.Println("Error Marking Last txn for review  ", err)
				}
			}
		}else {
			//if No adjusted txn sent - report as error
			fmt.Println("No adjustment  txnID is provided")
		 claim.Status =data.ClaimError
		}
	} else {
		iOopUpdateFlag,iOopBalance := InvokeRule("IIOOP", claimAmout, customerLimit.I_OOP_Limit,
			iOopBalance , &iOopAccumTrans )
		claim.Status = iOopAccumTrans.Status
		fmt.Println("Tx status : ",iOopAccumTrans.Status)
		if (iOopUpdateFlag){
			customerBalance.I_OOP_Balance = iOopBalance
			_,err := updateBalances(customerBalance, stub)
			if (err !=nil ){
				fmt.Println("Error updating balances ", err)
			}
			fmt.Println("claim.MemberID,claim.ContractID ", claim.MemberID,claim.ContractID)
			updateCarrierBalances(claim.MemberID,claim.ContractID,0,claimAmout,claim.Participant,stub)
		}
	}
	iOopAccumTrans.Status = claim.Status
	iOopAccumTrans.AccumBalance =	customerBalance.I_OOP_Balance
	claim.Transactions = append(claim.Transactions,iOopAccumTrans)
	_,err = createUpdateClaim(claim,newClaim,stub)
	if err != nil {
		fmt.Println("createUpdateClaim failed")
		return false, err
	}

	lastTxn.Transaction =iOopAccumTrans
	//if (claim.Status != data.ClaimError) {
	_,err =	updateLastIndTxn(lastTxn,stub)
	if (err != nil ){
		fmt.Println("Error adding Last Ind Txns")
	}
	//}
	fmt.Println("In services.ProcessIIOOPClaim end ")
	return false, nil
}

// Create New Claim or Update the existing claim with new Txn or Txn status modifications
func createUpdateClaim(claim data.Claims, newClaim bool,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.createUpdateClaim start ")
	claim.LastUpdatedDTTM = time.Now()
	_, err := CreateClaims(claim,stub)
	if err != nil {
		fmt.Println("Creating claims failed  " )
		return nil, errors.New("Error creating claims of contract  "+claim.ContractID )
	}
	if newClaim {
		_,err := AddClaimIDToMemContClaims(claim.MemberID, claim.ContractID ,claim.ClaimID , stub)
		if err != nil {
			fmt.Println("Error adding  claims to Contract Claims List  " )
			return nil, errors.New("Error adding  claims to Contract Claims List  "+claim.ContractID )
		}
	}
	fmt.Println("In services.createUpdateClaim end ")
	return nil, nil
}

// Update the Balances IIDED and IIOOP of the Customer
func updateBalances(customerBalance data.CustomerBalance,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.updateBalances start ")
	_, err := CreateCustomerBalance(customerBalance,stub)
	if err != nil {
		fmt.Println("Updating Customer Balance " )
	}
	fmt.Println("In services.updateBalances end ")
	return nil, nil
}

//Update Insurance Carrier balances. For UI display
func updateCarrierBalances(memberID string, contractID string, iDedAccumAmt float64, iOppAccumAmt float64,
	source string,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.updateCarrierBalances start ")

	fmt.Printf("ContractID : %v , MemberID : %v \n", contractID,memberID)
	var contMemCarrierBalance data.ContMemCarrierBalance
	var carrierBalance data.CarrierBalance
	contMemCarrierBalance,err := query.GetCarrierBalances(memberID,contractID,stub)
	if (err != nil ){
		fmt.Println("No Records found. initialize the carrier balances")
		contMemCarrierBalance.ContractID = contractID
		contMemCarrierBalance.MemberID = memberID
		contMemCarrierBalance.CarrierBalances = make(map[string]data.CarrierBalance)
		carrierBalance.CarrierName = strings.Title(strings.ToLower(source))//strings.ToUpper(source)
		carrierBalance.AccumBalance  = make(map[string]float64)
		carrierBalance.AccumBalance["IIDED"] = iDedAccumAmt
		carrierBalance.AccumBalance["IIOOP"] = iOppAccumAmt
		contMemCarrierBalance.CarrierBalances[strings.ToUpper(source)] = carrierBalance
	}else {
		fmt.Println("Carrier balances found for Member.")
		carrierBalance = contMemCarrierBalance.CarrierBalances[strings.ToUpper(source)]
		if (carrierBalance.CarrierName == ""){
			carrierBalance.AccumBalance  = make(map[string]float64)
			carrierBalance.CarrierName = strings.Title(strings.ToLower(source))//strings.ToUpper(source)
		}
		fmt.Printf("Carrier balances of Source %v is : %v \n",source,carrierBalance)
		carrierBalance.AccumBalance["IIDED"] = carrierBalance.AccumBalance["IIDED"] + iDedAccumAmt
		carrierBalance.AccumBalance["IIOOP"] = carrierBalance.AccumBalance["IIOOP"] + iOppAccumAmt
		contMemCarrierBalance.CarrierBalances[strings.ToUpper(source)] = carrierBalance
	}
	fmt.Println("Carrier Balances to be updated : ", contMemCarrierBalance)
	_,err = CreateCarrierBalances(contMemCarrierBalance, stub)
	if err != nil {
		fmt.Println("Error updating Member Carier  Balances")
	}
	fmt.Println("In services.updateCarrierBalances end")
	return nil, nil
}

// Find the last/previous txn from Family contract after the negative txn processed and update the status to Review - Adjustment
func MarkLastFamTxnReview(isOppClaimAdjusted bool, contractID string, memberID string, claimID string,
	contractDefinition data.ContractDefinition,
	customerLimit data.CustomerLimit,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.MarkLastFamTxnReview start ")
		var lastTxns []data.LastTxnOfContract
		lastTxns, err := query.GetLastFamilyTxns(contractID,stub)
		if (err != nil ){
			fmt.Println("No Last txns found . ")
		}
		if len(lastTxns) > 0 {
			for _, lastTxn := range lastTxns {
				fmt.Println("Found last txn is from claim    : ",lastTxn.ClaimID)
				if (claimID == lastTxn.ClaimID  &&  memberID == lastTxn.MemberID ){
					fmt.Printf("Adjustment has come for this txn from %v. Hence no need to review \n",memberID )
					continue
				}
				/** NOTE: Adding this to deal with  ERROR status txn - start **/
				if (lastTxn.Transaction.Status == data.ClaimError){
					fmt.Println("Last ind txn is error txn. Hence get the previous txn which is not of Error status")
					var allTxns []data.LastTxnOfContract
					allTxns,err := query.GetTxnsOfMemCont(contractID,"",memberID,"",0,stub)
					if (err != nil){
							fmt.Printf("Error getting  txns for contract %v , member %v \n", contractID,memberID)
							return nil,errors.New("Error getting txns of contract "+contractID +" for Member " +memberID)
					}
					for _,tempTxn := range allTxns{
						if ((lastTxn.Transaction.TransactionID == tempTxn.Transaction.TransactionID)||
								(tempTxn.Transaction.Status == data.ClaimError)){
							continue
						}
						fmt.Printf("Got the last family txn of member %v which is not an Error txn \n ",lastTxn.MemberID)
						lastTxn.ClaimID = tempTxn.ClaimID
						lastTxn.Transaction = tempTxn.Transaction
						lastTxn.Source = tempTxn.Source
						break
					}
				}
				/**  Adding this to deal with  ERROR status txn - end **/
				fmt.Println("Last Calim with txn need to be verified after balances adjustment is ", lastTxn)
				_,err =UpdateTxnStatusAfterAdj(isOppClaimAdjusted,lastTxn,contractDefinition,customerLimit,stub)
				if (err != nil ){
					fmt.Println("Error Updating the status after last txn after adjustment. ")
				}
			}
		}
		fmt.Println("In services.MarkLastFamTxnReview end ")
		return nil, nil
}

// Find the last/previous txn from Individual contract after the negative txn processed and update the status to Review - Adjustment
func MarkLastIndTxnReview(isOppClaimAdjusted bool,contractID string, memberID string, claimID string,
	customerLimit data.CustomerLimit,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.markLastIndTxnReview start ")
	lastTxn, err := query.GetLastIndTxn(contractID,stub)
	if (err != nil ){
		fmt.Println("No Last txn found. ")
	}
	fmt.Println("Last ind txn for contract  ",lastTxn )
	/** NOTE: Adding this to deal with  ERROR status txn - start **/
	//In case last txn is error txn
	if (lastTxn.Transaction.Status == data.ClaimError){
		fmt.Println("Last ind txn is error txn. Hence get the previous txn which is not of Error status")
		var allTxns []data.LastTxnOfContract
		allTxns,err := query.GetTxnsOfMemCont(contractID,"",memberID,"",0,stub)
		if (err != nil){
				fmt.Printf("Error getting  txns for contract %v , member %v \n", contractID,memberID)
				return nil,errors.New("Error getting txns of contract "+contractID +" for Member " +memberID)
		}
		for _,tempTxn := range allTxns{
			if ((lastTxn.Transaction.TransactionID == tempTxn.Transaction.TransactionID)||
					(tempTxn.Transaction.Status == data.ClaimError)){
				continue
			}
			lastTxn.ClaimID = tempTxn.ClaimID
			lastTxn.Transaction = tempTxn.Transaction
			lastTxn.Source = tempTxn.Source
			fmt.Println("Got the last txn which is not an Error txn : ", tempTxn)
			break
		}
	}
	/**  Adding this to deal with  ERROR status txn - end **/
	fmt.Println("Last Calim with txn need to be verified after balances adjustment is ", lastTxn)
	var contractDefinition data.ContractDefinition
	_,err =UpdateTxnStatusAfterAdj(isOppClaimAdjusted,lastTxn,contractDefinition,customerLimit,stub)
	if (err != nil ){
		fmt.Println("Error Updating the status after last txn after adjustment. ")
	}
	fmt.Println("In services.markLastIndTxnReview end ")
	return nil,nil
}

// Update the previous Individual Txn status from Processed to Review - Adjustment
func UpdateTxnStatusAfterAdj(isOppClaimAdjusted bool,lastTxn data.LastTxnOfContract,contractDefinition data.ContractDefinition,
		customerLimit data.CustomerLimit,stub shim.ChaincodeStubInterface)([]byte, error){
	fmt.Println("In services.UpdateTxnStatusAfterAdj start ")
	claim,err := query.GetClaim(lastTxn.MemberID,lastTxn.ContractID,lastTxn.ClaimID,stub)
	claimTxns := claim.Transactions
	updateFlag := false
	fmt.Printf("Last Claim %v found with txnID %v \n",lastTxn.ClaimID,lastTxn.Transaction.TransactionID)
	for index, txn := range claimTxns {
		fmt.Printf("Claim %v Transaction with txn :  %v \n",claim.ClaimID,txn)
		if (lastTxn.Transaction.TransactionID == txn.TransactionID) {
			fmt.Printf("Found Transaction with txnID %v  in claims History : \n",txn.TransactionID)
			if (!isOppClaimAdjusted && txn.AccumType == "IFDED"){
				fmt.Println("Found IFDED")
				if (txn.AccumBalance == contractDefinition.F_Ded_Limit && txn.Status == data.ClaimProcessed){
					fmt.Println("Updating IFDED txn to Review")
					 txn.Status = data.ClaimReviewAdjustment
					 txn.TxnUpdatedDate = time.Now()
					 claim.Transactions[index] = txn
					 updateFlag = true
				 }
			} else if ( txn.AccumType == "IFOOP"){
				fmt.Println("Found IFOOP")
				 if (txn.AccumBalance == contractDefinition.F_OOP_Limit && txn.Status == data.ClaimProcessed){
					fmt.Println("Updating IFOOP txn to Review")
					 txn.Status = data.ClaimReviewAdjustment
					 txn.TxnUpdatedDate = time.Now()
					 claim.Transactions[index] = txn
					 updateFlag = true
				 }
			} else if (!isOppClaimAdjusted && txn.AccumType == "IIDED"){
				fmt.Println("Found IIDED")
				 if (txn.AccumBalance == customerLimit.I_Ded_Limit && txn.Status == data.ClaimProcessed){
					 fmt.Println("Updating IIDED txn to Review")
					 txn.Status = data.ClaimReviewAdjustment
					 txn.TxnUpdatedDate = time.Now()
					 claim.Transactions[index] = txn
					 updateFlag = true
				 }
			} else if (txn.AccumType == "IIOOP"){
				fmt.Println("Found IIOOP")
				 if (txn.AccumBalance == customerLimit.I_OOP_Limit && txn.Status == data.ClaimProcessed){
					 fmt.Println("Updating IIOOP txn to Review")
					 txn.Status = data.ClaimReviewAdjustment
					 txn.TxnUpdatedDate = time.Now()
					 claim.Transactions[index] = txn
					 updateFlag = true
				 }
			}
		}
	}
	if (updateFlag){
		fmt.Println("Tx status has been modified.")
		claim.LastUpdatedDTTM = time.Now()
		_, err = CreateClaims(claim,stub)
		if err != nil {
			fmt.Println("Updating claim with transaction %s status marked for review \n",lastTxn.Transaction.TransactionID )
			return nil, errors.New("Error updating claims of contract  "+claim.ContractID )
		}
	}
	fmt.Println("In services.UpdateTxnStatusAfterAdj end ")
	return nil, nil
}

// Adjust the Individual Ded and OOP Limits of the Contract
func AdjustLimitsOfContract(args []string, stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.AdjustLimitsOfContract start ")
	subscriberID :=args[0]
	memberID :=args[1]
	contractID  :=args[2]
	var accumLimits []data.Accum
	err := json.Unmarshal([]byte(args[3]), &accumLimits)
	fmt.Println("accumLimits are ", accumLimits)
	if err!= nil {
		fmt.Println("Error unmarshal accumtypes  ", err)
	}
	var customerLimit data.CustomerLimit
	customerLimit,err = query.GetCustomerLimits(memberID,contractID,stub)
	if err !=nil {
		fmt.Println("Error retrieving Customer Limits " , memberID )
		return nil, errors.New("Error retrieving Customer Limits " + memberID)
	}
	fmt.Println("Customer Limits received : " ,customerLimit)

	contractDefinition,err := query.GetContract(contractID,stub)
	if err != nil {
		fmt.Println("Error receiving  the Customer Contract")
		return nil, errors.New("Error receiving  the Customer Contract")
	}
	if (contractDefinition.ContractType =="I"){
		fmt.Println("Contract Type : I" )
		for _, accum := range accumLimits {
			if (accum.Type == "IIDED"){
				fmt.Println("Set IIDED Limit" )
				customerLimit.I_Ded_Limit = accum.Amount
			}
			if (accum.Type == "IIOOP"){
				fmt.Println("Set IIOOP Limit" )
				customerLimit.I_OOP_Limit = accum.Amount
			}
		}
	}
	//NOTE: we are not handling this currently. Kept for further enhancements

	// if (contractDefinition.ContractType == "F"){
	// 	for _, accum := range accumLimits {
	// 		if (accum.Type == "IFDED"){
	// 				contractDefinition.F_Ded_Limit = accum.Amount
	// 		}
	// 		if (accum.Type == "IIDED"){
	// 				contractDefinition.F_OOP_Limit = accum.Amount
	// 		}
	// 		if (accum.Type == "IIDED"){
	// 				customerLimit. I_Ded_Limit = accum.Amount
	// 		}
	// 		if (accum.Type == "IIOOP"){
	// 				customerLimit. I_OOP_Limit = accum.Amount
	// 		}
	// 	}
	// 	_,err = CreateContractDefinition(contractDefinition,stub )
	// 	if err !=nil {
	// 		fmt.Println("Error Updating the limits of Contract  " , contractID )
	// 	}
	// }
	var membersOfContract data.ContractMembers
	membersOfContract,err = query.GetMembersOfContract(contractID,stub)
	if err !=nil {
		fmt.Println("Error retrieving Members of Contract  " , contractID )
		//return nil, errors.New("Error retrieving Members of Contract " + contractID)
	}
	memberIDs := membersOfContract.MemberIDs
	if (len(memberIDs) > 0 ){
		for _,memID := range memberIDs {
			customerLimit.MemberID = memID
			_,err = CreateCustomerLimit(customerLimit,stub )
			if err != nil {
				fmt.Println("Error updating the customer limits for member :  "+memID)
				return nil, errors.New("Error updating the Customer Limits")
			}
			fmt.Println("Limits updated. Check if any previous txns of the member effected due to limits ajdusted " )
			// Mark all previous txns till the balances fall in the limits after limits adjusted
			var lastTxns []data.LastTxnOfContract
			lastTxns, err = query.GetTxnsOfMemCont(contractID,subscriberID,memberID,"",0,stub)
			fmt.Println("No of previous txns for the member " ,len(lastTxns))
			//Create a map to add the already adjusted txnid of claim.
			//This helps not to adjust the same claim txn again
			claimTxns := make(map[string][]int)
			txnAlreadyAdj := false
			updateFlag := true

			if (len(lastTxns)>0){
				for _,txn := range lastTxns {
					txnAlreadyAdj = false
					if ( claimTxns[txn.ClaimID] != nil ){
						// search if txn is there in the list
						txnIDs := claimTxns[txn.ClaimID]
						fmt.Println("txnIDs : ",txnIDs)
						for _,txnID := range txnIDs{
							if (txnID ==   txn.Transaction.TransactionID){
								txnAlreadyAdj = true
								break
							}
						}
					}
					if !txnAlreadyAdj {
						//call UpdateTxnStatusAfterLimitAdj
						updateFlag,err = UpdateTxnStatusAfterILimitAdj(txn,contractDefinition,customerLimit,stub)
						fmt.Println("Is txn updated because limits adjustment is : ", updateFlag)
						fmt.Println("Txn is : ", txn)
						//add to map so that we can check that  this txn is already adjusted
						claimTxns[txn.ClaimID] = append(claimTxns[txn.ClaimID], txn.Transaction.TransactionID)
					}
				}
				//NOTE: Checking the scenarios where few previous txns are errors,
				// but there are more txns down which has balance is > new limit values. To handle this,
				// will go through all previous txns
				// if (updateFlag == false ){
				// 	break
				// }
			}
		}
	}
	fmt.Println("In services.AdjustLimitsOfContract end ")
	return nil,nil
}

//Update status of the previous Txns after the Individual Contract Limits updated
func UpdateTxnStatusAfterILimitAdj(lastTxn data.LastTxnOfContract,contractDefinition data.ContractDefinition,
		customerLimit data.CustomerLimit,stub shim.ChaincodeStubInterface)(bool, error){
	fmt.Println("In services.UpdateTxnStatusAfterILimitAdj start ")
	claim,err := query.GetClaim(lastTxn.MemberID,lastTxn.ContractID,lastTxn.ClaimID,stub)
	fmt.Printf("Last Claim %v found with txnID %v \n",lastTxn.ClaimID,lastTxn.Transaction.TransactionID)
	var txns_index_map = make(map[string]int)
	for index, txn := range claim.Transactions {
		fmt.Printf("Claim %v Transaction with txn : \n",claim.ClaimID,txn)
		if (lastTxn.Transaction.TransactionID == txn.TransactionID) {
			txns_index_map[txn.AccumType] = index
		}
	}
	fmt.Println("%v of txns found with txn id : %v ", len(txns_index_map) ,lastTxn.Transaction.TransactionID)
	reviewOOPFlag := false
	reviewDEDFlag := false
	if (txns_index_map!= nil && len(txns_index_map) == 2) {
		fmt.Println("IIDED claim. So verify both IIDED and IIOOP Txns. ")
		index := txns_index_map["IIOOP"]
		if (claim.Transactions[index].AccumBalance > customerLimit.I_OOP_Limit){
			fmt.Println("IIOOP Balance is >=  limit ")
			reviewOOPFlag = true
		}
		index = txns_index_map["IIDED"]
	  if (claim.Transactions[index].AccumBalance > customerLimit.I_Ded_Limit){
				fmt.Println("IIDED Balance is >=  limit ")
				reviewDEDFlag = true
		}
		if (reviewDEDFlag ||  reviewOOPFlag ){
			fmt.Printf("Set txns of claim %v of txns of ID %v as Review – Contract \n",claim.ClaimID,claim.Transactions[txns_index_map["IIDED"]].TransactionID )
			claim.Transactions[txns_index_map["IIDED"]].Status = data.ClaimReviewContract
			claim.Transactions[txns_index_map["IIDED"]].TxnUpdatedDate = time.Now()
			claim.Transactions[txns_index_map["IIOOP"]].Status = data.ClaimReviewContract
			claim.Transactions[txns_index_map["IIOOP"]].TxnUpdatedDate = time.Now()
		}
	}
	if (txns_index_map!= nil && len(txns_index_map) == 1) {
		index := txns_index_map["IIOOP"]
		fmt.Println("IIOOP claim.")
		if (claim.Transactions[index].AccumBalance > customerLimit.I_OOP_Limit){
			fmt.Printf("Set txns of claim %v of txns of ID %v as Review – Contract \n",claim.ClaimID,claim.Transactions[txns_index_map["IIOOP"]].TransactionID )
			claim.Transactions[txns_index_map["IIOOP"]].Status = data.ClaimReviewContract
			claim.Transactions[txns_index_map["IIOOP"]].TxnUpdatedDate = time.Now()
			reviewOOPFlag = true
		}
	}
	if (reviewDEDFlag || reviewOOPFlag){
		fmt.Println("Tx status has been modified. ")
		claim.LastUpdatedDTTM = time.Now()
		_, err = CreateClaims(claim,stub)
		if err != nil {
			fmt.Println("Updating claim with transaction %s status marked for review \n",lastTxn.Transaction.TransactionID )
			return false, errors.New("Error updating claims of contract  "+claim.ContractID )
		}
	}
	fmt.Println("In services.UpdateTxnStatusAfterILimitAdj end ")
	return true,nil
}

//Updates Txn status from Review - Adjustment to Processed
func UpdateTxnStatus(args []string,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.UpdateTxnStatus start ")
	subscriberID :=args[0]
	memberID :=args[1]
	contractID  :=args[2]

	var accumTypes []data.Accum
	err := json.Unmarshal([]byte(args[3]), &accumTypes)
	fmt.Println("accumTypes are ", accumTypes)
	if err!= nil {
		fmt.Println("Error unmarshal accumtypes  ", err)
	}
	claimID  :=args[4]
	txnID, err := strconv.Atoi(args[5])

	fmt.Println("Input params :")
	fmt.Printf("SubscriberID : %v , memberID %v ,contractID %v,claimID %v,txnID %v \n ",subscriberID,memberID,contractID,claimID,txnID)
	claim,err := query.GetClaim(memberID,contractID,claimID,stub)
	var txns_index_map = make(map[string]int)
	for index, txn := range claim.Transactions {
		if (txn.TransactionID ==  txnID) {
			txns_index_map[txn.AccumType] = index
		}
	}
	for _,accum := range accumTypes{
		//Review – Adjustment changes to  Adjusted
		 claim.Transactions[txns_index_map[accum.Type]].Status =data.ClaimProcessed
		 claim.Transactions[txns_index_map[accum.Type]].TxnUpdatedDate = time.Now()
	}
	_,err = CreateClaims(claim,stub)
	if (err != nil){
		fmt.Println("Error Updting the txn")
	}
	fmt.Println("In services.UpdateTxnStatus end ")
	return nil,nil
}

func AddMemberRelation(subscriberID string, memberID string, relation string,stub shim.ChaincodeStubInterface) ([]byte, error){
	fmt.Println("In services.AddMemberRelation start ")
	var memberRelations []data.MemberRelation
	memberRelations,err := query.GetSubscriberMemberRelation(subscriberID, stub)
	if err != nil {
		fmt.Println("Error receiving  the Customer Relationships")
		return nil, errors.New("Error receiving  the Customer Relationships")
	}

	if (len(memberRelations)>0){
		for _,memRel := range memberRelations {
			if (memRel.MemberID == memberID ){
				fmt.Printf("Relationship between subscriber %v and member %v is already exists \n",subscriberID,memberID)
				return nil, errors.New("Relationship between subscriber "+subscriberID + "and member "+memberID +" already exists")
			}
		}
	}
	fmt.Printf("Adding Relationship between subscriber %v and member %v  \n",subscriberID,memberID)
	var memRelation data.MemberRelation
	memRelation.SubscriberID = subscriberID
	memRelation.MemberID = memberID
	memRelation.Relationship = relation
	memberRelations = append(memberRelations, memRelation)
	_,err = CreateMemRelation(subscriberID,memberRelations,stub)
	if (err != nil ){
		fmt.Println("Error adding member Relatonship ")
		return nil,err
	}

	fmt.Println("In services.AddMemberRelation end ")
	return nil,nil
}
func SwitchPolicy(args []string,stub shim.ChaincodeStubInterface) ([]byte, error){
	fmt.Println("In services.SwitchPolicy start ")
	memberID :=args[0]
	contractID  :=args[1]
	newContractID := args[2]
	relation := args[3]
	fmt.Printf("memberID: %v , contractID: %v , newContractID: %v , relation: %v \n",memberID,contractID,newContractID,relation)

	// 1. Get New Contract Details so that we have subscriberID
	contractDefinition,err := query.GetContract(newContractID,stub)
	if err != nil {
		fmt.Println("Error receiving  the Customer Contract")
		return nil, errors.New("Error receiving  the Customer Contract")
	}
	fmt.Println("New Contract :", contractDefinition)

	//2. Add Relationship between subscriber and new member
	_,err = AddMemberRelation(contractDefinition.SubscriberID, memberID,relation,stub)
	if (err != nil){
		fmt.Printf("Unable to add Relatonship between subscriber %v and member %v  \n", contractDefinition.SubscriberID, memberID)
		fmt.Println("Error :" , err)
	}

	//3. Get IIDED and IIOOP  limits of the subscriber. The same limits applies to the new member
	var customerLimit data.CustomerLimit
	customerLimit,err = query.GetCustomerLimits(memberID,contractID,stub)
	if err !=nil {
		fmt.Println("Error retrieving Customer Limits " , memberID )
		return nil, errors.New("Error retrieving Customer Limits " + memberID)
	}
	fmt.Println("Customer Limits received : " ,customerLimit)

	//4. Add member to New Contract
	AddMemberToContract(newContractID,memberID,customerLimit.I_Ded_Limit,customerLimit.I_OOP_Limit,stub)

	//5. Get balances from Old contract
	var customerBalance data.CustomerBalance
	customerBalance,err = query.GetCustomerBalance(memberID,contractID,stub)
	if err != nil {
		fmt.Printf("Error retrieving Member Balance %v from contract %v \n " , memberID , contractID)
		return nil, errors.New("Error retrieving Customer Balance " + memberID)
	}
	fmt.Println("Old contract Balances  : ", customerBalance)

	//6. Update balances from old contract to New contract
	//Check if newContract is Family , update individual and Family
	if (contractDefinition.ContractType == "F"){
		var customerAllBalance data.CustomerAllBalance
		customerAllBalance,err := query.GetAllBalancesOfCustomer(contractDefinition.SubscriberID,newContractID,stub)
		if err != nil {
			fmt.Println("Error retrieving Customer Balance " , contractDefinition.SubscriberID )
			return nil, errors.New("Error retrieving Customer Balance " + contractDefinition.SubscriberID)
		}
		fmt.Println("Customer Balance received : " ,customerAllBalance)
		customerAllBalance.ContractID = newContractID
		customerAllBalance.MemberID = memberID
		customerAllBalance.F_Ded_Balance = customerAllBalance.F_Ded_Balance + customerBalance.I_Ded_Balance
		customerAllBalance.F_OOP_Balance = customerAllBalance.F_OOP_Balance + customerBalance.I_OOP_Balance
		customerAllBalance.I_Ded_Balance = customerBalance.I_Ded_Balance
		customerAllBalance.I_OOP_Balance = customerBalance.I_OOP_Balance

		_,err = UpdateCustomerAndFamilyBalances(customerAllBalance, stub)
		if (err !=nil ){
			fmt.Println("Error updating Individual & Family balances ", err)
		}

	}else {
		customerBalance.ContractID = newContractID
		_,err = updateBalances(customerBalance, stub)
		if (err !=nil ){
			fmt.Printf("Error updating balances of member to new contract \n", err)
		}
	}

	//7. Update Carrier balances from old contract to New contract
	var contMemCarrierBalance data.ContMemCarrierBalance
	contMemCarrierBalance,err = query.GetCarrierBalances(memberID,contractID,stub)
	fmt.Println("Customer Carrier Balances received from old contract : " ,contMemCarrierBalance)
	if (err != nil ){
		fmt.Println("No Customer Carrier Balances Found")
	}else {
		contMemCarrierBalance.ContractID = newContractID
	}
	CreateCarrierBalances(contMemCarrierBalance,stub)

	// 8. Remove member from old contract
	// _,err = RemoveMemberFromContract(contractID, memberID,stub)
	// if (err != nil){
	// 		fmt.Printf("Error removing member %v from contract %v  \n ", memberID,contractID)
	// }
	fmt.Println("In services.SwitchPolicy end ")
	return nil, nil
}
/////
