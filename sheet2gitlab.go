package main

import (
	"fmt"
	"strconv"

	gitlab "gitlab.botsunit.com/infra/sheet2gitlab/gitlab"
	sheetReader "gitlab.botsunit.com/infra/sheet2gitlab/googlespreadsheet"
)

func main() {

	gitlab.InitializeGitlab()

	gitlab.SynchronizeLabels("boobs")

	userStoryList := sheetReader.ReadSheet()
	fmt.Println(len(userStoryList))
	us := userStoryList[45]
	fmt.Println(us.DescriptionToTitle())

	manDayStr := strconv.FormatFloat(float64(us.ManDay), 'f', 1, 64)

	fmt.Println("Man day converted :", manDayStr)
	labelList := []string{manDayStr, "functionnal"}
	if len(us.Bot) > 1 {
		labelList = append(labelList, "tranversal")
	}
	gitlabLabels := gitlab.GetGitlabLabels(labelList, "boobs", "meta_boob")

	gitlab.PrintGitlabLabels(gitlabLabels)
	milestone, _ := gitlab.GetMilestone("0.5", "meta_boob")
	gitlab.CreateGitlabIssue(us.Bot, us.DescriptionToTitle(), us.Description, milestone.ID, gitlabLabels, us.Number)

}
