package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

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

var mainDomain string = "localhost" + portNumber

var dbAdminUsers = map[string]AdminUser{}
var dbSessions = map[string]string{}

func init() {
	portNumber = ":8080"
	tpl = template.Must(template.ParseGlob("./Views/*.gohtml"))

	bs, _ := bcrypt.GenerateFromPassword([]byte("Welcome1!"), bcrypt.MinCost)
	dbAdminUsers["tk@trevorknott.io"] = AdminUser{"tk@trevorknott.io", bs} /* Dummy User For Now */

}

func main() {
	fmt.Println("Server Initialized on port:", portNumber)

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
	mux.GET("/admin/signin", adminSignin)
	mux.GET("/about", about)
	mux.GET("/pc_langs", lang_call)
	mux.POST("/pc_langs/delete/:lang", lang_delete_call)
	mux.POST("/pc_langs/add", lang_add_call)
	mux.GET("/api/pc_langs/", lang_json_call)
	mux.GET("/contact", amigos)
	mux.ServeFiles("/assets/*filepath", http.Dir("/assets"))

	mux.POST("/admin/login", admin_login)

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

func adminSignin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fmt.Println("Blog Control Page Hit")

	tpl.ExecuteTemplate(w, "admin_signin.gohtml", nil)
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

// LANGUAGES //

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
	if os.Getenv("admin") == "true" {
		if req.Method == "POST" {
			req.ParseForm()
			var lang_to_del string
			lang_to_del = ps.ByName("lang") // req.FormValue("lang_del") ALternate way via. form

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

			fmt.Println(w, "DELETED RECORD", rows)

		}
	} else {
		fmt.Println(w, "MUST BE ADMIN TO DELETE LANGUAGES")
	}

	http.Redirect(w, req, "/pc_langs", 301)

}

func lang_add_call(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if os.Getenv("admin") == "true" {
		if req.Method == "POST" {
			req.ParseForm()

			var lang_to_add string
			lang_to_add = req.FormValue("lang_add")

			fmt.Println("Lang to add:", lang_to_add)

			stmt, err := db2.Prepare(`INSERT INTO pc_langs(lang_id, lang_name) VALUES(?, ?);`) // `INSERT INTO customer VALUES ("James");`
			check(err)

			if lang_to_add != "" {
				result, err := stmt.Exec(0, lang_to_add)
				check(err)

				fmt.Println(w, "ADD RECORD", result)
			} else {
				fmt.Println(w, "Unable to add NULL FIELDS!", lang_to_add)
			}

		}
	} else {
		fmt.Println(w, "MUST HAVE ADMIN ACCESS TO ADD LANGS!", nil)
	}

	http.Redirect(w, req, "/pc_langs", 301)

}

func admin_login(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == "POST" {
		var admin_name string
		var admin_password string

		admin_name = req.FormValue("admin-email")
		admin_password = req.FormValue("admin-password")

		rows, err := db2.Query(`SELECT admin_email FROM admin_users;`)
		fmt.Println(w, "Established admin_users db connection", nil)
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

		if admin_name == names[0] && admin_password == "admin" {
			os.Setenv("admin", "true")
			tpl.ExecuteTemplate(w, "admin_users.gohtml", names)
		} else {
			fmt.Fprintln(w, "YOU ARE NOT AN AUTHORIZED ADMIN.. GO BACK QUICKLY CONTACT TREVOR KNOTT... Geolocating in 20 seconds...", nil)
		}

	}
}
