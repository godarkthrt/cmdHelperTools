package main_test

import (
	"fmt"
	"net/url"
	"testing"
)

func TestEncodeUrlFunctionality(t *testing.T) {
	escapedUrl := url.QueryEscape("hollo { +how+ are you }")
	t.Log("Escaped url is : ", escapedUrl)
	fmt.Println("Escaped url is : ", escapedUrl)

}
