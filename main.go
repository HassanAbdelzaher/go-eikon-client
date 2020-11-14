package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var mainport int = 9000
var secoundaryport = 36036
var url = "http://localhost:9000/api/v1/data"

func main() {
	http.HandleFunc("/api/request", makeEikonRequest)
	http.HandleFunc("/api/screen", screen)
	http.HandleFunc("/api/exchange", exchange)
	http.HandleFunc("/api/country", country)
	http.HandleFunc("/", home)
	http.ListenAndServe(":2017", nil)
}
func home(w http.ResponseWriter, hreq *http.Request) {
	fmt.Fprintf(w, "welcom")
}
func makeEikonRequest(w http.ResponseWriter, hreq *http.Request) {
	fmt.Println("URL:>", url)
	ek := newEikonRequest([]string{"IBM"}, []string{"TR.InterimNumber", "TR.InterimNumber(Period=FQ0).Date", "TR.InterimNumber(Period=FY0).Date", "TR.InterimNumber.Date"}, "FQ-1")
	//var jsonStr = []byte(`{"Entity":{"E":"DataGrid_StandardAsync","W":{"requests":[{"instruments":["IBM"],"fields":[{"name":"TR.Ric","parameters":{}},{"name":"TR.InterimNumber","parameters":{}},{"name":"TR.InterimNumber.Date","parameters":{}}],"parameters":{}}]}}}`)
	jsonStr, err := json.Marshal(ek)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("x-tr-applicationid", "ecb7a4ef2ef3477087f551bfe9dc7d6476bc1919")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	fmt.Fprintf(w, string(body))
}

//SCREEN(U(IN(Equity(active or inactive,public,countryprimaryquote))/*UNV:Public*/),IN(TR.ExchangeMarketIdCode,{0}))
func screen(w http.ResponseWriter, hreq *http.Request) {
	fmt.Println("URL:>", url)
	rics, ok := hreq.URL.Query()["rics"]
	if !ok {
		w.WriteHeader(500)
		fmt.Fprintf(w, "please provide rics")
		return
	}
	fields, ok := hreq.URL.Query()["fields"]
	if !ok {
		w.WriteHeader(500)
		fmt.Fprintf(w, "please provide fields")
		return
		//fields = []string{"TR.RIC", "TR.CompanyName", "TR.ISIN", "TR.InterimNumber(Period=FQ0)", "TR.InterimNumber(Period=FQ0).Date", "TR.InterimNumber(Period=FY0).Date"}
	}
	period, ok := hreq.URL.Query()["period"]
	if !ok {
		w.WriteHeader(500)
		fmt.Fprintf(w, "please provide period")
		return
		//fields = []string{"TR.RIC", "TR.CompanyName", "TR.ISIN", "TR.InterimNumber(Period=FQ0)", "TR.InterimNumber(Period=FQ0).Date", "TR.InterimNumber(Period=FY0).Date"}
	}
	ek := newEikonRequest(rics, fields, period[0])
	//var jsonStr = []byte(`{"Entity":{"E":"DataGrid_StandardAsync","W":{"requests":[{"instruments":["IBM"],"fields":[{"name":"TR.Ric","parameters":{}},{"name":"TR.InterimNumber","parameters":{}},{"name":"TR.InterimNumber.Date","parameters":{}}],"parameters":{}}]}}}`)
	jsonStr, err := json.Marshal(ek)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Println(string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("x-tr-applicationid", "ecb7a4ef2ef3477087f551bfe9dc7d6476bc1919")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, err.Error())
		return
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	fmt.Fprintf(w, string(body))
}

func exchange(w http.ResponseWriter, hreq *http.Request) {
	fmt.Println("URL:>", url)
	code, ok := hreq.URL.Query()["exchangeCode"]
	if !ok {
		fmt.Fprintf(w, "please provide exchangeCode")
		return
	}
	rics := fmt.Sprintf(`SCREEN(U(IN(Equity(public,active,countryprimaryquote))),IN(TR.ExchangeMarketIdCode,'%s'))`, code[0])
	fields, ok := hreq.URL.Query()["fields"]
	if !ok {
		fields = []string{"TR.RIC"}
	}
	ek := newEikonRequest([]string{rics}, fields, "")
	//var jsonStr = []byte(`{"Entity":{"E":"DataGrid_StandardAsync","W":{"requests":[{"instruments":["IBM"],"fields":[{"name":"TR.Ric","parameters":{}},{"name":"TR.InterimNumber","parameters":{}},{"name":"TR.InterimNumber.Date","parameters":{}}],"parameters":{}}]}}}`)
	jsonStr, err := json.Marshal(ek)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("x-tr-applicationid", "ecb7a4ef2ef3477087f551bfe9dc7d6476bc1919")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 180 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	fmt.Fprintf(w, string(body))
}

func country(w http.ResponseWriter, hreq *http.Request) {
	fmt.Println("URL:>", url)
	code, ok := hreq.URL.Query()["exchangeCode"]
	if !ok {
		fmt.Fprintf(w, "please provide exchangeCode")
		return
	}
	rics := fmt.Sprintf(`SCREEN(U(IN(Equity(public,active))),IN(TR.CoRTradingCountryCode,%s))`, code[0])
	fields, ok := hreq.URL.Query()["fields"]
	if !ok {
		fields = []string{"TR.RIC"}
	}
	ek := newEikonRequest([]string{rics}, fields, "")
	//var jsonStr = []byte(`{"Entity":{"E":"DataGrid_StandardAsync","W":{"requests":[{"instruments":["IBM"],"fields":[{"name":"TR.Ric","parameters":{}},{"name":"TR.InterimNumber","parameters":{}},{"name":"TR.InterimNumber.Date","parameters":{}}],"parameters":{}}]}}}`)
	jsonStr, err := json.Marshal(ek)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("x-tr-applicationid", "ecb7a4ef2ef3477087f551bfe9dc7d6476bc1919")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 180 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	fmt.Fprintf(w, string(body))
}
