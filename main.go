package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

//go:embed index.html
var html string

var (
	df = "data.txt"
	cf = "cache.txt"
)

func main() {

	dataf, err := os.Open(df)

	if err != nil {
		fmt.Println("Not found data.txt")
		fmt.Println("Automatic shutdown after 10 seconds")
		time.Sleep(10 * time.Second)
		dataf.Close()
		return
	}

	check()
	txt, _ := ioutil.ReadAll(dataf)
	ck := strings.Split(string(txt), "\n")
	dataf.Close()

	var data []string

	for _, val := range ck {

		s := strings.TrimSpace(val)

		if s != "" && !find(data, s) {
			data = append(data, s)
		}
	}

	http.HandleFunc("/api/get", func(w http.ResponseWriter, r *http.Request) {

		check()

		file, err := os.Open(cf)
		if err != nil {
			w.Write([]byte("打开" + cf + "失败"))
			return
		}
		defer file.Close()

		old, _ := ioutil.ReadAll(file)

		var cache []string

		tidy := strings.Split(string(old), "\n")

		for _, val := range tidy {
			s := strings.TrimSpace(val)
			if s != "" {
				cache = append(cache, s)
			}

		}

		if len(cache) >= len(data) {
			w.WriteHeader(http.StatusNoContent)
			w.Write([]byte(""))
		} else {

			var word string

			for {
				word = data[randNum(len(data))]

				if !find(cache, word) {
					break
				}

			}

			cache = append(cache, word)

			os.WriteFile(cf, []byte(strings.Join(cache, "\n")), os.ModePerm)

			w.Write([]byte(word))

			return
		}

	})

	http.HandleFunc("/api/delete", func(w http.ResponseWriter, r *http.Request) {
		check()
		os.WriteFile(cf, []byte(""), os.ModePerm)
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// aa, _ := os.Open("index.html")
		// body, _ := ioutil.ReadAll(aa)
		// w.Write([]byte(body))

		w.Write([]byte(html))
	})

	fmt.Println("")
	fmt.Println("http://localhost:6788")
	http.ListenAndServe(":6788", nil)

}

func randNum(max int) int {
	var timeStamp = time.Now().Unix()
	r := rand.New(rand.NewSource(timeStamp))
	num := r.Intn(max)
	return num
}

func find(array []string, value string) bool {
	for _, txt := range array {

		if value == txt {
			return true
		}

	}
	return false
}

func check() {
	_, err := os.Stat(cf)

	if os.IsNotExist(err) || err != nil {
		f, _ := os.Create(cf)
		f.Close()
	}
}
