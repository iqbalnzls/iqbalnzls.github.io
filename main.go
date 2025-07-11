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
	select {}
}

func convertText(this js.Value, args []js.Value) interface{} {
	input := args[0].String()

	b := bytes.TrimSpace([]byte(input))
	result := UnescapeString(string(b))

	return js.ValueOf(result)
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

func join(str []string) string {
	return "[" + strings.Join(str, ",") + "]"
}
