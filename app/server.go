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
    Fname, Sname, Dname, Email string
	Id int
}

func helloWorld(w http.ResponseWriter, r *http.Request){
    name, err := os.Hostname()
	checkErr(err)
    fmt.Fprintf(w, "HOSTNAME : %s\n", name)
}

func dbConnect() (db *sql.DB) {
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

func dbSelect() []Employee{
	db := dbConnect()
	rows, err := db.Query("select * from employees")
	checkErr(err)
	
	employee := Employee{}
    employees := []Employee{}
	
	for rows.Next() {
		var first_name, last_name, department, email string
		var id int
		err = rows.Scan(&id, &first_name, &last_name, &department, &email)
		checkErr(err)
		employee.Id = id
		employee.Fname = first_name
		employee.Sname = last_name
		employee.Dname = department
		employee.Email = email
		employees = append(employees, employee)
		
	}
	defer db.Close()
	return employees
}

var tmpl = template.Must(template.ParseGlob("templates/*"))

func dbTableHtml(w http.ResponseWriter, r *http.Request){
	table := dbSelect()
	err := tmpl.ExecuteTemplate(w, "Index", table)
	checkErr(err)
}

func dbTable(w http.ResponseWriter, r *http.Request){
    table := dbSelect()
	for i := range(table) {
        emp := table[i]
        fmt.Fprintf(w,"YES|%5d|%12s|%12s|%12s|%20s|\n" ,emp.Id, emp.Fname ,emp.Sname ,emp.Dname ,emp.Email)
    }
}

func main() {
    http.HandleFunc("/", helloWorld)
    http.HandleFunc("/view", dbTableHtml) 
	http.HandleFunc("/raw", dbTable)
    http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
