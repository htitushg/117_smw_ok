package controllers

import (
	"117_smw_ok/assets"
	"117_smw_ok/models"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Ajouté le 24/02/2024 19h20
// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

	//return
}

// Ajout du 24/02/2024 19h20
var logs, _ = os.Create("logs/logs.log")
var jsonHandler = slog.NewJSONHandler(logs, &slog.HandlerOptions{
	Level:     slog.LevelDebug,
	AddSource: true,
}).WithAttrs([]slog.Attr{slog.Int("Info", 13)})
var Logger = slog.New(jsonHandler)
var LogId = 0

// Log is a models.Middleware that writes a series of information in logs/logs.log
// in JSON format: time, function name, request Id (incremented int),
// client IP, request Method, and request URL.
func Log() models.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("session_token")
			var pseudo string
			if err != nil {
				pseudo = ""
			} else {
				token := c.Value
				pseudo = assets.Sessions[token].Pseudo
			}
			start := time.Now()
			wrapped := wrapResponseWriter(w)
			//next.ServeHTTP(wrapped, r)

			/* "method", r.Method,
			 */

			LogId++
			log.Println("Log()")
			Logger.Info("Log() Middleware", slog.Int("reqId", LogId), slog.Duration("duration", time.Since(start)), slog.Int("status", wrapped.status), slog.String("path", r.URL.EscapedPath()), slog.String("clientIP", models.GetIP(r)), slog.String("pseudo", pseudo), slog.String("reqMethod", r.Method), slog.String("reqURL", r.URL.String()))
			handler.ServeHTTP(w, r)
		}
	}
}

// slog.String("Name", GetCurrentName(w, r)),
func Guard() models.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("Guard()")
			c, err := r.Cookie("session_token")
			resultat := false
			stoken := ""
			if err != nil {
				if err == http.ErrNoCookie {
					resultat = false
				}
				resultat = false
				http.Redirect(w, r, "/Login?err=restricted", http.StatusSeeOther)
				return
			}
			stoken = c.Value
			_, exists := assets.Sessions[stoken]
			if !exists {
				resultat = false
			}

			// If the previous session is valid, create a new session token for the current user
			// on peut utiliser google : "github.com/google/uuid"
			// ou bien pseudo_uuid() fonction ci dessus qui utilise "crypto/rand"

			newSessionToken := pseudo_uuid()
			maxAge := 120
			// Set the token in the session map, along with the user whom it represents
			assets.Sessions[newSessionToken] = assets.Session{
				Pseudo: assets.Sessions[stoken].Pseudo,
				MaxAge: maxAge,
			}
			fmt.Println(assets.Sessions)
			// Delete the older session token
			delete(assets.Sessions, stoken)
			fmt.Println(assets.Sessions)
			//http.SetCookie(w, &http.Cookie{Name: "session_token", MaxAge: -1})
			// Set the new token as the users `session_token` cookie
			var cookie = http.Cookie{
				Name:   "session_token",
				Value:  newSessionToken,
				MaxAge: maxAge,
			}
			http.SetCookie(w, &cookie)
			r.AddCookie(&cookie)
			resultat = true
			if resultat {
				log.Println("Guard()", assets.Sessions[newSessionToken].MaxAge)
				handler.ServeHTTP(w, r)
				return
			} else {
				http.Redirect(w, r, "/Login?err=restricted", http.StatusSeeOther)
				return
			}
		}
	}
}

func Foo() models.Middleware {
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println("Foo()")
			handler.ServeHTTP(w, r)
		}
	}
}

// Join is used to concatenate various middlewares, for better visibility.
// it takes the http.HandlerFunc corresponding to the route, and then
// any number of models.Middleware that will be concatenated in order like this:
// middlewares[0](middlewares[1](middlewares[2](handlerFunc))).
func Join(handlerFunc http.HandlerFunc, middlewares ...models.Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handlerFunc = middlewares[i](handlerFunc)
	}
	return handlerFunc
}
