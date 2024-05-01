package db
import (
	//"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	//"os"
)

func CheckPassword(userid,pass string) (bool) {
	statement, _ := db_obj.Prepare("SELECT * FROM user where id=? and password=?")
	rows,_ :=statement.Query(userid,pass)
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return false
}
func GetPassword(userid string) (string,error) {
	statement, _ := db_obj.Prepare("SELECT password FROM user where id=?")
	rows,_ :=statement.Query(userid)
	var tempC string
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&tempC)
		return tempC,nil
	}
	return "",fmt.Errorf("Invalid User Info")
}
func AddUser(userid,pass string){
	statement, _ := db_obj.Prepare("INSERT INTO user (id, password) VALUES (?, ?)")
	statement.Exec(userid,pass)
}

func ResetPassowrd(userid,pass string){
	statement, _ := db_obj.Prepare("update user set password=?  where id=?")
	statement.Exec(pass,userid)
}
