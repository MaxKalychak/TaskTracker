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

func checkSite(link string, wg *sync.WaitGroup, ch chan bool) {
	defer wg.Done()
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(link)
	if err != nil {
		fmt.Printf("%s - Error: %v\n", link, err)
		ch <- false
		return
	}

	defer resp.Body.Close()
	fmt.Println("\nFound URL:", link)
	fmt.Println("Status: ", resp.Status)
	fmt.Println("Code: ", resp.StatusCode)
	ch <- true
}

func readFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Can not open file: ", err)
		return
	}

	defer file.Close()

	var wg sync.WaitGroup
	result := make(chan bool)

	total, alive, errors := 0, 0, 0

	go func() {
		for res := range result {
			total++
			if res == true {
				alive++
			} else {
				errors++
			}
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}

		wg.Add(1)
		go checkSite(url, &wg, result)
	}

	wg.Wait()

	close(result)

	time.Sleep(100 * time.Millisecond)

	if err := scanner.Err(); err != nil {
		fmt.Println("Error while reading file: ", err)
	}
	fmt.Println("\n--- Статистика перевірки ---")
	fmt.Printf("Всього перевірено: %d\n", total)
	fmt.Printf("Працюють:        %d\n", alive)
	fmt.Printf("Помилок:         %d\n", errors)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go sites.txt")
		return
	}

	fname := os.Args[1]
	readFile(fname)
}
