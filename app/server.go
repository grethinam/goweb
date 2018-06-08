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
    fname, sname, dname, email string
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
		err = rows.Scan(&first_name, &last_name, &department, &email)
		checkErr(err)
		employee.fname = first_name
		employee.sname = last_name
		employee.dname = department
		employee.email = email
		employees = append(employees, employee)
		
	}
	defer db.Close()
	return employees
}

var tmpl = template.Must(template.ParseFiles("layout.html"))
//var tmpl = template.Must(template.ParseGlob("layout.html"))
func dbTableHtml(w http.ResponseWriter, r *http.Request){
	db := dbConnect()
	rows, err := db.Query("select * from employees")
	checkErr(err)
	
	emp := Employee{}
    employees := []Employee{}
	
	for rows.Next() {
		var first_name string
		var last_name string
		var department string
		var email string
		err = rows.Scan(&first_name, &last_name, &department, &email)
		checkErr(err)
		//fmt.Fprintf(w,"|%12s|%12s|%12s|%20s|\n" ,first_name ,last_name ,department ,email)
		emp.fname = first_name
		emp.sname = last_name
		emp.dname = department
		emp.email = email
		employees = append(employees, emp)
		
	}
	
	/*for i := range(res) {
        emp := res[i]
        fmt.Fprintf(w,"HA|%12s|%12s|%12s|%20s|\n" ,emp.fName ,emp.sName ,emp.dptName ,emp.eMail)
    }*/
	
	tmpl.ExecuteTemplate(w, "Index", employees)
	defer db.Close()
}

func dbTable(w http.ResponseWriter, r *http.Request){
    table := dbSelect()
	for i := range(table) {
        emp := table[i]
        fmt.Fprintf(w,"YES|%12s|%12s|%12s|%20s|\n" ,emp.fname ,emp.sname ,emp.dname ,emp.email)
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
