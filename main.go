package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-martini/martini"
)

var port = flag.Int("p", 9527, "port to listen")
var secret = flag.String("s", "", "github secret")
var script = flag.String("e", "", "script to execute")
var test = flag.Bool("t", false, "run script then exit")

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
		log.Println(msg)
		return msg
	}
	out, err := exec.Command("sh", *script).Output()
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	log.Println(string(out))
	return string(out)
}

func main() {
	flag.Parse()
	if *script == "" {
		flag.PrintDefaults()
		return
	}
	if *test == true {
		out, err := exec.Command("sh", *script).Output()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(out))
		return
	}
	if *secret == "" {
		flag.PrintDefaults()
		return
	}
	m := martini.Classic()
	m.Map(log.New(os.Stdout, "", log.LstdFlags))
	m.Post("/", index)
	addr := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(addr, m))
}
