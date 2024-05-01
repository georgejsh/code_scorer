package libs

import (
	"encoding/json"
	"net/http"
	//"fmt"
)

var input_select =[]Help {
	Help{"name","enum",false,[]string{"course","problem"}},
	Help{"id","string",false,[]string{}},
}
func Select(w http.ResponseWriter, r *http.Request) {
	var head string
	w.Header().Set("Content-Type", "application/json")
	head, r.URL.Path = shiftPath(r.URL.Path)
	if head == "help" {
		output, _ := json.Marshal(input_select)
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
		
		pass:=selectCourse(claims["user"].(string),r.FormValue("id"))
		//fmt.Printf("results %s\n",result)
		if pass != true {
			// Handle error
			data, _ := json.Marshal(Output{"","Wrong course ID"})
			w.Write(data)
			return
		}
		data, _ := json.Marshal(Output{"Selected Course", ""})
		w.Write(data)
		return
	}else{
		course:=getSelectedCourse(claims["user"].(string))
		if course=="Not Selected"{
			data, _ := json.Marshal(Output{"","No Course Selected yet!!"})
			w.Write(data)
			return
		}
		//fmt.Printf("results %s\n",course)
		pass:=selectProblem(claims["user"].(string),course,r.FormValue("id"))
		if pass != true {
			// Handle error
			data, _ := json.Marshal(Output{"","Wrong problem ID!!"})
			w.Write(data)
			return
		}
		data, _ := json.Marshal(Output{"Selected Problem", ""})
		w.Write(data)
		return
	}

}
