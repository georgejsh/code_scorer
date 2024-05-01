package libs

import (
	"encoding/json"
	"net/http"
	//"strings"
	//"fmt"
	//"log"
)
var input_addcourse =[]Help {
	Help{"id","string",false,[]string{}},
	Help{"name","string",false,[]string{}},
	Help{"year","string",false,[]string{}},

}
func AddCourse(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("got /hello request\n")
	//io.WriteString(w, "Hello, HTTP!\n")
	
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input_addcourse)
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
	// to do verify admin
	if isRoot(claims["user"].(string)) ==false{
		data, _ := json.Marshal(Output{"", "Not Authorised User"})
		w.Write(data)
		return
	}
	
	addCourse(r.FormValue ("id"),r.FormValue ("name"),r.FormValue ("year"))
	data, _ := json.Marshal(Output{"Course Added Successfully", ""})
	w.Write(data)
}

var input_addenroll =[]Help {
	Help{"userid","string",false,[]string{}},
	Help{"courseid","string",false,[]string{}},
	Help{"isadmin","string",false,[]string{}},
	
}
func AddEnroll(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("got /hello request\n")
	//io.WriteString(w, "Hello, HTTP!\n")
	
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input_addenroll)
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
	// to do verify admin
	if isAdmin(claims["user"].(string)) ==false{
		data, _ := json.Marshal(Output{"", "Not Authorised User"})
		w.Write(data)
		return
	}
	if !isRoot(claims["user"].(string)) && getSelectedCourse(claims["user"].(string))!=r.FormValue ("courseid"){
		data, _ := json.Marshal(Output{"", "Not Authorised User. Select the course"})
		w.Write(data)
		return
	}
	enroll(r.FormValue ("userid"),r.FormValue ("courseid"),r.FormValue ("isadmin"))
	data, _ := json.Marshal(Output{"Course Added Successfully", ""})
	w.Write(data)
}


