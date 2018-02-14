package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
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
	this.Ctx.WriteString(string("That's not the secret string!"))
}

type CmdController struct {
	beego.Controller
}

func (this *CmdController) Get() {
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}
	this.Ctx.WriteString(string(out))
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

func main() {
	cmd = os.Args[1:]
	user = "cmd2web"
	pass = RandString(8)
	path := RandString(30)
	port, err := GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nURL: https://%s:%d/%s\n", hostname, port, path)
	fmt.Printf("Username: %s\n", user)
	fmt.Printf("Password: %s\n\n", pass)
	beego.BConfig.RunMode = "prod"
	beego.BConfig.Listen.EnableHTTP = false
	//beego.BConfig.Listen.HTTPPort = port
	beego.BConfig.Log.AccessLogs = true
	beego.BConfig.Listen.EnableHTTPS = true
	beego.BConfig.Listen.HTTPSPort = port
	beego.BConfig.Listen.HTTPSCertFile = "ssl/cert.pem"
	beego.BConfig.Listen.HTTPSKeyFile = "ssl/cert.key"
	beego.BConfig.WebConfig.DirectoryIndex = true
	authPlugin := auth.NewBasicAuthenticator(SecretAuth, "cmd2web")
	//beego.InsertFilter("*", beego.BeforeRouter, authPlugin)
	//beego.InsertFilter("*", beego.BeforeExec, authPlugin)
	//Tried the above two variations, but the "beego.BeforeStatic" is the only
	//one the that will prompt for password for static content (files in this case)
	beego.InsertFilter("*", beego.BeforeStatic, authPlugin)
	beego.SetStaticPath(fmt.Sprintf("/%s/files", path), ".")
	beego.Router(path, &CmdController{})
	beego.Router("/*", &MainController{})
	beego.Run()
}
