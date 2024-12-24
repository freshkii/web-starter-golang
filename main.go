package main

import (
	"fmt"
	"net/http"
	"os"

	"path/filepath"
	"mime"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"

	"crypto/rand"
	"encoding/hex"

	"io"
	. "github.com/bitly/go-simplejson"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func sendFile(w http.ResponseWriter, filename string) {
	f, err := os.ReadFile(filename)
	if err != nil { return }

	w.Header().Set("Content-Type", mime.TypeByExtension( filepath.Ext(string(f)) ))
	w.Write(f)
}

func sendJson(w http.ResponseWriter, json *Json) {
	w.Header().Set("Content-type", "application/json")
	string_json,_ := json.Encode()
	w.Write(string_json)
}

func GenerateSecureToken(length int) string {
  b := make([]byte, length)
  if _, err := rand.Read(b); err != nil { return "" }
  return hex.EncodeToString(b)
}
 
func parseBody(requestBody io.ReadCloser) (*Json, error) {
	bodyContent, err := io.ReadAll(requestBody)
	if err != nil { return nil, err }
	return NewJson(bodyContent)
}

func hashPassword(password string) (string,error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password),[]byte(hash))
	return err == nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] - %s (%s)\n",r.Method,r.URL.Path,r.RemoteAddr)
	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/" {
			sendFile(w,"pages/home.html")
			return
		}

		if filepath.Ext(r.URL.Path) != "" {
			sendFile(w, "static"+r.URL.Path)
		} else {
			sendFile(w, "pages"+r.URL.Path+".html")
		}
		break
	case http.MethodPost:
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		body, err := parseBody(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}


		json,_ := NewJson([]byte("{}"))

		switch r.URL.Path {
		case "/login":
			// check if username and password are of correct type
			username, err := body.Get("username").String()
			if err != nil || username == "" {
				json.Set("error", "Username or password is incorrect")
				sendJson(w,json)
				return
			}

			password, err := body.Get("password").String()
			if err != nil || password == "" {
				json.Set("error", "Username or password is incorrect")
				sendJson(w,json)
				return
			}
			
			// check if username and password are valid
			var hashedPassword string
			err = db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&hashedPassword)
			if err == sql.ErrNoRows {
				// Username not found
				json.Set("error", "Invalid username or password.")
				w.WriteHeader(http.StatusUnauthorized)
				sendJson(w, json)
				return
			} else if err != nil {
				// Database issues
				json.Set("error", "An internal server error has occurred.")
				w.WriteHeader(http.StatusInternalServerError)
				sendJson(w, json)
				return
			}


			if checkPasswordHash(hashedPassword, password) {
				json.Set("error", "Username or password is incorrect")
				sendJson(w,json)
				return
			}

			// generate token
			token := GenerateSecureToken(40)

			// store it in the database
			db.Exec("UPDATE users SET token = ? WHERE username = ?", token, username)

			// send it to the user
			json.Set("token",token)
			sendJson(w,json)

			break
		case "/signin":
			// check if a user with the username are of the correct type
			username, err := body.Get("username").String()

			if err != nil || username == "" {
				json.Set("error", "Username or password is incorrect")
				sendJson(w,json)
				return
			}

			password, err := body.Get("password").String()
			if err != nil || password == "" {
				json.Set("error", "Username or password is incorrect")
				sendJson(w,json)
				return
			}

			var dbPassword string
			row := db.QueryRow("SELECT password FROM users WHERE username = ?", username)
			row.Scan(&dbPassword)

			if dbPassword != "" {
				json.Set("error", "Username already used")
				sendJson(w,json)
				return
			}

			hashedPassword,_ := hashPassword(password)

			// store username, password in database
			db.Exec("INSERT INTO users (username, password) VALUES (?,?)",username,hashedPassword)

			// generate token
			token := GenerateSecureToken(40)

			// store it in the database
			db.Exec("UPDATE users SET token = ? WHERE username = ?", token, username)

			// send it to the user
			json.Set("token",token)
			sendJson(w,json)

			break
		}
		break
	case http.MethodDelete:
		if r.Header.Get("Content-Type") != "application/json" { return }

		body,err := parseBody(r.Body)
		if err != nil { return }

		switch r.URL.Path {
		case "/logout":
			// check if token is valid, if it isn't return
			token, err := body.Get("token").String()
			if err != nil || len(token) < 40 { return }

			username, err := body.Get("username").String()
			if err != nil || username == "" { return }


			var dbToken string
			err = db.QueryRow("SELECT token FROM users WHERE username = ?",username).Scan(&dbToken)
			if err == sql.ErrNoRows || dbToken != token { return }
			if err != nil { return }

			// remove token from database
			_, err = db.Exec("UPDATE users SET token = '' WHERE token = ?", token)
			if err != nil { return }


			break
		}
		break
	}
}

func main() {
	// init database
	var err error
	db, err = sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Print("error: could not open database\n")
		return
	}


	http.HandleFunc("/",handler)
	http.ListenAndServe(":8080",nil)
}
