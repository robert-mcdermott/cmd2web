package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/auth"
)

var cmd []string
var user, pass string

func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Ctx.WriteString(fmt.Sprintf(four04html, "Sorry, correct URL required!"))
}

type CmdController struct {
	beego.Controller
}

func (this *CmdController) Get() {
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		out = []byte(err.Error())
	}
	this.Ctx.WriteString(fmt.Sprintf(cmdhtml, strings.Join(cmd, " "), time.Now().Format(time.RFC1123), string(out)))
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func SecretAuth(username, password string) bool {
	if username == user && password == pass {
		return true
	}
	return false
}

func createCerts(tempdir string, certKey, certPub []byte) {
	err := ioutil.WriteFile(fmt.Sprintf("%s/cert.key", tempdir), certKey, 0644)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/cert.pem", tempdir), certPub, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteCerts(tempdir string) {
	time.Sleep(time.Second * 10)
	err := os.Remove(fmt.Sprintf("%s/cert.key", tempdir))
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(fmt.Sprintf("%s/cert.pem", tempdir))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cmd = os.Args[1:]
	user = "cmd2web"
	pass = RandString(8)
	path := RandString(30)
	tempdir := os.TempDir()
	port, err := GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	createCerts(tempdir, certKey, certPub)

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Println("\nAccess Information")
	fmt.Println("-------------------------------------")
	fmt.Printf("URL: https://%s:%d/%s\n", hostname, port, path)
	fmt.Printf("Username: %s\n", user)
	fmt.Printf("Password: %s\n\n", pass)
	fmt.Printf("Easy Access URL: https://%s:%s@%s:%d/%s\n\n", user, pass, hostname, port, path)

	beego.BConfig.RunMode = "prod"
	beego.BConfig.Listen.EnableHTTP = false
	//beego.BConfig.Listen.HTTPPort = port
	beego.BConfig.Log.AccessLogs = true
	beego.BConfig.Listen.EnableHTTPS = true
	beego.BConfig.Listen.HTTPSPort = port
	beego.BConfig.Listen.HTTPSCertFile = fmt.Sprintf("%s/cert.pem", tempdir)
	beego.BConfig.Listen.HTTPSKeyFile = fmt.Sprintf("%s/cert.key", tempdir)
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.MaxMemory = 134217728 // 128MiB
	authPlugin := auth.NewBasicAuthenticator(SecretAuth, "cmd2web")
	//beego.InsertFilter("*", beego.BeforeRouter, authPlugin)
	//beego.InsertFilter("*", beego.BeforeExec, authPlugin)
	//Tried the above two variations, but the "beego.BeforeStatic" is the only
	//one the that will prompt for password for static content (files in this case)
	beego.InsertFilter("*", beego.BeforeStatic, authPlugin)
	beego.SetStaticPath(fmt.Sprintf("/%s/files", path), ".")
	beego.Router(path, &CmdController{})
	beego.Router("/*", &MainController{})
	go deleteCerts(tempdir)
	beego.Run()
}
