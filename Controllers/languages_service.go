// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"net/http"

// 	_ "github.com/go-sql-driver/mysql"
// )

// func lang_call(w http.ResponseWriter, req *http.Request) {
// 	rows, err := db2.Query(`SELECT language_name FROM languages;`)
// 	check(err)
// 	defer rows.Close()

// 	// data to be used in query
// 	var s, name string
// 	s = "RETRIEVED RECORDS:\n"

// 	// query
// 	for rows.Next() {
// 		err = rows.Scan(&name)
// 		check(err)
// 		s += name + "\n"
// 	}
// 	// fmt.Fprintln(w, s)

// 	tpl.ExecuteTemplate(w, "all_langs.gohtml", s) //admin
// }
