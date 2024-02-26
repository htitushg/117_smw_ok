package main

import (
	"117_smw_ok/assets"
	"117_smw_ok/controllers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// On relie le fichier css et le favicon au nom static
	fmt.Printf("Main Chemin= %s\n", assets.Chemin+"assets/") //
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(assets.Chemin+"assets/"))))
	// routes
	mux.HandleFunc("GET /{$}", controllers.HomeBundle)
	mux.HandleFunc("GET /Login", controllers.LoginBundle)
	mux.HandleFunc("POST /Signin", controllers.SigninBundle)
	mux.HandleFunc("GET /Index", controllers.IndexBundle)
	mux.HandleFunc("POST /Logout", controllers.LogoutBundle)
	mux.HandleFunc("GET /Register", controllers.RegisterBundle)
	mux.HandleFunc("POST /AfficheUserInfo", controllers.AfficheUserInfoBundle)

	// start the server
	fmt.Printf("http://localhost%v , Cliquez sur le lien pour lancer le navigateur", assets.Port)
	log.Fatal(http.ListenAndServe(assets.Port, mux))
}
