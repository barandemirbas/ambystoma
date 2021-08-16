package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	openwith "github.com/barandemirbas/open-with"
	"github.com/dietsche/rfsnotify"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//go:embed injected.html
var js string

func main() {
	port := flag.Int("p", 8080, "Set the port to serve")
	flag.Parse()
	SetPort(*port)

	openwith.Browser("localhost", *port)

	fmt.Println("[+]", fmt.Sprintf("Server running on http://localhost:%d", *port))

	http.HandleFunc("/", Server)
	http.HandleFunc("/reload", Reload)
	fmt.Println("[!]", http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func Inject(file []byte) string {
	m := regexp.MustCompile("</head>")
	res := m.ReplaceAllString(string(file), js)
	return res
}

func SetPort(port int) {
	m := regexp.MustCompile("var port = (.*);")
	res := m.ReplaceAllString(js, fmt.Sprintf("var port = %d;", port))
	js = res
}

// Server function is the web server
func Server(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, max-age=0")

	file, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Fprintf(w, "<html><h3><strong>index.html</strong> not found</h3></html>")
	}

	if strings.Contains(r.URL.Path, ".html") || strings.Contains(r.URL.Path, ".htm") || r.URL.Path == "/" {

		if r.URL.Path != "/" {
			file, err = ioutil.ReadFile(r.URL.Path[1:])
			if err != nil {
				fmt.Fprintf(w, "<html><h3><strong>%s</strong> not found</h3></html>", r.URL.Path[1:])
			}
		}

		fmt.Fprintf(w, Inject(file))
	} else {
		http.ServeFile(w, r, r.URL.Path[1:])
	}
}

// Reload function is detect file changes and reload the page with websockets
func Reload(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	watcher, err := rfsnotify.NewWatcher()
	if err != nil {
		fmt.Println("[!]", err)
		return
	}
	err = watcher.AddRecursive(".")
	if err != nil {
		fmt.Println("[!]", err)
		return
	}
	for {
		<-watcher.Events
		_ = conn.WriteMessage(1, []byte("reload"))
	}
}
