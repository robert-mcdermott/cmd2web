package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

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
