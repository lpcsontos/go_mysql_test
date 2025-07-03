package handlers

import (
	"db_test/utils"
   _ "database/sql"
	"fmt"
   "log"
   "net/http"
	_ "os"
	_ "github.com/joho/godotenv"
   _ "github.com/go-sql-driver/mysql"
	_ "golang.org/x/crypto/bcrypt"
)


func LoginHand( w http.ResponseWriter, r *http.Request){
	err := utils.LoginPage.Execute(w, nil)
	if err != nil {
		http.Error(w, "404", http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "only post method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")
	
	if len(name) > 20 || len(password) > 20{
		http.Error(w, "name or password should be maximum 20 character long", http.StatusBadRequest)
		return
	}

	var hash string
	Qerr := utils.DB.QueryRow("SELECT password FROM users WHERE name = ?", name).Scan(&hash)
	if Qerr != nil {
		w.WriteHeader(http.StatusUnauthorized)
    	fmt.Fprintf(w, `<script>alert("Wrong username or password");window.history.back()</script>`)
		log.Println("User not found:", Qerr)
		return
	}

	if utils.PassComp(hash, password){
		fmt.Printf("siker")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}else{
		w.WriteHeader(http.StatusUnauthorized)
    	fmt.Fprintf(w, `<script>alert("Wrong username or password");window.history.back()</script>`)
	}

}
