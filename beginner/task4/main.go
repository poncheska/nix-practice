package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	num := 5
	for i := 1; i < num + 1; i++{
		wg.Add(1)
		go func(ii int) {
			resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/" + strconv.Itoa(ii))
			if err != nil {
				log.Fatalln(err)
			}
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println(string(bs))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
