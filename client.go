package main

import (
	"errors"
	"bufio"
	"github.com/kr/pretty"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"fmt"
	"github.com/gwuhaolin/lightsocks/cmd"
	"net/url"
	"time"
	"io/ioutil"
	"github.com/mitchellh/go-homedir"
	"net/http"
	"path"
	"strings"
	"net"
	"os"
	"encoding/json"
)
const (
	DefaultListenAddr = ":1080"
	VERSION = "1"
)

var(
 gLocalConn net.Listener
 configPath string
)

type AutoUpdate struct {
	Url string 
	Softname string 
	CurrVer string 
	CurPath string 
}

const authserver = "http://111.231.82.173:9000/auth"
//const authserver = "http://127.0.0.1:9000/auth"

type PublicIp struct{
	RemoteAddr  string `json:"remote"`
	//ListenAddr string `json:"listen"`
	Msg  string `json:"msg"`
	Code int `json:"code"`
}

func Auth() {
	config := &cmd.Config{
		ListenAddr: DefaultListenAddr,
	}
	config.ReadConfig()
	//config.SaveConfig()
	ext := get_external()
	var public PublicIp
	for {
		resp, err := http.PostForm(authserver, url.Values{"name":{config.Auth},"ext":{ext},"auth":{config.Auth}})
		if err != nil {
			fmt.Printf("请检查网络")
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err := json.Unmarshal([]byte(body), &public); err == nil {
			if public.Code == -1 {
				fmt.Println(`Error:`, public.Msg)
			}
		}
		SaveConfig(public.RemoteAddr)
		time.Sleep(5 * time.Second)
	}
}
func  SaveConfig(DynamicIp string) {
	home, _ := homedir.Dir()
	configFilename := "dynamic"
	configPath = path.Join(home, configFilename)
	//configJson, _ := json.MarshalIndent(, "", "	")
	err := ioutil.WriteFile(configPath, []byte(DynamicIp), 0644)
	if err != nil {
		fmt.Errorf("保存配置到文件 %s 出错: %s", configPath, err)
	}
}
func getIp () string {
	home, _ := homedir.Dir()
	configFilename := "dynamic"
	configPath = path.Join(home, configFilename)
	fileIn, fileInErr := os.Open(configPath)
	if fileInErr != nil{
		fmt.Println(fileInErr)
	}
	finReader := bufio.NewReader(fileIn)
	inputString, _ := finReader.ReadString('\n')
	newString := strings.Replace(inputString,"\n","",-1)
	return newString
}


func get_external() string {
	resp, err := http.Get("http://icanhazip.com")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

func handle() {
	remoteIp := make(chan string)
	//for {
		sourceConn, err := gLocalConn.Accept()
		if err != nil {
			log.Println("server err:", err.Error())
		}
		go func() {
			for {
				select {
				case ip := <-remoteIp:
					targetConn, err := net.DialTimeout("tcp", ip, 10*time.Second)
					go func() {
						//defer targetConn.Close()
						_, err = io.Copy(targetConn, sourceConn)
						if err != nil {
							//fmt.Println("io.Copy 1 failed：", err.Error())
						}
					}()

					go func() {
						//defer targetConn.Close()
						_, err = io.Copy(sourceConn, targetConn)
						if err != nil {
							//fmt.Println("io.Copy 2 failed：", err.Error())
						}
					}()
				}
			}
		}()
		remoteIp <-getIp()
	//}
}

func forever(fn func()) {
	f := func() {
		defer func() {
			if r := recover(); r != nil {
				debug.PrintStack()
				pretty.Println("Recover from error:", r)
			}
		}()
		fn()
	}
	for {
		f()
	}
}


func main(){
	pretty.Println("GOOS:", runtime.GOOS, "GOARCH:", runtime.GOARCH)
	au := &AutoUpdate{
		Url: "http://111.231.82.173/file/",
		CurrVer: VERSION,
	}
	switch runtime.GOOS {
	case "windows":
		au.Softname = "client.exe"
		au.CurPath,_ = filepath.Abs(filepath.Dir(au.Softname))
		if err :=au.WUpdate(); err!=nil{
			fmt.Println(err)
		}
	case "linux":
		au.Softname = "linux_client-x64"
		au.CurPath,_ = filepath.Abs(filepath.Dir(au.Softname))
		if err :=au.WUpdate(); err!=nil{
			fmt.Println(err)
		}
	case "darwin":
		au.Softname = "darwin_client-x64"
		au.CurPath,_ = filepath.Abs(filepath.Dir(au.Softname))
		if err :=au.WUpdate(); err!=nil{
			fmt.Println(err)
		}
	case "freebsd":
		au.Softname = "freebsd_client-x64"
		au.CurPath,_ = filepath.Abs(filepath.Dir(au.Softname))
		if err :=au.WUpdate(); err!=nil{
			fmt.Println(err)
		}
	}
	go Auth()
	config := &cmd.Config{
		ListenAddr: DefaultListenAddr,
	}
	config.ReadConfig()
	fmt.Println("服务启动成功，服务地址：", config.ListenAddr)
	fmt.Println("授权",config.Auth)
	fmt.Println("Ready to connnected ...")
	localConn, err := net.Listen("tcp", config.ListenAddr) 
	if err != nil {
		fmt.Println(err.Error())
	}
	gLocalConn = localConn
	forever(handle)
}

func (au AutoUpdate) WUpdate() error {
	resp, err := http.Get(au.Url + "version.txt")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	newVer := strings.Replace(string(body),"\n","",-1)
	if newVer != au.CurrVer {
		fmt.Println("有新版本请去下载对应系统版本:",au.Url+au.Softname)
	}
	return nil
}

func (au AutoUpdate) getNewVer() error {
	client := http.Client{Timeout: 900 * time.Second}
	resp, err := client.Get(au.Url + "client.exe")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.Status == "200 OK" {
		newFile, err := os.Create("update.exe")
		if err != nil {
			return err
		}
		defer newFile.Close()
		_, err = io.Copy(newFile, resp.Body)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New(resp.Status)
	}
}

func (au AutoUpdate) copyFile() bool {
	os.Rename("client.exe","clientbak.exe")
	error := os.Rename("update.exe","client.exe")
	if error != nil {
		return false
	} else {
		return true
	}
}

