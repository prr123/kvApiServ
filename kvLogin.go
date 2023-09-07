// kvLogin
// test server for login
//
// Author: prr
// Date: 7 Sept 2023
// copyright 2023 prr, azul software
//


package main

import (
//    "io"
	"fmt"
    "log"
//	"bytes"
	"os"
    "net/http"
	"io/ioutil"
	"strings"

	util "github.com/prr123/utility/utilLib"
	"github.com/prr123/azulkv2/azulkvLib"
)

type myHandler struct{}

type apiKvObj struct {
	KvDb *azulkv2.KvObj
}

//func (*myHandler) ServeApi(w http.ResponseWriter, r *http.Request) {
//func (kv *apiKvObj)LoginApi(w http.ResponseWriter, r *http.Request) {
func LoginApi(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("received %s from %s\n", r.Method, r.RemoteAddr)

	urlDat := []byte(r.URL.String())
	log.Printf("urlDat: %s\n", string(urlDat))

	//todo replace with hash table
	if string(urlDat[:6]) != "/login" {
		log.Printf("path %s\n", string(urlDat))
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintf(w, "501 no match: \n")
		return
	}
	log.Printf("path %s\n", string(urlDat))
	fmt.Printf("server: headers:\n")
	for headerName, headerValue := range r.Header {
		  fmt.Printf("\t%s = %s\n", headerName, strings.Join(headerValue, ", "))
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("server: could not read request body: %s\n", err)
	}
	fmt.Printf("server: request body: %s\n", reqBody)
}


func main(){

    numarg := len(os.Args)
    dbg := false
    flags:=[]string{"dbg","db", "port"}

    useStr := "./kvLogin [/db=dbfolder] [/port=] [/dbg]"
    helpStr := "http Server to test login\n"

    if numarg > 4 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage is: %s\n", useStr)
        os.Exit(-1)
    }

    if numarg > 1 && os.Args[1] == "help" {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage is: %s\n", useStr)
        os.Exit(1)
    }

    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    _, ok := flagMap["dbg"]
    if ok {dbg = true}
    if dbg {
        fmt.Printf("dbg -- flag list:\n")
        for k, v :=range flagMap {
            fmt.Printf("  flag: /%s value: %s\n", k, v)
        }
    }

	dbDir := "$HOME/dbtest"
    dbval, ok := flagMap["db"]
    if !ok {
        log.Printf("default db: %s\n", dbDir)
    } else {
        if dbval.(string) == "none" {log.Fatalf("error: no dir provided!")}
        dbDir = dbval.(string)
        log.Printf("csrList: %s\n", dbDir)
    }

	portStr := "10901"
    ptval, ok := flagMap["port"]
    if !ok {
        log.Printf("default port: %s\n", portStr)
    } else {
        if ptval.(string) == "none" {log.Fatalf("error: no port provided!")}
        portStr = ptval.(string)
        log.Printf("port: %s\n", portStr)
    }



    mux := http.NewServeMux()

    mux.HandleFunc("/", LoginApi)
//    mux.Handle("/",&myHandler{})

    //http.ListenAndServe uses the default server structure.
	log.Printf("listening on port: %s!\n", portStr)
	port:= ":" + portStr
    err = http.ListenAndServe(port, mux)
    if err != nil {log.Fatalf("ListenAndServe: %v\n", err)}

}
