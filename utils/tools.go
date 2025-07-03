package utils

import(
	"database/sql"
	_ "fmt"
   _ "log"
   _ "net/http"
	_ "os"
	_ "github.com/joho/godotenv"
   _ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"html/template"
)


var IndPage = template.Must(template.ParseFiles("pages/index.html"))
var RegPage = template.Must(template.ParseFiles("pages/register.html"))
var LoginPage = template.Must(template.ParseFiles("pages/login.html"))

var DB *sql.DB

func PassComp(hashpass string, normpass string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashpass), []byte(normpass))
	return err == nil
}

