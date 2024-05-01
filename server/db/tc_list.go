package db
import (
	//"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)


func AddTCFile(userid,probid,sample string)int64 {
	statement, _ := db_obj.Prepare("INSERT INTO tc_list (user_id, prob_id,sample) VALUES (?, ?,?)")
	res, _ :=statement.Exec(userid,probid,sample)
	id, _ := res.LastInsertId()
    return id
}

func UpdateTCFile(tc_list,id string) {
	statement, _ := db_obj.Prepare("update tc_list set tc_id_list=?  where id=?")
	statement.Exec(tc_list,id)
}
func GetTCList(id int64) [][2]string {
	statement, _ := db_obj.Prepare("select tc_id_list from tc_list where id=?")
	rows,_ :=statement.Query(id)
	defer rows.Close()
	var tcs [][2]string
	for rows.Next() {
		var tempC string
		rows.Scan(&tempC)
		s1 := strings.Split(tempC, ",")
		for i := 0; i < len(s1); i++{
			tc:=GetTC(s1[i])
			m := [2]string{tc.input,tc.output}
			tcs=append(tcs,m)
		}
		return tcs
	}
	return [][2]string{}
}
func IsActiveSampleTC(userid,probid string) int64 {
	statement, _ := db_obj.Prepare("select id from tc_list where user_id=? and prob_id=? and sample=\"true\"")
	rows,_ :=statement.Query(userid,probid)
	defer rows.Close()
	for rows.Next() {
		var tempC int64
		rows.Scan(&tempC)
		return tempC
	}
	return -1
}

func IsActiveFullTC(userid,probid string) int64 {
	statement, _ := db_obj.Prepare("select id from tc_list where user_id=? and prob_id=? and sample=\"false\" and creation>(select datetime('now','-10 minutes'))")
	rows,_ :=statement.Query(userid,probid)
	defer rows.Close()
	for rows.Next() {
		var tempC int64
		rows.Scan(&tempC)
		return tempC
	}
	return -1
}


