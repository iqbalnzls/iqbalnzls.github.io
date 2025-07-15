package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"syscall/js"
)

type ConversionRequest struct {
	Text string `json:"text"`
}

type ConversionResponse struct {
	Result string `json:"result"`
}

func main() {
	js.Global().Set("convertText", js.FuncOf(convertText))
	js.Global().Set("convertTextK8s", js.FuncOf(convertTextK8s))
	select {}
}

func convertText(this js.Value, args []js.Value) interface{} {
	input := args[0].String()
	b := bytes.TrimSpace([]byte(input))
	result := UnescapeString(string(b))
	return js.ValueOf(result) // NO PROMISE NEEDED HERE
}

func convertTextK8s(this js.Value, args []js.Value) interface{} {
	input := args[0].String()
	b := bytes.TrimSpace([]byte(input))
	result := UnescapeStringKube(string(b))
	return js.ValueOf(result) // NO PROMISE NEEDED HERE
}

func UnescapeString(str string) string {
	var (
		temp  = make([]string, 0)
		temp2 = make([][]string, 0)
	)

	if str[:1] != "[" && str[len(str)-1:] != "]" {
		str = "[" + str + "]"
	}

	if err := json.Unmarshal([]byte(str), &temp); err != nil {
		temp = nil
		if err = json.Unmarshal([]byte(str), &temp2); err != nil {
			panic(err)
		}
	}

	switch l := len(temp); {
	case l > 0:
		return join(temp)
	case l == 0:
		tempResult := make([]string, 0)
		for _, v := range temp2 {
			tempResult = append(tempResult, join(v))
		}
		return join(tempResult)
	}
	return ""
}

func UnescapeStringKube(str string) string {
	type mes struct {
		Message string
	}

	var (
		temp  = make([]string, 0)
		temp2 = make([][]string, 0)
	)

	if str[:1] != "[" && str[len(str)-1:] != "]" {
		str = "[" + str + "]"
	}

	if err := json.Unmarshal([]byte(str), &temp); err != nil {
		temp = nil
		if err = json.Unmarshal([]byte(str), &temp2); err != nil {
			panic(err)
		}
	}

	switch l := len(temp); {
	case l > 0:
		res := make([]string, 0)
		for _, s := range temp {
			mTemp := mes{}
			if err := json.Unmarshal([]byte(s), &mTemp); err != nil {
				panic(err)
			}
			res = append(res, mTemp.Message)
		}

		return join(res)
	case l == 0:
		tempResult := make([]string, 0)
		for _, v := range temp2 {
			tem := make([]string, 0)
			for _, i := range v {
				mTemp := mes{}
				if err := json.Unmarshal([]byte(i), &mTemp); err != nil {
					panic(err)
				}
				tem = append(tem, mTemp.Message)
			}

			tempResult = append(tempResult, join(tem))
		}

		return join(tempResult)
	}
	return ""
}

func join(str []string) string {
	return "[" + strings.Join(str, ",") + "]"
}
