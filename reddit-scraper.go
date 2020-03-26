package main

import (
	"encoding/json"
	"fmt"
	"github.com/devfacet/gocmd"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var client = &http.Client{}
var verboseExecution = false

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flags := struct {
		Help      bool `short:"h" long:"help" description:"Display usage" global:"true"`
		Version   bool `short:"v" long:"version" description:"Display version"`
		VersionEx bool `long:"vv" description:"Display version (extended)"`
		Url struct {
			SubReddit     	string `short:"s" long:"sub" required:"true" description:"Subreddit to scrap"`
			Output 			string `short:"o" long:"output" required:"false" description:"Output file. If not present, send to STDOUT"`
			Limit 			int `short:"l" long:"limit" required:"false" description:"Limit the number of urls to retrieve. If not present or 0, fetch everything."`
			Verbose 		bool`short:"v" long:"verbose" required:"false" description:"Show extra information during the process and stats"`
		} `command:"urls" description:"fetch urls from post"`
	}{}

	// Echo command
	gocmd.HandleFlag("Url", func(cmd *gocmd.Cmd, args []string) error {
		verboseExecution = flags.Url.Verbose
		getURLs(flags.Url.SubReddit, flags.Url.Output, flags.Url.Limit)
		return nil
	})

	// Init the app
	gocmd.New(gocmd.Options{
		Name:        "Reddit URLs Scrapper",
		Version:     "1.0.0",
		Description: "A basic Reddit Scrapper to get all the url from posts in a given SubReddit",
		Flags:       &flags,
		ConfigType:  gocmd.ConfigTypeAuto,
	})
}

func getURLs(subreddit string, outfile string, limit int) {
	allUrls := make([]string,0,100)
	afterCode := ""
	for afterCode != "END" {
		results, newAfter := doRequest(afterCode, subreddit)

		sendVerboseMsg("Urls fetched in batch: "+ strconv.Itoa(len(results)))

		allUrls = append(allUrls, results...)
		afterCode = newAfter

		if limit != 0 && limit < len(allUrls) {
			allUrls = allUrls[0:limit]
			break
		}
	}

	var f *os.File
	var err error
	var usingStdout = false

	if outfile != "" {
		f, err = os.Create(outfile)
		if err != nil {
			f =	os.Stdout
			usingStdout = true
		}
	}else {
		f =	os.Stdout
		usingStdout = true
	}

	for _, aUrl := range allUrls {
		_, err = f.WriteString(aUrl +"\n")
		check(err)
	}

	sendVerboseMsg("Total Urls fetched: "+ strconv.Itoa(len(allUrls)))

	_ = f.Sync()
	if !usingStdout {
		_ = f.Close()
	}

}

func sendVerboseMsg(msg string) {
	if verboseExecution {
		fmt.Println(msg)
	}
}

func doRequest(afterCode string, subreddit string) ([]string, string) {
	urlToFetch := "https://www.reddit.com/r/" + subreddit + "/.json?limit=100&after=" + afterCode
	sendVerboseMsg("Request to: "+ urlToFetch)
	req, err := http.NewRequest("GET", urlToFetch, nil)

	req.Header.Add("User-Agent", `random"`)
	resp, err := client.Do(req)

	check(err)
	body, err := ioutil.ReadAll(resp.Body)

	//unmarshal json response ina generic way using maps
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	data := result["data"].(map[string]interface{})
	newAfterCode := data["after"]
	children := data["children"].([]interface{})
	var resultUrls []string
	resultUrls = make([]string,0,100)
	for _, child := range children {
		url := child.(map[string]interface{})["data"].(map[string]interface{})["url"].(string)
		resultUrls = append(resultUrls, url)
	}

	resp.Body.Close()
	if newAfterCode == nil {
		return  resultUrls, "END"
	}else {
		return resultUrls, newAfterCode.(string)
	}
}