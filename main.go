package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
	"waf-tester/client"
	"waf-tester/config"
	"waf-tester/domain"
)

func main() {
	//loading config file
	var c config.Config
	c.LoadConfig("./config/config-local.yml")

	//loading sample file
	file, _ := os.Open("./data/sample.txt")
	defer file.Close()

	//reading sample texts
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	//making requests concurrently
	maxConcurrency := 64
	limiter := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup

	wafCounter := make(chan int, len(lines))
	var cc client.Client

	target := domain.GetTestTargetInstance()
	start := time.Now()
	for _, line := range lines {
		limiter <- struct{}{}
		wg.Add(1)
		go func(line string) {
			defer func() {
				<-limiter
				wg.Done()
			}()
			_, statusCode, _ := cc.DoRequestWithoutBody(target.Method, target.GetUrl()+"/"+line)
			if statusCode == 403 {
				wafCounter <- 1
			}
		}(line)
	}

	go func() {
		wg.Wait()
		close(wafCounter)
	}()
	end := time.Now()

	var i int8
	for range wafCounter {
		i++
	}

	fmt.Println(i)
	fmt.Println("elapsed: ", end.Sub(start))
}
