package libs

import (
	"encoding/json"
	"net/http"
	"fmt"
)

var input_list =[]Help {
	Help{"name","enum",false,[]string{"course","problem"}},
}
func GetList(w http.ResponseWriter, r *http.Request) {
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input_list)
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
	if r.FormValue ("name") =="course"{
		
		result :=getCourses(claims["user"].(string))
		//fmt.Printf("results %s\n",result)
		
		data, _ := json.Marshal(Output{result, ""})
		w.Write(data)
		return
	}else{
		result :=getSelectedCourse(claims["user"].(string))
		if result == "Not Selected"{
			data, _ := json.Marshal(Output{"", "No Courses Selected"})
			w.Write(data)
			return
		}
		fmt.Printf("results %s\n",result)
		output:= getProblems(claims["user"].(string),result)
		data, _ := json.Marshal(Output{output, ""})
		w.Write(data)
		return
	}

}
