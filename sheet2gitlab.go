package main

import gitlab "gitlab.botsunit.com/infra/sheet2gitlab/gitlab"

func main() {
	// userStoryList := sheetReader.ReadSheet()
	// for _, us := range userStoryList {
	// 	us.Print()
	// }

	gitlab.InitializeGitlab()
	gitlab.CreateGitlabIssue()

}
