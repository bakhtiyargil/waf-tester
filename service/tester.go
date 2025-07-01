package service

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"waf-tester/client"
	"waf-tester/model"
	"waf-tester/utility"
)

type TesterService struct {
	client *client.Client
	wp     utility.Worker
}

func NewTesterService(client *client.Client, pool *utility.WorkerPool) *TesterService {
	return &TesterService{client: client, wp: pool}
}

func (t *TesterService) StartInjectionTest(testRequest *model.TestRequest) error {
	t.wp.Start()

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

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			var routineFunc utility.RoutineFunction = t.processMethod
			task := utility.NewTask(scanner.Text(), model.FromRequest(testRequest), routineFunc)
			t.wp.Submit(task)
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

func (t *TesterService) processMethod(paramStatic interface{}, param interface{}) error {
	prs := paramStatic.(*model.Target)
	pr := param.(string)
	body, i, err := t.client.DoRequestWithoutBody(prs.Method, prs.GetUrl()+pr)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	if i != 403 {
		file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		if _, err := file.WriteString(prs.GetUrl() + pr + "\n" + string(body) + "\n"); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
	return nil
}
