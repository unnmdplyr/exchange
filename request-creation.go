package main

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func CreateRequest(service *ServiceData, envData *EnvironmentData) (*http.Request, error) {

	replacedBody, err := executeTemplateForBody(service.Body, envData)
	if err != nil {
		fmt.Println("An error happened during traversing body", err)
		return nil, err
	}

	body := createBodyString(&replacedBody, &service.Headers)

	serviceUrl, err := ExecuteTemplate(service.Url, envData)
	if err != nil {
		return nil, err
	}
	payload := strings.NewReader(body)

	logRequest(service.Method, serviceUrl, body)

	request, err := http.NewRequest(service.Method, serviceUrl, payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, v := range service.Headers {
		for mk, mv := range v {
			mv2, err2 := ExecuteTemplate(mv, envData)
			if err2 != nil {
				return nil, err2
			}
			request.Header.Add(mk, mv2)
		}
	}

	logHeaders(request.Header)

	return request, nil
}

func logRequest(method, serviceUrl, body string) {
	fmt.Println("\nRequest created with:\nMethod:", method, "\nUrl:", serviceUrl, "\nBody:\n   ", body)
}

func logHeaders(headers http.Header) {
	fmt.Println("Headers:")
	for hk, hv := range headers {
		fmt.Println("   ", hk, ":", hv)
	}
}

// Create body string if header "Content-Type" with value "application/x-www-form-urlencoded" is present.
func createBodyString(replacedBody *interface{}, headers *Headers) string {

	v := reflect.ValueOf(replacedBody)
	if v.Kind() == reflect.Map && len(v.MapKeys()) == 0 {
		return ""
	}
	if v.Kind() == reflect.Slice && v.Len() == 0 {
		return ""
	}

	var builder BodyBuilder
	builder = &UnknownBuilder{}

	if isFormUrlEncoded(headers) {
		builder = &ToFormUrlEncodedBuilder{}
	} else if isApplicationJson(headers) {
		builder = &ToJsonBuilder{}
	} else if isRaw(headers) {
		builder = &RawBuilder{}
	}

	return builder.build(replacedBody)
}

const ContentType = "content-type"

func isFormUrlEncoded(headers *Headers) bool {
	return isMatchingHeader(headers, ContentType, "application/x-www-form-urlencoded")
}

func isApplicationJson(headers *Headers) bool {
	return isMatchingHeader(headers, ContentType, "application/json")
}

func isRaw(headers *Headers) bool {
	return isMatchingHeader(headers, ContentType, "text/plain")
}

func isMatchingHeader(headers *Headers, header string, headerValue string) bool {
	for _, v := range *headers {
		for mk, mv := range v {
			if strings.ToLower(mk) == header &&
				strings.HasPrefix(strings.ToLower(strings.TrimSpace(mv)), headerValue) {
				return true
			}
		}
	}
	return false
}

// It traverses the body and execute the template for all the string values in the body.
func executeTemplateForBody(data interface{}, envData *EnvironmentData) (interface{}, error) {

	switch v := reflect.ValueOf(data); v.Kind() {

	case reflect.Map:
		newMap := make(map[interface{}]interface{})
		for _, key := range v.MapKeys() {
			var err error
			newMap[key.Interface()], err = executeTemplateForBody(v.MapIndex(key).Interface(), envData)
			if err != nil {
				fmt.Println("Error happened at map key", key.Interface())
				return nil, err
			}
		}
		return newMap, nil

	case reflect.Slice:
		newSlice := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			var err error
			newSlice[i], err = executeTemplateForBody(v.Index(i).Interface(), envData)
			if err != nil {
				fmt.Println("Error happened at slice index", v.Index(i).Interface())
				return nil, err
			}
		}
		return newSlice, nil

	case reflect.String:
		replacedValue, err := ExecuteTemplate(v.String(), envData)
		if err != nil {
			fmt.Println("Error happened at template execution", v.String())
			return nil, err
		}
		if v.String() != replacedValue {
			return replacedValue, nil
		}

		return v.String(), nil

	default:
		return data, nil
	}
}
