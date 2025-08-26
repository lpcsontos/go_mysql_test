package utils

import(
	"database/sql"
	_"fmt"
   "log"
   "net/http"
	"time"
	_ "os"
	_ "github.com/joho/godotenv"
   _ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"crypto/rand"
	"encoding/base64"
	"html/template"
	_"errors"
)


var IndPage = template.Must(template.ParseFiles("pages/index.html"))
var RegPage = template.Must(template.ParseFiles("pages/register.html"))
var LoginPage = template.Must(template.ParseFiles("pages/login.html"))
var TestPage = template.Must(template.ParseFiles("pages/test.html"))

var DB *sql.DB

func PassComp(hashpass string, normpass string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(normpass))
	return err == nil
}

func GenToken(length int) string{
	bytes := make([]byte, length)
	if _,err := rand.Read(bytes); err!=nil {
		log.Fatalf("failed to gen token %v", err)
	}

	return base64.RawURLEncoding.EncodeToString(bytes);
}

func IsLoggedIn(r *http.Request) bool {	
	user, err := r.Cookie("user")
	if err != nil{
		log.Println("User cookie not found:", err)
		return false
	}
	username := user.Value

	var db int
	qerr := DB.QueryRow("SELECT COUNT(*) as count FROM users WHERE name = ?", username).Scan(&db)
	if qerr != nil || db == 0{
		log.Println("User not found:", qerr)
		return false
	}

	st, err := r.Cookie("session_token")
	if err != nil || st.Value == ""{
		return false
	}

	var sessionExpiresSTR string
	qerr = DB.QueryRow("SELECT sessionExpires FROM tokens WHERE id = (select id from users where name = ?) AND sessionToken = ?", username, st.Value).Scan(&sessionExpiresSTR)
	if qerr != nil{
		log.Println("User not found:", qerr)
		return false
	}
	sessionExpires, err := time.Parse("2006-01-02 15:04:05", sessionExpiresSTR)
	if err != nil {
   	log.Println("Time parse error:", err)
   	return false
	}

	if time.Now().After(sessionExpires) {
      log.Println("Session expired for user:", username)
		return false
   }
	return true
}


func Auth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if !IsLoggedIn(r) {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }
        next(w, r)
    }
}

const csrfCookieName = "csrf_token"
const csrfHeaderName = "X-CSRF-Token"

func CSRFMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
			user, err := r.Cookie("user")
			if err != nil{
				log.Println("User cookie not found:", err)
            http.Error(w, "User cookie not found", http.StatusForbidden)
				return
			}
			username := user.Value

         cookie, err := r.Cookie(csrfCookieName)
         if err != nil || cookie.Value == "" {
            http.Error(w, "Missing CSRF cookie", http.StatusForbidden)
         	return
         }

			var csrfExpiresSTR string
			qerr := DB.QueryRow("SELECT csrfExpires FROM tokens WHERE id = (select id from users where name = ?) AND csrfToken = ?", username, cookie.Value).Scan(&csrfExpiresSTR)
			if qerr != nil{
				log.Println("Token not found:", qerr)
            http.Error(w, "Invalid CSRF cookie", http.StatusForbidden)
				return
			}
			csrfExpires, err := time.Parse("2006-01-02 15:04:05", csrfExpiresSTR)
			if err != nil {
    			log.Println("Time parse error:", err)
    			return
			}

			if time.Now().After(csrfExpires) {
      		log.Println("csrf expired for user:", username)
            http.Error(w, "Expired CSRF token", http.StatusForbidden)
				return
   		}

         headerToken := r.Header.Get(csrfHeaderName)
			if headerToken == "" || headerToken != cookie.Value{
				log.Println(headerToken + " : " + cookie.Value) 
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
        }

        next(w, r)
    }
}
