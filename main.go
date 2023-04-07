package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/ping", Pong).Methods("GET")
	r.HandleFunc("/email-validator", EmailValidator).Methods("GET")

	http.ListenAndServe(":3000", r)

}

func EmailValidator(w http.ResponseWriter, r *http.Request) {
	domain := r.URL.Query().Get("domain")
	// domain := "a.in"
	response := make(map[string]interface{})
	response["domain"] = domain
	mxRecords, _ := net.LookupMX(domain)

	if len(mxRecords) > 0 {
		response["has_mx_records"] = true
		fmt.Println("Has MX records true")
	}
	txtRecords, _ := net.LookupTXT(domain)

	for _, rec := range txtRecords {
		// fmt.Println(fmt.Sprintf("%v", rec))
		if strings.HasPrefix(rec, "v=spf1") {
			fmt.Println("Has SPF true")
			response["has_spf_records"] = true
			response["spf_records"] = fmt.Sprintf("%v", rec)

			break
		}
	}

	dmarcRecords, _ := net.LookupTXT("_dmarc." + domain)

	for _, rec := range dmarcRecords {
		fmt.Println(fmt.Sprintf("%v", rec))
		if strings.HasPrefix(rec, "v=DMARC1") {
			fmt.Println("Has DMARC true")
			response["has_dmarc_records"] = true
			response["dmarc_records"] = fmt.Sprintf("%v", rec)
			// fmt.Println(fmt.Sprintf("%v", rec))
			// break
		}
	}
	resp, _ := json.Marshal(response)
	w.Write(resp)
	w.WriteHeader(http.StatusOK)
}

func Pong(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["message"] = "Success"
	resp, _ := json.Marshal(response)
	w.Write(resp)
	w.WriteHeader(http.StatusOK)

}
