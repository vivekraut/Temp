package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hcsc/claims/data"
	"github.com/hcsc/claims/query"
	"encoding/json"
)

func ProcessClaimAdjust(args []string,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.ProcessClaimAdjust start ")
	subscriberID :=args[0]
	memberID 	:=args[1]
	contractID  :=args[2]
	var accumTypes []data.Accum

	err := json.Unmarshal([]byte(args[3]), &accumTypes)
	fmt.Println("accumTypes are ", accumTypes)
	if err!= nil {
		fmt.Println("Error unmarshal accumtypes  ", err)

	}
	claimID :=args[4]
	var txnID int =0
	txnID, err  = strconv.Atoi(args[5])
	if (err != nil ){
		fmt.Println("Please pass integer as args[5]")
	}


	fmt.Println("subscriberID  ", subscriberID)
	fmt.Println("memberID  ", memberID)

	var customerLimit data.CustomerLimit
	customerLimit,err = query.GetCustomerLimits(memberID,contractID,stub)
	if err !=nil {
		fmt.Println("Error retrieving Customer Limits " , memberID )
		return nil, errors.New("Error retrieving Customer Limits " + memberID)
	}
	fmt.Println("Customer Limits received : " ,customerLimit)

	claim, err := query.GetClaim(memberID,contractID,claimID,stub)

	contractDefinition,err := query.GetContract(contractID,stub)
	if err != nil {
		fmt.Println("Error receiving  the Customer Contract")
		return nil, errors.New("Error receiving  the Customer Contract")
	}
	fmt.Println("Contract :", contractDefinition)
	if (contractDefinition.ContractType == "I"){
		fmt.Println("Recevied Individual Claim for Contract ", contractID)
		if accumTypes != nil {
			for _, accum := range accumTypes {
				if (accum.Type == "IIDED"){
					_,err :=ProcessIIDEDClaim(false,true,accumTypes,claim,customerLimit,txnID,stub)
					if (err != nil){
						fmt.Println("Error processing IIDED claims")
					}
					break
				}

				if (accum.Type == "IIOOP"){
					_,err :=ProcessIIOOPClaim(false,true,accum.Amount,claim,customerLimit,txnID,stub)
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
		if accumTypes != nil {
			for _, accum := range accumTypes {

				if (accum.Type == "IIDED"){

					_,err :=ProcessIFDEDClaim(false,accumTypes,claim, contractDefinition ,
						customerLimit ,stub)
					if (err != nil){
						fmt.Println("Error processing IFDED claims")
					}
					break
				}

				if (accum.Type == "IIOOP"){
					_,err :=ProcessIFOOPClaim(false ,accum.Amount, claim, contractDefinition ,
						customerLimit ,stub)

					if (err != nil){

						fmt.Println("Error processing IFDED claims")
					}
					break
				}
			}
		}
	}
	fmt.Println("In services.ProcessClaimAdjust end ")

	return nil,nil
}


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
	//claimAmout,err := strconv.ParseFloat(args[5], 64)
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
	fmt.Println("No of Members in contract %v is %v ", contractID, len(memberIDs))
	if (len(memberIDs) > 0 ){
		for _,memID := range memberIDs {
			fmt.Printf("memID : %v , input memberID : %v \n", memID,memberID)
			if (memID == memberID){
				memberFound = true
				fmt.Printf("Found Member %v in Contract %v ",memberID,contractID )
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
		claim.Participant =source
		claim.Status = ""
		claim.TotalClaimAmount = 0
		claim.CreateDTTM = time.Now()
		newClaim = true

		fmt.Println("New Claim : ",claim)
	}
	//TODO - update later
	var txnID int = 0

	lastTxn, err := query.GetLastIndTxn(contractID,stub)
	if (err != nil ){
		fmt.Println("No Last txn found. Set txnID as : 1")
		txnID = 1
	}else {
		txnID = lastTxn.Transaction.TransactionID + 1
	}

	fmt.Println("Last ind txn for contract  ",lastTxn )

	contractDefinition,err := query.GetContract(contractID,stub)
	if err != nil {
		fmt.Println("Error receiving  the Customer Contract")
		return nil, errors.New("Error receiving  the Customer Contract")
	}
	fmt.Println("Contract :", contractDefinition)
	if (contractDefinition.ContractType == "I"){
		fmt.Println("Recevied Individual Claim for Contract ", contractID)
		if accumTypes != nil {
			for _, accum := range accumTypes {
				if (accum.Type == "IIDED"){
					_,err :=ProcessIIDEDClaim(newClaim,false,accumTypes,claim,customerLimit,txnID,stub)
					if (err != nil){
						fmt.Println("Error processing IIDED claims")
					}
					break
				}

				if (accum.Type == "IIOOP"){
					_,err :=ProcessIIOOPClaim(newClaim,false,accum.Amount,claim,customerLimit,txnID,stub)
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
		if accumTypes != nil {
			for _, accum := range accumTypes {

				if (accum.Type == "IIDED"){

					_,err :=ProcessIFDEDClaim(newClaim ,accumTypes, claim, contractDefinition ,
						customerLimit ,stub)
					if (err != nil){
						fmt.Println("Error processing IFDED claims")
					}
					break
				}

				if (accum.Type == "IIOOP"){
					_,err :=ProcessIFOOPClaim(newClaim ,accum.Amount, claim, contractDefinition ,
						customerLimit ,stub)

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

func ProcessIFDEDClaim(newClaim bool, accums []data.Accum,claim data.Claims,
	contractDefinition data.ContractDefinition,
	customerLimit data.CustomerLimit,stub shim.ChaincodeStubInterface)(bool, error){

	fmt.Println("In services.ProcessIFDEDClaim start ")

		var customerAllBalance data.CustomerAllBalance
		customerAllBalance,err := query.GetAllBalancesOfCustomer(claim.MemberID,claim.ContractID,stub)

		if err != nil {
			fmt.Println("Error retrieving Customer Balance " , claim.MemberID )
			return false, errors.New("Error retrieving Customer Balance " + claim.MemberID)
		}

		fmt.Println("Customer Balance received : " ,customerAllBalance)

		txnID := 1
		var transLength int = len(claim.Transactions)
		fmt.Println("Claim transactions length : ",transLength)
		if (!newClaim && transLength > 0) {
			txnID = claim.Transactions[transLength-1].TransactionID + 1
		}

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

		if (fDedAccumTrans.AccumAmount < 0 ){
			fmt.Printf("Adjustment claim has come in with amout %v.  Processing .... \n",fDedAccumTrans.AccumAmount)
			customerAllBalance.F_Ded_Balance = customerAllBalance.F_Ded_Balance + fDedAccumTrans.AccumAmount
			customerAllBalance.F_OOP_Balance = customerAllBalance.F_OOP_Balance + fOopAccumTrans.AccumAmount
			customerAllBalance.I_Ded_Balance = customerAllBalance.I_Ded_Balance + iDedAccumTrans.AccumAmount
			customerAllBalance.I_OOP_Balance = customerAllBalance.I_OOP_Balance + iOopAccumTrans.AccumAmount

			claim.Status = "Processed"
			_,err := UpdateCustomerAndFamilyBalances(customerAllBalance, stub)
			if (err !=nil ){
				fmt.Println("Error updating balances ", err)
			}
			_,err = MarkLastFamTxnReview(false,claim.ContractID,claim.MemberID,claim.ClaimID,contractDefinition,customerLimit,stub )
			if (err != nil ){
					fmt.Println("Error Marking Last txn for review  ", err)
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

/*   //commenting because again all 4 txs has to be added irrespective of status
		 if (fDedUpdateFlag && fOopUpdateFlag  && iDedUpdateFlag && iOopUpdateFlag){
			claim.Transactions = append(claim.Transactions,fDedAccumTrans)
			claim.Transactions = append(claim.Transactions,fOopAccumTrans)
			claim.Transactions = append(claim.Transactions,iDedAccumTrans)
			claim.Transactions = append(claim.Transactions,iOopAccumTrans)

		 } else {

				fDedAccumTrans.Status = claim.Status
				claim.Transactions = append(claim.Transactions,fDedAccumTrans)

		 }
****/

		_,err = createUpdateClaim(claim,newClaim,stub)
		if err != nil {
			fmt.Println("createUpdateClaim failed")
			return false, err
		}

		lastTxn.Transaction = fDedAccumTrans
		_,err = UpdateLastFamMemTxn(lastTxn,stub)
		_,err = AddToLastTxnsList(lastTxn,stub)
		lastTxn.Transaction =fOopAccumTrans
		_,err = AddToLastTxnsList(lastTxn,stub)
		lastTxn.Transaction =iDedAccumTrans
		_,err = AddToLastTxnsList(lastTxn,stub)
		lastTxn.Transaction =iOopAccumTrans
		_,err = AddToLastTxnsList(lastTxn,stub)

		if (err != nil ){
			fmt.Println("Error adding to Last 10 Txns")
		}

		 fmt.Println("In services.ProcessIFDEDClaim end ")
		 return true,nil
}

func ProcessIFOOPClaim(newClaim bool, claimAmout float64,claim data.Claims,contractDefinition data.ContractDefinition,
	customerLimit data.CustomerLimit,stub shim.ChaincodeStubInterface)(bool, error){
		fmt.Println("In services.ProcessIFOOPClaim start ")

		fmt.Println("Accum type is IFOOP")
		var customerAllBalance data.CustomerAllBalance
		customerAllBalance,err := query.GetAllBalancesOfCustomer(claim.MemberID,claim.ContractID,stub)

		if err != nil {
			fmt.Println("Error retrieving Customer Balance " , claim.MemberID )
			return false, errors.New("Error retrieving Customer Balance " + claim.MemberID)
		}
		fmt.Println("Customer Balance received : " ,customerAllBalance)
		txnID := 0
		var transLength int = len(claim.Transactions)
		fmt.Println("Claim transactions length : ",transLength)
		if (!newClaim && transLength > 0) {
			txnID = claim.Transactions[transLength-1].TransactionID
		}
		fOopAccumTrans := data.Transaction{AccumType:"IFOOP", AccumAmount:claimAmout,
			TransactionDate:time.Now(),TxnUpdatedDate:time.Now(),TransactionID:txnID+1}
		iOopAccumTrans := data.Transaction{AccumType:"IIOOP", AccumAmount:claimAmout,
			TransactionDate:time.Now(),TxnUpdatedDate:time.Now(),TransactionID:txnID+1}

		fOopBalance := customerAllBalance.F_OOP_Balance
		iOopBalance := customerAllBalance.I_OOP_Balance

		lastTxn :=data.LastTxnOfContract{ClaimID:claim.ClaimID,ContractID:claim.ContractID,
				MemberID:claim.MemberID,SubscriberID:claim.SubscriberID,
					Source:claim.Participant}
		if (claimAmout < 0 ){
			fmt.Printf("Adjustment claim has come in with amout %v.  Processing .... \n",claimAmout)
			customerAllBalance.F_OOP_Balance = customerAllBalance.F_OOP_Balance + claimAmout
			customerAllBalance.I_OOP_Balance = customerAllBalance.I_OOP_Balance + claimAmout
			claim.Status ="Processed"
			_,err := UpdateCustomerAndFamilyBalances(customerAllBalance, stub)
			if (err !=nil ){
				fmt.Println("Error updating balances ", err)
			}
			_,err = MarkLastFamTxnReview(true,claim.ContractID,claim.MemberID,claim.ClaimID,contractDefinition,customerLimit,stub )
			if (err != nil ){
					fmt.Println("Error Marking Last txn for review  ", err)
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

/*  //commenting because again all 2 txs has to be added irrespective of status
		if (fOopUpdateFlag && iOopUpdateFlag) {
			fmt.Println("Both IFOOP &IIOOP are success")
			claim.Transactions = append(claim.Transactions,fOopAccumTrans)
			claim.Transactions = append(claim.Transactions,iOopAccumTrans)
			claim.TotalClaimAmount = claim.TotalClaimAmount+claimAmout
		} else {
			fmt.Println("Claim made for review")
			fOopAccumTrans.Status = claim.Status
			claim.Transactions = append(claim.Transactions,fOopAccumTrans)
		}
*/

	_,err = createUpdateClaim(claim,newClaim,stub)
	if err != nil {
		fmt.Println("createUpdateClaim failed")
		return false, err
	}
	lastTxn.Transaction =fOopAccumTrans
	_,err = UpdateLastFamMemTxn(lastTxn,stub)
	_,err = AddToLastTxnsList(lastTxn,stub)
	lastTxn.Transaction =iOopAccumTrans
	_,err = AddToLastTxnsList(lastTxn,stub)



	if (err != nil ){
		fmt.Println("Error adding to Last 10 Txns")
	}
	fmt.Println("In services.ProcessIFOOPClaim end ")
	return true,nil
}

func ProcessIIDEDClaim(newClaim bool, adjustTxn bool,accums []data.Accum,claim data.Claims,
	customerLimit data.CustomerLimit, txnID int, stub shim.ChaincodeStubInterface)(bool, error){
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

	if (iDedAccumTrans.AccumAmount < 0 || iOopAccumTrans.AccumAmount < 0){
		fmt.Printf("Adjustment claim has come in with amout %v.  Processing .... \n",iDedAccumTrans)
		customerBalance.I_Ded_Balance = customerBalance.I_Ded_Balance  + iDedAccumTrans.AccumAmount
		customerBalance.I_OOP_Balance = customerBalance.I_OOP_Balance  + iOopAccumTrans.AccumAmount

		claim.Status = "Processed"
		_,err := updateBalances(customerBalance,stub)
		if err != nil {
			fmt.Println("updateBalances failed")
			return false, err
		}
		_,err = MarkLastIndTxnReview(false,claim.ContractID,claim.MemberID,claim.ClaimID,customerLimit,stub )
		if (err != nil ){
				fmt.Println("Error Marking Last txn for review  ", err)
		}
	} else {

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

	if (adjustTxn){
		fmt.Println("Adjust the old txn after reprocessing.  ")
		for index,txn := range claim.Transactions{
			if (txn.TransactionID == txnID){
				fmt.Println("Txn match found. Update the txn.  ")
				if (claim.Transactions[index].AccumType == "IIDED"){
					claim.Transactions[index].AccumAmount = iDedAccumTrans.AccumAmount
					claim.Transactions[index].AccumBalance = iDedAccumTrans.AccumBalance
					claim.Transactions[index].Overage = iDedAccumTrans.Overage
					claim.Transactions[index].Status = iDedAccumTrans.Status
					claim.Transactions[index].TxnUpdatedDate = time.Now()
				}

				if (claim.Transactions[index].AccumType == "IIOOP"){
					claim.Transactions[index].AccumAmount = iOopAccumTrans.AccumAmount
					claim.Transactions[index].AccumBalance = iOopAccumTrans.AccumBalance
					claim.Transactions[index].Overage = iOopAccumTrans.Overage
					claim.Transactions[index].Status = iOopAccumTrans.Status
					claim.Transactions[index].TxnUpdatedDate = time.Now()
				}
			}
		}
	}else {
		claim.Transactions = append(claim.Transactions,iDedAccumTrans)
		claim.Transactions = append(claim.Transactions,iOopAccumTrans)

	}

	/*  //commenting because again all 2 txs has to be added irrespective of status
	if (iDedUpdateFlag && iOopUpdateFlag ){
			fmt.Println("Adding IIDED and IIOOP txs")
			claim.Transactions = append(claim.Transactions,iDedAccumTrans)
			claim.Transactions = append(claim.Transactions,iOopAccumTrans)
	} else {
		fmt.Println("Only IIDED success. Adding IIDED Tx")
		iDedAccumTrans.Status = claim.Status
		claim.Transactions = append(claim.Transactions,iDedAccumTrans)
	}

*/
	_,err = createUpdateClaim(claim,newClaim,stub)
	if err != nil {
		fmt.Println("createUpdateClaim failed")
		return false, err
	}

	lastTxn.Transaction =iDedAccumTrans
	_,err =	updateLastIndTxn(lastTxn,stub)
	_,err = AddToLastTxnsList(lastTxn,stub)

	lastTxn.Transaction =iOopAccumTrans
	_,err = AddToLastTxnsList(lastTxn,stub)

	if (err != nil ){
		fmt.Println("Error adding to Last 10 Txns")
	}

	fmt.Println("In services.ProcessIIDEDClaim end ")
	 return true,nil
}


func ProcessIIOOPClaim(newClaim bool, adjustTxn bool,claimAmout float64,claim data.Claims,
	customerLimit data.CustomerLimit,txnID int,stub shim.ChaincodeStubInterface)(bool, error){

	fmt.Println("In services.ProcessIIOOPClaim start ")

	var customerBalance data.CustomerBalance
	customerBalance,err := query.GetCustomerBalance(claim.MemberID,claim.ContractID,stub)

	if err != nil {
		fmt.Println("Error retrieving Customer Balance " , claim.MemberID )
		return false, errors.New("Error retrieving Customer Balance " + claim.MemberID)
	}
	fmt.Println("Customer Balance received : " ,customerBalance)

	// txnID := 0
	// var transLength int = len(claim.Transactions)
	// fmt.Println("Claim transactions length : ",transLength)
	// if (!newClaim && transLength > 0) {
	// 	txnID = claim.Transactions[transLength-1].TransactionID
	// }

	iOopAccumTrans := data.Transaction{AccumType:"IIOOP", AccumAmount:claimAmout,
		TxnUpdatedDate:time.Now(),TransactionDate:time.Now(),TransactionID:txnID}

	iOopBalance := customerBalance.I_OOP_Balance

	lastTxn :=data.LastTxnOfContract{ClaimID:claim.ClaimID,ContractID:claim.ContractID,
		MemberID:claim.MemberID,SubscriberID:claim.SubscriberID, Source:claim.Participant}
 	if (claimAmout <  0 ){
		fmt.Println("Adjustment claim has come in with amout %v.  Processing .... \n",claimAmout)
		customerBalance.I_OOP_Balance = customerBalance.I_OOP_Balance + claimAmout
		claim.Status ="Processed"
		_,err := updateBalances(customerBalance, stub)
		if (err !=nil ){
			fmt.Println("Error updating balances ", err)
		}
		_,err = MarkLastIndTxnReview(true,claim.ContractID,claim.MemberID,claim.ClaimID,customerLimit,stub )
		if (err != nil ){
				fmt.Println("Error Marking Last txn for review  ", err)
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

		}
	}

	iOopAccumTrans.Status = claim.Status
	iOopAccumTrans.AccumBalance =	customerBalance.I_OOP_Balance

	if (adjustTxn){
		fmt.Println("Adjust the old txn after reprocessing.  ")
		for index,txn := range claim.Transactions{
			if (txn.TransactionID == txnID){
				fmt.Println("Txn match found. Update the txn.  ")

				if (claim.Transactions[index].AccumType == "IIOOP"){
					claim.Transactions[index].AccumAmount = iOopAccumTrans.AccumAmount
					claim.Transactions[index].AccumBalance = iOopAccumTrans.AccumBalance
					claim.Transactions[index].Overage = iOopAccumTrans.Overage
					claim.Transactions[index].Status = iOopAccumTrans.Status
					claim.Transactions[index].TxnUpdatedDate = time.Now()
				}
			}
		}
	}else {
			claim.Transactions = append(claim.Transactions,iOopAccumTrans)
	}


	_,err = createUpdateClaim(claim,newClaim,stub)
	if err != nil {
		fmt.Println("createUpdateClaim failed")
		return false, err
	}

	lastTxn.Transaction =iOopAccumTrans
	_,err =	updateLastIndTxn(lastTxn,stub)
	_,err = AddToLastTxnsList(lastTxn,stub)
	if (err != nil ){
		fmt.Println("Error adding to Last 10 Txns")
	}

	fmt.Println("In services.ProcessIIOOPClaim end ")
	return false, nil
}

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

func updateBalances(customerBalance data.CustomerBalance,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.updateBalances start ")

	_, err := CreateCustomerBalance(customerBalance,stub)
	if err != nil {
		fmt.Println("Updating Customer Balance " )
	}
	fmt.Println("In services.updateBalances end ")
	return nil, nil
}


func updateBalancesAndClaim(updateflag bool,claim data.Claims,customerBalance data.CustomerBalance,
	stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.updateBalancesAndClaim start ")
	_, err := CreateClaims(claim,stub)
	if err != nil {
		fmt.Println("Creating claims failed  " )
	}
	if (updateflag){
		_, err = CreateCustomerBalance(customerBalance,stub)
		if err != nil {
			fmt.Println("Updating Customer Balance " )
		}
	}
	fmt.Println("In services.updateBalancesAndClaim end ")
	return nil, nil
}

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
				if (txn.AccumBalance == contractDefinition.F_Ded_Limit && txn.Status == "Processed"){
					fmt.Println("Updating IFDED txn to Review")
					 txn.Status = "Review. Balances / Limits Adjusted"
					 txn.TxnUpdatedDate = time.Now()
					 claim.Transactions[index] = txn
					 updateFlag = true
					 lastTxn.Transaction = txn
					 _,err = AddToLastTxnsList(lastTxn,stub)

				 }
			} else if ( txn.AccumType == "IFOOP"){
				fmt.Println("Found IFOOP")
				 if (txn.AccumBalance == contractDefinition.F_OOP_Limit && txn.Status == "Processed"){
					fmt.Println("Updating IFOOP txn to Review")
					 txn.Status = "Review. Balances / Limits Adjusted"
					 txn.TxnUpdatedDate = time.Now()
					 claim.Transactions[index] = txn
					 updateFlag = true
					 lastTxn.Transaction = txn
					 _,err = AddToLastTxnsList(lastTxn,stub)
				 }
			} else if (!isOppClaimAdjusted && txn.AccumType == "IIDED"){
				fmt.Println("Found IIDED")
				 if (txn.AccumBalance == customerLimit.I_Ded_Limit && txn.Status == "Processed"){
					 fmt.Println("Updating IIDED txn to Review")
					 txn.Status = "Review. Balances / Limits Adjusted"
					 txn.TxnUpdatedDate = time.Now()
					 claim.Transactions[index] = txn
					 updateFlag = true
					 lastTxn.Transaction = txn
					 _,err = AddToLastTxnsList(lastTxn,stub)
				 }
			} else if (txn.AccumType == "IIOOP"){
				fmt.Println("Found IIOOP")
				 if (txn.AccumBalance == customerLimit.I_OOP_Limit && txn.Status == "Processed"){
					 fmt.Println("Updating IIOOP txn to Review")
					 txn.Status = "Review. Balances / Limits Adjusted"
					 txn.TxnUpdatedDate = time.Now()
					 claim.Transactions[index] = txn
					 updateFlag = true
					 lastTxn.Transaction = txn
					 _,err = AddToLastTxnsList(lastTxn,stub)
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
func MarkLastIndTxnReview(isOppClaimAdjusted bool,contractID string, memberID string, claimID string,
	customerLimit data.CustomerLimit,stub shim.ChaincodeStubInterface) ([]byte, error) {
	fmt.Println("In services.markLastIndTxnReview start ")
	lastTxn, err := query.GetLastIndTxn(contractID,stub)
	if (err != nil ){
		fmt.Println("No Last txn found. ")
	}
	fmt.Println("Last ind txn for contract  ",lastTxn )
	var contractDefinition data.ContractDefinition
	_,err =UpdateTxnStatusAfterAdj(isOppClaimAdjusted,lastTxn,contractDefinition,customerLimit,stub)
	if (err != nil ){
		fmt.Println("Error Updating the status after last txn after adjustment. ")
	}

	fmt.Println("In services.markLastIndTxnReview end ")
	return nil,nil
}

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
	contractDefinition,err := query.GetContract(contractID,stub)
	if err != nil {
		fmt.Println("Error receiving  the Customer Contract")
		return nil, errors.New("Error receiving  the Customer Contract")
	}

	var customerLimit data.CustomerLimit
	customerLimit,err = query.GetCustomerLimits(memberID,contractID,stub)
	if err !=nil {
		fmt.Println("Error retrieving Customer Limits " , memberID )
		return nil, errors.New("Error retrieving Customer Limits " + memberID)
	}
	fmt.Println("Customer Limits received : " ,customerLimit)


	if (contractDefinition.ContractType =="I"){
		fmt.Println("Contract Type : I" )
		for _, accum := range accumLimits {
			if (accum.Type == "IIDED"){
				fmt.Println("Set IIDED Limit" )
				customerLimit.I_Ded_Limit = accum.Amount
			}
			if (accum.Type == "IIOOP"){
				fmt.Println("Set IIDED Limit" )
				customerLimit.I_OOP_Limit = accum.Amount
			}
		}

	}

	if (contractDefinition.ContractType == "F"){

		for _, accum := range accumLimits {
			if (accum.Type == "IFDED"){
					contractDefinition.F_Ded_Limit = accum.Amount
			}
			if (accum.Type == "IIDED"){
					contractDefinition.F_OOP_Limit = accum.Amount
			}
			if (accum.Type == "IIDED"){
					customerLimit. I_Ded_Limit = accum.Amount
			}
			if (accum.Type == "IIDED"){
					customerLimit. I_OOP_Limit = accum.Amount
			}
		}

		_,err = CreateContractDefinition(contractDefinition,stub )
		if err !=nil {
			fmt.Println("Error Updating the limits of Contract  " , contractID )

		}
	}

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

			// Mark all previous txns till the balances fall in the limits after limits adjusted
			var lastTxns []data.LastTxnOfContract
			lastTxns, err = query.GetLastTxnsOfMem(contractID,subscriberID,memberID,stub)

			//Create a map to add the already adjusted txnid of claim.
			//This helps not to adjust the same claim again
			claimTxns := make(map[string][]int)
			txnAlreadyAdj := false
			updateFlag := true

			if (len(lastTxns)>0){
				for _,txn := range lastTxns {
					if ( claimTxns[txn.ClaimID] != nil ){
						// search if txn is there in the list
						txnIDs := claimTxns[txn.ClaimID]
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
						//add to map so that we can check that  this txn is already adjusted
						claimTxns[txn.ClaimID] = append(claimTxns[txn.ClaimID], txn.Transaction.TransactionID)
					}

				}
				if (updateFlag == false ){
					break
				}
			}
		}
	}

	fmt.Println("In services.AdjustLimitsOfContract end ")
	return nil,nil
}
func UpdateTxnStatusAfterILimitAdj(lastTxn data.LastTxnOfContract,contractDefinition data.ContractDefinition,
customerLimit data.CustomerLimit,stub shim.ChaincodeStubInterface)(bool, error){
	fmt.Println("In services.UpdateTxnStatusAfterILimitAdj start ")

	claim,err := query.GetClaim(lastTxn.MemberID,lastTxn.ContractID,lastTxn.ClaimID,stub)
	fmt.Printf("Last Claim %v found with txnID %v \n",lastTxn.ClaimID,lastTxn.Transaction.TransactionID)

	var txns_index_map = make(map[string]int)
	for index, txn := range claim.Transactions {
		fmt.Printf("Claim %v Transaction with txn : ",claim.ClaimID,txn)
		if (lastTxn.Transaction.TransactionID == txn.TransactionID) {
			txns_index_map[txn.AccumType] = index
		}
	}

	fmt.Println("%v of txns found with txn id : %v \n", len(txns_index_map) ,lastTxn.Transaction.TransactionID)
	reviewOOPFlag := false
	reviewDEDFlag := false

	if (txns_index_map!= nil && len(txns_index_map) == 2) {
		fmt.Println("IIDED claim. So verify both IIDED and IIOOP Txns. ")
		index := txns_index_map["IIOOP"]
		if (claim.Transactions[index].AccumBalance >= customerLimit.I_OOP_Limit){
			fmt.Println("IIOOP Balance is >=  limit ")
			reviewOOPFlag = true
		}
		index = txns_index_map["IIDED"]
	  if (claim.Transactions[index].AccumBalance >= customerLimit.I_Ded_Limit){
				fmt.Println("IIDED Balance is >=  limit ")
				reviewDEDFlag = true
		}

		if (reviewDEDFlag ||  reviewOOPFlag ){
			claim.Transactions[txns_index_map["IIDED"]].Status = "Review – Exceeds Limit"
			claim.Transactions[txns_index_map["IIDED"]].TxnUpdatedDate = time.Now()
			claim.Transactions[txns_index_map["IIOOP"]].Status = "Review – Exceeds Limit"
			claim.Transactions[txns_index_map["IIOOP"]].TxnUpdatedDate = time.Now()
		}
	}

	if (txns_index_map!= nil && len(txns_index_map) == 1) {
		index := txns_index_map["IIOOP"]
		fmt.Println("IIOOP claim.")
		if (claim.Transactions[index].AccumBalance >= customerLimit.I_OOP_Limit){
			claim.Transactions[txns_index_map["IIOOP"]].Status = "Review – Exceeds Limit"
			claim.Transactions[txns_index_map["IIOOP"]].TxnUpdatedDate = time.Now()
			reviewOOPFlag = true

		}
	}

	if (reviewDEDFlag || reviewOOPFlag){
		fmt.Println("Tx status has been modified.")
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
/**   Backup before adjust claim *

func ProcessIIDEDClaim(newClaim bool, accums []data.Accum,claim data.Claims,
	customerLimit data.CustomerLimit, txnID string, stub shim.ChaincodeStubInterface)(bool, error){
	fmt.Println("In services.ProcessIIDEDClaim start ")

	var customerBalance data.CustomerBalance
	customerBalance,err := query.GetCustomerBalance(claim.MemberID,claim.ContractID,stub)

	if err != nil {
		fmt.Println("Error retrieving Customer Balance " , claim.MemberID )
		return false, errors.New("Error retrieving Customer Balance " + claim.MemberID)
	}
	fmt.Println("Customer Balance received : " ,customerBalance)

	txnID := 1
	var transLength int = len(claim.Transactions)
	fmt.Println("Claim transactions length : ",transLength)
	if (!newClaim && transLength > 0) {
		txnID = claim.Transactions[transLength-1].TransactionID +1
	}

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

	if (iDedAccumTrans.AccumAmount < 0 || iOopAccumTrans.AccumAmount < 0){
		fmt.Printf("Adjustment claim has come in with amout %v.  Processing .... \n",iDedAccumTrans)
		customerBalance.I_Ded_Balance = customerBalance.I_Ded_Balance  + iDedAccumTrans.AccumAmount
		customerBalance.I_OOP_Balance = customerBalance.I_OOP_Balance  + iOopAccumTrans.AccumAmount

		claim.Status = "Processed"
		_,err := updateBalances(customerBalance,stub)
		if err != nil {
			fmt.Println("updateBalances failed")
			return false, err
		}
		_,err = MarkLastIndTxnReview(false,claim.ContractID,claim.MemberID,claim.ClaimID,customerLimit,stub )
		if (err != nil ){
				fmt.Println("Error Marking Last txn for review  ", err)
		}
	} else {

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

	//   //commenting because again all 2 txs has to be added irrespective of status
	// if (iDedUpdateFlag && iOopUpdateFlag ){
	// 		fmt.Println("Adding IIDED and IIOOP txs")
	// 		claim.Transactions = append(claim.Transactions,iDedAccumTrans)
	// 		claim.Transactions = append(claim.Transactions,iOopAccumTrans)
	// } else {
	// 	fmt.Println("Only IIDED success. Adding IIDED Tx")
	// 	iDedAccumTrans.Status = claim.Status
	// 	claim.Transactions = append(claim.Transactions,iDedAccumTrans)
	// }


	_,err = createUpdateClaim(claim,newClaim,stub)
	if err != nil {
		fmt.Println("createUpdateClaim failed")
		return false, err
	}

	lastTxn.Transaction =iDedAccumTrans
	_,err =	updateLastIndTxn(lastTxn,stub)
	_,err = AddToLastTxnsList(lastTxn,stub)

	lastTxn.Transaction =iOopAccumTrans
	_,err = AddToLastTxnsList(lastTxn,stub)

	if (err != nil ){
		fmt.Println("Error adding to Last 10 Txns")
	}

	fmt.Println("In services.ProcessIIDEDClaim end ")
	 return true,nil
}

*/
