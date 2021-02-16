package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/barandemirbas/open-with"
	"github.com/dietsche/rfsnotify"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
	Id      string `xml:"id,attr"`
	Style   string `xml:"style,attr"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	port := flag.Int("p", 8080, "Set the port to serve")
	flag.Parse()
	openwith.Browser("localhost", *port)
	log.Println(fmt.Sprintf("Server running on http://localhost:%d", *port))
	http.HandleFunc("/", Server)
	http.HandleFunc("/reload", Reload)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

// Server function is the web server
func Server(w http.ResponseWriter, r *http.Request) {
	h := html{}
	w.Header().Set("Cache-Control", "no-store, max-age=0")
	if strings.Contains(r.URL.Path, ".css") == false && strings.Contains(r.URL.Path, ".js") == false {
		file, err := ioutil.ReadFile("index.html")
		if r.URL.Path != "/" {
			file, err = ioutil.ReadFile(r.URL.Path[1:] + ".html")
		}
		err = xml.NewDecoder(bytes.NewBuffer(file)).Decode(&h)
		h.Body.Content += "<script>\nvar socket = new WebSocket(\"ws://localhost:8080/reload\");\nsocket.onopen = function () {\nconsole.log(\"Status: Connected.\");\n};\nsocket.onmessage = function (e) {\nlocation.reload();\n};\n</script>\n"
		fmt.Fprintf(w, "<html lang='"+h.Lang+"'><head>"+h.Head.Content+"</head>\n"+"<body class='"+h.Body.Class+"' id='"+h.Body.Id+"' style='"+h.Body.Style+"'>"+h.Body.Content+"</body>\n</html>")
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
	} else {
		http.ServeFile(w, r, r.URL.Path[1:])
	}
}

// Reload function is detect file changes and reload the page with websockets
func Reload(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	watcher, _ := rfsnotify.NewWatcher()
	err = watcher.AddRecursive(".")
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
