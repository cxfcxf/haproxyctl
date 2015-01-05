package haproxyctl

import (
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

const (
	SOCKET_TYPE = "unix"
)

type HaProxy struct {
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

func (h *HaProxy) Loadenv(cfg string) {
	h.Cfg = cfg
	//Load sock and pid to HaProxy struct
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
	//Load binary path to HaProxy struct
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

func (h *HaProxy) Exec(command string) []string {
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

func (h *HaProxy) Showstatus() string {
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

func (h *HaProxy) Showhealth() string {
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

func (h *HaProxy) ShowRegexp(reg string) string {
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

func (h *HaProxy) DisableServer(px string, sv string) string {
	c := fmt.Sprintf("disable server %s/%s", px, sv)
	h.Exec(c)
	return fmt.Sprintf("Server %s/%s has been disabled\n", px, sv) + h.ShowRegexp(sv)
}

func (h *HaProxy) DisableAll(server string) {
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

func (h *HaProxy) EnableServer(px string, sv string) string {
	c := fmt.Sprintf("enable server %s/%s", px, sv)
	h.Exec(c)
	return fmt.Sprintf("Server %s/%s has been enabled\n", px, sv) + h.ShowRegexp(sv)
}

func (h *HaProxy) EnableAll(server string) {
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

func (h *HaProxy) Configcheck() {
	shell := fmt.Sprintf("%s -c -f %s", h.Bin, h.Cfg)
	cmd := exec.Command("sh", "-c", shell)
	res, _ := cmd.CombinedOutput()
	fmt.Println(strings.Trim(string(res), "\n"))
}
