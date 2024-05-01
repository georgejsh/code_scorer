package libs

import (
	"encoding/json"
	"net/http"
	"github.com/jedib0t/go-pretty/v6/table"
	"time"
	//"fmt"
)

var input_status =[]Help {
	
}
func GetStatus(w http.ResponseWriter, r *http.Request) {
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input_status)
		data, _ := json.Marshal(Output{string(output), ""})
		
		w.Write(data)
		return
	}
	err := r.ParseMultipartForm(0)

	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
	}
	claims,err:=verifyLogin(r)
	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
	}
	t := table.NewWriter()
	t.AppendRow([]interface{}{"User ID",claims["user"]})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"IP Logged",claims["ip"]})
	t.AppendSeparator()
	tm := time.Unix(0,int64(claims["exp"].(float64)))
	t.AppendRow([]interface{}{"Token Expire",tm})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"Selected Course",getSelectedCourse(claims["user"].(string))})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"Selected Problem",getSelectedProblem(claims["user"].(string))})
	t.AppendSeparator()
	data, _ := json.Marshal(Output{t.Render(), ""})
	w.Write(data)
	return
}
