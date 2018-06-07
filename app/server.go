package main

import (
    "fmt"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "os"
    "html/template"
)

type Employee struct {
    fName  string
    sName  string
    dptName   string
    eMail string
}

func helloWorld(w http.ResponseWriter, r *http.Request){
    name, err := os.Hostname()

    if err != nil {
        panic(err)
    }
    
    fmt.Fprintf(w, "Hello World!!!!\n")
    fmt.Fprintf(w, name)
}

func dbConnect(db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "supersecret"
    dbHost := "mysql.go"
    dbPort := "3306"
    dbName := "company"
    db, err := sql.Open(dbDriver, dbUser +":"+ dbPass +"@tcp("+ dbHost +":"+ dbPort +")/"+ dbName +"?charset=utf8")
    checkErr(err)
    return db
}

tmpl := template.Must(template.ParseFiles("form/*"))

func dbTableHtml(w http.ResponseWriter, r *http.Request){
	db := dbConn()
	rows, err := db.Query("select * from employees")
	checkErr(err)
	
	emp := Employee{}
    res := []Employee{}
	
	for rows.Next() {
		var first_name string
		var last_name string
		var department string
		var email string
		err = rows.Scan(&first_name, &last_name, &department, &email)
		checkErr(err)
		//fmt.Fprintf(w,"|%12s|%12s|%12s|%20s|\n" ,first_name ,last_name ,department ,email)
		emp.fName = first_name
		emp.sName = last_name
		emp.dptName = department
		emp.eMail = email
		res = append(res, emp)
		
	}
	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}


func main() {
    http.HandleFunc("/", helloWorld)
    http.HandleFunc("/view", dbTableHtml) 
    http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
