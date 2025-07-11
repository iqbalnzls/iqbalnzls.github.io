package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type ConversionRequest struct {
	Text string `json:"text"`
}

type ConversionResponse struct {
	Result string `json:"result"`
}

func main() {
	// Serve static files
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Handle conversion endpoint
	http.HandleFunc("/convert", handler)

	log.Println("Server starting on :8282")
	if err := http.ListenAndServe(":8282", nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ConversionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b := bytes.TrimSpace([]byte(req.Text))

	result := UnescapeString(string(b))

	response := ConversionResponse{
		Result: result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
