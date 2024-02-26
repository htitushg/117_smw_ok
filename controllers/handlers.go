package controllers

import (
	"117_session_ok/assets"
	"117_session_ok/models"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"time"
)

// Ajouté le 02/02/2024
// Note - NOT RFC4122 compliant

// Génère un UUID (Jeton de session : Token)
func pseudo_uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}

// Fin de l'Ajout du 02/02/2024

// Si la session est valide, renvoie le Token et true, sinon nil et false
func SessionValide(w http.ResponseWriter, r *http.Request) (stoken string, resultat bool) {
	c, err := r.Cookie("session_token")
	resultat = false
	stoken = ""
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return stoken, resultat
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return stoken, resultat
	}
	stoken = c.Value
	_, exists := assets.Sessions[stoken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return stoken, resultat
	}
	// If the previous session is valid, create a new session token for the current user
	// on peut utiliser google : "github.com/google/uuid"
	// ou bien pseudo_uuid() fonction ci dessus qui utilise "crypto/rand"

	/* newSessionToken := uuid.NewString() */
	newSessionToken := pseudo_uuid()
	maxAge := 120

	// Set the token in the session map, along with the user whom it represents
	assets.Sessions[newSessionToken] = assets.Session{
		Pseudo: assets.Sessions[stoken].Pseudo,
		MaxAge: maxAge,
	}
	// Delete the older session token
	delete(assets.Sessions, stoken)
	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  newSessionToken,
		MaxAge: maxAge,
	})
	resultat = true
	return newSessionToken, resultat
}

// Controlleur Home: Affiche le Page publique(home) si la session n'est pas valide, sinon affiche la page privée(index)
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Home log: UrlPath: %#v\n", r.URL.Path) // testing
	var data assets.Data
	var t *template.Template
	var err error
	stoken, exists := SessionValide(w, r)
	if !exists {
		t, err = template.ParseFiles(assets.Chemin + "templates/home.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		// Il nous faut ici rassembler les infos utilisateur
		DJour := time.Now().Format("2006-01-02")
		data.CSessions = assets.Sessions[stoken]
		data.Date_jour = DJour
		data.SToken = stoken
		/* data.Email:       credsR.Email
		data.Firstname:   credsR.Firstname
		data.Lastname:   credsR.Lastname */

		t, err = template.ParseFiles(assets.Chemin + "templates/index.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	if err = t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Controlleur Login: Affiche la page de connexion
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Login log: UrlPath: %#v\n", r.URL.Path)
	var message template.HTML
	switch r.URL.Query().Get("err") {
	case "pass":
		message = "<div class=\"message\">Wrong username or password!</div>"
	case "restricted":
		message = "<div class=\"message\">You need to sign in to access to this resource!</div>"
	}
	t, err := template.ParseFiles(assets.Chemin + "templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, message); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Controlleur Signin: Traite les informations de connexin fournies par login,
// Si les informations sont incorrectes, renvoie vers home
// Sinon crée la session et renvoie vers index
func Signin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Signin log: UrlPath: %#v\n", r.URL.Path)
	var creds assets.Credentials
	/* var data assets.Data
	var t *template.Template
	var err error */
	creds.Pseudo = r.FormValue("pseudo")
	creds.Password = r.FormValue("passid")
	// Get the expected password from our in memory map
	expectedPassword, ok := assets.Users[creds.Pseudo]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/Login?err=pass", http.StatusSeeOther)
		return
	} else {
		// Create a new random session token
		// on peut utiliser google : "github.com/google/uuid"
		// ou bien pseudo_uuid() fonction ci dessus qui utilise "crypto/rand"

		/* newSessionToken := uuid.NewString() */
		sessionToken := pseudo_uuid()
		maxAge := 120
		// Set the token in the session map, along with the user whom it represents
		assets.Sessions[sessionToken] = assets.Session{
			Pseudo: creds.Pseudo,
			MaxAge: maxAge,
		}
		// Finally, we set the client cookie for "session_token" as the session token we just generated
		// we also set an expiry time of 120 seconds
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  sessionToken,
			MaxAge: maxAge,
		})
		http.Redirect(w, r, "/Index", http.StatusSeeOther)
		/* //DatedeCreation := assets.Sessions[sessionToken].Expiry.Format("“2006-01-02 15h04 05 secondes”")
			DJour := time.Now().Format("2006-01-02")
			data.CSessions = assets.Sessions[sessionToken]
			data.Date_jour = DJour
			data.SToken = sessionToken

			t, err = template.ParseFiles(assets.Chemin + "templates/index.html")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}
		if err := t.Execute(w, data); err != nil {
			http.Error(w, err.Error(), 500)
			return*/
	}
}
func Index(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session_token")
	var sessionToken = cookie.Value
	var data = assets.Data{
		SToken:    sessionToken,
		CSessions: assets.Sessions[sessionToken],
		Date_jour: time.Now().Format("2006-01-02"),
	}
	t, err := template.ParseFiles(assets.Chemin + "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Controlleur Logout: Si la session est valide, Ferme la session et renvoie vers home
func Logout(w http.ResponseWriter, r *http.Request) {
	sessionToken, exists := SessionValide(w, r)
	if exists {
		// remove the users session from the session map
		delete(assets.Sessions, sessionToken)
		// We need to let the client know that the cookie is expired
		// In the response, we set the session token to an empty
		// value and set its expiry as the current time
		http.SetCookie(w, &http.Cookie{
			Name:   "session_token",
			Value:  "",
			MaxAge: -1,
		})
	}
	t, err := template.ParseFiles(assets.Chemin + "templates/home.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Controlleur AfficheUserInfo: Si la session est valide renvoie vers afficheuserinfo
// Sinon renvoie vers home
func AfficheUserInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("AfficheUserInfo log: UrlPath: %#v\n", r.URL.Path) // testing
	var data assets.Data
	var err error
	var t *template.Template
	sessionToken, exists := SessionValide(w, r)
	if exists {
		DJour := time.Now().Format("2006-01-02")
		data.CSessions = assets.Sessions[sessionToken]
		data.Date_jour = DJour
		data.SToken = sessionToken
		t, err = template.ParseFiles(assets.Chemin + "templates/afficheuserinfo.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		t, err = template.ParseFiles(assets.Chemin + "templates/home.html")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Ajouté le 02/02/2024
// for GET

// Controlleur Register: Renvoie vers register pour enregistrement
func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Register log: UrlPath: %#v\n", r.URL.Path) // testing
	var err error
	var t *template.Template
	t, err = template.ParseFiles(assets.Chemin + "templates/register.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
func IndexHandlerNoMeth(w http.ResponseWriter, r *http.Request) {
	log.Println(models.GetCurrentFuncName())
	log.Println("HTTP Error", http.StatusMethodNotAllowed)
	w.WriteHeader(http.StatusMethodNotAllowed)
	Logger.Warn("indexHandlerNoMeth", slog.String("reqURL", r.URL.String()), slog.Int("HttpStatus", http.StatusMethodNotAllowed))
	//http.Redirect(w, r, "/Home?err=pass", http.StatusSeeOther)
	http.Redirect(w, r, "/Home", http.StatusSeeOther)
	//w.Write([]byte("Error " + fmt.Sprint(http.StatusMethodNotAllowed) + " !"))
}
func IndexHandlerOther(w http.ResponseWriter, r *http.Request) {
	log.Println(models.GetCurrentFuncName())
	log.Println("HTTP Error", http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	Logger.Warn("indexHandlerOther", slog.String("reqURL", r.URL.String()), slog.Int("HttpStatus", http.StatusNotFound))
	//http.Redirect(w, r, "/Home?err=pass", http.StatusSeeOther)
	http.Redirect(w, r, "/Home", http.StatusSeeOther)
	//w.Write([]byte("Error " + fmt.Sprint(http.StatusNotFound) + " !"))
}
