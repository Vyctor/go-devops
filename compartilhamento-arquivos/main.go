package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	auth "github.com/abbot/go-http-auth"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Print("Uso: go run main.go <diretorio> <porta>\n")
		os.Exit(1)
	}
	httpDir := os.Args[1]
	httpPort := os.Args[2]

	authenticator := auth.NewBasicAuthenticator("meuserver.com", Secret)
	http.HandleFunc("/", authenticator.Wrap(func(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
		http.FileServer(http.Dir(httpDir)).ServeHTTP(w, &r.Request)
	}))

	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}

func Secret(user, real string) string {
	if user == "vyctor" {
		return "$1$I7x8VHOM$gpizjN3X0CZrjFC18Beg7/"
	}
	return ""
}
