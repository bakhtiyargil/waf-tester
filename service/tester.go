package service

import (
	"bufio"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"waf-tester/client"
	"waf-tester/logger"
	"waf-tester/model"
	"waf-tester/utility"
)

type TesterService struct {
	client *client.Client
	wp     utility.Worker
	logger *logger.AppLogger
}

func NewTesterService(client *client.Client, wp utility.Worker, logger *logger.AppLogger) *TesterService {
	return &TesterService{
		client: client,
		wp:     wp,
		logger: logger,
	}
}

func (t *TesterService) StartInjectionTest(testRequest *model.TestRequest) error {
	err := filepath.Walk("./data", func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return t.logger.ErrorR(walkErr)
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() != "payloads.txt" {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return t.logger.ErrorR(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var routineFunc utility.RoutineFunction = t.processMethod
		for scanner.Scan() {
			task := utility.NewTask(scanner.Text(), model.FromRequest(testRequest), routineFunc)
			t.wp.Submit(task)
		}
		if err := scanner.Err(); err != nil {
			return t.logger.ErrorR(err)
		}
		return nil
	})

	if err != nil {
		return t.logger.ErrorR(err)
	}
	return nil
}

func (t *TesterService) processMethod(paramStatic interface{}, param interface{}) {
	prs := paramStatic.(*model.Target)
	escPr := url.PathEscape(param.(string))
	body, httpStatus, err := t.client.DoRequestWithoutBody(prs.Method, prs.GetUrl()+"/"+escPr)
	if err != nil {
		t.logger.Error(err)
	}
	if strconv.Itoa(httpStatus) != prs.Criteria[1] {
		if !strings.Contains(string(body), prs.Criteria[0]) {
			file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				t.logger.Error(err)
				return
			}
			defer file.Close()

			if _, err := file.WriteString(prs.GetUrl() + "/" + escPr + "\n" + string(body) + "\n"); err != nil {
				t.logger.Error(err)
				return
			}
		}
	}
}
