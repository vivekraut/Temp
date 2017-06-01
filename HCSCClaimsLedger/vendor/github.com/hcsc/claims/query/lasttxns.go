package query

import (
    "fmt"
    "sort"
    "errors"
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
  lastTxns,err := GetLastTxnsOfMem(contractID,"",memberID,stub)
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
func GetLastTxnsOfMem(contractID string, subscriberID string, memberID string,
  stub shim.ChaincodeStubInterface) ([]data.LastTxnOfContract, error) {

	fmt.Println("In lasttxns.GetLastTxnsOfMem start ")
  var claimIDs []string
  claimIDs,err := GetMemberClaimIDsOfContract(memberID,contractID,stub)
  fmt.Println("Claim IDs : " ,claimIDs)
  var txns_data_map = make(map[int]data.LastTxnOfContract)
  txns_sorted := make(timeSlice, 0, len(txns_data_map))
  txnIndex := 0
  if err !=nil {
      fmt.Printf("Error retrieving Claim Ids of Member %v for contract %v  \n" , memberID ,contractID)
  }else {
    fmt.Printf("No of Claims for member %v : \n" , memberID,len(claimIDs))
    for _,claimID := range claimIDs {
	     claim, _ := GetClaim(memberID,contractID,claimID,stub)
       if (len(txns_data_map) > 0 ){
         txnIndex = len(txns_data_map)
       }
       claimTxns := claim.Transactions
       if (len(claimTxns) > 0 ){
         fmt.Printf("No of txns for claim :  %v  is %v \n" , claimID,len(claimTxns))
         for _,txn := range claimTxns{

           var lastTxn data.LastTxnOfContract
           lastTxn.ContractID = claim.ContractID
           lastTxn.MemberID = claim.MemberID
           lastTxn.SubscriberID = claim.SubscriberID
           lastTxn.ClaimID = claim.ClaimID
           lastTxn.Transaction = txn

           txns_data_map[txnIndex] =lastTxn
           txnIndex = txnIndex +1

         }
       }

    }
    fmt.Println("Total Txns : ", len(txns_data_map))

    txns_sorted = make(timeSlice, 0, len(txns_data_map))
    for _, d := range txns_data_map {
        txns_sorted = append(txns_sorted, d)
    }
    sort.Sort(txns_sorted)
    fmt.Println("Sorted Txns :  ")
    fmt.Println(txns_sorted)
    //return nil, errors.New("Error retrieving Members of Contract " + contractID)
  }
  fmt.Println("In lasttxns.GetLastTxnsOfMem end ")
  return txns_sorted,nil
}

func GetAllLastTxns(noOftxns int, stub shim.ChaincodeStubInterface)([]data.LastTxnOfContract, error) {
  fmt.Println("In lasttxns.GetLastTxnsOfContract start ")

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
  			lastTxns, err = GetLastTxnsOfMem(contractID,"",memID,stub)
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

  if (len(txns_sorted) >= noOftxns){
    return allLastTxns[:noOftxns+1]
  }
  fmt.Println(txns_sorted)
  fmt.Println("In lasttxns.GetLastTxnsOfContract end ")
  return allLastTxns,nil
}
/*
func main() {

// claim1 := data.Claims{ ClaimID:"100", data.Transaction{AccumType:"IFDED",TransactionDate:time.Now(),TransactionID:1}
// claim2 := data.Claims{ ClaimID:"101", data.Transaction{AccumType:"IFOOP",TransactionDate:time.Now(),TransactionID:2}
// claim3 := data.Claims{ ClaimID:"102", data.Transaction{AccumType:"IFOOP",TransactionDate:time.Now(),TransactionID:3}

tempTxn := data.Transaction{AccumType:"IFDED",TransactionDate:time.Now().Add(12 * time.Hour),TransactionID:1}
txn1 := data.LastTxnOfContract{ClaimID:"100", Transaction:tempTxn }
tempTxn = data.Transaction{AccumType:"IFDED",TransactionDate:time.Now(),TransactionID:1}
txn2 := data.LastTxnOfContract{ClaimID:"102", Transaction:tempTxn }
tempTxn = data.Transaction{AccumType:"IFDED",TransactionDate:time.Now().Add(24 * time.Hour),TransactionID:1}
txn3 := data.LastTxnOfContract{ClaimID:"103",Transaction:tempTxn }

var txns_data_map = make(map[int]data.LastTxnOfContract)
txns_data_map[0] = txn1//append(txns_data_map,txn1)
txns_data_map[1] = txn2//append(txns_data_map,txn2)
txns_data_map[2] = txn3//append(txns_data_map,txn3)

fmt.Println("Len of map : ", len(txns_data_map))
fmt.Println("txns_data_map[0] of map : ", txns_data_map[0])
  //Sort the map by date
  date_sorted_reviews := make(timeSlice, 0, len(txns_data_map))
  for _, d := range txns_data_map {
      date_sorted_reviews = append(date_sorted_reviews, d)
  }
  fmt.Println(date_sorted_reviews)
  sort.Sort(date_sorted_reviews)
  fmt.Println("")
  fmt.Println(date_sorted_reviews)
}
*/
