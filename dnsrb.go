package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	//g='golang.org/x/sys' ; go get ${g}/unix
	//mkdir -p ~/go/src/${g} ; git clone https://github.com/golang/sys ~/go/src/${g}
	//g='golang.org/x/net' ; go get ${g}/ipv4
	//mkdir -p ~/go/src/${g} ; git clone https://github.com/golang/net ~/go/src/${g}
	//g='github.com/miekg/dns' ; go get ${g}
	//mkdir -p ~/go/src/${g} ; git clone https://github.com/miekg/dns ~/go/src/${g}
	"github.com/miekg/dns"
)

var dns_sec = new(int64)
var dns_rec = map[string][]string{}

func loadRecords(records map[string][]string) {
	fmt.Printf("Reading file [%d]...\n", *dns_sec)
	file, err := os.Open("dns.txt")
	if err == nil {
		for k, _ := range records {
			delete(records, k)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			list := strings.Fields(line)
			if len(list) > 1 {
				host := list[0]
				adrs := strings.Split(list[1], ",")
				records[host] = []string{}
				for x := 0; x < len(adrs); x += 1 {
					records[host] = append(records[host], adrs[x])
				}
			}
		}
	}
	defer file.Close()
}

func formatAnswer(host string, addr string) string {
	return fmt.Sprintf("%s A %s", host, addr)
}

func parseQuery(dns_msg *dns.Msg) {
	for i, query := range dns_msg.Question {
		if query.Qtype == dns.TypeA {
			fmt.Printf("[%d]A-Query: %s\n", i, query.Name)
			ipfs := []string{}
			for dns_key, dns_val := range dns_rec {
				matched, _ := regexp.Match(dns_key+"\\.?", []byte(query.Name))
				if matched {
					for x := 0; x < len(dns_val); x += 1 {
						ipfs = append(ipfs, dns_val[x])
					}
				}
			}
			if len(ipfs) < 1 {
				ans, _ := net.LookupIP(query.Name)
				for _, ip := range ans {
					addr := ip.String()
					if strings.Contains(addr, ".") {
						ipfs = append(ipfs, addr)
					}
				}
			}
			for j, ip := range ipfs {
				fmt.Printf("[%d]A-Reply: %s\n", j, ip)
				rr, _ := dns.NewRR(formatAnswer(query.Name, ip))
				dns_msg.Answer = append(dns_msg.Answer, rr)
			}
		}
	}
}

func handleRequest(dns_out dns.ResponseWriter, dns_req *dns.Msg) {
	var now_sec = time.Now().Unix()
	var new_sec int64 = 15
	if now_sec > (*dns_sec + new_sec) {
		*dns_sec = now_sec
		loadRecords(dns_rec)
	}

	dns_ans := new(dns.Msg)
	dns_ans.SetReply(dns_req)
	dns_ans.Compress = false

	if dns_req.Opcode == dns.OpcodeQuery {
		parseQuery(dns_ans)
	}
	if len(dns_ans.Answer) > 0 {
		dns_out.WriteMsg(dns_ans)
	}
}

func main() {
	// request handler
	dns.HandleFunc(".", handleRequest)

	// server options
	port := 53053
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	fmt.Printf("Starting: %d\n", port)

	// start server
	err := server.ListenAndServe()
	defer server.Shutdown()
	fmt.Printf("Errors: %s\n", err)
}
