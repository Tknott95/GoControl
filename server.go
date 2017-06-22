package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"


	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
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
	tpl = template.Must(template.ParseGlob("./Views/*.gohtml"))

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

	mux := httprouter.New()
	mux.GET("/", index)
	mux.GET("/about", about)
	mux.GET("/pc_langs", lang_call)
	mux.POST("/pc_langs/delete/:lang", lang_delete_call)
	mux.POST("/pc_langs/add", lang_add_call)
	mux.GET("/api/pc_langs/", lang_json_call)
	mux.GET("/contact", amigos)

	// mux.Handler("GET", "/assets/", http.StripPrefix("/assets", justFilesFilesystem{http.Dir("Public")}))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	http.ListenAndServe(portNumber, mux)
}

func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Println("Index Page Hit...")

	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}

func about(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Println("About Page hit...")

	tpl.ExecuteTemplate(w, "about.gohtml", nil)
}

func blogControl(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Println("Blog Control Page Hit")
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func amigos(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

func lang_call(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	rows, err := db2.Query(`SELECT lang_name FROM pc_langs;`)
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

func lang_json_call(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := db2.Query(`SELECT lang_name FROM pc_langs;`)
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

func lang_delete_call(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	if req.Method == "POST" {
		req.ParseForm()
		var lang_to_del string
		lang_to_del = ps.ByName("lang"); // req.FormValue("lang_del") ALternate way via. form

		if lang_to_del == "C" { // BUG FIX TO ALLOW C# DELETION (C WONT BE ADDED REGARDLESS)
			lang_to_del = "C#"
		}

		fmt.Println("Lang to delete:", lang_to_del)

		stmt, err := db2.Prepare(`DELETE FROM pc_langs WHERE lang_name= ?;`)
		check(err)
		defer stmt.Close()

		rows, err := stmt.Query(lang_to_del)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			// ...
		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Fprintln(w, "DELETED RECORD", rows)

	}

}

func lang_add_call(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == "POST" {
		req.ParseForm()

		var lang_to_add string
		lang_to_add = req.FormValue("lang_add")

		fmt.Println("Lang to add:", lang_to_add)

		stmt, err := db2.Prepare(`INSERT INTO pc_langs(lang_id, lang_name) VALUES(?, ?);`) // `INSERT INTO customer VALUES ("James");`
		check(err)
		
		result, err := stmt.Exec(0, lang_to_add)
		check(err)

		fmt.Fprintln(w, "ADD RECORD", result)

	}
}