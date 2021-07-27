<h1 align="center">GOHTTP</h1>


---
Vulnerabilities that this tool helps to check:

- Broken Link Hijacking
- Subdomain Takeover

```
----GOHTTP 1.1----

COMMAND  DESCRIPTION               REQUIRED
-------  -----------               --------
-u       Host URL.                 Yes
-l       List of hosts.            Yes
-title   Grab title banner.        No
-server  Grab server banner.       No
-ip      Get host IP.              No
-cnm     Get host Canonical Name.  No
-cl      Content length.           No
-fh      Specified header.         No
-c       Crawler.                  No
-H       Set request headers.      No
-h       Help Menu.

```

Output Example:
```
╰─$ go run ghttp.go -u http://testphp.vulnweb.com/ -c -cnm -server -title -cl -H "User-Agent:Mozilla"                                                                                       
6 Links Found.
1 [200]https://www.acunetix.com/ [Acunetix | Web Application Security Scanner][acunetix.com][www.acunetix.com.][111857b]
2 [200]https://www.acunetix.com/vulnerability-scanner/ [Vulnerability Scanner - Web Application Security | Acunetix][acunetix.com][www.acunetix.com.][146129b]
3 [200]http://www.acunetix.com [Acunetix | Web Application Security Scanner][acunetix.com][www.acunetix.com.][111857b]
4 [200]https://www.acunetix.com/vulnerability-scanner/php-security-scanner/ [PHP Security Scanner | Acunetix][acunetix.com][www.acunetix.com.][144610b]
5 [200]https://www.acunetix.com/blog/articles/prevent-sql-injection-vulnerabilities-in-php-applications/ [Prevent SQL injection vulnerabilities in PHP applications and fix them][acunetix.com][www.acunetix.com.][174700b]
6 [403]http://www.eclectasy.com/Fractal-Explorer/index.html [403 Forbidden][Apache/2][www.eclectasy.com.][236b]

```
