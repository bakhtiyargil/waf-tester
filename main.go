package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"waf-tester/client"
	"waf-tester/config"
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

	var cc client.Client

	start := time.Now()
	var i int8
	for _, line := range lines {
		resp, _ := cc.DoRequest("GET", "http://waffy.xyz/DVWA/"+line, "")
		if resp.StatusCode == 403 {
			i++
		}
	}
	end := time.Now()
	//64
	fmt.Println(i)
	fmt.Println("elapsed: ", end.Sub(start))
}
