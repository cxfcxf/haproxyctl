#This is a partial rewrite of https://github.com/flores/haproxyctl in go

the programing is still under construction. what you see is not completed yet.
#current status
```
1. need to finish flag
2. more function will be added
3. seprate the library for exec for library usage
4. bindding it to a port for remote access and execution
```

```
i am not a expert of haproxy, but i find its actually annonying to gather data or enable/disable servers
through stat socket haproxy provides. so this is a tool or wrapper to get around with that. 
```
###further
```
i am thinking ofbinding the socket to a port with net/http to provide remote accessesbility in the future.
```

thanks for cflores' great work on the original haproxyctl
i am working towards to complete this when i have time, so we will have compiled binary in the future.

###compile haproxyctl.go
```
go build haproxyctl
```

###example
you can either point the haproxy.cfg file with -f or program will detect it by default

show running status, 
```
[root@haproxy haproxyctl]# go run haproxyctl.go -status
haproxy is running on pid 4554.
these ports are used and guys are connected:
TCP *:commplex-main
UDP *:56381
```
disable app/app1
```
[root@haproxy haproxyctl]# go run haproxyctl.go -disable app/app1
Server app/app1 has been disabled

now printing app1 Health Check
# pxname   svname     status     weight
app        app1       DOWN       1
```

showhealth
```
[root@haproxy haproxyctl]# go run haproxyctl.go -showhealth

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
[root@haproxy haproxyctl]# go run haproxyctl.go -showbackend

now printing BACKEND Health Check
# pxname   svname     status     weight
static     BACKEND    DOWN       0
app        BACKEND    DOWN       0
```

#License
This code is released under the MIT License. You should feel free to do whatever you want with it. 

#contact
you can contact me for any suggestion or request @ siegfried.chen@gmail.com
