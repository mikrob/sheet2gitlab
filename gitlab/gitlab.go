package gitlab

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"
)

const (
	gitlabToken = "SAyvT7LLPozwLd5sorSx"
	gitlabURL   = "https://gitlab.botsunit.com/api/v3"
)

var (
	git            *gitlab.Client
	projectMapping = map[string]string{
		"chat":      "boob_chat",
		"broadcast": "boob_broadcast",
		"html":      "boob_html",
		"identity":  "boob_identity",
		"payment":   "payment",
		"user":      "boob_user",
		"shop":      "boob_shop",
		"realtime":  "boob_realtime",
	}
)

// InitializeGitlab initialize GitlabClient
func InitializeGitlab() {
	git = gitlab.NewClient(nil, gitlabToken)
	git.SetBaseURL(gitlabURL)
}

//CreateGitlabIssue allow to create a gitlab issue
func CreateGitlabIssue() {

	users, _, err := git.Users.ListUsers(nil)
	if err != nil {
		panic(err.Error())
	}
	for _, user := range users {
		fmt.Println(user.Name)
	}
}
