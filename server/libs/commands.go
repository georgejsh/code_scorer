package libs
import (
	"net/http"
)
type fn func(http.ResponseWriter, *http.Request)

var commands map[string]fn
var commands_access map[string]bool
func init() {
	commands =map[string] fn{}
	commands_access =map[string] bool{}
}
func AddCommand(name string, v fn ,admin bool){
	commands[name]=v
	commands_access[name]=admin
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)
	_, ok := commands[head]
	if !ok {
		return
	}
	commands[head](w, r)
}