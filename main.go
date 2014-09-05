package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/go-martini/martini"
)

var port = flag.Int("p", 9527, "port to listen")
var secret = flag.String("s", "", "github secret")
var script = flag.String("e", "", "script to execute")

func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.Form)
	out, err := exec.Command(*script).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}

func main() {
	flag.Parse()
	if *script == "" || *secret == "" {
		flag.PrintDefaults()
		return
	}
	m := martini.Classic()
	m.Post("/", index)
	addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(addr, m))
}
