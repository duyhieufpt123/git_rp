package util

import (
	"log"
	"time"

	"github.com/xanzy/go-gitlab"
)

// get Date with h,m,s < 0 for previous
func GetDate(h, m, s int64) *time.Time {
	timein := time.Now().Add(time.Hour*time.Duration(h) +
		time.Minute*time.Duration(m) +
		time.Second*time.Duration(s))
	return &timein
}

func ListBranchNameProject(pid interface{}) []string {
	res := []string{}
	config, err := LoadConfig(".")
	if err != nil {
		log.Panic(err)
	}

	git, err := gitlab.NewClient(config.TOKEN, gitlab.WithBaseURL(config.SERVER_GIT))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	opt := &gitlab.ListBranchesOptions{}
	branches, _, err := git.Branches.ListBranches(pid, opt)
	if err != nil {
		log.Fatal(err)
	}

	for _, branch := range branches {
		res = append(res, branch.Name)
	}

	return res
}
func ListSubGroupsID(gid interface{}) []int {
	//LOAD CONFIG
	config, err := LoadConfig(".")
	var ListSubGrID []int
	if err != nil {
		log.Panic(err)
	}
	//
	git, err := gitlab.NewClient(config.TOKEN, gitlab.WithBaseURL(config.SERVER_GIT))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	group_sub_opt := &gitlab.ListSubGroupsOptions{AllAvailable: gitlab.Bool(true)}
	group_sub, _, err := git.Groups.ListSubGroups(gid, group_sub_opt)
	if err != nil {
		log.Fatal(err)
	}
	for _, mem := range group_sub {
		ListSubGrID = append(ListSubGrID, mem.ID)

	}
	return ListSubGrID
}
func ListGroupProjectsID(gid interface{}) []int {
	var ListProjectID []int
	// LOAD CONFIG
	config, err := LoadConfig(".")
	if err != nil {
		log.Panic(err)
	}
	git, err := gitlab.NewClient(config.TOKEN, gitlab.WithBaseURL(config.SERVER_GIT))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	//
	group_proj_opt := &gitlab.ListGroupProjectsOptions{}
	group_project, _, err := git.Groups.ListGroupProjects(gid, group_proj_opt)
	if err != nil {
		log.Fatal(err)
	}
	for _, proj := range group_project {
		ListProjectID = append(ListProjectID, proj.ID)
	}
	return ListProjectID
}
