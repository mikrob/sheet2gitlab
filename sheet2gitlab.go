package main

import (
	"fmt"

	gitlab "gitlab.botsunit.com/infra/sheet2gitlab/gitlab"
)

func main() {
	// userStoryList := sheetReader.ReadSheet()
	// for _, us := range userStoryList {
	// 	us.Print()
	// }

	gitlab.InitializeGitlab()

	gitlab.SynchronizeLabels("boobs")
	// labelList := []string{
	// 	"functionnal", "0.5", "tranversal",
	// }
	// gitlabLabels := gitlab.GetGitlabLabels(labelList, "boobs", "meta_boob")
	// gitlab.PrintGitlabLabels(gitlabLabels)
	//
	// gitlab.CreateGitlabIssue("meta_boob", "issue title", "issue description", nil, gitlabLabels)

	milestone, _ := gitlab.GetMilestone("0.5", "meta_boob")
	fmt.Println(milestone.ID)
	fmt.Println(milestone.Title)
	fmt.Println(milestone.Description)
}
