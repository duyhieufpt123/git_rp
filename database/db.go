package database

import (
	"gitlabapi/models"
	"gitlabapi/util"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/xanzy/go-gitlab"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

type Database struct {
}

/*
created an init function to load the environment variables from the config.env file
and connect the Golang server to the Postgres database
*/
//init la ham de tao
func (database *Database) Init() {
	config, err := util.LoadConfig(".")
	//nil = null

	if err != nil {
		log.Fatal("cannot load config:", err) //err = null thi chay bth ko co loi
	}

	dbURL := config.DBSource //config goi DB
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	//su dung Migrator de tao table
	db.Migrator().CreateTable(&models.Project{})
	db.Migrator().CreateTable(&models.User{})
	db.Migrator().CreateTable(&models.Branch{})
	db.Migrator().CreateTable(&models.Commit{})
	db.Migrator().CreateTable(&models.FileCommit{})
}

// Connect database
func (database *Database) Connect() *gorm.DB {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	dbURL := config.DBSource
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func (database *Database) DropTable() {
	gorm := database.Connect()
	//Execute table
	gorm.Exec(`drop table "file_commits" cascade;`)
	gorm.Exec(`drop table "commits" cascade;`)
	gorm.Exec(`drop table "branches" cascade;`)
	gorm.Exec(`drop table "users" cascade;`)
	gorm.Exec(`drop table "projects" cascade;`)
}
func (database *Database) ImportDB(pid interface{}, branch string) {
	var exist bool
	// load database and config
	db := database.Connect()
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Panic(err)
	}
	git, err := gitlab.NewClient(config.TOKEN, gitlab.WithBaseURL(config.SERVER_GIT))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// access gitlab commit models
	opt_commit := &gitlab.ListCommitsOptions{RefName: gitlab.String(branch), WithStats: gitlab.Bool(true)}
	commit, _, err := git.Commits.ListCommits(pid, opt_commit)
	if err != nil {
		log.Fatal(err)
	}
	// access gitlab commit diff models
	opt_diff := &gitlab.GetCommitDiffOptions{PerPage: 10}
	Commits := []models.Commit{}
	// modify commit and file_commit models
	for _, com := range commit {
		commitinfo := models.Commit{}
		user := models.User{}
		// check user exist
		_ = db.Model(user).
			Select("count(*) > 0").
			Where("Name = ?", com.AuthorName).
			Find(&exist).Error
		if exist { //  exist
			db.First(&user, "Name = ?", com.AuthorName)
		} else { // non exist
			// create user models
			user.Email = com.AuthorEmail
			user.Name = com.AuthorName
			user.Id = uuid.New()
			db.Create(&user)
		}
		//
		diffs, _, err := git.Commits.GetCommitDiff(pid, com.ShortID, opt_diff)
		if err != nil {
			log.Fatal(err)
		}
		for _, diff := range diffs {
			fileinfo := models.FileCommit{Name: diff.OldPath}
			if diff.DeletedFile {
				fileinfo.Status = "Deleted"
			} else if diff.NewFile {
				fileinfo.Status = "Added"
			} else {
				fileinfo.Status = "Modified"
			}
			commitinfo.FileCommits = append(commitinfo.FileCommits, fileinfo)
		}
		commitinfo.Line_added = com.Stats.Additions
		commitinfo.Line_removed = com.Stats.Deletions
		commitinfo.User_name = com.AuthorName
		commitinfo.User_Id = user.Id
		commitinfo.CommitDate = com.CommittedDate
		commitinfo.Message = com.Message
		commitinfo.CommitID = com.ID
		Commits = append(Commits, commitinfo)
	}
	// access gitlab project models
	opt_project := &gitlab.GetProjectOptions{}
	project, _, err := git.Projects.GetProject(pid, opt_project)
	if err != nil {
		log.Fatal(err)
	}
	// create project models
	// project exist ?
	project_model := models.Project{}
	_ = db.Model(&models.Project{}).
		Select("count(*) > 0").
		Where("ID = ?", project.ID).
		Find(&exist).Error
	if exist { // project exist
		db.First(&project_model, "ID = ?", project.ID)
	} else { // non exist
		project_model.Name = project.Name
		project_model.Id = project.ID
		db.Create(&project_model)
	}
	//create branch models
	branch_model := models.Branch{}
	branch_model.Id = uuid.New()
	branch_model.Name = branch
	branch_model.Project_Id = project_model.Id
	db.Create(&branch_model)
	// create commit models
	for _, commit := range Commits {
		commit.Id = uuid.New()
		temp := commit.Id
		commit.Project_Id = project_model.Id
		commit.Branch_Id = branch_model.Id
		db.Create(&commit)
		//create file models
		for _, file := range commit.FileCommits {
			file.Commit_Id = temp
			file.File_Id = uuid.New()
			db.Create(&file)
		}
	}
}
func (database *Database) UpdateDB(pid interface{}, branch string, time *time.Time) {
	var exist bool
	// load database and config
	db := database.Connect()
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Panic(err)
	}
	git, err := gitlab.NewClient(config.TOKEN, gitlab.WithBaseURL(config.SERVER_GIT))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// access gitlab commit models
	opt_commit := &gitlab.ListCommitsOptions{RefName: gitlab.String(branch), WithStats: gitlab.Bool(true), Since: time}
	commit, _, err := git.Commits.ListCommits(pid, opt_commit)
	if err != nil {
		log.Fatal(err)
	}
	// access gitlab commit diff models
	opt_diff := &gitlab.GetCommitDiffOptions{PerPage: 10}
	Commits := []models.Commit{}
	// modify commit and file_commit models
	for _, com := range commit {
		commitinfo := models.Commit{}
		user := models.User{}
		// check user exist
		_ = db.Model(user).
			Select("count(*) > 0").
			Where("Name = ?", com.AuthorName).
			Find(&exist).Error
		if exist { //  exist
			db.First(&user, "Name = ?", com.AuthorName)
		} else { // non exist
			// create user models
			user.Email = com.AuthorEmail
			user.Name = com.AuthorName
			user.Id = uuid.New()
			db.Create(&user)
		}
		//
		diffs, _, err := git.Commits.GetCommitDiff(pid, com.ShortID, opt_diff)
		if err != nil {
			log.Fatal(err)
		}
		for _, diff := range diffs {
			fileinfo := models.FileCommit{Name: diff.OldPath}
			if diff.DeletedFile {
				fileinfo.Status = "Deleted"
			} else if diff.NewFile {
				fileinfo.Status = "Added"
			} else {
				fileinfo.Status = "Modified"
			}
			commitinfo.FileCommits = append(commitinfo.FileCommits, fileinfo)
		}
		commitinfo.Line_added = com.Stats.Additions
		commitinfo.Line_removed = com.Stats.Deletions
		commitinfo.User_name = com.AuthorName
		commitinfo.User_Id = user.Id
		commitinfo.CommitDate = com.CommittedDate
		commitinfo.CommitID = com.ID
		commitinfo.Message = com.Message
		Commits = append(Commits, commitinfo)
	}
	// access gitlab project models
	opt_project := &gitlab.GetProjectOptions{}
	project, _, err := git.Projects.GetProject(pid, opt_project)
	if err != nil {
		log.Fatal(err)
	}
	// create project models
	// project exist ?
	project_model := models.Project{}
	_ = db.Model(&models.Project{}).
		Select("count(*) > 0").
		Where("ID = ?", project.ID).
		Find(&exist).Error
	if exist { // project exist
		db.First(&project_model, "ID = ?", project.ID)
	} else { // non exist
		project_model.Name = project.Name
		project_model.Id = project.ID
		db.Create(&project_model)
	}

	//create branch models
	branch_model := models.Branch{}
	branch_model.Id = uuid.New()
	branch_model.Name = branch
	branch_model.Project_Id = project_model.Id
	db.Create(&branch_model)

	// create commit models
	for _, commit := range Commits {
		commit.Id = uuid.New()
		temp := commit.Id
		commit.Project_Id = project_model.Id
		commit.Branch_Id = branch_model.Id
		db.Create(&commit)
		//create file models
		for _, file := range commit.FileCommits {
			file.Commit_Id = temp
			file.File_Id = uuid.New()
			db.Create(&file)
		}
	}
}
func (database *Database) HandleDB(isUpdate bool, timeIn ...*time.Time) {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Panic(err)
	}
	if !isUpdate {
		database.DropTable()
		database.Init()
	}
	var list_proj_id []int
	list_sub_group := util.ListSubGroupsID(config.GROUP_ID)
	for _, group := range list_sub_group {
		list := util.ListGroupProjectsID(group)
		list_proj_id = append(list_proj_id, list...)
	}
	for _, pid := range list_proj_id {
		branches := util.ListBranchNameProject(pid)
		for _, branch := range branches {
			//fmt.Printf("Project: %d, Branch: %s\n", pid, branch)
			if isUpdate {
				database.UpdateDB(pid, branch, timeIn[0])
			} else {
				database.ImportDB(pid, branch)
			}
		}
	}
}
