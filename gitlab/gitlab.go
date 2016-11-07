package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"
)

const (
	gitlabToken = "SAyvT7LLPozwLd5sorSx"
	gitlabURL   = "https://gitlab.botsunit.com/api/v3"
)

var (
	git            *gitlab.Client
	projectMapping = map[string]string{
		"chat":         "boob_chat",
		"broadcast":    "boob_broadcast",
		"html":         "boob_html",
		"identity":     "boob_identity",
		"payment":      "boob_payment",
		"user":         "boob_user",
		"shop":         "boob_shop",
		"realtime":     "boob_realtime",
		"meta_boob":    "meta_boob",
		"notification": "boob_notification",
		"flash":        "boob_flash",
	}

	labels = map[string]string{
		"functionnal": "#CB4335",
		"tech":        "#1F618D",
		"bug":         "#FF5733",
		"tranversal":  "#8E44AD",
		"0.2":         "#34495E",
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
		fmt.Printf("label: %s\n", label)
	}
}

//GetMilestone find a milestone with name milestoneName for project with pid
func GetMilestone(milestoneName string, project string) (gitlab.Milestone, error) {
	pid := getProjectID(project)
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

func getProjectID(projectName string) int {
	//fmt.Printf("Project name is : [%s] ", projectName)
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

func issueExist(title string, pid int, issueNumber int32) bool {
	issueList, _, err := git.Issues.ListProjectIssues(pid, &gitlab.ListProjectIssuesOptions{})
	if err != nil {
		fmt.Println("Error while testing if issue exists, prefrer crash instead of creating duplicates")
		panic(err.Error())
	}
	for _, issue := range issueList {
		if strings.Contains(issue.Title, fmt.Sprintf("#%d", issueNumber)) {
			return true
		}
	}
	return false
}

//ListAllProjectAccessibles give access to all project printed on stoud
func ListAllProjectAccessibles() {
	falseBool := false
	optAll := &gitlab.ListProjectsOptions{Archived: &falseBool}
	projectsAll, _, err := git.Projects.ListProjects(optAll)
	if err != nil {
		panic(err.Error())
	}
	for _, project := range projectsAll {
		fmt.Println(project.Name)
	}
}

//CreateGitlabIssue allow to create a gitlab issue
func CreateGitlabIssue(projects []string, title string, description string, labels []string, number int32, milestone string) {

	// title, description, AssigneeID, MilestoneID, Labels
	for _, p := range projects {
		pid := getProjectID(projectMapping[p])
		milestone, _ := GetMilestone(milestone, projectMapping[p])
		if !issueExist(title, pid, number) {
			options := &gitlab.CreateIssueOptions{
				Title:       &title,
				Description: &description,
				MilestoneID: &milestone.ID,
				Labels:      labels,
			}
			git.Issues.CreateIssue(pid, options)
		}
	}
}

//SearchIssueWithoutMileStone list issue without milestone affected
func SearchIssueWithoutMileStone() {
	labs := &gitlab.Labels{"functionnal"}
	listIssuesOpts := &gitlab.ListIssuesOptions{Labels: *labs}
	issues, _, _ := git.Issues.ListIssues(listIssuesOpts)
	for _, issue := range issues {
		if issue.Milestone.ID == 0 {
			fmt.Printf("Url : %s \n", issue.WebURL)
			fmt.Printf("Pid : %d", issue.ProjectID)
			fmt.Printf("Title : %s", issue.Title)
		}
	}
}

//GetIssuesForMilestone retrieve all the issue for a milestone
func GetIssuesForMilestone(milestoneName string, group string) map[string][]gitlab.Issue {
	falseBool := false
	optAll := &gitlab.ListProjectsOptions{Archived: &falseBool}
	projectsAll, _, err := git.Projects.ListProjects(optAll)
	if err != nil {
		panic("Cannot retrieve project list!")
	}
	var idxDelete int
	for i, proj := range projectsAll {
		if proj.Name == "payment" {
			idxDelete = i
		}
	}
	projectsAll = append(projectsAll[:idxDelete], projectsAll[idxDelete+1:]...)
	issuesPerProject := map[string][]gitlab.Issue{}

	for _, project := range projectsAll {
		pName := project.Name
		milestone, err := GetMilestone(milestoneName, pName)
		if err != nil {
			panic(fmt.Sprintf("error while retrieving milestone %s for project : %s", milestoneName, pName))
		}
		opt := &gitlab.ListProjectIssuesOptions{
			Milestone: &milestone.Title,
		}
		issuesForProject, _, err := git.Issues.ListProjectIssues(project.ID, opt)
		if err != nil {
			panic(fmt.Sprintf("error while retrieving issues for project : %s", project.Name))
		}
		for _, issue := range issuesForProject {
			// if len(issuesPerProject[project.Name]) == 0 {
			// 	issuesPerProject[project.Name] = []gitlab.Issue{}
			// }
			issuesPerProject[project.Name] = append(issuesPerProject[project.Name], *issue)
		}
	}
	return issuesPerProject
}

//GetLabelManDay return label as float32 to retrieve man day
func GetLabelManDay(labels []string) float32 {
	result := 0.0
	for _, label := range labels {
		if strings.Contains(label, ".") {
			result, _ = strconv.ParseFloat(label, 32)
		}
	}
	return float32(result)
}
