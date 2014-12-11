package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

var (
	showhealth = flag.Bool("showhealth", false, "switch for health check")
	showbackend = flag.Bool("showbackend", false, "switch for backend info")
	status     = flag.Bool("status", false, "switch for current running status")
	enable     = flag.String("enable", "", "enable px/sv / put px/sv into ACTIV mode")
	disable    = flag.String("disable", "", "disable px/sv / put px/sv into MAINT mode")
	f	   = flag.String("f", "/etc/haproxy/haproxy.cfg", "point configration file, default /etc/haproxy/haproxy.cfg")
)

const (
	SOCKET_TYPE = "unix"
)

type haProxy struct {
	Pid  string
	Sock string
	Bin  string
}

func (h *haProxy) Loadenv(cfg string) {
	//Load sock and pid to haProxy struct
	repf := regexp.MustCompile(`pidfile`)
	reso := regexp.MustCompile(`stats\ socket`)

	c, _ := ioutil.ReadFile(cfg)
	cn := strings.Split(string(c), "\n")
	for _, line := range cn {
		if repf.MatchString(line) {
			p, _ :=  ioutil.ReadFile(strings.Fields(line)[1])
			h.Pid = strings.Trim(string(p), "\n")
		}
		if reso.MatchString(line) {
			h.Sock = strings.Fields(line)[2]
		}
	}
	//Load binary path to 
	
	

}

func (h *haProxy) Showstatus() {
	if len(h.Pid) > 0 {
		fmt.Printf("haproxy is running on pid %s.\nthese ports are used and guys are connected:\n", h.Pid)
		shell := fmt.Sprintf("lsof -ln -i |awk '$2 ~ /%s/ {print $8\" \"$9}'", h.Pid)
		cmd := exec.Command("sh", "-c", shell)
		res, _ := cmd.CombinedOutput()
		fmt.Println(string(res))
	} else {
		fmt.Printf("haproxy is not running\n")
	}
}

func (h *haProxy) Exec(command string) string {
	sock, err := net.Dial(SOCKET_TYPE, h.Sock)
	if err != nil {
		panic(err)
	}
	defer sock.Close()

	cmd := fmt.Sprintf("%s\r\n", command)
	_, err = sock.Write([]byte(cmd))
	if err != nil {
		panic(err)
	}

	res, err := ioutil.ReadAll(sock)
	if err != nil {
		panic(err)
	}
	return string(res)
}

func (h *haProxy) Showhealth() {
	output := h.Exec("show stat")
	out := strings.Split(output, "\n")
	fmt.Printf("\nnow printing Health Check...\n\n")
	for _, line := range out {
		if line != "" {
			data := strings.Split(line, ",")
			fmt.Printf("%-10s %-10s %-10s %-10s\n", data[0], data[1], data[17], data[18])
		}
	}
}

func (h *haProxy) ShowRegexp(reg string) {
	re := regexp.MustCompile(reg)

	output := h.Exec("show stat")
	out := strings.Split(output, "\n")
	fmt.Printf("\nnow printing %s Health Check\n", reg)
	for i, line := range out {
		if i == 0 || re.MatchString(line) {
			data := strings.Split(line, ",")
			fmt.Printf("%-10s %-10s %-10s %-10s\n", data[0], data[1], data[17], data[18])
		}
	}
}

func (h *haProxy) DisableServer(px string, sv string) {
	c := fmt.Sprintf("disable server %s/%s", px, sv)
	h.Exec(c)
	fmt.Printf("Server %s/%s has been disabled\n", px, sv)
	h.ShowRegexp(sv)
}

func (h *haProxy) EnableServer(px string, sv string) {
	c := fmt.Sprintf("enable server %s/%s", px, sv)
	h.Exec(c)
	fmt.Printf("Server %s/%s has been enabled\n", px, sv)
	h.ShowRegexp(sv)
}

func main() {
	flag.Parse()
	haproxy := new(haProxy)
	if len(*f) > 0 {
		haproxy.Loadenv(*f)
	} else {
		haproxy.Loadenv("/etc/haproxy/haproxy.cfg")
	}

	if *status {
		haproxy.Showstatus()
	}

	if *showbackend {
		haproxy.ShowRegexp("BACKEND")
	}

	if *showhealth {
		haproxy.Showhealth()
	}
	if *enable != "" {
		en := strings.Split(*enable, "/")
		haproxy.EnableServer(en[0], en[1])
	}
	if *disable != "" {
		dis := strings.Split(*disable, "/")
		haproxy.DisableServer(dis[0], dis[1])
	} 
	// flag.PrintDefaults()
}
