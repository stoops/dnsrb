# dnsrb
```
$ go run dnsrb.go

Starting: 53053
Reading file [1616610565]...
[0]A-Query: facebook.com.
[0]A-Reply: 127.0.0.1
[0]A-Query: www.facebook.com.
[0]A-Reply: 127.0.0.1
[0]A-Query: blah.facebook.com.
[0]A-Reply: 127.0.0.1
[0]A-Query: fb.me.
[0]A-Reply: 127.0.0.1
[0]A-Query: amazon.ca.
[0]A-Reply: 54.239.18.172
[1]A-Reply: 54.239.19.238
[2]A-Reply: 52.94.225.242

$ echo ; for d in facebook.com www.facebook.com blah.facebook.com fb.me amazon.ca ; do echo "[$d] -> $(dig @127.0.0.1 -p 53053 $d +short)" ; done ; echo

[facebook.com] -> 127.0.0.1
[www.facebook.com] -> 127.0.0.1
[blah.facebook.com] -> 127.0.0.1
[fb.me] -> 127.0.0.1
[amazon.ca] -> 54.239.18.172 54.239.19.238 52.94.225.242
```
