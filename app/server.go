package main

import (
    "fmt"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "os"
    "html/template"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func layout(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("layout.html"))
		data := TodoPageData{
			PageTitle: "My TODO list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
			},
		}
		tmpl.Execute(w, data)
}

func helloWorld(w http.ResponseWriter, r *http.Request){
    name, err := os.Hostname()

    if err != nil {
        panic(err)
    }
    
    fmt.Fprintf(w, "Hello World!!!!")
    fmt.Fprintf(w, name)
}

func dbTable(w http.ResponseWriter, r *http.Request){
    db, err := sql.Open("mysql", "root:supersecret@tcp(mysql.go:3306)/company?charset=utf8")
	checkErr(err)
	rows, err := db.Query("select * from employees")
	checkErr(err)
	
	for rows.Next() {
		var first_name string
		var last_name string
		var department string
		var email string
		err = rows.Scan(&first_name, &last_name, &department, &email)
		checkErr(err)
		fmt.Fprintf(w,"|%12s|%12s|%12s|%20s|\n" ,first_name ,last_name ,department ,email)
	}
	
	db.Close()

}

func main() {
    http.HandleFunc("/view", helloWorld)
    http.HandleFunc("/", dbTable)
    http.HandleFunc("/html", layout) 
    http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
