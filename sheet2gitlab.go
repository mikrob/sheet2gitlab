package main

import gitlab "gitlab.botsunit.com/infra/sheet2gitlab/gitlab"

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
	gitlab.SearchIssueWithoutMileStone()
}
