package main

import (
	"crypto/hmac"
	"crypto/sha1"
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

func ComputeHmac(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	h.Write([]byte(message))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func index(w http.ResponseWriter, r *http.Request) string {
	sig := r.Header.Get("X-Hub-Signature")
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	sig1 := fmt.Sprintf("sha1=%s", ComputeHmac(string(body), *secret))
	if sig != sig1 {
		msg := fmt.Sprintln("signature not match", sig, sig1)
		fmt.Println(msg)
		return msg
	}
	out, err := exec.Command(*script).Output()
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	fmt.Println(out)
	return string(out)
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
