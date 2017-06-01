package query

import (
    "fmt"
    "sort"
    "errors"
    "encoding/json"
    "strings"
  	"github.com/hyperledger/fabric/core/chaincode/shim"
  	"github.com/hcsc/claims/data"
)

type timeSlice []data.LastTxnOfContract

func (p timeSlice) Len() int {
    return len(p)
}
func (p timeSlice) Less(i, j int) bool {
    return p[i].Transaction.TransactionDate.After(p[j].Transaction.TransactionDate)
}
func (p timeSlice) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

//Get last one txn for a given member from the contract
func GetLastTxnOfMem(contractID string, memberID string,
  stub shim.ChaincodeStubInterface) (data.LastTxnOfContract, error) {
	fmt.Println("In lasttxns.GetLastTxnOfMem start ")
  var lastTxn data.LastTxnOfContract
  lastTxns,err := GetTxnsOfMemCont(contractID,"",memberID,"",0,stub)
  if (err != nil){
    fmt.Println("Error getting last txns of Member  %v for contract %v ", memberID, contractID)
  }
  if (len(lastTxns) > 0){
    lastTxn = lastTxns[0]
  }
  fmt.Println("In lasttxns.GetLastTxnOfMem end ")
  return lastTxn,nil
}

//Get all txns in sorted order for a given member from the contract
func GetTxnsOfMemCont(contractID string, subscriberID string, memberID string,source string,noOftxns int,
    stub shim.ChaincodeStubInterface) ([]data.LastTxnOfContract, error) {
	fmt.Println("In lasttxns.GetTxnsOfMemCont start ")

  fmt.Printf("Input: contractID %v ,subscriberID %v ,memberID %v \n", contractID,subscriberID,memberID)
  var allTxnsOfCont []data.LastTxnOfContract
  var claimIDs []string
  claimIDs,err := GetMemberClaimIDsOfContract(memberID,contractID,stub)
  fmt.Println("Claim IDs : " ,claimIDs)
  //var txns_data_map = make(map[int]data.LastTxnOfContract)
  txns_sorted := make(timeSlice, 0, len(allTxnsOfCont))
  //  txnIndex := 0
  if err !=nil {
      fmt.Printf("Error retrieving Claim Ids of Member %v for contract %v  \n" , memberID ,contractID)
  }else {
    fmt.Printf("No of Claims for member %v : \n" , memberID,len(claimIDs))
    for _,claimID := range claimIDs {
	     claim, _ := GetClaim(memberID,contractID,claimID,stub)
      //  if (len(txns_data_map) > 0 ){
      //    txnIndex = len(txns_data_map)
      //  }
       claimTxns := claim.Transactions
       if (len(claimTxns) > 0 ){
         fmt.Printf("No of txns for claim :  %v  is %v \n" , claimID,len(claimTxns))
         for _,txn := range claimTxns{
           fmt.Printf("Source  %v Txn Source %v \n", source,claim.Participant )
           if ((source == "") || (source != "" && strings.ToUpper(source) == strings.ToUpper(claim.Participant))){
             fmt.Println("Source match ")
             var lastTxn data.LastTxnOfContract
             lastTxn.ContractID = claim.ContractID
             lastTxn.MemberID = claim.MemberID
             lastTxn.SubscriberID = claim.SubscriberID
             lastTxn.ClaimID = claim.ClaimID
             lastTxn.Source = strings.Title(strings.ToLower(claim.Participant))
             lastTxn.Transaction = txn
            allTxnsOfCont = append(allTxnsOfCont,lastTxn)
           }
         }
       }
    }

    fmt.Println("Total Txns : ", len(allTxnsOfCont))
    txns_sorted = make(timeSlice, 0, len(allTxnsOfCont))
    for _, d := range allTxnsOfCont {
        txns_sorted = append(txns_sorted, d)
    }
    sort.Sort(txns_sorted)
    fmt.Printf("Sorted txns : %v \n",len(txns_sorted))
  }
  fmt.Println("In lasttxns.GetTxnsOfMemCont end ")
  if (noOftxns == 0 ){
    //if sent as 0 , will return all txns
      fmt.Printf("Return all txns  %v \n",len(txns_sorted))
      return txns_sorted,nil
  } else if (len(txns_sorted) >= noOftxns){
      fmt.Printf("Return %v txns from %v \n",noOftxns,len(txns_sorted))
    // if request for   > 0  txns, send those many txns if available
    return txns_sorted[:noOftxns],nil
  }

  fmt.Printf("Less: Return %v txns from %v \n",len(txns_sorted), noOftxns)
  return txns_sorted,nil
}

// Get transactions from Ledger irrespective of contract in the order of date
func GetTxnsFromLedger(noOftxns int,source string, stub shim.ChaincodeStubInterface)([]data.LastTxnOfContract, error) {
  fmt.Println("In lasttxns.GetTxnsFromLedger start ")
  fmt.Printf("Input: noOftxns %v , source %v \n", noOftxns, source)
  var allLastTxns []data.LastTxnOfContract
  var contractsList []string
  // 1. Get all Contracts from ledger
  contractsList,err := GetContractsList(stub)
  if (err != nil){
    fmt.Println("No Contracts Found in Ledger")
    return allLastTxns, errors.New("No Contracts Found in Ledger")
  }
  for _,contractID :=range contractsList {
    //2. Get Members of the Contract
    fmt.Println("Get Members of the Contract ", contractID)
    membersOfContract,err := GetMembersOfContract(contractID,stub)
  	fmt.Printf("Members of Contract %v is : %v \n",contractID,membersOfContract)
  	if err !=nil {
  		fmt.Println("NO Members Found for the contract   " , contractID )
  	}
    memberIDs := membersOfContract.MemberIDs
  	fmt.Printf("%v Members found in contract %v  \n",  len(memberIDs),contractID,)
  	if (len(memberIDs) > 0 ){
  		for _,memID := range memberIDs {
        //3. Get Claims of each member
        var lastTxns []data.LastTxnOfContract
  			lastTxns, err = GetTxnsOfMemCont(contractID,"",memID,source,0,stub)
        fmt.Printf("%v Txns found for Member  %v  \n",  len(lastTxns),memID)
        allLastTxns = append(allLastTxns,lastTxns...)
  		 }
     }
  }
  txns_sorted := make(timeSlice, 0, len(allLastTxns))
  for _, d := range allLastTxns {
      txns_sorted = append(txns_sorted, d)
  }
  sort.Sort(txns_sorted)
  fmt.Println("Total No. of Txns is :  ", len(txns_sorted))
  if (noOftxns == 0 ){
    //if sent as 0 , will return all txns
      return txns_sorted,nil
  } else if (len(txns_sorted) >= noOftxns){
    // if request for   > 0  txns, send those many txns if available
    return txns_sorted[:noOftxns],nil
  }
  fmt.Println("In lasttxns.GetTxnsFromLedger end ")
  // if request for > 0 , but not those many txns available, send whatever available
  return txns_sorted,nil
}

func GetClaimsOfContract(contractID string,stub shim.ChaincodeStubInterface)([]data.Claims, error){
	fmt.Println("In query.GetClaimsOfContract start ")
	var claimsList []data.Claims
	var contractMembers data.ContractMembers
  contractMembers,err := GetMembersOfContract(contractID, stub)
  if err != nil {
    fmt.Println("Error receiving  the Members of Contract")
    return nil, errors.New("Error receiving  Members of Contract")
  }
  fmt.Println("Members of contract : ",contractMembers)
  //2. Get claims of each member for given contract
  memberIDs :=  contractMembers.MemberIDs
  for _,memberID := range memberIDs {
    memContClaims, err := GetMemberClaimIDsOfContract(memberID, contractID,stub)
  	if err != nil {
  		fmt.Println("Error retrieving claims of contract  ",contractID)
  		return claimsList, errors.New("Error retrieving claims of contract  "+contractID )
  	}
    if (len(memContClaims) > 0 ){
  		for _, claimID := range memContClaims {
  			var claim data.Claims
  			claim, err = GetClaim(memberID,contractID,claimID,stub)
  			if err != nil {
  				return claimsList,err
  			}
  			fmt.Println("Claim : " , claim)
   			claimsList = append(claimsList,claim)
  		}
  	}else {
  		fmt.Println("No Claims found for contract : "+contractID )
  	}
  }
	fmt.Println("In query.GetClaimsOfContract end ")
	return claimsList, err
}

func GetTxnsOfContract(contractID string,noOftxns int,source string,stub shim.ChaincodeStubInterface)([]data.LastTxnOfContract, error) {
	fmt.Println("In query.GetTxnsOfContract start ")
  var allTxnsOfCont []data.LastTxnOfContract
	//1. Get all the memberIDs of the contract
	var contractMembers data.ContractMembers
	contractMembers,err := GetMembersOfContract(contractID, stub)
	if err != nil {
		fmt.Println("No members found for the Contract ",contractID)
		return nil, errors.New("No members found for the Contract  "+ contractID)
	}
	fmt.Println("Members of contract : ",contractMembers)
	//2. Get claims of each member for given contract
  memberIDs :=  contractMembers.MemberIDs
  for _,memberID := range memberIDs {
    memContClaims, err := GetMemberClaimIDsOfContract(memberID, contractID,stub)
  	if err != nil {
  		fmt.Println("No claims found for member  ",memberID)
  		//return allTxnsOfCont, errors.New("Error retrieving claims of contract  "+contractID )
      continue
  	}
    fmt.Printf("%v claims found for Member %v from Contract %v ", len(memContClaims), memberID,contractID)
    if (len(memContClaims) > 0 ){
  		for _, claimID := range memContClaims {
  			var claim data.Claims
  			claim, err = GetClaim(memberID ,contractID ,claimID,stub)
  			if err != nil {
          fmt.Println("No claim found for claimID  ",claimID)
  				//return allTxnsOfCont,err
          continue
  			}
        for _,txn := range claim.Transactions {
          fmt.Printf("Source  %v Txn Source %v \n", source,claim.Participant )
          if ((source == "") || (source != "" && strings.ToUpper(source) == strings.ToUpper(claim.Participant))){
            fmt.Println("Source match ")
            var lastTxn data.LastTxnOfContract
            lastTxn.ContractID = claim.ContractID
            lastTxn.SubscriberID =claim.SubscriberID
            lastTxn.MemberID = claim.MemberID
            lastTxn.ClaimID = claim.ClaimID
            lastTxn.Source = strings.Title(strings.ToLower(claim.Participant))//strings.Title(claim.Participant)
            lastTxn.Transaction = txn
            allTxnsOfCont = append(allTxnsOfCont, lastTxn)
          }
        }
  		}
  	}else {
  		fmt.Println("No Claims found for contract : "+contractID )
  	}
  }
  fmt.Println("No of Txns of Contract : ", len(allTxnsOfCont))
  txns_sorted := make(timeSlice, 0, len(allTxnsOfCont))
  for _, d := range allTxnsOfCont {
      txns_sorted = append(txns_sorted, d)
  }
  sort.Sort(txns_sorted)
  if (noOftxns == 0 ){
    //if sent as 0 , will return all txns
      return txns_sorted,nil
  } else if (len(txns_sorted) >= noOftxns){
    // if request for   > 0  txns, send those many txns if available
    return txns_sorted[:noOftxns],nil
  }
	fmt.Println("In query.GetTxnsOfContract end ")
	return txns_sorted, err
}

func GetTxnsOfMem(memberID string,noOftxns int,source string,stub shim.ChaincodeStubInterface) ([]data.LastTxnOfContract, error) {
	fmt.Println("In lasttxns.GetTxnsOfMem start ")
  var allTxnsOfCont []data.LastTxnOfContract
  var memberContracts []string
  //1. Get member's contracts
  var contractsList []string
  contractsList,err := GetContractsList(stub)
  if err != nil {
    fmt.Println("Error getting  contracts from ledger  ")
    return nil, errors.New("Error getting  contracts from ledger ")
  }

  for _,contractID := range contractsList {
    //Get members of contract
    var membersOfContract data.ContractMembers
    membersOfContract,err := GetMembersOfContract(contractID, stub)
    if err != nil {
      fmt.Println("No members found")
      continue
    }
    fmt.Printf("Members of Contract %v is %v \n", contractID,membersOfContract)
  	memberIDs := membersOfContract.MemberIDs
  	for _,memID := range memberIDs{
      if (memID == memberID){
        fmt.Printf("Member %v found in contract %v  \n  ",memberID, contractID)
        memberContracts = append(memberContracts, contractID)
        break
      }
    }
  }
  fmt.Printf("Contracts for the member %v are : %v \n  ",memberID, memberContracts)

  //Get claims of member for each contract
  var claimIDs []string
  if (len(memberContracts) > 0 ){
    for _,contractID := range memberContracts{
        claimIDs,err = GetMemberClaimIDsOfContract(memberID,contractID,stub)
        fmt.Printf("Claim IDs for Member contract %v \n  " ,contractID, claimIDs)
        if err !=nil {
            fmt.Printf("Error retrieving Claim Ids of Member %v for contract %v  \n" , memberID ,contractID)
        }else {
          for _,claimID := range claimIDs {
      	     claim, _ := GetClaim(memberID,contractID,claimID,stub)
             claimTxns := claim.Transactions
             if (len(claimTxns) > 0 ){
               for _,txn := range claimTxns{
                 fmt.Printf("Source  %v Txn Source %v \n", source,claim.Participant )
                 if ((source == "") || (source != "" && strings.ToUpper(source) == strings.ToUpper(claim.Participant))){
                   fmt.Println("Source match ")
                   var lastTxn data.LastTxnOfContract
                   lastTxn.ContractID = claim.ContractID
                   lastTxn.MemberID = claim.MemberID
                   lastTxn.SubscriberID = claim.SubscriberID
                   lastTxn.ClaimID = claim.ClaimID
                   lastTxn.Source = strings.Title(strings.ToLower(claim.Participant))//strings.Title(claim.Participant)
                   lastTxn.Transaction = txn
                   allTxnsOfCont = append(allTxnsOfCont,lastTxn)
                 }
               }
             }
           }
        }
    }
  }
  fmt.Println("Total Txns : ", len(allTxnsOfCont))
  txns_sorted := make(timeSlice, 0, len(allTxnsOfCont))
  for _, d := range allTxnsOfCont {
        txns_sorted = append(txns_sorted, d)
  }
  sort.Sort(txns_sorted)
  fmt.Println("In lasttxns.GetTxnsOfMem end ")
  if (noOftxns == 0 ){
    //if sent as 0 , will return all txns
      return txns_sorted,nil
  } else if (len(txns_sorted) >= noOftxns){
    // if request for   > 0  txns, send those many txns if available
    return txns_sorted[:noOftxns],nil
  }
  fmt.Println("In lasttxns.GetTxnsOfMem end ")
  return txns_sorted,nil
}

//NOTE: This method is used in initialize.AddToLastTxnsList which is not being used.
// will remove after confirmaton
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
