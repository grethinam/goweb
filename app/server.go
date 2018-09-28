package main

import (
    "fmt"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "os"
    "html/template"
	"log"
	"strings"
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
    dbHost := "mysqldb"
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

func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConnect()
    emp := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM employees WHERE id=?")
	checkErr(err)
    delForm.Exec(emp)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/view", 301)
}

func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConnect()
    if r.Method == "POST" {
        fname := r.FormValue("fname")
		sname := r.FormValue("sname")
		dname := r.FormValue("dname")
        email := r.FormValue("email")
        insForm, err := db.Prepare("INSERT INTO employees(first_name, last_name, department, email) VALUES(?,?,?,?)")
	    checkErr(err)
        insForm.Exec(fname, sname, dname, email)
        log.Println("INSERT: First Name: " + fname + " | LAST_NAME: " + sname+ " | DEPARTMENT: " + dname+ " | EMAIL: " + email)
    }
    defer db.Close()
    http.Redirect(w, r, "/view", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    db := dbConnect()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM employees WHERE id=?", nId)
	checkErr(err)
	employee := Employee{}
    for selDB.Next() {
		var first_name, last_name, department, email string
		var id int
		err = selDB.Scan(&id, &first_name, &last_name, &department, &email)
		checkErr(err)
		employee.Id = id
		employee.Fname = first_name
		employee.Sname = last_name
		employee.Dname = department
		employee.Email = email
    }
    tmpl.ExecuteTemplate(w, "Edit", employee)
    defer db.Close()
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConnect()
    if r.Method == "POST" {
        fname := r.FormValue("fname")
		sname := r.FormValue("sname")
		dname := r.FormValue("dname")
        email := r.FormValue("email")
        id := r.FormValue("uid")
        updForm, err := db.Prepare("UPDATE employees SET first_name=?, last_name=?, department=?, email=? WHERE id=?")
	    checkErr(err)
        updForm.Exec(fname, sname, dname, email, id)
        log.Println("UPDATE: First Name: " + fname + " | LAST_NAME: " + sname+ " | DEPARTMENT: " + dname+ " | EMAIL: " + email)
    }
    defer db.Close()
    http.Redirect(w, r, "/view", 301)
}

func main() {

	checkTableExist()
    http.HandleFunc("/", dbTableHtml)
	http.HandleFunc("/raw", dbTable)
	http.HandleFunc("/host", helloWorld)
	http.HandleFunc("/new", New)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)
    http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func checkTableExist() {
	
	db := dbConnect()
	defer db.Close()
	// make sure connection is available
	err := db.Ping()
    checkErr(err)
	
	dbName := "company"
	dbTable	:= "employees"
	
	infoTablestmt, err := db.Query("SELECT TABLE_NAME FROM information_schema.tables WHERE table_schema = '" + dbName + "' AND table_name = '" + dbTable + "' LIMIT 1;")
	checkErr(err)
	
	var queryTable string
	
	for infoTablestmt.Next() {
		var table string
        err = infoTablestmt.Scan(&table)
        checkErr(err)
        queryTable = table
    }
	if strings.TrimSpace(queryTable) != dbTable {
		log.Println("INFO: Employee Table creation started....")
		createT, err := db.Prepare("CREATE TABLE employees ( id int(6) unsigned NOT NULL AUTO_INCREMENT, first_name varchar(25) NOT NULL, last_name  varchar(25) NOT NULL, department varchar(15) NOT NULL, email  varchar(50) NOT NULL, PRIMARY KEY (id))ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;")
		
		checkErr(err)
		_, err = createT.Exec()
		if err != nil {
			panic(err)
		} else {
			log.Println("INFO: Employee Table successfully created....")
		}
	}else{
		log.Println("INFO: Employee Table exist....")
	}
}
