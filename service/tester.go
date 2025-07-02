package service

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"waf-tester/client"
	"waf-tester/model"
	"waf-tester/utility"
)

type TesterService struct {
	client *client.Client
}

func NewTesterService(client *client.Client) *TesterService {
	return &TesterService{client: client}
}

func (t *TesterService) StartInjectionTest(testRequest *model.TestRequest) error {
	wp := utility.NewWorkerPoolExecutor(testRequest.Host, 128)
	wp.Start()
	defer func() {
		go wp.Stop()
	}()

	err := filepath.Walk("./data", func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking to file %s: %w", path, walkErr)
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() != "payloads.txt" {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var routineFunc utility.RoutineFunction = t.processMethod
		for scanner.Scan() {
			task := utility.NewTask(scanner.Text(), model.FromRequest(testRequest), routineFunc)
			wp.Submit(task)
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("scanner error on file %s: %w", path, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to start injection test: %w", err)
	}
	return nil
}

func (t *TesterService) processMethod(paramStatic interface{}, param interface{}) {
	prs := paramStatic.(*model.Target)
	escPr := url.PathEscape(param.(string))
	body, httpStatus, err := t.client.DoRequestWithoutBody(prs.Method, prs.GetUrl()+"/"+escPr)
	if err != nil {
		fmt.Printf("failed to do request: %s", err)
		return
	}
	if strconv.Itoa(httpStatus) != prs.Criteria[1] {
		if !strings.Contains(string(body), prs.Criteria[0]) {
			file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("failed to open file: %s", err)
				return
			}
			defer file.Close()

			if _, err := file.WriteString(prs.GetUrl() + "/" + escPr + "\n" + string(body) + "\n"); err != nil {
				fmt.Printf("failed to write to file: %s", err)
				return
			}
		}
	}
}
