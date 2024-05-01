package libs
import (
    "crypto/rsa"
    "encoding/hex"
	"crypto/x509"
    "encoding/pem"
	"log"
	"io/ioutil"
)
var privateKey *rsa.PrivateKey
func rSA_OAEP_Decrypt(cipherText string) (string,error) {
    //ct, _ := base64.StdEncoding.DecodeString(cipherText)
	ct,err := hex.DecodeString(cipherText)
	if err != nil {
		return "",err
	}
    plaintext, err := rsa.DecryptPKCS1v15(nil, privateKey, ct)
    if err != nil {
		return "",err
	}
    //fmt.Println("Plaintext:", string(plaintext))
    return string(plaintext),nil
}

func init(){
	pemString, err := ioutil.ReadFile("./libs/id_rsa")     // the file is inside the local directory
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(pemString)
	privPemBytes := block.Bytes
    
    var parsedKey interface{}
    //PKCS1
    if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
        //If what you are sitting on is a PKCS#8 encoded key
        if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			log.Fatal(err)
            
        }
    }

    
    var ok bool
    privateKey, ok = parsedKey.(*rsa.PrivateKey)
	if !ok {
        log.Fatal(err)
    }
	//fmt.Printf("%s",RSA_OAEP_Decrypt("5cd5c085e4b938d6731d483075bfe9e74dc87db00a779104adbf83a3f759a0546b89cf27e1c2f587d35b3ebba49f2ab34a96d7049bddc98eeba7375d686724461eca7b1383a989eb8fd4cc872018ea09adc086372e8335426a7aa69d524894a601bf563cd497f6a2dd9bffeccb5555ed37cb1b93cb7daa70a030da5ebd7bfced218794de59f762f61d6ce9c3d640d4a3e689b005d5181d6a0d21635797f45da499dc4ae53674be3a2e1bf71f1215de745e788ece3f2b5a42d998daca0179ff388ebf624db9be35aaa170ee18c0586f4874e285942108f61d02aa20e39f2f8e7b0b334da647c1a647ce3b7ce1f288327c22108959b3ac13c9419b928b6441d0550dc0cd192cbaa6f78ab30d8b72029f00aa7db4150307ba1fcd183bc7192f62367aec3ad17b77047802f2be9d8e6b7d8c42066f44268f82583eaff2b8e838840e60ab50d74a841d62a6c16de0cd0fb523465e5dedff36d1d07688012107fb9f5a67abfced03ec4dc86e58cb8676000ce82d1042702ed8b0cb8939e36f1dd8c442"))
}