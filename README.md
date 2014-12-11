#This is a partial rewrite of https://github.com/flores/haproxyctl in go

###the programing is still under construction. what you see is not compeleted yet.

```
i am not a expert of haproxy, but i find its actually annonying to gather data or enable/disable servers
through stat socket haproxy provides. so this is a tool or wrapper to get around with that. i am thinking of
binding the socket to a port with net/http to provide remote accessesbility in the future.
```

thanks for flores' great work on the original haproxyctl
i am working towards to compelete this, so we will have compiled binary in the future.

```
go build haproxyctl
```

example
```
[root@Centos ~]# go run haproxyctl.go -disable app/app1
Server app/app1 has been disabled
now printing app1 Health Check
=====================================
# pxname   svname     status     weight
app        app1       DOWN       1
=====================================


[root@Centos ~]# go run haproxyctl.go -enable app/app1
Server app/app1 has been enabled
now printing app1 Health Check
=====================================
# pxname   svname     status     weight
app        app1       DOWN       1
=====================================


[root@Centos ~]# go run haproxyctl.go -stats
now printing Health Check
=====================================
# pxname   svname     status     weight
main       FRONTEND   OPEN
static     static     DOWN       1
static     BACKEND    DOWN       0
app        app1       DOWN       1
app        app2       DOWN       1
app        app3       DOWN       1
app        app4       DOWN       1
app        BACKEND    DOWN       0
=====================================

```

i dont actually have servers i can test now, so its not maintain mode instead of down, it will never be up lol

#License
This code is released under the MIT License. You should feel free to do whatever you want with it. 
