package db
import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//"fmt"
	"log"
	"io/ioutil"
	
)
var db_obj *sql.DB
func init(){
	var err  error
	db_obj, err = sql.Open("sqlite3", "./db/code.db")

	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("db initialized: \n")
	content, err := ioutil.ReadFile("./db/start.sql")     // the file is inside the local directory
	if err != nil {
		log.Fatal(err)
	}
	_, err = db_obj.Exec(string(content))
	if err != nil {
		log.Fatal(err)
	}
}

func Run(){

}
