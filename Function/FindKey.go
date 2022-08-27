package Function

import (
	"ShiroKeyCheck/AES_Encrypt"
	"ShiroKeyCheck/GlobalVar"
	"encoding/base64"
	"fmt"
)

func FindTheKey(Shirokeys string, Content []byte) bool {
	key, _ := base64.StdEncoding.DecodeString(Shirokeys)
	RememberMe1 := AES_Encrypt.AES_CBC_Encrypt(key, Content) //AES CBC加密
	RememberMe2 := AES_Encrypt.AES_GCM_Encrypt(key, Content) //AES GCM加密
	if HttpRequset2(RememberMe1, Shirokeys, "CBC") {
		fmt.Println("Find the Key!")
		fmt.Println("[+] CBC-KEY:", Shirokeys)
		fmt.Printf("[+] %s=%s\n", GlobalVar.ReqHeader, RememberMe1)
		return true
	}
	if HttpRequset2(RememberMe2, Shirokeys, "GCM") {
		fmt.Println("Find the Key!")
		fmt.Println("[+] GCM-KEY:", Shirokeys)
		fmt.Printf("[+] %s=%s\n", GlobalVar.ReqHeader, RememberMe2)
		return true
	}
	return false
}
