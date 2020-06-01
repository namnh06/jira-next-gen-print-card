package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Issue struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Fields struct {
		Issuetype struct {
			Name string `json:"name"`
		} `json:"issuetype"`
		Created  string `json:"created"`
		Priority struct {
			Name string `json:"name"`
		} `json:"priority"`
		Customfield10016 float64 `json:"customfield_10016"`
		Assignee         struct {
			EmailAddress string `json:"emailAddress"`
			DisplayName  string `json:"displayName"`
		} `json:"assignee"`
		Status struct {
			Self           string `json:"self"`
			Description    string `json:"description"`
			IconURL        string `json:"iconUrl"`
			Name           string `json:"name"`
			ID             string `json:"id"`
			StatusCategory struct {
				Self      string `json:"self"`
				ID        int    `json:"id"`
				Key       string `json:"key"`
				ColorName string `json:"colorName"`
				Name      string `json:"name"`
			} `json:"statusCategory"`
		} `json:"status"`
		Summary string `json:"summary"`
	} `json:"fields"`
}

type Data struct {
	Expand     string `json:"expand"`
	StartAt    int    `json:"startAt"`
	MaxResults int    `json:"maxResults"`
	Total      int    `json:"total"`
	Issues     []Issue
}

var err = godotenv.Load()
var sprint = os.Args[1]
var apiEndpoint = "https://yuriqa.atlassian.net/rest/agile/1.0/board/2/sprint/" + sprint + "/issue?maxResults=100"
var requestType = "GET"
var username string = os.Getenv("JIRA_EMAIL")
var password string = os.Getenv("JIRA_API_TOKEN")

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	return nil
}

func main() {
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	req, err := http.NewRequest(requestType, apiEndpoint, nil)
	if err != nil {
		log.Fatalln("[DEV] 78,", err)
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("[DEV] 84, ", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	data := Data{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalln(err)
	}

	// filter tasks to get lastPrintDateTime and currentDateTime
	tz, _ := time.LoadLocation("Australia/Brisbane")
	currentDateTime := time.Now().In(tz)
	var printDate time.Time
	if currentDateTime.Weekday() == time.Monday {
		printDate = currentDateTime.AddDate(0, 0, -3)
	} else {
		printDate = currentDateTime.AddDate(0, 0, -1)
	}
	lastPrintDateTime := time.Date(printDate.Year(), printDate.Month(), printDate.Day(), 11, 30, 0, 0, tz)
	layout := "2006-01-02T15:04:05.000Z"
	dataPrint := []Issue{}
	for i, v := range data.Issues {
		re := regexp.MustCompile(`\+(11|10)00`)
		str := re.ReplaceAllString(data.Issues[i].Fields.Created, "Z")
		lastDateTime, err := time.Parse(layout, str)

		if err != nil {
			log.Println("[DEV] Error ", err)
		}

		if lastDateTime.After(lastPrintDateTime) && lastDateTime.Before(currentDateTime) {
			dataPrint = append(dataPrint, v)
		}
	}

	sort.Slice(dataPrint, func(i, j int) bool {
		first, _ := strconv.ParseInt(dataPrint[i].ID, 10, 32)
		second, _ := strconv.ParseInt(dataPrint[j].ID, 10, 32)
		return first < second
	})

	file, _ := json.MarshalIndent(dataPrint, "", " ")
	err = ioutil.WriteFile("data.json", file, 0644)
}
