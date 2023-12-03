package main

import (
	"fmt"
	"net/http"
	"strings"
)

func CreateRequest(service *ServiceData, envData *EnvironmentData) (*http.Request, error) {
	for k, v := range service.Body {
		replacedValue, err := ExecuteTemplate(v, envData)
		if err != nil {
			return nil, err
		}
		if v != replacedValue {
			service.Body[k] = replacedValue
		}
	}

	body := createBodyString(service)

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
	splitBody := strings.ReplaceAll(body, "&", "\n    ")
	fmt.Println("\nRequest created with:\nMethod:", method, "\nUrl:", serviceUrl, "\nBody:\n   ", splitBody)
}

func logHeaders(headers http.Header) {
	fmt.Println("Headers:")
	for hk, hv := range headers {
		fmt.Println("   ", hk, ":", hv)
	}
}

// Create body string if header "Content-Type" with value "application/x-www-form-urlencoded" is present.
func createBodyString(service *ServiceData) string {
	if len(service.Body) == 0 {
		return ""
	}

	var builder BodyBuilder
	builder = &UnknownBuilder{}

	if isFormUrlEncoded(&service.Headers) {
		builder = &ToFormUrlEncodedBuilder{}
	} else if isApplicationJson(&service.Headers) {
		builder = &ToJsonBuilder{}
	}

	return builder.build(&service.Body)
}

func isFormUrlEncoded(headers *Headers) bool {
	return isMatchingHeader(headers, "content-type", "application/x-www-form-urlencoded")
}

func isApplicationJson(headers *Headers) bool {
	return isMatchingHeader(headers, "content-type", "application/json")
}

func isMatchingHeader(headers *Headers, header string, headerValue string) bool {
	for _, v := range *headers {
		for mk, mv := range v {
			if strings.ToLower(mk) == header &&
				strings.ToLower(mv) == headerValue {
				return true
			}
		}
	}
	return false
}
