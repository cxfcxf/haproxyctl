package main

import (
	"flag"
	"fmt"
	"strings"
	"net/http"
	haproxyctl "github.com/cxfcxf/haproxyctl/lib"
)

var (
	configcheck  = flag.Bool("configcheck", false, "check config file")
	action	     = flag.String("action", "", "tell me what action you wanna do?")
	execution    = flag.String("execution", "", "parameters to action")
	binding      = flag.String("binding", "", "http port you want program to bind to")
	client	     = flag.String("client", "", "only one client ip are allowed")
	f            = flag.String("f", "/etc/haproxy/haproxy.cfg", "point configration file, default /etc/haproxy/haproxy.cfg")
)

func handler(w http.ResponseWriter, r *http.Request, h *haproxyctl.HaProxy, c string) {
	usage := "please use /haproxyctl?action=xxxx&execution=yyyy"
	if strings.Split(r.RemoteAddr, ":")[0] != c {
		fmt.Fprintf(w, "sorry, you are not in the whitelist in order to exec api")
	} else if r.URL.Path != "/haproxyctl" {
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
	haproxy := new(haproxyctl.HaProxy)

	if len(*f) > 0 {
		haproxy.Loadenv(*f)
	} else {
		haproxy.Loadenv("/etc/haproxy/haproxy.cfg")
	}

	if len(*binding) > 0 && len(*client) > 0 {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			handler(w, r, haproxy, *client)
		})
		http.ListenAndServe(*binding, nil)
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
