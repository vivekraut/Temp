package services
import (
	"fmt"
	"github.com/hcsc/claims/data"
	"strconv"
)

func InvokeRule(accumType string, claimAmout float64, accumLimit float64,
	accumBalance float64, transaction *data.Transaction) (bool,float64) {
	fmt.Println("In services.InvokeRule for "+accumType +" start  ")
	fmt.Println("claimAmout : %f , accumLimit %f, accumBalance %f" ,claimAmout, accumLimit,accumBalance )
	var updateFlag bool =false

	if (accumBalance < accumLimit) {
		if ((claimAmout + accumBalance) <= accumLimit){
			fmt.Println("ClaimAmout + accumBalance is less than accumLimit")
			accumBalance = accumBalance + claimAmout
			transaction.Status = "Processed"
			updateFlag = true
		}else if ((claimAmout + accumBalance) > accumLimit){
			fmt.Println(accumType + ": transaction lead to Overage ")
			transaction.Overage = claimAmout + accumBalance - accumLimit
			transaction.Status = "Review. "+ accumType +" Overage : " +strconv.FormatFloat(transaction.Overage, 'f',2, 64)
		}
	}else {
		fmt.Println(accumType + ": Limit Reached ")
		transaction.Status = "Review. " +accumType + " Limit Reached"
	}

	fmt.Println("In services.InvokeRule Transaction :  ", transaction)
	fmt.Println("In services.InvokeRule end ")

	return updateFlag,accumBalance

}
