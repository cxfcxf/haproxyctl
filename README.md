#Haproxyctl in golang
#####This is a partial rewrite of https://github.com/flores/haproxyctl in go
Thanks for cflores' great work on the original haproxyctl
Thanks for cflores' advise of building on haproxyctl

i added more functions to the haproxyctl

####Library is under github.com/cxfcxf/haproxyctl/lib

####Things need to be improved
1. it is not robust enough, i need people to test it in different situation, so i can improve it by problem people run into
2. ARGV prob need to match the original one or not, i am not sure
3. Need more user cases for haproxy to improve haproxyctl.

####feature
support of multiple cores
support binding haproxyctl to a port

###how to compile haproxyctl.go
```
go get github.com/cxfcxf/haproxyctl/lib
git clone https://github.com/cxfcxf/haproxyctl.git
cd haproxyctl
go build haproxyctl
```

###example for running as a command
you can either point the haproxy.cfg file with -f or 
program will detect it at /etc/haproxy/haproxy.cfg  by default

####show running status
```
[root@haproxy haproxyctl]# ./haproxyctl -action="showstatus"
haproxy is running on pid 1434.
these ports are used and guys are connected:
TCP *:commplex-main
UDP *:48266
haproxy is running on pid 1435.
these ports are used and guys are connected:
TCP *:commplex-main
UDP *:48266
```

####disable server app/app2
```
[root@haproxy haproxyctl]# ./haproxyctl -action="disable" -execution="app/app2"
Server app/app2 has been disabled

now printing app2 Health Check
# pxname   svname     status     weight
app        app2       MAINT      1

now printing app2 Health Check
# pxname   svname     status     weight
app        app2       MAINT      1
```

####showhealth
```
[root@haproxy haproxyctl]# ./haproxyctl -action="showhealth"

now printing Health Check...

# pxname   svname     status     weight
main       FRONTEND   OPEN
static     static     DOWN       1
static     BACKEND    DOWN       0
app        app1       DOWN       1
app        app2       MAINT      1
app        app3       DOWN       1
app        app4       DOWN       1
app        BACKEND    DOWN       0

now printing Health Check...

# pxname   svname     status     weight
main       FRONTEND   OPEN
static     static     DOWN       1
static     BACKEND    DOWN       0
app        app1       DOWN       1
app        app2       MAINT      1
app        app3       DOWN       1
app        app4       DOWN       1
app        BACKEND    DOWN       0
```

####show backend
```
[root@haproxy haproxyctl]# ./haproxyctl -action="showbackend"

now printing BACKEND Health Check
# pxname   svname     status     weight
static     BACKEND    DOWN       0
app        BACKEND    DOWN       0

now printing BACKEND Health Check
# pxname   svname     status     weight
static     BACKEND    DOWN       0
app        BACKEND    DOWN       0
```

####disable/enable a server in all backend
```
./haproxyctl -action="disableall" -execution="app2"
./haproxyctl -action="enableall" -execution="app1"
```

####exectution of any socket command (directly execution socket command)
```
[root@haproxy haproxyctl]# ./haproxyctl -action="socketexec" -execution="get weight app/app2"
1 (initial 1)

1 (initial 1)
```

###Example for binding to web port for RESTful API
```
[root@haproxy haproxyctl]# ./haproxyctl --binding="8888"
```
now you can do something like
```
####showstatus
curl http://youripaddress:8888/haproxyctl?action=showstatus
####socketexec
curl http://youripaddress:8888/haproxyctl?action=socketexec&execution=get%20weight%20app/app1
####disable server app/app2
curl http://youripaddress:8888/haproxyctl?action=disable&exectution=app/app2
```
web haproxyctl bascilly support all kind of command through your request
so you can disable server remotely, it shares the same return asd command line


#License
This code is released under the MIT License. You should feel free to do whatever you want with it. 

#contact
you can contact me for any suggestion or request @ siegfried.chen@gmail.com
