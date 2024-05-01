package libs

import (
	"encoding/json"
	"net/http"
	"strings"
	//"fmt"
	//"log"
)
var input_add =[]Help {
	Help{"user_name","string",false,[]string{}},
	Help{"password","string",true,[]string{}},
}/*
func init(){
	claims,err:=getClaims("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTIzMzM2MjU2ODA5MTgxMDAsImlwIjoiMTI3LjAuMC4xIiwidXNlciI6ImNzMTVtMDIxIn0.FBePmtm01-tyQeArkyv6AIhtaWqN0OLTEOLkYUiS5Xc")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s \n",claims["user"],claims["ip"])
}*/
func AddUser(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("got /hello request\n")
	//io.WriteString(w, "Hello, HTTP!\n")
	
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input_add)
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
	pass_c,err := rSA_OAEP_Decrypt(r.FormValue ("password"))
	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", "Invalid Password Info"})
		w.Write(data)
		return
	}

	i := strings.Index(pass_c, "_")+1
	at_pass_utc:= pass_c[i:] 
	pass := checkTimeStr(at_pass_utc,claims["utc"].(string),500*1000*1000);
	if pass != true {
		// Handle error
		data, _ := json.Marshal(Output{"", "Time Expired"})
		w.Write(data)
		return
	}

	pass_c =pass_c[:i-1]
	addUser(r.FormValue ("user_name"),pass_c)
	data, _ := json.Marshal(Output{"User Added Successfully", ""})
	w.Write(data)
}
