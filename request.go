package main

import (
	"errors"
	"fmt"
)

type ExchangePerformable interface {
	extendEnvironmentData(
		environmentData *EnvironmentData,
		environmentDataExtenderFilePath string) error

	serviceExchange(
		environmentData *EnvironmentData,
		requestDataFilePath string,
		responseDataFilePath string) error
}

type ExchangePerformer struct {
	previousExchange    *ExchangePerformer
	exchangePerformable ExchangePerformable

	environmentData                 *EnvironmentData
	environmentDataExtenderFilePath string
	requestDataFilePath             string
	responseDataFilePath            string
}

func (exchangePerformer *ExchangePerformer) extendEnvironmentDataAndPerformServiceExchange() error {

	if exchangePerformer.previousExchange != nil {
		err := (*exchangePerformer.previousExchange).extendEnvironmentDataAndPerformServiceExchange()
		if err != nil {
			return err
		}
	}

	err := exchangePerformer.exchangePerformable.extendEnvironmentData(
		exchangePerformer.environmentData,
		exchangePerformer.environmentDataExtenderFilePath)
	if err != nil {
		return err
	}

	err = exchangePerformer.exchangePerformable.serviceExchange(
		exchangePerformer.environmentData,
		exchangePerformer.requestDataFilePath,
		exchangePerformer.responseDataFilePath)
	if err != nil {
		return err
	}

	return nil
}

type TokenExchangePerformer struct {
	exchangePerformer ExchangePerformer
}

func (tokenExchangePerformer *TokenExchangePerformer) extendEnvironmentData(
	environmentData *EnvironmentData,
	environmentDataExtenderFilePath string) error {

	//  no-op
	return nil
}

func (tokenExchangePerformer *TokenExchangePerformer) serviceExchange(
	environmentData *EnvironmentData,
	requestDataFilePath string,
	responseDataFilePath string) error {

	err := serviceExchange(
		environmentData,
		requestDataFilePath,
		responseDataFilePath)
	if err != nil {
		return err
	}
	return nil
}

func serviceExchange(
	environmentData *EnvironmentData,
	requestDataFilePath string,
	responseDataFilePath string) error {

	serviceData := ServiceData{}
	err := ReadData(requestDataFilePath, &serviceData)
	if err != nil {
		return err
	}

	request, err := CreateRequest(&serviceData, environmentData)
	if err != nil {
		return err
	}

	body, statusCode, err := Exchange(request)
	if err != nil {
		return err
	}

	err = Write(responseDataFilePath, body)
	if err != nil {
		return err
	}

	if statusCode >= 300 {
		errorMessage := fmt.Sprint("Status code is equal to or greater than 300.\nError:", statusCode)
		return errors.New(errorMessage)
	}

	return nil
}

func extendEnvironmentData(environmentData *EnvironmentData, tokenDataFilePath string) error {
	tokenData := TokenData{}
	err := ReadData(tokenDataFilePath, &tokenData)
	if err != nil {
		return err
	}

	if len(tokenData.AccessToken) < 1 && len(tokenData.Token) < 1 {
		fmt.Println(" ********* accessToken:", tokenData.AccessToken)
		fmt.Println(" ********* token:", tokenData.Token)
		fmt.Println("Token is invalid.")
		return err
	}

	if len(tokenData.AccessToken) > 0 {
		fmt.Println("AccessToken is used.")
		(*environmentData)["access_token"] = tokenData.AccessToken
	} else if len(tokenData.Token) > 0 {
		fmt.Println("Token is used.")
		(*environmentData)["access_token"] = tokenData.Token
	}

	return nil
}

type ServiceExchangePerformer struct {
	exchangePerformer ExchangePerformer
}

func (serviceExchangePerformer *ServiceExchangePerformer) extendEnvironmentData(
	environmentData *EnvironmentData,
	environmentDataExtenderFilePath string) error {

	err := extendEnvironmentData(
		environmentData,
		environmentDataExtenderFilePath)
	if err != nil {
		return err
	}
	return nil
}

func (serviceExchangePerformer *ServiceExchangePerformer) serviceExchange(
	environmentData *EnvironmentData,
	requestDataFilePath string,
	responseDataFilePath string) error {
	err := serviceExchange(
		environmentData,
		requestDataFilePath,
		responseDataFilePath)
	if err != nil {
		return err
	}
	return nil
}
