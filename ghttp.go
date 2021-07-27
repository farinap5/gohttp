package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/cheynewallace/tabby"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Global Functions
var n int = 1
var setHeader string
var version string = "1.1"
// ----------------

type OBJ struct {
	gTitle   bool
	gServer  bool
	gLength  bool
	gIP      bool
	findHead string
	gCNAME   bool
	//--//--//--//
	url    string
	domain string
	statuc string
	title  string
	server string
	length string
	ip     string
	headCH string
	header string
	cname  string
}

func (o OBJ)outputList() {
	fmt.Println(n,o.statuc+o.url+" "+o.title+o.server+o.cname+o.length+o.header)
	n++
}

func (o OBJ)output() {
	println("Host:",o.url)
	println("Status Code:",o.statuc)

	if o.gTitle {
		println("Title:",o.title)
	}

	if o.gServer {
		println("Server:",o.server)
	}

	if o.gCNAME {
		println("CNAME:",o.cname)
	}

	if o.gLength {
		println("Content Length:",o.length,)
	}

	if len(o.findHead) >= 1 {
		println("Header:",o.headCH,":",o.header)
	}
}


func (fr OBJ)get() *http.Response {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fr.url, nil)
	req.Header.Set("User-Agent","GHTTP 1.1")
	if setHeader != "false" {
		h1 := strings.Replace(setHeader," ","",-1)
		h2 := strings.Split(h1,",")
		for i := range(h2) {
			z := strings.Split(h2[i],":")
			req.Header.Set(z[0],z[1])
		}
	}
	res,err := client.Do(req)
	if err != nil {
		println(err.Error())
	}
	return res
}

// This function will grab the title from web page.
func getTitle(r *http.Response) string {
	body,err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	reTitle := regexp.MustCompile(`(?im)<\s*title.*>(.*?)<\s*/\s*title>`)
	t := reTitle.FindString(string(body))
	if len(t) > 3 {
		tS := strings.Split(t,">")
		tSS := strings.Split(tS[1],"</")
		return "[\033[1;36m"+tSS[0]+"\033[0;0m]"
	} else {
		return "[\033[1;31mNo Title Found!\033[0;0m]"
	}
}

// This function will grab the banner from the server.
func getServer(r *http.Response) string {
	head := r.Header["Server"]
	server := strings.Join(head, " ")
	return "[\033[1;33m"+server+"\033[0;0m]"
}

// This function count the length.
func getLength(r *http.Response) string {
	body,err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	cL := len(body)
	strCL := strconv.Itoa(cL)
	return "[\033[1;31m"+strCL+"\033[0;0mb]"
}

// This function will grab the ip.
func getIP(domain string) {
	add,err := net.LookupIP(domain)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	for _, ip := range add {
		fmt.Printf("-> "+domain+". IN A %s\n", ip.String())
	}
}

func getCNAME(domain string) string {
	add,err := net.LookupCNAME(domain)
	if err != nil {
		println(err.Error())
		//os.Exit(1)
	}
	return "[\033[1;35m"+add+"\033[0;0m]"
}

// This function will get the response header specified
// by the user.
func specificHeader(h string,r *http.Response) string {
	head := r.Header[h]
	newH := strings.Join(head, " ")
	if len(newH) >= 1 {
		return "["+newH+"]"
	} else {
		return "[No Data!]"
	}

}

// This function returns status code
func stsCode(r *http.Response) string {
	var resp string
	sts := r.StatusCode
	strStatus := strconv.Itoa(sts)
	if strStatus == "200" {
		resp = "[\033[1;32m"+strStatus+"\033[0;0m]"
	} else {
		resp = "[\033[1;31m"+strStatus+"\033[0;0m]"
	}

	return resp
}

// Request unique host.
func req(u* string,gT* bool,gS* bool,nh* string,gIp* bool,cl* bool,cnm* bool) {
	b := new(OBJ)
	b.url = *u
	b.domain = strings.Split(b.url, "/")[2]
	b.gTitle = *gT
	b.gServer = *gS
	b.findHead = *nh
	b.gIP = *gIp
	b.gLength = *cl
	b.gCNAME = *cnm
	b.statuc = stsCode(b.get())
	if *gT {
		b.title = getTitle(b.get())
	}
	if *gS {
		b.server = getServer(b.get())
	}
	if len(*nh) > 1 {
		b.headCH = *nh
		b.header = specificHeader(*nh,b.get())
	}
	if *cl {
		b.length = getLength(b.get())
	}
	if *cnm {
		b.cname = getCNAME(b.domain)
	}

	b.output()
	if *gIp {
		getIP(b.domain)
	}
}

// Request list of hosts
func reqList(u string,gT* bool,gS* bool,nh* string,gIp* bool,cl* bool,cnm* bool) {
	b := new(OBJ)
	b.url = u
	b.domain = strings.Split(b.url, "/")[2]
	b.statuc = stsCode(b.get())
	if *gT {
		b.title = getTitle(b.get())
	}
	if *gS {
		b.server = getServer(b.get())
	}
	if len(*nh) > 1 {
		b.header = specificHeader(*nh,b.get())
	}
	if *cl {
		b.length = getLength(b.get())
	}

	if *cnm {
		b.cname = getCNAME(b.domain)
	}

	b.outputList()
	if *gIp {
		getIP(b.domain)
	}
}



// This function is the crawler,
// all returned URL will be requested
// to test
func crawl(u string) []string {
	var list []string
	client := &http.Client{}
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent","GHTTP 1.1")
	if setHeader != "false" {
		h1 := strings.Replace(setHeader," ","",-1)
		h2 := strings.Split(h1,",")
		for i := range(h2) {
			z := strings.Split(h2[i],":")
			req.Header.Set(z[0],z[1])
		}
	}
	r,err := client.Do(req)
	if err != nil {
		println(err.Error())
	}

	body,err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	links := regexp.MustCompile(`href="https?.*?"`)
	match := links.FindAllString(string(body),-1)
	for i := range match {
		insideM := 0
		m := strings.Split(match[i],"\"")[1]
		for j := range list {
			if list[j] == m{
				insideM++
			}
		}
		if insideM == 0 {
			list = append(list, m)
		}
	}

	srcs := regexp.MustCompile(`src=".*?"`)
	srcsm := srcs.FindAllString(string(body),-1)
	for i := range srcsm {
		insideM := 0
		m := strings.Split(match[i],"\"")[1]
		for j := range list {
			if list[j] == m{
				insideM++
			}
		}
		if insideM == 0 {
			list = append(list, m)
		}
	}
	fmt.Println(len(list),"Links Found.")
	return list
}

// This function is the base of threads.
// Its name doest mean it will just request when using crawl
// function, is used for list too.
func reqByCrawl(v chan int,list []string,gT* bool,gS* bool,nH* string,gIp* bool,cl* bool,cnm* bool) {
	var nrq int = 0
	for i := range v {
		reqList(list[nrq],gT,gS,nH,gIp,cl,cnm)
		println(i)
		nrq++
	}
}

func main() {
	var list []string
	// Unique host url.
	var url  = flag.String("u","false","URL that will be tested.")
	// Title
	var gT   = flag.Bool("title",false,"Grab title from web page.")
	// Server
	var gS   = flag.Bool("server",false,"Grab banner from server used by the web application")
	var nH   = flag.String("fh","","Grab specified header.")
	// Get host ip.
	var gIp  = flag.Bool("ip",false,"Show ip.")
	var cnm  = flag.Bool("cnm",false,"Shou Canonical Name.")
	var cL   = flag.Bool("cl",false,"Content Length.")
	// List containing hosts to test
	var hL   = flag.String("l","false","List containing URLs.")
	var crwl = flag.Bool("c",false,"Crawler.")
	//var thr  = flag.Int("t",5,"Threads.")
	var h    = flag.Bool("h",false,"Help Menu.")
	// Set request headers: "User-Agent: ghttp 1.1"
	var sH   = flag.String("H","false","Set Headers.")
	flag.Parse()

	// Pass Header to a global variable.
	setHeader = *sH

	if *h {
		_help()
		os.Exit(1)
	}
	if *url == "false" && *hL == "false" {
		println("No URL or list of hosts!")
		println("type argument -h")
		os.Exit(1)
	}

	// If a file is requested.
	if *hL != "false" {
		// Open file;
		file,err := os.Open(*hL)
		if err != nil {
			println("[\u001B[1;31m!\u001B[0;0m]- Error while trying to open file!")
			os.Exit(1)
		}
		println("[\u001B[1;32mOK\u001B[0;0m]- File:",*hL)
		scanner := bufio.NewScanner(file)
		// Split content by lines.
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			list = append(list, scanner.Text())
		}
		file.Close()
		// Make a request for each url in the file.
		for i := range list{
			reqList(list[i],gT,gS,nH,gIp,cL,cnm)
		}
		/*//------------- threads is not working well.
		requestId := make(chan int)
		conn := *thr
		for i := 0; i <= conn; i++ {
			go reqByCrawl(requestId,list,gT,gS,nH,gIp,cL,cnm)
		}
		for i := 0; i < len(list); i++ {
			requestId <- i
		}
		//-------------*/
	} else {
		if *crwl {
			LinkSrc := crawl(*url)

			for i := range LinkSrc{
				reqList(LinkSrc[i],gT,gS,nH,gIp,cL,cnm)
			}
			/*//------------- threads is not working well.
			requestId := make(chan int)
			conn := *thr

			for i := 1; i <= conn; i++ {
				go reqByCrawl(requestId,LinkSrc,gT,gS,nH,gIp,cL,cnm)
			}
			for i := 0; i < len(LinkSrc); i++ {
				requestId <- i
			}
			//-------------*/
		} else {
			req(url,gT,gS,nH,gIp,cL,cnm)
		}

	}

}

func _help() {
	fmt.Println("----\033[1;36mGOHTTP \u001B[1;35m"+version+"\033[0;0m----\n")
	t := tabby.New()
	t.AddHeader("COMMAND","DESCRIPTION","REQUIRED")
	t.AddLine("-u","Host URL.","Yes")
	t.AddLine("-l","List of hosts.","Yes")
	t.AddLine("-title","Grab title banner.","No")
	t.AddLine("-server","Grab server banner.","No")
	t.AddLine("-ip","Get host IP.","No")
	t.AddLine("-cnm","Get host Canonical Name.","No")
	t.AddLine("-cl","Content length.","No")
	t.AddLine("-fh","Specified header.","No")
	t.AddLine("-c","Crawler.","No")
	//t.AddLine("-t","Threads - Default: 5","No")
	t.AddLine("-H","Set request headers.","No")
	t.AddLine("-h","Help Menu.")
	t.Print()
}
