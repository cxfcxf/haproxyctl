package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"net/http"
)

var (
	configcheck  = flag.Bool("configcheck", false, "check config file")
	action	     = flag.String("action", "", "tell me what action you wanna do?")
	execution    = flag.String("execution", "", "parameters to action")
	binding      = flag.String("binding", "", "http port you want program to bind to")
	f            = flag.String("f", "/etc/haproxy/haproxy.cfg", "point configration file, default /etc/haproxy/haproxy.cfg")
)

const (
	SOCKET_TYPE = "unix"
)

type haProxy struct {
	Pid  []string
	Sock []string
	Bin  string
	Cfg  string
}

func appendifuniq(slice []string, s string) []string {
    for _, ele := range slice {
        if ele == s {
            return slice
        }
    }
    return append(slice, s)
}

func (h *haProxy) Loadenv(cfg string) {
	h.Cfg = cfg
	//Load sock and pid to haProxy struct
	repf := regexp.MustCompile(`pidfile`)
	reso := regexp.MustCompile(`stats\ socket`)

	c, _ := ioutil.ReadFile(cfg)
	cn := strings.Split(string(c), "\n")
	for _, line := range cn {
		if repf.MatchString(line) {
			p, _ := ioutil.ReadFile(strings.Fields(line)[1])
			q := strings.Split(string(p), "\n")
			for _, l := range q[0:len(q)-1] {
				h.Pid = appendifuniq(h.Pid, string(l))
			}
		}
		if reso.MatchString(line) {
			h.Sock = appendifuniq(h.Sock, strings.Fields(line)[2])
		}
	}
	//Load binary path to haProxy struct
	rewh := regexp.MustCompile(`no haproxy in`)

	shell := fmt.Sprintf("which haproxy")
	cmd := exec.Command("sh", "-c", shell)
	res, _ := cmd.CombinedOutput()
	if rewh.MatchString(string(res)) {
		fmt.Println("your haproxy binary is not in the $PATH")
	} else {
		h.Bin = strings.Trim(string(res), "\n")
	}
}

func (h *haProxy) Showstatus() string {
	if len(h.Pid) > 0 {
	var status string
		for _, p := range h.Pid {
			head := fmt.Sprintf("haproxy is running on pid %s.\nthese ports are used and guys are connected:\n", p)
			shell := fmt.Sprintf("lsof -ln -i |awk '$2 ~ /%s/ {print $8\" \"$9}'", p)
			cmd := exec.Command("sh", "-c", shell)
			res, _ := cmd.CombinedOutput()
			status += head + string(res)
		}
		return status
	} else {
		return "haproxy is not running"
	}
}

func (h *haProxy) Exec(command string) []string {
	var result []string
	for _, socket := range h.Sock {
        	sock, err := net.Dial(SOCKET_TYPE, socket)
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
		result = append(result, string(res))
	}
	return result
}

func (h *haProxy) Showhealth() string {
	s := h.Exec("show stat")
	var result string
	for _, output := range s {
		out := strings.Split(output, "\n")
		res := fmt.Sprintf("\nnow printing Health Check...\n\n")
		for _, line := range out {
			if line != "" {
				data := strings.Split(line, ",")
				res += fmt.Sprintf("%-10s %-10s %-10s %-10s\n", data[0], data[1], data[17], data[18])
			}
		}
		result += res
	}
	return result
}

func (h *haProxy) ShowRegexp(reg string) string {
	re := regexp.MustCompile(reg)

	s := h.Exec("show stat")
	var result string
	for _, output := range s {
		out := strings.Split(output, "\n")
		res := fmt.Sprintf("\nnow printing %s Health Check\n", reg)
		for i, line := range out {
			if i == 0 || re.MatchString(line) {
				data := strings.Split(line, ",")
				res += fmt.Sprintf("%-10s %-10s %-10s %-10s\n", data[0], data[1], data[17], data[18])
			}
		}
		result += res
	}
	return result
}

func (h *haProxy) DisableServer(px string, sv string) string {
	c := fmt.Sprintf("disable server %s/%s", px, sv)
	h.Exec(c)
	return fmt.Sprintf("Server %s/%s has been disabled\n", px, sv) + h.ShowRegexp(sv)
}

func (h *haProxy) DisableAll(server string) {
	re := regexp.MustCompile(server)

	s := h.Exec("show stat")
	output := s[0]
	out := strings.Split(output, "\n")
	for i, line := range out {
		if i != 0 && line != "" {
			data := strings.Split(line, ",")
			if re.MatchString(data[1]) && data[17] == "UP" {
				c := fmt.Sprintf("disable server %s/%s", data[0], server)
				h.Exec(c)
			}
		}
	}
}

func (h *haProxy) EnableServer(px string, sv string) string {
	c := fmt.Sprintf("enable server %s/%s", px, sv)
	h.Exec(c)
	return fmt.Sprintf("Server %s/%s has been enabled\n", px, sv) + h.ShowRegexp(sv)
}

func (h *haProxy) EnableAll(server string) {
	re := regexp.MustCompile(server)
	rest := regexp.MustCompile(`(?i)Down|MAINT`)

	s := h.Exec("show stat")
	output := s[0]
	out := strings.Split(output, "\n")
	for i, line := range out {
		if i != 0 && line != "" {
			data := strings.Split(line, ",")
			if re.MatchString(data[1]) && rest.MatchString(data[17]) {
				c := fmt.Sprintf("enable server %s/%s", data[0], server)
				h.Exec(c)
			}
		}
	}
}

func (h *haProxy) Configcheck() {
	shell := fmt.Sprintf("%s -c -f %s", h.Bin, h.Cfg)
	cmd := exec.Command("sh", "-c", shell)
	res, _ := cmd.CombinedOutput()
	fmt.Println(strings.Trim(string(res), "\n"))
}

//Binding Handler Section
func handler(w http.ResponseWriter, r *http.Request, h *haProxy) {
	usage := "please use /haproxyctl?action=xxxx&exec=yyyy"
	if r.URL.Path != "/haproxyctl" {
		fmt.Fprintf(w, usage)
	} else {
		err := r.ParseForm()
		if err != nil { panic(err) }
		uri := r.Form
		if uri.Get("action") == "" {
			fmt.Fprintf(w, "you need to specify action")
		} else {
			action := uri.Get("action")
			execution := uri.Get("execution")
			switch action {
			case "showstatus":
				fmt.Fprintf(w, h.Showstatus())
			case "socketexec":
				res := h.Exec(execution)
				for _, r := range res {
                                	fmt.Fprintf(w, r)
				}
			case "showhealth":
				res := h.Showhealth()
				fmt.Fprintf(w, res)
			case "showbackend":
				res := h.ShowRegexp("BACKEND")
				fmt.Fprintf(w, res)
			case "showregexp":
				res := h.ShowRegexp(execution)
				fmt.Fprintf(w, res)
			case "enable":
				en := strings.Split(execution, "/")	
				res := h.EnableServer(en[0], en[1])
				fmt.Fprintf(w, res)
			case "enableall":
				h.EnableAll(execution)
			case "disable":
				dis := strings.Split(execution, "/")
				res := h.DisableServer(dis[0], dis[1])
				fmt.Fprintf(w, res)
			case "disableall":
				h.DisableAll(execution)
			default:
				fmt.Fprintf(w, usage)
			}
		}
	}
}

func main() {
	flag.Parse()
	haproxy := new(haProxy)
	if len(*f) > 0 {
		haproxy.Loadenv(*f)
	} else {
		haproxy.Loadenv("/etc/haproxy/haproxy.cfg")
	}

	if len(*binding) > 0 {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, haproxy)
		})
		http.ListenAndServe(":" + *binding, nil)
	}

	if *configcheck {
		haproxy.Configcheck()
	}

	if len(*action) > 0 {
		switch *action {
		case "socketexec":
                	res := haproxy.Exec(*execution)
			for _, r := range res {
                		fmt.Print(r)
			}
		case "showstatus":
                	res := haproxy.Showstatus()
                	fmt.Print(res)
		case "showbackend":
			res := haproxy.ShowRegexp("BACKEND")
                	fmt.Print(res)
		case "showhealth":
                	res := haproxy.Showhealth()
                	fmt.Print(res)
		case "showregexp":
			res := haproxy.ShowRegexp(*execution)
			fmt.Print(res)
		case "enable":
                	en := strings.Split(*execution, "/")
                	res := haproxy.EnableServer(en[0], en[1])
                	fmt.Printf(res)
		case "enableall":
			haproxy.EnableAll(*execution)
		case "disable":
                	dis := strings.Split(*execution, "/")
                	res := haproxy.DisableServer(dis[0], dis[1])
                	fmt.Printf(res)
		case "disableall":
			haproxy.DisableAll(*execution)
		default:
			flag.PrintDefaults()
		}
	}
}
