package db
import (
	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//"fmt"
)


func CreateSubmit(tcid int64)int64 {
	statement, _ := db_obj.Prepare("INSERT INTO submit (tc_id) VALUES (?)")
	res, _ :=statement.Exec(tcid)
	id, _ := res.LastInsertId()
    return id
}
func UpdateScore(tcid,score int64) {
	statement, _ := db_obj.Prepare("update submit set score=?  where id=?")
	statement.Exec(score,tcid)
}

