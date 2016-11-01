package gitlab

import (
	"fmt"
	"log"

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
		"meta_boob": "meta_boob",
	}

	labels = map[string]string{
		"functionnal": "#CB4335",
		"tech":        "#1F618D",
		"bug":         "#FF5733",
		"tranversal":  "#8E44AD",
		"0.5":         "#34495E",
		"1.0":         "#34495E",
		"2.0":         "#34495E",
		"3.0":         "#34495E",
		"4.0":         "#34495E",
		"5.0":         "#34495E",
		"6.0":         "#34495E",
		"7.0":         "#34495E",
	}
)

// InitializeGitlab initialize GitlabClient
func InitializeGitlab() {
	git = gitlab.NewClient(nil, gitlabToken)
	git.SetBaseURL(gitlabURL)
}

// ListUsers allow to test if gitlab connection work by listing and printing users
func ListUsers() {
	users, _, err := git.Users.ListUsers(nil)
	if err != nil {
		panic(err.Error())
	}
	for _, user := range users {
		fmt.Println(user.Name)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// GetGitlabLabels allow to pass a list of label in slice of string and to retrieve a list of Real struct gitlab label
func GetGitlabLabels(list []string, group string, project string) []string {
	var gitlabLabels []gitlab.Label
	var computedResult []string
	foundLabels, _, err := git.Labels.ListLabels(fmt.Sprintf("%s/%s", group, project))
	if err != nil {
		log.Println(err)
	}
	for _, label := range foundLabels {
		if contains(list, label.Name) {
			gitlabLabels = append(gitlabLabels, *label)
		}
	}
	for _, strLabel := range gitlabLabels {
		computedResult = append(computedResult, strLabel.Name)
	}
	return computedResult
}

//PrintGitlabLabels allow to print a list of gitlab labels
func PrintGitlabLabels(list []string) {
	for _, label := range list {
		log.Printf("label: %s", label)
	}
}

//GetMilestone find a milestone with name milestoneName for project with pid
func GetMilestone(milestoneName string, project string) (gitlab.Milestone, error) {
	pid := getProjectId(project)
	mileStoneList, _, err := git.Milestones.ListMilestones(pid, &gitlab.ListMilestonesOptions{})
	for _, milestone := range mileStoneList {
		if milestone.Title == milestoneName {
			return *milestone, nil
		}
	}
	return gitlab.Milestone{}, err
}

//SynchronizeLabels allow to ensure all labels exists in all related project in gitlab
func SynchronizeLabels(group string) {
	// List all labels
	var gitlabLabels []string
	for _, repo := range projectMapping {
		// retrive gitlab labels per projects
		//fmt.Printf("Labels for bot %s gitlab project : %s", bot, repo)
		foundLabels, _, err := git.Labels.ListLabels(fmt.Sprintf("%s/%s", group, repo))
		if err != nil {
			log.Println(err)
		}
		// build labels as array of string
		for _, label := range foundLabels {
			//log.Printf("Found label: %s", label.Name)
			gitlabLabels = append(gitlabLabels, label.Name)
		}
		// create missings labels
		for label, color := range labels {
			if !contains(gitlabLabels, label) {
				fmt.Println("Gitlab doesnt have label", label)
				l := &gitlab.CreateLabelOptions{
					Name:  gitlab.String(label),
					Color: gitlab.String(color),
				}
				label, _, err := git.Labels.CreateLabel(fmt.Sprintf("%s/%s", group, repo), l)
				if err != nil {
					fmt.Println(err)
				}
				log.Printf("Created label: %s\nWith color: %s\n", label.Name, label.Color)
				//log.Printf("Created label: %s\nWith color: %s\n", *l.Name, *l.Color)
			}
		}
	}

}

func getProjectId(projectName string) int {
	fmt.Println("looking for project")
	opt := &gitlab.ListProjectsOptions{Search: gitlab.String(projectName)}
	projects, _, err := git.Projects.ListProjects(opt)
	var pid int
	if err != nil {
		panic(err.Error())
	}
	if len(projects) != 1 {
		panic(fmt.Sprintf("Cannot find project %s", projectName))
	} else {
		pid = projects[0].ID
	}
	return pid
}

//CreateGitlabIssue allow to create a gitlab issue
func CreateGitlabIssue(project string, title string, description string, milestoneID int, labels []string) {
	// title, description, AssigneeID, MilestoneID, Labels
	pid := getProjectId(project)
	options := &gitlab.CreateIssueOptions{
		Title:       &title,
		Description: &description,
		MilestoneID: &milestoneID,
		Labels:      labels,
	}

	git.Issues.CreateIssue(pid, options)
}
