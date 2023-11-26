package db

import
(
	"database/sql"
    "log"
    "fmt"

    _ "github.com/lib/pq"
)

var DB *sql.DB

const Connstr="user=postgres password=novell@123 dbname=postgres sslmode=disable"

func InitDB(){
	var err error

	DB,err=sql.Open("postgres",Connstr)
	if err!=nil{
		fmt.Printf("error in connecting to db %v\n",err)
		return
	}
	createTable()

}
func createTable() {
    query := `
    CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL
    );`

    _, err := DB.Exec(query)
    if err != nil {
        log.Fatalf("Error creating table: %s", err)
    }
}
func CloseDB(){
	DB.Close()
}