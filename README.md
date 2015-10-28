# stxy
Haproxy Stats to StatsD written in Golang

### Halp!
```
NAME:
   stxy - haproxy stats to statsd

USAGE:
   stxy [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   help, h	Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --haproxy-url "localhost:22002/;csv"	host:port of redis servier
   --statsd-url, -s "localhost:8125"	host:port of statsd server
   --prefix, -p "haproxy"		statsd namespace
   --interval, -i "5000"		time in milliseconds to periodically check redis
   --help, -h				show help
   --version, -v			print the version
```

### Stats Tracked:
```
admin.FRONTEND.scur:1|g
admin.FRONTEND.smax:1|g
admin.FRONTEND.ereq:0|g
admin.FRONTEND.econ:0|g
admin.FRONTEND.rate:1|g
admin.FRONTEND.bin:1881|g
admin.FRONTEND.bout:1666554|g
admin.FRONTEND.hrsp_1xx:0|g
admin.FRONTEND.hrsp_2xx:11|g
admin.FRONTEND.hrsp_3xx:0|g
admin.FRONTEND.hrsp_4xx:0|g
admin.FRONTEND.hrsp_5xx:0|g
admin.FRONTEND.qtime:0|g
admin.FRONTEND.ctime:0|g
admin.FRONTEND.rtime:0|g
admin.FRONTEND.ttime:0|g
```

### Adding more stats:
```
0 pxname
1 svname
2 qcur
3 qmax
4 scur
5 smax
6 slim
7 stot
8 bin
9 bout
10 dreq
11 dresp
12 ereq
13 econ
14 eresp
15 wretr
16 wredis
17 status
18 weight
19 act
20 bck
21 chkfail
22 chkdown
23 lastchg
24 downtime
25 qlimit
26 pid
27 iid
28 sid
29 throttle
30 lbtot
31 tracked
32 type
33 rate
34 rate_lim
35 rate_max
36 check_status
37 check_code
38 check_duration
39 hrsp_1xx
40 hrsp_2xx
41 hrsp_3xx
42 hrsp_4xx
43 hrsp_5xx
44 hrsp_other
45 hanafail
46 req_rate
47 req_rate_max
48 req_tot
49 cli_abrt
50 srv_abrt
51 comp_in
52 comp_out
53 comp_byp
54 comp_rsp
55 lastsess
56 last_chk
57 last_agt
58 qtime
59 ctime
60 rtime
61 ttime
```
