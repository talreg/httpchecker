package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

type errorLog struct {
	errorMessage string
	errorTime    time.Time
}

type trancactionResults struct {
	transactionError error
	transactionCode  int
}

var (
	good_hits              int
	bad_hits               int
	connection_errors      int
	current_env            string
	listening_port         int
	connection_errors_data []errorLog
	testing_address        string
	serverStart            time.Time
)

const (
	serverLoadFactor=2
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintf(w, "<html><head><title>http checker</title></head><body><H4>Http checker<br>You're running in %s mode.\r\n</H4>", current_env)
	fmt.Fprintf(w, "<br>running from: %s<br>", serverStart)
	fmt.Fprintf(w, "<br>we're checking: %s<br>", testing_address)
	fmt.Fprintf(w, "<h2>Total checks: %d<br> good ones: %d <br> bad ones: %d <br> connection errors: %d</h2>", good_hits+bad_hits, good_hits, bad_hits, connection_errors)
	if connection_errors > 0 {
		fmt.Fprint(w, "<br>connection errors:<br>")
		for index, error_data := range connection_errors_data {
			fmt.Fprintf(w, "%d:(%s) %s<br>", index, error_data.errorTime, error_data.errorMessage)
		}
	}
	fmt.Fprintf(w, "</body></html>")
	r.Body.Close()
}

func testAddress(address string, transChannel chan trancactionResults) {
	for {
		resp, err := http.Get(address)
		if err != nil {
			transChannel <- trancactionResults{transactionError: err, transactionCode: 0}
		}else{
			transChannel <- trancactionResults{transactionError: nil, transactionCode: resp.StatusCode}
			resp.Body.Close()
		}
	}
}

func loopAddresses(address string) {
	transChannel := make(chan trancactionResults)
	var result trancactionResults
	for i := 0; i < serverLoadFactor; i++ {
		go testAddress(address, transChannel)
	}
	for result = range transChannel {
		if result.transactionError != nil {
			connection_errors++
			connection_errors_data = append(connection_errors_data, errorLog{result.transactionError.Error(), time.Now()})
		} else {
			if result.transactionCode < 300 && result.transactionCode > 199 {
				good_hits++
			} else {
				bad_hits++
			}
		}
	}
}

func main() {
	current_env = os.Getenv("RUN_ENV")
	if current_env == "" {
		current_env = "development"
	}
	address := flag.String("address", "", "address to check, full path")
	flag.Parse()
	if *address == "" {
		fmt.Printf("bad address value")
		os.Exit(1)
	}
	testing_address = *address
	connection_errors_data = make([]errorLog, 0, 100)
	go loopAddresses(*address)
	listening_port = 3000
	fmt.Printf("starting server at %d\r\n", listening_port)
	serverStart = time.Now()
	http.HandleFunc("/version", func(arg1 http.ResponseWriter, arg2 *http.Request) {
		fmt.Fprint(arg1, "version 0.0.0.1")
	})
	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", listening_port), nil)
}
