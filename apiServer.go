// apiServer
//
// Author: prr
// Date: 2 Sept 2023
// copyright 2023 prr, azul software
//


package main

import (
//    "io"
	"fmt"
    "log"
	"bytes"
    "net/http"
)

type myHandler struct{}

//func (*myHandler) ServeApi(w http.ResponseWriter, r *http.Request) {
func ServeApi(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("received %s from %s\n", r.Method, r.RemoteAddr)

	urlDat := []byte(r.URL.String())
	log.Printf("urlDat: %s\n", string(urlDat))
	dbStr := string(urlDat[:4])
	if dbStr != "/db/" {
		w.WriteHeader(http.StatusNotImplemented)
    	fmt.Fprintf(w, "501 no match: %s\n", dbStr)
		return
	}
	parPos := bytes.IndexByte(urlDat, '?')
	cmdDat := []byte{}
	parDat := []byte{}
	if parPos == -1 {
		cmdDat = urlDat[4:]
	} else {
		cmdDat = urlDat[4:parPos]
		parDat = urlDat[parPos+1:]
	}
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Proto:  %s\n", r.Proto)
	fmt.Fprintf(w, "Request URI:  %s\n", r.RequestURI)

	for k,v := range r.Header {
		fmt.Fprintf(w, "key: %s value: %s\n", k, v)
	}

	fmt.Fprintf(w, "cmd: %s par: %s\n", string(cmdDat), string(parDat))
//    io.WriteString(w, "URL: " + r.URL.String())




}

/*
func Api(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("received %s from %s\n", r.Method, r.RemoteAddr)

	fmt.Fprintf(w, "hello fron db handler\n")
	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Proto:  %s\n", r.Proto)
	fmt.Fprintf(w, "Request URI:  %s\n", r.RequestURI)

	for k,v := range r.Header {
		fmt.Fprintf(w, "key: %s value: %s\n", k, v)
	}

    io.WriteString(w, "URL: " + r.URL.String())

    io.WriteString(w, "Api")
}
*/

func main(){


    mux := http.NewServeMux()

    mux.HandleFunc("/", ServeApi)
//    mux.Handle("/",&myHandler{})

    //http.ListenAndServe uses the default server structure.
	log.Printf("listening on port 10800!\n")
    err := http.ListenAndServe(":10800", mux)
    if err != nil {log.Fatalf("ListenAndServe: %v\n", err)}

}
