package main

import (
	"fmt"
	"strings"
)

type BodyBuilder interface {
	build(b *Body) string
}

type ToJsonBuilder struct {
}

type ToFormUrlEncodedBuilder struct {
}

type UnknownBuilder struct {
}

func (builder *ToJsonBuilder) build(b *Body) string {
	//  TODO: marshal the body
	var keyValuePairs []string

	for mk, mv := range *b {
		keyValuePairs = append(keyValuePairs, "\""+mk+"\": \""+mv+"\"")
	}

	return "{" + strings.Join(keyValuePairs, ",") + "}"
}

func (builder *ToFormUrlEncodedBuilder) build(b *Body) string {
	var keyValuePairs []string

	for mk, mv := range *b {
		keyValuePairs = append(keyValuePairs, mk+"="+mv)
	}

	return strings.Join(keyValuePairs, "&")
}

func (builder *UnknownBuilder) build(b *Body) string {
	// no-op
	fmt.Println("Warning: Undetermined body builder is used.",
		"Provide an \"application/json\" or similar \"Content-Type\" header.")
	return ""
}
