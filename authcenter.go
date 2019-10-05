package main

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

func ClearPic() string {
	fileIn, fileInErr := os.Open("clearpic")
	if fileInErr != nil {
		fmt.Println("killhostiderror")
	}
	defer fileIn.Close()
	finReader := bufio.NewReader(fileIn)
	inputString, _ := finReader.ReadString('\n')
	newString := strings.Replace(inputString, "\n", "", -1)
	return newString
}
func Killip() string {
	fileIn, fileInErr := os.Open("killip")
	if fileInErr != nil {
		fmt.Println("killhostiderror")
	}
	defer fileIn.Close()
	finReader := bufio.NewReader(fileIn)
	inputString, _ := finReader.ReadString('\n')
	newString := strings.Replace(inputString, "\n", "", -1)
	return newString
}
func Allkill() string {
	fileIn, fileInErr := os.Open("allkill")
	if fileInErr != nil {
		fmt.Println("allkillerror")
	}
	defer fileIn.Close()
	finReader := bufio.NewReader(fileIn)
	inputString, _ := finReader.ReadString('\n')
	newString := strings.Replace(inputString, "\n", "", -1)
	return newString
}
func getContent() []string {
	fileIn, fileInErr := os.Open("text")
	if fileInErr != nil {
		fmt.Println("get auth order error")
	}
	defer fileIn.Close()
	finReader := bufio.NewReader(fileIn)
	var fileList []string
	for {
		inputString, readerError := finReader.ReadString('\n')
		if readerError == io.EOF {
			break
		}
		fileList = append(fileList, strings.Replace(inputString, "\n", "", -1))
	}
	return fileList
}
func getIp() string {
	fileIn, fileInErr := os.Open("ip")
	if fileInErr != nil {
		fmt.Println("get dynamic ip error")
	}
	defer fileIn.Close()
	finReader := bufio.NewReader(fileIn)
	inputString, _ := finReader.ReadString('\n')
	newString := strings.Replace(inputString, "\n", "", -1)
	return newString
}
func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, fmt.Errorf("not in array")
}

func main() {
	gin.DisableConsoleColor()
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r := gin.Default()
	m := melody.New()
	r.POST("/auth", func(c *gin.Context) {
		list := getContent()
		ip := getIp()
		message := c.PostForm("name")
		ext := c.PostForm("ext")
		auth := c.PostForm("auth")
		log.Println(auth, ext)
		_, error := Contain(message, list)
		if error != nil {
			c.JSON(http.StatusOK, gin.H{"msg": "授权失败,联系开发者付款", "code": -1})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "认证通过", "code": 0, "remote": ip})
	})

	r.POST("/killip", func(c *gin.Context) {
		t := c.PostForm("killip")
		err := ioutil.WriteFile("killip", []byte(t), 0644)
		if err != nil {
			fmt.Errorf("killip保存失败: %s", err)
		}
		c.JSON(http.StatusOK, gin.H{"msg": "杀掉IP", "code": 0, "remote": t})
	})
	r.POST("/allkill", func(c *gin.Context) {
		t := c.PostForm("allkill")
		err := ioutil.WriteFile("allkill", []byte(t), 0644)
		if err != nil {
			fmt.Errorf("allkill保存失败: %s", err)
		}
		c.JSON(http.StatusOK, gin.H{"msg": "allkill", "code": 0, "remote": t})
	})
	r.POST("/save", func(c *gin.Context) {
		t := c.PostForm("ip")
		err := ioutil.WriteFile("ip", []byte(t), 0644)
		if err != nil {
			fmt.Errorf("保存配置到文件出错: %s", err)
		}
		c.JSON(http.StatusOK, gin.H{"msg": "认证通过", "code": 0, "remote": t})
	})
	r.POST("/ClearPic", func(c *gin.Context) {
		t := c.PostForm("clearpic")
		err := ioutil.WriteFile("clearpic", []byte(t), 0644)
		if err != nil {
			fmt.Errorf("allkill保存失败: %s", err)
		}
		c.JSON(http.StatusOK, gin.H{"msg": "allkill", "code": 0, "remote": t})
	})
	r.GET("/clearpic", func(c *gin.Context) {
		status := ClearPic()
		c.JSON(http.StatusOK, gin.H{"msg": "认证通过", "code": 0, "hostid": status})
	})
	r.GET("/Killip", func(c *gin.Context) {
		status := Killip()
		c.JSON(http.StatusOK, gin.H{"msg": "认证通过", "code": 0, "hostid": status})
	})
	r.GET("/Allkill", func(c *gin.Context) {
		status := Allkill()
		c.JSON(http.StatusOK, gin.H{"msg": "认证通过", "code": 0, "hostid": status})
	})
	r.GET("/", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		//fmt.Println(len(msg))
		m.Broadcast(msg)
	})

	r.Run(":9000")
}
