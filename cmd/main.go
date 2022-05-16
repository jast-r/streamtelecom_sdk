package main

import (
	"context"
	"fmt"

	streamtelecom "github.com/jast-r/streamtelecom_sdk"
)

func main() {
	test, err := streamtelecom.NewClient("assi", "dQc20d6ybH")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(test.GetBalance(context.Background()))
	fmt.Println(test.GetSenderList(context.Background()))
	test.GetTariffList(context.Background())
	// fmt.Println(test.SendSingleSMS(context.Background(), streamtelecom.SingleSMSRequest{
	// 	DestinationAddress: "79997600375",
	// 	Text:               "test",
	// 	SourceAddress:      "ASSI",
	// 	TTL:                "5",
	// 	SendDate:           time.Now().Add(5 * time.Minute),
	// }))
}
