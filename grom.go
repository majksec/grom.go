package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
	"net/http"
	"strconv"
	"time"
	"sync"
	"io/ioutil"
)

func main() {
	// Colors used for terminal output
	var Reset  = "\033[0m"
	var Red    = "\033[31m"
	var Green  = "\033[32m"
	var Yellow = "\033[33m"
	headers := "\n ################################################################### \n #Description : Domain response status code checker \n #Author      : majksec (twitch.tv/majksec)\n ################################################################### \n"
	

	var uncheckedSubdomainList []string
	var checkedSubdomainList []string

	// Defining timeout if request takes too long to complete
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	var inputFilePath = flag.String("f", "", "Domain text file path")
	var outputFilePath = flag.String("o", "", "Path and name where to output file")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, string(Green) + "Basic usage: go run grom.go -f=subdomains.txt -o=output.txt \n      -f flag requires raw subdomain list \n      -o flag is not required, specifies output file name \n" + string(Reset))
	}

	flag.Parse()

	time.Sleep(3*time.Second)

	// Reading input file
	inputFile, err := os.Open(*inputFilePath)
	if err != nil {
		fmt.Println(string(Red) + "Input file not specified, use -h for help" + string(Reset))
		return
	}
	defer inputFile.Close()
	
	fmt.Println(string(Yellow) + headers + string(Reset))

	data := bufio.NewScanner(inputFile)

	// Formatting url
	for data.Scan() {
		uncheckedSubdomainList = append(uncheckedSubdomainList, "https://" + data.Text())
		uncheckedSubdomainList = append(uncheckedSubdomainList, "http://" + data.Text())
	}

	var wg sync.WaitGroup

	// Checking subdomain status code and appending ones that are alive to output variable (checkedSubdomainList)
	for _, subdomain := range uncheckedSubdomainList {
		wg.Add(1)
		go func(subdomain string) {
			defer wg.Done()
			var request, err = client.Get(subdomain)

			if err != nil {
				fmt.Println(string(Red) + "[-] Domain " + subdomain + " is unreachable, continuing..." + string(Reset))
			} else {
				var statusCode = request.StatusCode
	
				if statusCode >= 500 && statusCode <= 599 {
					fmt.Println(string(Red) + "[" + strconv.Itoa(statusCode) + "] " + subdomain + string(Reset))
				} else if statusCode >= 200 && statusCode <= 299 {
					fmt.Println(string(Green) + "[" + strconv.Itoa(statusCode) + "] " + subdomain + string(Reset))
					checkedSubdomainList = append(checkedSubdomainList, "[" + strconv.Itoa(statusCode) + "]" + subdomain)
				} else {
					fmt.Println(string(Yellow) + "[" + strconv.Itoa(statusCode) + "] " + subdomain + string(Reset))
					checkedSubdomainList = append(checkedSubdomainList, "[" + strconv.Itoa(statusCode) + "]" + subdomain)
				}
			}
		}(subdomain)
	}
	wg.Wait()


	// Writing result to the output file
	resultList := ""

	for _, subdomain := range checkedSubdomainList {
		resultList += subdomain
		resultList += "\n"
	}

	ioutil.WriteFile(*outputFilePath, []byte(resultList), 0644)

	fmt.Println(string(Green) + "\n Script successfully finished" + string(Reset))
}