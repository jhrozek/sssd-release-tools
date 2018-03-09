package main

import "encoding/json"
import "flag"
import "fmt"
import "log"
import "net/http"
import "net/url"
import "os"

const IssuesURLFmt = "https://pagure.io/api/0/%s/issues"
const issueBaseURL = "https://pagure.io/SSSD/sssd/issue"

type IssueList struct {
	Status string
	Issues []Issue
}

type Issue struct {
	Id    int
	Title string
}

func SearchIssues(issuesUrl string, milestone string, status string) (*IssueList, error) {
	v := url.Values{"milestones": {milestone}, "status": {status}}
	resp, err := http.Get(issuesUrl + "?" + v.Encode())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http.Get returned %s", resp.Status)
	}

	var result IssueList
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func FormatIssue(i *Issue) string {
	issueUrl := fmt.Sprintf("%s/%d", issueBaseURL, i.Id)
	return fmt.Sprintf(" * `%d <%s>`_ - %s\n", i.Id, issueUrl, i.Title)
}

func main() {
	var project = flag.String("project", "", "The project to print issues for")
	var milestone = flag.String("milestone", "", "The milestone to print issues for")
	var status = flag.String("status", "closed", "Only return tickets with this status")

	flag.Parse()

	if *project == "" {
		fmt.Println("Please specify the project with --project")
		os.Exit(1)
	}

	if *milestone == "" {
		fmt.Println("Please specify the project with --milestone")
		os.Exit(1)
	}

	issuesUrl := fmt.Sprintf(IssuesURLFmt, *project)
	issuesList, err := SearchIssues(issuesUrl, *milestone, *status)
	if err != nil {
		log.Fatal(err)
	}

	for _, issue := range issuesList.Issues {
		fmt.Print(FormatIssue(&issue))
	}
}
