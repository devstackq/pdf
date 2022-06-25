package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"pdf_reader/pkg/pdf"
)

// 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
func test(w http.ResponseWriter, r *http.Request) {
	pdf := pdf.NewPdf()
	resp := pdf.Read(nil)

	err := json.NewEncoder(w).Encode(&resp)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println(resp)
}

// todo : integrate with arcana; rest/ service/ db; then test
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", test)
	http.ListenAndServe(":6969", mux)
}
