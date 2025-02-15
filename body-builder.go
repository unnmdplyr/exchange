package main

import (
	"fmt"
	"reflect"
	"strings"
)

type BodyBuilder interface {
	build(b *interface{}) string
}

type ToJsonBuilder struct {
}

type ToFormUrlEncodedBuilder struct {
}

type RawBuilder struct {
}

type UnknownBuilder struct {
}

func (builder *ToJsonBuilder) build(b *interface{}) string {
	body := marshal(*b)
	return body
}

func (builder *ToFormUrlEncodedBuilder) build(b *interface{}) string {
	var keyValuePairs []string

	// Check if the input is a map
	if reflect.ValueOf(*b).Kind() != reflect.Map {
		fmt.Println("!!!! The body is not a map in the FormUrlEncodedBuilder.")
		return ""
	}

	iter := reflect.ValueOf(*b).MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		pair := fmt.Sprintf("%v=%v", k, v)
		keyValuePairs = append(keyValuePairs, pair)
	}

	return strings.Join(keyValuePairs, "&")
}

func (builder *RawBuilder) build(b *interface{}) string {
	body := marshal(*b)
	return body
}

func (builder *UnknownBuilder) build(b *interface{}) string {
	// no-op
	fmt.Println("Warning: Undetermined body builder is used.",
		"Provide an \"application/json\" or similar \"Content-Type\" header.")
	return ""
}

// It traverses the body and marshals the body into a string.
func marshal(data interface{}) string {

	switch v := reflect.ValueOf(data); v.Kind() {

	case reflect.Map:
		var mapContent []string
		for _, key := range v.MapKeys() {
			mapContent = append(mapContent, "\""+fmt.Sprintf("%v", key)+"\": "+marshal(v.MapIndex(key).Interface()))
		}
		return "{ " + strings.Join(mapContent, ",") + " }"

	case reflect.Slice:
		var sliceContent []string
		for i := 0; i < v.Len(); i++ {
			sliceContent = append(sliceContent, marshal(v.Index(i).Interface()))
		}
		return "[ " + strings.Join(sliceContent, ",") + " ]"

	case reflect.String:
		return "\"" + v.String() + "\""

	default:
		return fmt.Sprintf("%v", data)
	}
}
