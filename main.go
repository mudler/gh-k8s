package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var start = flag.Bool("start", false, "Start VPN")
var generate = flag.Bool("generate", false, "Generate VPN token")

var ipCIDR = flag.String("ipCIDR", "", "IP in CIDR format")
var ip = flag.String("ip", "", "IP Address")
var luetVersion = flag.String("luetVersion", "0.20.10", "luetVersion")
var edgeVPNToken = flag.String("token", "", "EdgeVPN token")

var arch = "amd64"
var srcDir = "/home/src"

func main() {
	flag.Parse()

	//	utils.RunSH("dependencies", "apk add curl")
	//	utils.RunSH("dependencies", "apk add docker")
	//	utils.RunSH("dependencies", "apk add jq")
	RunSH("dependencies", "curl -L https://github.com/mudler/luet/releases/download/"+*luetVersion+"/luet-"+*luetVersion+"-linux-"+arch+" --output luet")
	RunSH("dependencies", "chmod +x luet")
	RunSH("dependencies", "mv luet /usr/bin/luet && mkdir -p /etc/luet/repos.conf.d/")
	RunSH("dependencies", "curl -L https://raw.githubusercontent.com/mocaccinoOS/repository-index/master/packages/mocaccino-extra.yml --output /etc/luet/repos.conf.d/mocaccino-extra.yml")
	RunSH("dependencies", "curl -L https://raw.githubusercontent.com/mocaccinoOS/repository-index/master/packages/luet.yml --output /etc/luet/repos.conf.d/luet.yml")
	RunSH("dependencies", "luet install -y system/luet utils/edgevpn container/k3s")

	switch {
	case *start:
		startVPN()
	case *generate:
		generateVPN()
	}
}

func generateVPN() {
	token, err := RunSHOUT("edgevpn", "edgevpn -g -b")
	checkErr(err)
	checkErr(ioutil.WriteFile("token", token, os.ModePerm))
}

func startVPN() {
	os.Setenv("EDGEVPNTOKEN", *edgeVPNToken)
	os.Setenv("ADDRESS", *ipCIDR)
	checkErr(RunSH("edgevpn", "edgevpn --api"))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
