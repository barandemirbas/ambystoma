package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/barandemirbas/open-with"
	"github.com/dietsche/rfsnotify"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type html struct {
	Lang string `xml:"lang,attr"`
	Head head   `xml:"head"`
	Body body   `xml:"body"`
}

type head struct {
	Content string `xml:",innerxml"`
}

type body struct {
	Content string `xml:",innerxml"`
	Class   string `xml:"class,attr"`
	ID      string `xml:"id,attr"`
	Style   string `xml:"style,attr"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// colors
var green = color.New(color.FgGreen).Add(color.Bold)
var yellow = color.New(color.FgYellow).Add(color.Bold)
var red = color.New(color.FgRed).Add(color.Bold)

func main() {
	port := flag.Int("p", 8080, "Set the port to serve")
	flag.Parse()
	openwith.Browser("localhost", *port)
	green.Println("[+]", fmt.Sprintf("Server running on http://localhost:%d", *port))
	http.HandleFunc("/", Server)
	http.HandleFunc("/reload", Reload)
	red.Println("[!]", http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

// Server function is the web server
func Server(w http.ResponseWriter, r *http.Request) {
	h := html{}
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	if !strings.Contains(r.URL.Path, ".css") && !strings.Contains(r.URL.Path, ".js") {
		file, _ := ioutil.ReadFile("index.html")
		if r.URL.Path != "/" {
			file, _ = ioutil.ReadFile(r.URL.Path[1:] + ".html")
		}
		err := xml.NewDecoder(bytes.NewBuffer(file)).Decode(&h)
		h.Body.Content += "<script>\nvar socket = new WebSocket(\"ws://localhost:8080/reload\");\nsocket.onopen = function () {\nconsole.log(\"Status: Connected.\");\n};\nsocket.onmessage = function (e) {\nlocation.reload();\n};\n</script>\n"
		fmt.Fprintf(w, "<html lang='"+h.Lang+"'><head>"+h.Head.Content+"</head>\n"+"<body class='"+h.Body.Class+"' id='"+h.Body.ID+"' style='"+h.Body.Style+"'>"+h.Body.Content+"</body>\n</html>")
		if err != nil && err != io.EOF {
			yellow.Println("[-]", err)
			return
		}
	} else {
		http.ServeFile(w, r, r.URL.Path[1:])
	}
}

// Reload function is detect file changes and reload the page with websockets
func Reload(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	watcher, err := rfsnotify.NewWatcher()
	if err != nil {
		red.Println("[!]", err)
		return
	}
	err = watcher.AddRecursive(".")
	if err != nil {
		red.Println("[!]", err)
		return
	}
	for {
		<-watcher.Events
		_ = conn.WriteMessage(1, []byte("reload"))
	}
}
