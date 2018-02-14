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
)

var cmd []string

func init() {
	rand.Seed(time.Now().UnixNano())
}


func RandString(n int) string {
        chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
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
	this.Ctx.WriteString(string("nice try!"))
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

func main() {
	cmd = os.Args[1:]
	path := RandString(30)
        port, err := GetFreePort()
        if err != nil {
            log.Fatal(err)
	}
        fmt.Printf("http://<this-server>:%d/%s\n", port, path)
        beego.BConfig.RunMode = "prod"
	beego.BConfig.Listen.HTTPPort = port 
        beego.BConfig.WebConfig.DirectoryIndex = true 
        beego.SetStaticPath(fmt.Sprintf("/%s/files", path),".")
	beego.Router(path, &CmdController{})
	beego.Router("/*", &MainController{})
	beego.Run()
}
