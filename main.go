package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/go-martini/martini"
)

var port = flag.Int("p", 9527, "port to listen")
var secret = flag.String("s", "", "github secret")
var script = flag.String("e", "", "script to execute")

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	sig := r.Header.Get("X-Hub-Signature")
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
	sig1 := ComputeHmac256(string(body), *secret)
	fmt.Println(sig, sig1)
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
