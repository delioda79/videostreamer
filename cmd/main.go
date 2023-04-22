package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/video/{url}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		url := vars["url"]
		decoded, err := base64.StdEncoding.DecodeString(url)
		if err != nil {
			fmt.Println("decode error:", err)
			return
		}

		rsp, err := http.Get(string(decoded))
		if err != nil {
			log.Fatal(fmt.Errorf("An error: %v", err))
		}

		defer rsp.Body.Close()
		fmt.Printf("%d %s %v\n", rsp.StatusCode, rsp.Status, rsp.Header)

		//w.Header().Set("Content-Type", "video/mp2t")
		//w.Header().Set("Content-Type", "application/x-mpegURL")
		w.Header().Set("Content-Type", rsp.Header.Get("Content-Type"))
		//w.Write([]byte("HELLO"))

		reader := bufio.NewReader(rsp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF{
				log.Fatal(fmt.Errorf("This happened %v",err))
			}

			//log.Println("ERROR: ", line)
			w.Write([]byte(string(line)))
		}
	})


	log.Fatal(http.ListenAndServe(":8080", r))
}

