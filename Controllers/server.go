package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type AdminUser struct {
	Email    string
	Password []byte
}

type PClanguage struct {
	Lang string
}

var portNumber string

var tpl *template.Template

var db2 *sql.DB
var db *sql.DB
var err error

var dbAdminUsers = map[string]AdminUser{}
var dbSessions = map[string]string{}

func init() {
	portNumber = ":8080"
	tpl = template.Must(template.ParseGlob("../Views/*.gohtml"))

	bs, _ := bcrypt.GenerateFromPassword([]byte("Welcome1!"), bcrypt.MinCost)
	dbAdminUsers["tk@trevorknott.io"] = AdminUser{"tk@trevorknott.io", bs} /* Dummy User For Now */

}

func main() {
	fmt.Println("Server Initialized on port :8080")

	db2, err = sql.Open("mysql", "tknott95:Welcome1!@tcp(godbinstance.cfchdss74ohb.us-west-1.rds.amazonaws.com:3306)/myGOdb?charset=utf8")
	check(err)
	defer db2.Close()

	err = db2.Ping()
	check(err)

	db, err = sql.Open("mysql", "tknott95:Welcome1!@tcp(godbinstance.cfchdss74ohb.us-west-1.rds.amazonaws.com:3306)/test02?charset=utf8")
	check(err)
	defer db.Close()

	err = db.Ping()
	check(err)

	http.HandleFunc("/", index)
	http.HandleFunc("/about", about)
  	http.HandleFunc("/pc_langs", lang_call)
  	http.HandleFunc("/api/pc_langs/", lang_json_call)
	http.HandleFunc("/contact", amigos)

  	http.HandleFunc("/pc_langs/delete/4", lang_delete)


	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("Public"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServe(portNumber, nil)
}

func index(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Index Page Hit...")

	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func about(w http.ResponseWriter, req *http.Request) {
	fmt.Println("About Page hit...")

	tpl.ExecuteTemplate(w, "about.gohtml", nil)
}

func blogControl(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Blog Control Page Hit")
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func amigos(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Query(`SELECT amigoname FROM amigos;`)
	check(err)
	defer rows.Close()

	// data to be used in query
	var s, name string

	// query
	for rows.Next() {
		err = rows.Scan(&name)
		check(err)
		s += name + "\n"
	}

	json, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")

	fmt.Println(string(json))
}

func lang_call(w http.ResponseWriter, req *http.Request) {
	rows, err := db2.Query(`SELECT language_name FROM languages;`) 
	check(err)
	defer rows.Close()

	// data to be used in query
	var name string
	var names []string

	// query
	for rows.Next() {
		err = rows.Scan(&name)
		check(err)
	
		names = append(names, name)

		
}
	tpl.ExecuteTemplate(w, "all_langs.gohtml", names)
}

func lang_json_call(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db2.Query(`SELECT language_name FROM languages;`) 
	check(err)
	defer rows.Close()

	// data to be used in query
	var name string
	var names []string

	// query
	for rows.Next() {
		err = rows.Scan(&name)
		check(err)
		// name := PClanguage{
		// 	Lang: name,
		// }

		// json.Marshal(name)

		names = append(names, name)

		
}

bs, err := json.Marshal(names)
	if err != nil {
		fmt.Println("error: ", err)
	}

	
		w.Write(bs)

	// tpl.ExecuteTemplate(w, "all_langs.gohtml", bs)//admin
}

func lang_delete(w http.ResponseWriter, req *http.Request) {
	stmt, err := db2.Prepare(`DELETE FROM languages WHERE languages_id=?;`)
	check(err)
	defer stmt.Close()

	r, err := stmt.Exec()
	check(err)

	n, err := r.RowsAffected()
	check(err)

	fmt.Fprintln(w, "Deleted Lang ID 4", n)
}
