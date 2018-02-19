package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/auth"
)

var cmd []string
var user, pass string

// generates a random string of chars of n lenght
func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

// mainController provides attempts to access the server withouth the correct path
// an error message
type mainController struct {
	beego.Controller
}

func (c *mainController) Get() {
	c.Ctx.WriteString(fmt.Sprintf(four04html, "Sorry, correct URL/Accesskey required!"))
}

// cmdController provides the command results output page
type cmdController struct {
	beego.Controller
}

func (c *cmdController) Get() {
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		out = []byte(err.Error())
	}
	if *rawFlag {
		c.Ctx.WriteString(string(out))
	} else {
		refresh := strconv.Itoa(*refreshFlag)
		if refresh != "0" {
			refresh = fmt.Sprintf("<meta http-equiv=\"refresh\" content=\"%s\">", refresh)
		} else {
			refresh = ""
		}
		c.Ctx.WriteString(fmt.Sprintf(cmdhtml, refresh, strings.Join(cmd, " "), time.Now().Format(time.RFC1123), string(out)))
	}
}

// find a free TCP port on the system that we can use
func getFreePort() (int, error) {
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

// validate the provided username and password
func secretAuth(username, password string) bool {
	if username == user && password == pass {
		return true
	}
	return false
}

// write the SSL cert/key to a temp location
// I couldn't find a way to provide the cert/key as a string so this is
// a workaround so I can still have a single file tool
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

// Delete the SSL cert/key from the the temp loacation after the server has started up
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
	//cmd = os.Args[1:]
	cmd = flag.Args()
	if len(cmd) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	user = "cmd2web"
	pass = randString(8)
	accessKey := randString(32)
	tempdir := os.TempDir()
	port, err := getFreePort()
	if err != nil {
		log.Fatal(err)
	}

	// put the certs in place
	createCerts(tempdir, certKey, certPub)

	// get the hostname of the system
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// beego server configuration
	beego.BConfig.RunMode = "prod"
	beego.BConfig.Listen.EnableHTTP = false
	//beego.BConfig.Listen.HTTPPort = port
	beego.BConfig.ServerName = hostname
	beego.BConfig.AppName = "cmd2web"
	beego.BConfig.Log.AccessLogs = true
	beego.BConfig.Listen.EnableHTTPS = true
	beego.BConfig.Listen.HTTPSPort = port
	beego.BConfig.Listen.HTTPSCertFile = fmt.Sprintf("%s/cert.pem", tempdir)
	beego.BConfig.Listen.HTTPSKeyFile = fmt.Sprintf("%s/cert.key", tempdir)
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.MaxMemory = 134217728 // 128MiB
	authPlugin := auth.NewBasicAuthenticator(secretAuth, "cmd2web")
	//beego.InsertFilter("*", beego.BeforeRouter, authPlugin)
	//beego.InsertFilter("*", beego.BeforeExec, authPlugin)
	//Tried the above two variations, but the "beego.BeforeStatic" is the only
	//one the that will prompt for password for static content (file or dir in this case)
	beego.InsertFilter("*", beego.BeforeStatic, authPlugin)
	// if the user provided an "--expose </path>" flag, expose the provided directory as a filesystem
	if *exposeFlag != "" {
		// make sure that the path to the directory or file the user provided exists
		if _, err := os.Stat(*exposeFlag); os.IsNotExist(err) {
			log.Fatal(err)
		}
		beego.SetStaticPath(fmt.Sprintf("/%s/file", accessKey), *exposeFlag)
	}
	// set the router that provides the command output
	beego.Router(accessKey, &cmdController{})
	// set the default router that doesn't match the static or command path/route
	beego.Router("/*", &mainController{})

	// delete the certs in the backgroup
	go deleteCerts(tempdir)

	// if an "--expire <int>" set a timer running in the background to exit after the desired number of minutes
	if *expireFlag != 0 {
		go func(t int) {
			<-time.After(time.Minute * time.Duration(t))
			fmt.Printf("\nTimeout of %d minutes expired, shutting down...\n\n", t)
			os.Exit(0)
		}(*expireFlag)
	}

	// provide user some information on stderr on how to access the server
	fmt.Fprintf(os.Stderr, "\nAccess Information\n")
	fmt.Fprintf(os.Stderr, "-------------------------------------\n")
	fmt.Fprintf(os.Stderr, "Command output: https://%s:%d/%s\n", hostname, port, accessKey)
	if *exposeFlag != "" {
		fmt.Fprintf(os.Stderr, "Exposed directory: https://%s:%d/%s/file/\n", hostname, port, accessKey)
	}
	fmt.Fprintf(os.Stderr, "Username: %s\n", user)
	fmt.Fprintf(os.Stderr, "Password: %s\n\n", pass)
	fmt.Fprintf(os.Stderr, "Easy Access URL: https://%s:%s@%s:%d/%s\n", user, pass, hostname, port, accessKey)
	fmt.Fprintf(os.Stderr, "-------------------------------------\n")
	// start the server
	beego.Run()
}
