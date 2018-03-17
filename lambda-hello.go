package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func Handler() (string, error) {
	// printing with both fmt and log
	fmt.Printf(">>>> Inside Lambda Handler\n")
	log.Printf(">>>> Inside Lambda Handler\n")
	return "All Done", nil
}

func main() {
	lambda.Start(Handler)
}

