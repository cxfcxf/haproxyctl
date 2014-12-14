#Haproxyctl in golang
#####This is a partial rewrite of https://github.com/flores/haproxyctl in go
thanks for cflores' great work on the original haproxyctl
i am working towards to complete this when i have time, so we will have compiled binary in the future.

###update
1. completed binding to port function
2. work with multiple process module

#current status
although most of original haproxyctl function are usable, 
the programing is still under actively developing. what you see is not completed yet.

#future plan
```
1. need to finish flag and usage setting
2. more functions will be added
3. seprate the library and command execution
```

###compile haproxyctl.go
```
go build haproxyctl
```

###example
you can either point the haproxy.cfg file with -f or program will detect it by default

show running status, 
```
[root@haproxy haproxyctl]# go run haproxyctl.go -action="showstatus"
haproxy is running on pid 1662.
these ports are used and guys are connected:
TCP *:commplex-main
UDP *:41747
```
disable app/app2
```
[root@haproxy haproxyctl]# go run haproxyctl.go -action="disable" -execution="app/app2"
Server app/app2 has been disabled

now printing app2 Health Check
# pxname   svname     status     weight
app        app2       DOWN       1
```

showhealth
```
[root@haproxy haproxyctl]# go run haproxyctl.go -action="showhealth"

now printing Health Check...

# pxname   svname     status     weight
main       FRONTEND   OPEN
static     static     DOWN       1
static     BACKEND    DOWN       0
app        app1       DOWN       1
app        app2       DOWN       1
app        app3       DOWN       1
app        app4       DOWN       1
app        BACKEND    DOWN       0
```
show backend
```
[root@haproxy haproxyctl]# go run haproxyctl.go -action="showbackend"

now printing BACKEND Health Check
# pxname   svname     status     weight
static     BACKEND    DOWN       0
app        BACKEND    DOWN       0
```

disable/enable a server in all backend
```
go run haproxyctl.go -action="disableall" -execution="app2"
go run haproxyctl.go -action="enableall" -execution="app1"
```

exectution socket command (directly execution socket command)
```
[root@haproxy haproxyctl]# go run haproxyctl.go -action="socketexec" -execution="get weight app/app2"
1 (initial 1)
```


#License
This code is released under the MIT License. You should feel free to do whatever you want with it. 

#contact
you can contact me for any suggestion or request @ siegfried.chen@gmail.com
