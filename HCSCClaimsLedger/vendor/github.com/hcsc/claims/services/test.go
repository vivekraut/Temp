package services

import (
	// "errors"
	// "fmt"
	// //"strconv"
	// "time"
	// "github.com/hyperledger/fabric/core/chaincode/shim"
	// "github.com/hcsc/claims/data"
	// "github.com/hcsc/claims/query"
	// "encoding/json"
)
import "fmt"

// type ClaimTxns struct {
// 	txnIDs []
// }



func main() {


  claimID :="CL10"
  txnID := "tx1"
  x := make(map[string][]string)

   x[claimID] = append(x[claimID], txnID)
   txnID = "tx2"
   x[claimID] = append(x[claimID], txnID)
   claimID ="CL11"
   txnID = "tx3"

   key :="CL12"
   value := "tx3"
   if (len(x) > 0){
     fmt.Println("Length is : ", len(x))
     if (x[key] != nil ){
       fmt.Println("x[key] is : ", x[key])
       txnIDs :=x[key]
      // if (txnIDs != nil && len(txnIDs) > 0){
         for index,txnID := range txnIDs{
           fmt.Println("x[key][index]: ",x[key][index])
           fmt.Println("txnID: ",txnID)
           if (txnID == value){
                fmt.Println("Found match: ",x[key][index])
           }
         }
       //}
     }

   }
   if(x["CL04"] != nil ){
      fmt.Println("Found x[CL04]is : ", x["CL04"])
   }


}
