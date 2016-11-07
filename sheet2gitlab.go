package main

import (
	gitlab "gitlab.botsunit.com/infra/sheet2gitlab/gitlab"
	sheetReader "gitlab.botsunit.com/infra/sheet2gitlab/googlespreadsheet"
)

func main() {

	gitlab.InitializeGitlab()
	//gitlab.ListAllProjectAccessibles()

	// gitlab.SynchronizeLabels("boobs")
	//
	// userStoryList := sheetReader.ReadSheet()
	// fmt.Printf("We have %d user stories", len(userStoryList))
	// for _, us := range userStoryList {
	// 	fmt.Println("####################################################################")
	// 	fmt.Printf("Issue #%d", us.Number)
	// 	fmt.Println(us.DescriptionToTitle())
	// 	fmt.Println(us.Bot)
	// 	manDayStr := strconv.FormatFloat(float64(us.ManDay), 'f', 1, 64)
	// 	labelList := []string{manDayStr, "functionnal"}
	// 	if len(us.Bot) > 1 {
	// 		labelList = append(labelList, "tranversal")
	// 	}
	// 	gitlabLabels := gitlab.GetGitlabLabels(labelList, "boobs", "meta_boob")
	// 	gitlab.PrintGitlabLabels(gitlabLabels)
	// 	gitlab.CreateGitlabIssue(us.Bot, us.DescriptionToTitle(), us.Description, gitlabLabels, us.Number, "0.5")
	// }
	//gitlab.SearchIssueWithoutMileStone()

	// issuesPerProject := gitlab.GetIssuesForMilestone("0.5", "boobs")
	// issueCounter := 0
	// for project, issues := range issuesPerProject {
	// 	fmt.Println("=================================================================")
	// 	fmt.Printf("Project : %s \n", project)
	// 	for _, issue := range issues {
	// 		fmt.Println(issue.Title)
	// 		manDays := gitlab.GetLabelManDay(issue.Labels)
	// 		userStory := sheetReader.UserStory{
	// 			Priority:    int32(issue.Milestone.ID),
	// 			Number:      int32(issueCounter),
	// 			Description: issue.Description,
	// 			Bot:         []string{project},
	// 			ManDay:      manDays,
	// 		}
	// 		issueCounter++
	// 		userStory.Print()
	// 	}
	// 	fmt.Println("=================================================================")
	// }
	sheetReader.WriteSheetAsCSV("0.5")
}
