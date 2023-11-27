package main

import (
	"flag"
	"fmt"
	"path/filepath"
)

//  E.g.: ./exchanger --service rtdialog --environment vcp-dev --request getAvailableIds

func main() {
	var service string
	flag.StringVar(&service, "service", "not-provided-service",
		"Service like, rtlmm, rtdialog or metering.")

	var request string
	flag.StringVar(&request, "request", "not-provided-request",
		"Request like, getAvailableIds or getDialogDetails.")

	var environment string
	flag.StringVar(&environment, "environment", "not-provided-environment",
		"Environment like, vcp-dev, prod/inno-prod")

	flag.Parse()

	fmt.Println("User provided service: ", service)
	fmt.Println("User provided request: ", request)
	fmt.Println("User provided environment: ", environment)

	if len(service) < 1 || len(environment) < 1 || len(request) < 1 {
		fmt.Println("One of the service, environment and request parameters has zero length. "+
			"service:", service, "; environment:", environment, "; request:", request)
		return
	}

	environmentDataFilePath := filepath.Join("environment", service, environment+".json")
	//	input
	serviceRequestDataFilePath := filepath.Join("service", service, request+".json")
	//  output
	serviceResponseDataFilePath := filepath.Join("service", service, request+"-output.json")

	//	input
	tokenRequestDataFilePath := filepath.Join("service", service, "token-request.json")
	//	output then input
	tokenResponseDataFilePath := filepath.Join("service", service, "token-response.json")

	var environmentData EnvironmentData
	err := ReadData(environmentDataFilePath, &environmentData)
	if err != nil {
		return
	}

	// var tokenRequest Request = &TokenRequest{
	// 	generalRequestData: GeneralRequestData{
	// 		extendEnvironmentData: false,
	// 		environmentData:       &environmentData,
	// 		requestDataFilePath:   tokenRequestDataFilePath,
	// 		responseDataFilePath:  tokenResponseDataFilePath,
	// 	},
	// }
	//
	// serviceRequest := ServiceRequest{
	// 	generalRequestData: GeneralRequestData{
	// 		extendEnvironmentData: true,
	// 		environmentData:       &environmentData,
	// 		requestDataFilePath:   serviceRequestDataFilePath,
	// 		responseDataFilePath:  serviceResponseDataFilePath,
	// 	},
	// 	environmentDataExtenderFilePath: tokenResponseDataFilePath,
	// 	previousRequest:       			 &tokenRequest,
	// }
	//
	// err = serviceRequest.Request()
	// if err != nil {
	// 	return
	// }

	tExchangePerformer := ExchangePerformer{
		exchangePerformable: &TokenExchangePerformer{},

		environmentData:      &environmentData,
		requestDataFilePath:  tokenRequestDataFilePath,
		responseDataFilePath: tokenResponseDataFilePath,
	}

	sExchangePerformer := ExchangePerformer{
		previousExchange:    &tExchangePerformer,
		exchangePerformable: &ServiceExchangePerformer{},

		environmentData:                 &environmentData,
		environmentDataExtenderFilePath: tokenResponseDataFilePath,
		requestDataFilePath:             serviceRequestDataFilePath,
		responseDataFilePath:            serviceResponseDataFilePath,
	}
	err = sExchangePerformer.extendEnvironmentDataAndPerformServiceExchange()
	if err != nil {
		return
	}
}
