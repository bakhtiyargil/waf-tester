package service

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
	"waf-tester/client"
	"waf-tester/model"
)

type TesterService struct {
	client *client.Client
}

func NewTesterService(client *client.Client) *TesterService {
	return &TesterService{client: client}
}

func (t *TesterService) StartInjectionTest(testRequest *model.TestRequest) (*model.Result, error) {
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

	target := model.Target{
		Host:   testRequest.Host,
		Path:   testRequest.Path,
		Method: testRequest.Method,
	}
	start := time.Now()
	for _, line := range lines {
		limiter <- struct{}{}
		wg.Add(1)
		go func(line string) {
			defer func() {
				<-limiter
				wg.Done()
			}()
			_, statusCode, _ := t.client.DoRequestWithoutBody(target.Method, target.GetUrl()+"/"+line)
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

	var i int
	for range wafCounter {
		i++
	}

	fmt.Println("Time elapsed: ", end.Sub(start))
	return &model.Result{TotalRequests: len(lines), BlockedRequests: i}, nil
}
