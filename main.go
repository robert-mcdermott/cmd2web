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
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/auth"
)

var cmd []string
var user, pass string

// generates a random string of chars of n length
func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
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
func createCerts(certFileBase string, certKey, certPub []byte) {
	err := ioutil.WriteFile(fmt.Sprintf("%s-cert.key", certFileBase), certKey, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s-cert.pem", certFileBase), certPub, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Delete the SSL cert/key from the the temp location after the server has started up
func deleteCerts(certFileBase string) {
	time.Sleep(time.Second * 10)
	err := os.Remove(fmt.Sprintf("%s-cert.key", certFileBase))
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove(fmt.Sprintf("%s-cert.pem", certFileBase))
	if err != nil {
		log.Fatal(err)
	}
}

// checkCmdExists makes sure that the provided command is in the users path
func checkCmdExists(command string) error {
	_, err := exec.LookPath(command)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	//cmd = os.Args[1:]
	cmd = flag.Args()
	if len(cmd) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	// check to see if the command exists before starting webserver
	if err := checkCmdExists(cmd[0]); err != nil {
		// command was not found; let the user know and exit
		fmt.Fprintf(os.Stderr, "\nError: command \"%s\" not found\n\n", cmd[0])
		os.Exit(2)
	}

	user = "cmd2web"
	pass = randString(8)
	accessKey := randString(32)
	certFileBase := fmt.Sprintf("%s/%s", os.TempDir(), randString(12))

	var port int
	if *portFlag != 0 {
		lport := *portFlag
		port = lport
	} else {
		lport, err := getFreePort()
		if err != nil {
			log.Fatal(err)
		}
		port = lport
	}

	// put the certs in place
	createCerts(certFileBase, certKey, certPub)

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
	beego.BConfig.Listen.HTTPSCertFile = fmt.Sprintf("%s-cert.pem", certFileBase)
	beego.BConfig.Listen.HTTPSKeyFile = fmt.Sprintf("%s-cert.key", certFileBase)
	beego.BConfig.WebConfig.DirectoryIndex = true
	beego.BConfig.MaxMemory = 134217728 // 128MiB
	authPlugin := auth.NewBasicAuthenticator(secretAuth, "cmd2web")
	//beego.InsertFilter("*", beego.BeforeRouter, authPlugin)
	//beego.InsertFilter("*", beego.BeforeExec, authPlugin)
	//Tried the above two variations, but the "beego.BeforeStatic" is the only
	//one the that will prompt for password for static content (file or dir in this case)
	if !*noauthFlag {
		beego.InsertFilter("*", beego.BeforeStatic, authPlugin)
	}
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
	// set the router that allows user to remotely stop the cmd2web server.
	beego.Router(fmt.Sprintf("/%s/stop", accessKey), &stopController{})
	// set the default router that doesn't match the static, command or stop route
	beego.Router("/*", &mainController{})

	// delete the certs in the backgroup
	go deleteCerts(certFileBase)

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
	fmt.Fprintf(os.Stderr, "Command output:    https://%s:%d/%s\n", hostname, port, accessKey)
	fmt.Fprintf(os.Stderr, "Remote stop:       https://%s:%d/%s/stop\n", hostname, port, accessKey)
	if *exposeFlag != "" {
		fmt.Fprintf(os.Stderr, "Exposed directory: https://%s:%d/%s/file\n", hostname, port, accessKey)
	}
	if !*noauthFlag {
		fmt.Fprintf(os.Stderr, "\nCredentials:\n")
		fmt.Fprintf(os.Stderr, "  Username: %s\n", user)
		fmt.Fprintf(os.Stderr, "  Password: %s\n", pass)
		fmt.Fprintf(os.Stderr, "\nEasy Access URL:   https://%s:%s@%s:%d/%s\n", user, pass, hostname, port, accessKey)
	}
	fmt.Fprintf(os.Stderr, "-------------------------------------\n")
	// start the server
	beego.Run()
}
