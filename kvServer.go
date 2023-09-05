// kvServer
// server that interacts with kv
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
	"os"
    "net/http"

	util "github.com/prr123/utility/utilLib"
	"github.com/prr123/azulkv2/azulkvLib"
)

type myHandler struct{}

type apiKvObj struct {
	KvDb *azulkv2.KvObj
}

//func (*myHandler) ServeApi(w http.ResponseWriter, r *http.Request) {
func (kv *apiKvObj)ServeKvApi(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("received %s from %s\n", r.Method, r.RemoteAddr)
	azulkv2.PrintDb(kv.KvDb)

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

	keyDat := []byte{}
	valDat := []byte{}
	errStr :=""
	if len(parDat) >0 {
		aIdx := bytes.IndexByte(parDat, '&')
		if aIdx == -1 {
			keyStPos :=  bytes.IndexByte(parDat, '=')
			if keyStPos == -1 {
				//error
				errStr = "no key value"
			} else {
				keyDat = parDat[keyStPos+1:]
			}
		} else {
			keyStPos :=   bytes.IndexByte(parDat[:aIdx], '=')
			if keyStPos == -1 {
				//error
				errStr = "no key value"
			} else {
				keyDat = parDat[keyStPos+1:aIdx]
			}
			valbData :=parDat[aIdx+1:]
			valStPos := bytes.IndexByte(valbData, '=')
			if valStPos == -1 {
				//error
				errStr = "no value"
			} else {
				valDat = valbData[valStPos+1:]
			}
		}
		if len(errStr) > 0 {
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 errer parameter: %s\n", errStr)
			return
		}
	}

	if len(errStr) == 0 {errStr = "none"}
	fmt.Fprintf(w, "cmd: %s err: %s key: %s value: %s\n", string(cmdDat), errStr, string(keyDat), string(valDat))

	switch string(cmdDat) {
	case "add":
		if len(keyDat) <1 || len(valDat) < 1 {
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 add: no key val parameters!\n")
			return
		}
		fmt.Fprintf(w, "cmd: add key: %s val: %s\n", string(keyDat), string(valDat))

	case "upd":
		if len(keyDat) <1 || len(valDat) < 1 {
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 upd: no key val parameters!\n")
			return
		}
		fmt.Fprintf(w, "cmd: upd key: %s val: %s\n", string(keyDat), string(valDat))

	case "del":
		if len(keyDat) <1{
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 add: no key val parameter!\n")
			return
		}
		if len(valDat) >0{
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 del: no val parameter required!\n")
			return
		}
		fmt.Fprintf(w, "cmd: del key: %s\n", string(keyDat))

	case "get":
		if len(keyDat) <1{
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 add: no key val parameter!\n")
			return
		}
		if len(valDat) >0{
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 del: no val parameter required!\n")
			return
		}
		fmt.Fprintf(w, "cmd: del key: %s\n", string(keyDat))

	case "list":
		if len(keyDat) >0 || len(valDat) >0 {
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 list: key val parameters not valid\n")
			return
		}
		fmt.Fprintf(w, "cmd: list\n")

	case "entries":
		if len(keyDat) >0 || len(valDat) >0 {
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 entries: key val parameters not valid\n")
			return
		}
		fmt.Fprintf(w, "cmd: entries\n")

//	case "set":

	case "info":
		if len(keyDat) >0 || len(valDat) >0 {
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 info: key val parameters not valid\n")
			return
		}
		fmt.Fprintf(w, "cmd: info\n")

	default:
			w.WriteHeader(http.StatusNotImplemented)
    		fmt.Fprintf(w, "502 error unknown command: %s\n", string(cmdDat))
			return
	}


//    io.WriteString(w, "URL: " + r.URL.String())

}


func main(){

    numarg := len(os.Args)
    dbg := false
    flags:=[]string{"dbg","db", "port"}

    useStr := "./kvServer [/db=dbfolder] [/port=] [/dbg]"
    helpStr := "http Server for azulkv kv-store\n"

    if numarg > 4 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: ./template [/flag1=] [/flag2]\n", useStr)
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

	portStr := "10800"
    ptval, ok := flagMap["port"]
    if !ok {
        log.Printf("default port: %s\n", portStr)
    } else {
        if ptval.(string) == "none" {log.Fatalf("error: no port provided!")}
        portStr = ptval.(string)
        log.Printf("port: %s\n", portStr)
    }


	kvDb, err :=  azulkv2.InitKV(dbDir ,dbg)
	if err != nil {log.Fatalf("error InitKv: %v", err)}
	kv := apiKvObj{
		KvDb: kvDb,
	}

    mux := http.NewServeMux()

    mux.HandleFunc("/", kv.ServeKvApi)
//    mux.Handle("/",&myHandler{})

    //http.ListenAndServe uses the default server structure.
	log.Printf("listening on port: %s!\n", portStr)
	port:= ":" + portStr
    err = http.ListenAndServe(port, mux)
    if err != nil {log.Fatalf("ListenAndServe: %v\n", err)}

}
