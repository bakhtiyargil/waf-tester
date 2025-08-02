package service

import (
	"bufio"
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"waf-tester/client"
	"waf-tester/domain"
	"waf-tester/domain/model"
	"waf-tester/logger"
	"waf-tester/utility"
)

type Tester interface {
	Start(testRequest *model.TestRequest) (testId string, err error)
	Terminate(testId string) error
}

type InjectionTester struct {
	client  client.Client
	logger  logger.Logger
	useCase domain.TestUseCase
}

func NewInjectionTester(client client.Client, logger logger.Logger, useCase domain.TestUseCase) Tester {
	return &InjectionTester{
		client:  client,
		logger:  logger,
		useCase: useCase,
	}
}

func (t *InjectionTester) Start(testRequest *model.TestRequest) (string, error) {
	wp := utility.NewWorkerPoolExecutor(testRequest.GetApi(), 128, t.logger)
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
			task := utility.NewTask(scanner.Text(), model.FromRequest(wp.GetId(), testRequest), routineFunc)
			wp.Submit(task)
		}
		if err := scanner.Err(); err != nil {
			return t.logger.ErrorR(err)
		}
		return nil
	})
	if err != nil {
		return "", t.logger.ErrorR(err)
	}

	err = wp.Start()
	if err != nil {
		return "", t.logger.ErrorR(err)
	}
	defer func() {
		go wp.Finish()
	}()

	return wp.GetId(), nil
}

func (t *InjectionTester) Terminate(testId string) error {
	var wp, err = utility.PlContext.Get(testId)
	if err != nil {
		return err
	}
	err = wp.TerminateGracefully()
	if err != nil {
		return err
	}
	return nil
}

func (t *InjectionTester) processMethod(paramStatic interface{}, param interface{}) {
	var (
		tst domain.Test
		err error
	)

	prs := paramStatic.(*model.Target)
	escPr := prs.GetUrl() + "/" + url.PathEscape(param.(string))
	body, httpStatus, elapsed, err := t.client.DoRequestWithoutBody(prs.Method, escPr)
	if err != nil {
		t.logger.Error(err)
		return
	}

	if strconv.Itoa(httpStatus) != prs.Criteria[1] {
		if !strings.Contains(string(body), prs.Criteria[0]) {
			tst = domain.Test{
				Host:           prs.Host,
				Path:           escPr,
				Method:         prs.Method,
				ResponseBdy:    string(body),
				ResponseStatus: httpStatus,
				ResponseTime:   elapsed.Milliseconds(),
				TestID:         prs.Id,
			}

			_, err = t.useCase.InsertOne(context.Background(), &tst)
			if err != nil {
				t.logger.Error(err)
				return
			}
		}
	}
}
