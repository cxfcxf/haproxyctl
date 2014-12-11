package main

import (
	"fmt"
	"net"
	"flag"
	"strings"
	"regexp"
	"io/ioutil"
)

var (
	stats = flag.Bool("stats", false, "check status")
	enable = flag.String("enable", "", "enable px/sv / put px/sv into ACTIV mode")
	disable = flag.String("disable", "", "disable px/sv / put px/sv into MAINT mode")
)

const (
	SOCKET_TYPE	= "unix"
)

type haProxy struct {
	Sock	string
}

func (h haProxy) Exec(command string) string {
	sock, err := net.Dial(SOCKET_TYPE, h.Sock)
    if err != nil { panic(err) }
    defer sock.Close()

    cmd := fmt.Sprintf("%s\r\n", command)
    _, err = sock.Write([]byte(cmd))
    if err != nil { panic(err) }

    res, err := ioutil.ReadAll(sock)
    if err != nil { panic(err) }
    return string(res)
}

func (h haProxy) Showhealth() {
    output := h.Exec("show stat")
    out := strings.Split(output, "\n")
	fmt.Println("now printing Health Check")
	fmt.Printf("=====================================\n")
    for _, line := range out {
        if line != "" {
            data := strings.Split(line, ",")
            fmt.Printf("%-10s %-10s %-10s %-10s\n", data[0], data[1], data[17], data[18])
        }
    }
	fmt.Printf("=====================================\n\n\n")
}

func (h haProxy) ShowRegexp(reg string) {
	re := regexp.MustCompile(reg)

    output := h.Exec("show stat")
    out := strings.Split(output, "\n")
	fmt.Printf("now printing %s Health Check\n", reg)
	fmt.Printf("=====================================\n")
    for i, line := range out {
        if i==0 || re.MatchString(line) {
            data := strings.Split(line, ",")
            fmt.Printf("%-10s %-10s %-10s %-10s\n", data[0], data[1], data[17], data[18])
        }
    }
	fmt.Printf("=====================================\n\n\n")

}

func (h haProxy) DisableServer(px string, sv string) {
	c := fmt.Sprintf("disable server %s/%s", px, sv)
	h.Exec(c)
	fmt.Printf("Server %s/%s has been disabled\n", px, sv)
	h.ShowRegexp(sv)
}

func (h haProxy) EnableServer(px string, sv string) {
    c := fmt.Sprintf("enable server %s/%s", px, sv)
    h.Exec(c)
	fmt.Printf("Server %s/%s has been enabled\n", px, sv)
    h.ShowRegexp(sv)
}

func main() {
	flag.Parse()
   	haproxy := new(haProxy)
   	haproxy.Sock = "/var/lib/haproxy/stats"
	
	if *stats {
		haproxy.Showhealth()
	} else if *enable != "" {
		en := strings.Split(*enable, "/")
		haproxy.EnableServer(en[0], en[1])
	} else if *disable != "" {
		dis := strings.Split(*disable, "/")
		haproxy.DisableServer(dis[0], dis[1])
	} else {
		flag.PrintDefaults()
	}

}
