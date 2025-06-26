package service

import (
	"bufio"
	"fmt"
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
create chunk based parallel tasks and execute
*/
func (t *TesterService) StartInjectionTest(testRequest *model.TestRequest) (*model.Result, error) {
	target := model.Target{
		Host:   testRequest.Host,
		Path:   testRequest.Path,
		Method: testRequest.Method,
	}
	result := &model.Result{}

	err := filepath.Walk("./data", func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking to file %s: %w", path, walkErr)
		}
		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("scanner error on file %s: %w", path, err)
		}

		total, blocked, err := t.makeConcurrentTest(lines, &target)
		if err != nil {
			return fmt.Errorf("concurrent test failed for file %s: %w", path, err)
		}

		result.BlockedRequests += blocked
		result.TotalRequests += total
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("injection test failed: %w", err)
	}
	return result, nil
}

func (t *TesterService) makeConcurrentTest(paths []string, target *model.Target) (total int, blocked int, e error) {
	maxConcurrency := 128
	pathsLen := len(paths)
	limiter := make(chan struct{}, maxConcurrency)
	wafBlock := make(chan int, pathsLen)

	var wg sync.WaitGroup
	start := time.Now()
	for _, line := range paths {
		limiter <- struct{}{}
		wg.Add(1)
		go func(line string) {
			defer func() {
				<-limiter
				wg.Done()
			}()
			_, statusCode, _ := t.client.DoRequestWithoutBody(target.Method, target.GetUrl()+"/"+line)
			if statusCode == 403 {
				wafBlock <- 1
			}
		}(line)
	}

	go func() {
		wg.Wait()
		close(wafBlock)
	}()
	end := time.Now()

	for range wafBlock {
		blocked++
	}
	total += pathsLen

	fmt.Println("time elapsed: ", end.Sub(start))
	return total, blocked, nil
}
