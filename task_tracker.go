package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func checkSite(link string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(link)
	if err != nil {
		fmt.Printf("%s - Error: %v\n", link, err)
		return
	}

	defer resp.Body.Close()
	fmt.Println("\nFound URL:", link)
	fmt.Println("Status: ", resp.Status)
	fmt.Println("Code: ", resp.StatusCode)
}

func readFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Can not open file: ", err)
		return
	}

	defer file.Close()

	var wg sync.WaitGroup

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}

		wg.Add(1)
		go checkSite(url, &wg)
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error while reading file: ", err)
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go sites.txt")
		return
	}

	fname := os.Args[1]
	readFile(fname)
}
