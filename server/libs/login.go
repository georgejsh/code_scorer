package libs

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	//"fmt"
	//"reflect"
)
var input =[]Help {
	Help{"user_name","string",false,[]string{}},
	Help{"password","string",true,[]string{}},
}
func Getlogin(w http.ResponseWriter, r *http.Request) {
	//fmt.Printf("got /hello request\n")
	//fmt.Printf("Url %s\n",r.URL.Path)
	//io.WriteString(w, "Hello, HTTP!\n")
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input)
		data, _ := json.Marshal(Output{string(output), ""})
		
		w.Write(data)
		return
	}
	//fmt.Printf("Url %s\n",r.URL.Path)
	err := r.ParseMultipartForm(0)

	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
	}
	//fmt.Println("Username ",r.FormValue ("user_name"))
	//fmt.Println("Password ",r.FormValue ("password"))
	//fmt.Println("UTC ",r.FormValue ("utc"))
	utc,err := rSA_OAEP_Decrypt(r.FormValue ("utc"))
	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", "Time Input Is Wrong"})
		w.Write(data)
		return
	}
	//fmt.Println("UTC ",utc)
	pass := checkTime(utc,time.Now().UnixNano(),200*1000*1000);
	if pass != true {
		// Handle error
		data, _ := json.Marshal(Output{"", "TIme Expired"})
		w.Write(data)
		return
	}
	pass_c,err := rSA_OAEP_Decrypt(r.FormValue ("password"))
	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", "Invalid User Info"})
		w.Write(data)
		return
	}
	pass_db,err :=getPassword(r.FormValue ("user_name"))
	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", "Invalid User Info"})
		w.Write(data)
		return
	}
	if strings.HasPrefix(pass_c,pass_db) ==false {
		data, _ := json.Marshal(Output{"", "Invalid User Info"})
		w.Write(data)
		return
	}
	i := strings.Index(pass_c, "_") +1
	at_pass_utc:= pass_c[i:]
	//fmt.Println("passc ",pass_c)
	//fmt.Println("passc ",pass_db)
	//fmt.Println("passc ",at_pass_utc)
	
	//getClaims()
	pass = checkTimeStr(at_pass_utc,utc,500*1000*1000);
	if pass != true {
		// Handle error
		data, _ := json.Marshal(Output{"", "Time Expired"})
		w.Write(data)
		return
	}
	ip,err := getIP(r)
	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
	}
	user:=r.FormValue ("user_name")
	
	//fmt.Println("user ",user)
	output,err:=generateJWT(user,ip)
	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
	}
	//fmt.Println("login user ",user)
	clearSelection(user);
	createSelection(user)
	data, _ := json.Marshal(Output{output, ""})
	
	w.Write(data)
}
