// speak.go
package main


import(
	"fmt"
	"os"
	"io/ioutil"
	"strings"
	"net/url"
)

// variables
var environ []string
var args string
var response_body string
var error_message string = "<h4>Oops, something went wrong</h4><p>Please send the exact URL and a description what you were doing to karlyan.kamerer (at) gmail.com"
var env_var map[string]string
var query_var map[string]string
var post_path = "/home/nurgling/speak/"

func buildPage() string {
	var response_body string
	content,err := ioutil.ReadFile(post_path + "posts.txt")
	if err != nil {
		response_body = error_message
		return response_body
	}
	response_body = strings.Replace(string(content), "\n", "<br>", -1)
	return response_body
}

func updatePage(speak string) {
        content,err := ioutil.ReadFile(post_path + "posts.txt")
        if err != nil {
                response_body += "<h2>go fuck yourself error 1</h2>"
        } else {
                err = ioutil.WriteFile(post_path + "posts.txt", []byte(speak + "\n" + string(content)), 0644)
		if err != nil {
			response_body += "<h2>go fuck yourself error 2</h2>"
		}
        }

}


func main() {
	env_var = make(map[string]string)
	query_var = make(map[string]string)
	environ = os.Environ()
	for _,n := range(environ) {
		split := strings.SplitN(n, "=", 2)
		env_var[split[0]] = split[1]
	}
	if env_var["QUERY_STRING"] != "" {
		split_query := strings.Split(env_var["QUERY_STRING"], "&")
		for _,n := range(split_query) {
			split := strings.SplitN(n, "=", 2)
			query_var[split[0]] = split[1]
		}
	}
	response_body = "<!doctype html>\n<html><meta charset=\"utf-8\">\n"
	response_body += "<header><title>SprachRohr Blog</title></header>\n<body>\n" 
	response_body += "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n"
	response_body += "</ul><br>\n"
	response_body += "<form action=\"" + env_var["REQUEST_URI"] + "\" method=\"POST\">"
	response_body += "<input type =\"text\" name=\"p\">"
	response_body += "<input type=\"submit\" value=\"Speak\">"
	response_body += "</form>"
	response_body += buildPage()
	if env_var["REQUEST_METHOD"] == "POST" {
		if len(os.Args) == 2 && os.Args[1] != "" {
			args,_ = url.QueryUnescape(os.Args[1])
			arg_var := make(map[string]string)
			split_args := strings.Split(args, "&")
			for _,n := range(split_args) {
				split := strings.SplitN(n, "=", 2)
				arg_var[split[0]] = split[1]
			}
			if arg_var["p"] != "" {
				clean_p := strings.Replace(arg_var["p"], "&", "&amp; ", -1)
				clean_p = strings.Replace(clean_p, ">", "&gt; ", -1)
				clean_p = strings.Replace(clean_p, "<", "&lt; ", -1)
				clean_p = strings.Replace(clean_p, "\"", "&quot; ", -1)
				clean_p = strings.Replace(clean_p, "'", "&apos; ", -1)
				if len(clean_p) > 1000 {
					clean_p = string([]rune(clean_p)[0:1000])
				}
				updatePage(clean_p)
			}
		}
	}
	response_body += "<p><a href=/faq.html>FAQ</a> For bugs, ideas, suggestion and other spam: karlyan.kamerer (at) gmail.com </p>\n"
	response_body += "</body>\n</html>\n"

	fmt.Printf("Content-Type: text/html; charset=utf-8\r\nContent-Length: %v\r\n\r\n" + response_body, len(response_body))
}
