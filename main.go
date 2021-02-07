package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/barandemirbas/open-with"
	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type html struct {
	Head head `xml:"head"`
	Body body `xml:"body"`
}

type head struct {
	Content string `xml:",innerxml"`
}

type body struct {
	Content string `xml:",innerxml"`
	Class   string `xml:"class,attr"`
	Id      string `xml:"id,attr"`
	Style   string `xml:"style,attr"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	openwith.Browser("localhost", 8080)
	log.Println(fmt.Sprintf("Server running on http://localhost%s", ":8080"))
	http.HandleFunc("/", Server)
	http.HandleFunc("/reload", Websocket)
	http.ListenAndServe(":8080", nil)
}

// Server function is the web server
func Server(w http.ResponseWriter, r *http.Request) {
	h := html{}
	file, err := ioutil.ReadFile("index.html")
	if r.URL.Path != "/" {
		file, err = ioutil.ReadFile(r.URL.Path[1:] + ".html")
		if strings.Contains(r.URL.Path, ".css") || strings.Contains(r.URL.Path, ".js") {
			http.ServeFile(w, r, r.URL.Path[1:])
		}
	}
	err = xml.NewDecoder(bytes.NewBuffer(file)).Decode(&h)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}
	h.Body.Content += "<script>\nvar socket = new WebSocket(\"ws://localhost:8080/reload\");\nsocket.onopen = function () {\nconsole.log(\"Status: Connected.\");\n};\nsocket.onmessage = function (e) {\nlocation.reload();\n};\n</script>\n"
	fmt.Fprintf(w, "<html><head>"+h.Head.Content+"</head>\n"+"<body class='"+h.Body.Class+"' id='"+h.Body.Id+"' style='"+h.Body.Style+"'>"+h.Body.Content+"</body></html>")
}

// Websocket function is reload the page with websockets
func Websocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	watcher, _ := fsnotify.NewWatcher()
	err = watcher.Add(".")
	if err != nil {
		return
	}
	for {
		select {
		case <-watcher.Events:
			err = conn.WriteMessage(1, []byte("reload"))
			if err != nil {
				return
			}
		}
	}
}
