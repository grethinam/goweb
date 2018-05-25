package main

import (
    "fmt"
    "net/http"
	_ "github.com/go-sql-driver/mysql"
    "database/sql"
    //"os"
)

/*func helloWorld(w http.ResponseWriter, r *http.Request){
    name, err := os.Hostname()

    if err != nil {
        panic(err)
    }
    
    fmt.Fprintf(w, "Hello World!!!!")
    fmt.Fprintf(w, name)
}*/

func dbTable(w http.ResponseWriter, r *http.Request){
    db, err := sql.Open("mysql", "root:supersecret@tcp(35.192.3.225:3306)/company?charset=utf8")
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
		/*fmt.Fprintf(w, last_name)
		fmt.Fprintf(w, department)
		fmt.Fprintf(w, email)*/
	}
	
	db.Close()

}

func main() {
    //http.HandleFunc("/", helloWorld)
    http.HandleFunc("/", dbTable)
    http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
