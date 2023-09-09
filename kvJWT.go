// kvJWT
// test server for login
//
// Author: prr
// Date: 8 Sept 2023
// copyright 2023 prr, azul software
//


package main

import (
//    "io"
	"fmt"
    "log"
	"bytes"
	"os"
	"time"
    "net/http"
	"io/ioutil"
	"strings"

	util "github.com/prr123/utility/utilLib"
	"github.com/prr123/azulkv2/azulkvLib"

    "github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v4"

)

type myHandler struct{}

type apiKvObj struct {
	KvDb *azulkv2.KvObj
}

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"peter": "pass1",
	"user1": "password1",
	"user2": "password2",
}


//func (*myHandler) ServeApi(w http.ResponseWriter, r *http.Request) {

func WelcomeApi(w http.ResponseWriter, r *http.Request) {
	urlDat := []byte(r.URL.String())
	log.Printf("Welcome: %s\n", string(urlDat))


	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))

//	fmt.Fprintf(w, "Welcome\n")

}


func RefreshApi(w http.ResponseWriter, r *http.Request) {
	urlDat := []byte(r.URL.String())
	log.Printf("Refresh: %s\n", string(urlDat))

	fmt.Fprintf(w, "Refresh\n")

//func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code until this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code until this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}


func LogoutApi(w http.ResponseWriter, r *http.Request) {
	urlDat := []byte(r.URL.String())
	log.Printf("Logout: %s\n", string(urlDat))

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
	fmt.Fprintf(w, "Logout\n")

}

//func (kv *apiKvObj)LoginApi(w http.ResponseWriter, r *http.Request) {
func LoginApi(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("received %s from %s\n", r.Method, r.RemoteAddr)

	urlDat := []byte(r.URL.String())
	log.Printf("Login: %s\n", string(urlDat))

	idx := bytes.IndexByte(urlDat,'?')
	cmdLen :=0
	if idx == -1 {cmdLen = len(urlDat)} else { cmdLen = idx}

	//todo replace with hash table


	if string(urlDat[:cmdLen]) != "/signin" {
		log.Printf("cmd %s\n", string(urlDat[:cmdLen]))
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

	creds := Credentials{}
	// Get the JSON body and decode into credentials
//	err = json.NewDecoder(r.Body).Decode(&creds)
//    NewUser:= user{}
    err = json.Unmarshal(reqBody, &creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		fmt.Printf("NewDecoder Error: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("creds: %v\n", creds)

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}


func main(){

    numarg := len(os.Args)
    dbg := false
    flags:=[]string{"dbg","db", "port"}

    useStr := "./kvJWT [/db=dbfolder] [/port=] [/dbg]"
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

    mux.HandleFunc("/signin", LoginApi)
    mux.HandleFunc("/welcome", WelcomeApi)
    mux.HandleFunc("/refresh", RefreshApi)
    mux.HandleFunc("/logout", LogoutApi)

    //http.ListenAndServe uses the default server structure.
	log.Printf("listening on port: %s!\n", portStr)
	port:= ":" + portStr
    err = http.ListenAndServe(port, mux)
    if err != nil {log.Fatalf("ListenAndServe: %v\n", err)}

}
