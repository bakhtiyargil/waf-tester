package service

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

/*
1. refactor
*/
func (t *TesterService) StartInjectionTest(testRequest *model.TestRequest) (*model.Result, error) {
	result := &model.Result{}
	filepath.Walk("./data", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()

			var lines []string
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				lines = append(lines, line)
			}
			totalReq, blockedReq, err := t.makeConcurrentRequest(lines, testRequest)
			result.BlockedRequests = result.BlockedRequests + blockedReq
			result.TotalRequests = result.TotalRequests + totalReq
			if err != nil {
				log.Printf("Error reading file %s: %v", path, err)
				return nil
			}
		}
		return nil
	})
	return result, nil
}

func (t *TesterService) makeConcurrentRequest(lines []string, testRequest *model.TestRequest) (totalRequest int, blockedRequest int, e error) {
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

	for range wafCounter {
		blockedRequest++
	}
	totalRequest += len(lines)

	fmt.Println("Time elapsed: ", end.Sub(start))
	return totalRequest, blockedRequest, nil
}
