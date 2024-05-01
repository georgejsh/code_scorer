package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"server/libs"
	
)




func init() {
	/*liCommands = map[string] fn{
		"list":  libs.GetList,
		//"login": libs.getlogin,
		//	"logout": logout,
	}*/
	libs.AddCommand("help",libs.GetHelp,false)
	libs.AddCommand("login",libs.Getlogin,false)
	libs.AddCommand("reset",libs.Reset,false)
	libs.AddCommand("adduser",libs.AddUser,true)
	libs.AddCommand("list",libs.GetList,false)
	libs.AddCommand("select",libs.Select,false)
	libs.AddCommand("status",libs.GetStatus,false)
	libs.AddCommand("problem",libs.GetProblem,false)
	libs.AddCommand("tc",libs.GetTC,false)
	libs.AddCommand("submit",libs.Submit,false)
	
}



// shiftPath splits the given path into the first segment (head) and
// the rest (tail). For example, "/foo/bar/baz" gives "foo", "/bar/baz".

func main() {
	//db.Run()
	http.HandleFunc("/", libs.Serve)
	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
