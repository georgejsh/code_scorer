package libs

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"net/http"
	"net"
	"time"
	"github.com/golang-jwt/jwt"
)

type Output struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
type Help struct {
	Name   string `json:"name"`
	Type   string `json:"typ"`
	Secret bool   `json:"secret"`
	Values []string `json:"values"`
}

func checkTime(x string,y int64 ,diff int64) bool {
	
	i, _ := strconv.ParseInt(x, 10, 64)
	//fmt.Println("x ",x)
	//fmt.Println("y ",y)
	//fmt.Println("diff ",diff)
	if y-i > diff {
		return false
	}
	return true

}
func checkTimeStr(x string,y string ,diff int64) bool {
	
	i, _ := strconv.ParseInt(x, 10, 64)
	j,_:= strconv.ParseInt(y, 10, 64)
	//fmt.Println("x ",x)
	//fmt.Println("y ",y)
	//fmt.Println("diff ",diff)
	if j-i > diff {
		return false
	}
	return true

}

var sampleSecretKey = []byte("BD05D63F0CFB6D88E87B1BA72F3065267DC5818097158918425ADA25500D44AF")

func generateJWT(user, ip string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(120 * time.Minute).UnixNano()
	claims["user"] = user
	claims["ip"] = ip
	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's an error with the signing method")
		}
		return sampleSecretKey, nil
	})

	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("token data not valid")
}

func verifyJWT(tokenString, ip,r_utc string) (jwt.MapClaims, error) {
	claims, err := getClaims(tokenString)
	if err != nil {
		return nil, err
	}
	exp_string := claims["exp"].(float64)
	//exp_time, _ := strconv.ParseInt(exp_string, 10, 64)
	//if err != nil {
	//	return nil, err
	//}
	if float64(time.Now().UnixNano()) > exp_string {
		return nil, fmt.Errorf("Token expired")
	}
	if ip != claims["ip"].(string) {
		return nil, fmt.Errorf("Wrong ip found")
	}
	utc,err := rSA_OAEP_Decrypt(r_utc)
	if err != nil {
		// Handle error
		return nil, err
	}
	pass := checkTime(utc,time.Now().UnixNano(),200*1000*1000);
	if pass != true {
		// Handle error
		return nil, fmt.Errorf("Time Expired")
	}
	claims["utc"] = utc
	return claims, nil
}

func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}


func getIP(r *http.Request) (string, error) {
    //Get IP from the X-REAL-IP header
    ip := r.Header.Get("X-REAL-IP")
    netIP := net.ParseIP(ip)
    if netIP != nil {
        return ip, nil
    }

    //Get IP from X-FORWARDED-FOR header
    ips := r.Header.Get("X-FORWARDED-FOR")
    splitIps := strings.Split(ips, ",")
    for _, ip := range splitIps {
        netIP := net.ParseIP(ip)
        if netIP != nil {
            return ip, nil
        }
    }

    //Get IP from RemoteAddr
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
        return "", err
    }
    netIP = net.ParseIP(ip)
    if netIP != nil {
        return ip, nil
    }
    return "", fmt.Errorf("No valid ip found")
}

func verifyLogin(r *http.Request)(jwt.MapClaims, error){
	ip,err := getIP(r)
	if err != nil {
		// Handle error
		return nil,err
	}
	reqToken := r.Header.Get("Authorization")
	//fmt.Printf("Req Token%s \n",reqToken)
	return verifyJWT(reqToken,ip,r.FormValue ("utc"))
	/*claims,err:=getClaims(reqToken)
	if err != nil {
		// Handle error
		data, _ := json.Marshal(Output{"", err.Error()})
		w.Write(data)
		return
	}
	fmt.Printf("%s %s \n",claims["user"],claims["ip"])
	if ip!=claims["ip"]{
		data, _ := json.Marshal(Output{"", "Ip Mismatch"})
		w.Write(data)
		return
	}*/
}

func getProblemID(r * http.Request) (jwt.MapClaims,error) {
	err := r.ParseMultipartForm(0)

	if err != nil {
		return nil ,err
	}
	claims,err:=verifyLogin(r)
	if err != nil {
		return nil ,err
	}
	course:=getSelectedCourse(claims["user"].(string))
	if course=="Not Selected"{
		return nil ,fmt.Errorf("No Course Selected yet!!")
	}
	desc:=getSelectedProblem(claims["user"].(string))
	if desc=="Not Selected"{
		return nil ,fmt.Errorf("No Problem Selected yet!!")
	}
	claims["course"]=course
	claims["problem"]=desc
	return claims,nil
}