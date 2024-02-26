# Programme en Go avec authentification des utilisateurs avec Sessions et Cookies

 _Tiré de:(<https://www.sohamkamani.com/golang/session-cookie-authentication/>)_

## Démarrer cette application

To run this application, build and run the Go binary:

```sh
go run .

```

Ce programme met en oeuvre un système de controle d'accès pour pouvoir accéder à un contenu.
Il utilise un système de session associé à des cookies ("github.com/google/uuid")

Un cookie HTTP (également appelé cookie web ou cookie de navigateur) est une donnée de petite taille envoyée automatiquement par le serveur au navigateur web de l'utilisateur. Le navigateur peut alors enregistrer le cookie et le renvoyer au serveur lors des requêtes ultérieures.

Généralement, un cookie HTTP sert à indiquer que deux (ou plusieurs) requêtes proviennent du même navigateur où une personne est connectée. Il permet de mémoriser des informations d'état alors que le protocole HTTP est sans état.
Les cookies ont trois usages principaux :

La gestion de session
Connexions aux sites, chariots d'achats, scores de jeux, ou toute autre chose que le serveur devrait mémoriser

La personnalisation
Les préférences et autres éléments de configuration

Le pistage
L'enregistrement et l'analyse du comportement de la personne visitant le site

### Remarque

Les cookies ont été un outil général de stockage côté client. Bien que cela était pertinent lorsque c'était la seule façon de stocker des données côté client, il est désormais recommandé d'utiliser des API modernes dédiées à cet usage. Les cookies sont envoyés avec chaque requête et peuvent alourdir les performances (notamment pour les connexions mobiles). Les API modernes pour le stockage de données client sont :

* L'API Web Storage (localStorage et sessionStorage)
* IndexedDB.

Note : Pour observer les cookies enregistrés (et les autres types de stockage utilisés par une page web), vous pouvez activer l'inspecteur de stockage dans les outils de développement de Firefox et ouvrir le niveau Cookies dans la hiérarchie de l'onglet Stockage.

#### Créer un cookie

Après avoir reçu une requête HTTP, un serveur peut envoyer un ou plusieurs en-têtes Set-Cookie avec la réponse. Le navigateur enregistre alors généralement le ou les cookies et les renvoie via l'en-tête HTTP Cookie (en-US) pour les requêtes envers le même serveur. Il est possible d'indiquer une date d'expiration ou une durée de vie après laquelle le cookie ne devrait plus être envoyé. Il est également possible d'ajouter des restrictions supplémentaires pour le domaine et les chemins pour lesquels le cookie peut être envoyé. Pour plus de détails sur les attributs des en-têtes mentionnés plus tôt, consultez la page de référence pour Set-Cookie("<https://developer.mozilla.org/fr/docs/Web/HTTP/Headers/Set-Cookie>").

### Caractéristiques du programme

// Une map contient les sessions des utilisateurs.

```go

var Sessions = map[string]Session{}

```

// Chaque session contient le Pseudo de l'utilisateur et la date d'expiration du jeton (Token) associé à la session
dans cette structure on peut ajouter les informations dont on a besoin en plus des deux informations requises.

```go
type Session struct {
 Pseudo    string
 Expiry    time.Time
 Email     string
 Firstname string
 Lastname  string
 Address   string
 Town      string
 ZipCode   string
 Country   string
}
```

// On crée une map pour stocker les paires (Pseudo, password). Dans une application réelle on pourrait utiliser une base de données.
Par exemple avec deux utilisateurs:

```go
var Users = map[string]string{
 "henry": "1nhri96p",
 "marie": "marie123",
}
```

** "session_token" est le code secret qui permet d'accéder au cookie (r est de type http.Request)

```go
c, err := r.Cookie("session_token")
```

** On accède au jeton en cours :

```go
stoken = c.Value
```

** On vérifie si la date de fin du token est atteinte :

```go
if assets.Sessions[stoken].Expiry.Before(time.Now())
```

## Fonctions utilisées

### SessionValide(w http.ResponseWriter, r *http.Request) (stoken string, resultat bool)

 vérifie si la session est valide et renvoie un booléen et le jeton (token)

### Home(w http.ResponseWriter, r *http.Request)

Ce controller affiche la page d'accès à l'pplication avant connexion de l'utilisateur
Si la session est valide on oriente l'utilisateur vers la page index. html
sinon on lui présente la page home.html

### Login(w http.ResponseWriter, r *http.Request)

Ce controlleur affiche la page de connexion

### Signin(w http.ResponseWriter, r *http.Request)

Ce controlleur traite les éléments entrés par l'utilisateur:
Il vérifie lexactitude du nom (Pseudo) et du mot de passe(Password),

### Si l'entrée est correcte

création du jeton (sessionToken) et de sa date d'expiration (expiresAt)

```go
    sessionToken := uuid.NewString()
 maxAge := 120
```

Initialisation dans la map du Pseudo est de la date d'expiration:

```go
    assets.Sessions[sessionToken] = assets.Session{
  Pseudo: creds.Pseudo,
  MaxAge: maxAge,
 }
```

Mise à jour du Cookie :

```go
    http.SetCookie(w, &http.Cookie{
  Name:    "session_token",
  Value:   sessionToken,
  MaxAge: maxAge,
 })
```

Appel du template index.html

### Si l'entrée est incorrecte

Appel de la page home.html

### SessionValide(w http.ResponseWriter, r *http.Request)

On verifie si la session est valide

### Si non valide

on affiche le template home.html

### Si valide

On crée un nouveau jeton(newSessionToken) et une nouvelle date d'expiration

```go
    newSessionToken := uuid.NewString()
 MaxAge := 120
```

On stocke dans la map le pseudo et la date d'expiration

```go
    assets.Sessions[newSessionToken] = assets.Session{
  Pseudo: assets.Sessions[sessionToken].Pseudo,
  MaxAge: 120,
 }
```

On efface l'ancien jeton dans la map

```go
    delete(assets.Sessions, sessionToken)
```

On met à jour le cookie

```go
    http.SetCookie(w, &http.Cookie{
  Name:    "session_token",
  Value:   newSessionToken,
  MaxAge: 120,
 })
```

Appel du template index.html

### Logout(w http.ResponseWriter, r *http.Request)

On vérifie que la session est valide

#### si valide

On efface l'ancien jeton dans la map

```go
    delete(assets.Sessions, sessionToken)
```

On met à jour le cookie

```go
    http.SetCookie(w, &http.Cookie{
  Name:    "session_token",
  Value:   "",
  MaxAge: -1,
 })
```

Appel du template home.html

### AfficheUserInfo(w http.ResponseWriter, r *http.Request)

Controller qui affiche la page afficheuserinfo.html" avec les éléments de la session
