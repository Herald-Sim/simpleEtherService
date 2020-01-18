package main

import (
	"api"
	"fmt"
	"log"
	"net/http"
	"router"
)

func main() {
	// API routing table
	webRouter := &router.Router{make(map[string]map[string]router.HandlerFunc)}

	// For test..
	webRouter.HandleFunc("GET", "/users/:user_id/addresses/:address_id", func(c *router.Context) {
		fmt.Fprintf(c.ResponseWriter, "Retrive user's address\nuser:%s\naddress:%s\n", c.Params["user_id"], c.Params["address_id"])
	})

	// Sign-up
	webRouter.HandleFunc("POST", "/user", api.SignUp)
	// Sign-in
	webRouter.HandleFunc("POST", "/user/:user_id", api.SignIn)

	// Add user ether wallet
	webRouter.HandleFunc("GET", "/addWallet/:user_id", api.JWTauth(api.AddWallet))

	// Read contract
	webRouter.HandleFunc("GET", "/contract", api.JWTauth(api.ReadContract))
	// For test without JWT
	//webRouter.HandleFunc("GET", "/contract", api.ReadContract)

	// Write contract
	webRouter.HandleFunc("POST", "/contract", api.JWTauth(api.WriteContract))
	// For test without JWT
	//webRouter.HandleFunc("POST", "/contract", api.WriteContract)

	// Refresh JWT Toekn
	webRouter.HandleFunc("GET", "/token", api.Refresh)

	/*
		// HTTPS
		log.Fatal(http.ListenAndServeTLS(":443",
			"/etc/letsencrypt/live/appserver.acewallet.net/fullchain.pem",
			"/etc/letsencrypt/live/appserver.acewallet.net/privkey.pem",
			webRouter))

	*/

	// GET user wallet meta info
	//webRouter.HandleFunc("GET", "/wallet/:walletAddress", api.JWTauth(api.WalletInfo))

	// POST token (user -> user)
	//webRouter.HandleFunc("POST", "/HRToken", api.JWTauth(api.TransferToken))

	// HTTP
	log.Fatal(http.ListenAndServe(":8080", webRouter))
}
