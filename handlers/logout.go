package handlers

import (
	"db_test/utils"
   _ "database/sql"
	"fmt"
   "log"
   "net/http"
	"time"
	_ "os"
	_ "github.com/joho/godotenv"
   _ "github.com/go-sql-driver/mysql"
	_"golang.org/x/crypto/bcrypt"
)

func LogoutHand( w http.ResponseWriter, r *http.Request){
	
	user, err := r.Cookie("user")
	if err != nil{
		log.Println("User cookie not found:", err)
		return
	}
	name := user.Value
	st, err := r.Cookie("session_token")
	if err != nil{
		log.Println("session cookie not found:", err)
		return
	}

	dbQuery := `
		DELETE FROM tokens
		WHERE sessionToken = ? AND id = (SELECT id FROM users WHERE name = ?)
	`

	_, err = utils.DB.Exec(dbQuery, st.Value, name)
	if err != nil {
		log.Printf("Delete query failed with error: %v", err)
		log.Printf("Error type: %T", err)
		http.Error(w, fmt.Sprintf("Database update failed: %v", err), http.StatusInternalServerError)
		return
	}
	
	expired := time.Now().Add(-1 * time.Hour)
	http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", Expires: expired, HttpOnly: true, Secure: true})
	http.SetCookie(w, &http.Cookie{Name: "csrf_token", Value: "", Expires: expired, HttpOnly: false, Secure: true})
	http.SetCookie(w, &http.Cookie{Name: "user", Value: "", Expires: expired, HttpOnly: true, Secure: true})

	
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
