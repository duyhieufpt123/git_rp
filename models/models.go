package models

import (
	"time"

	"github.com/google/uuid"
)

//insert value

type FileCommit struct {
	File_Id   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Commit_Id uuid.UUID `gorm:"type:uuid"`
	Name      string    `gorm:"type:varchar(2048)"`
	Status    string    `gorm:"type:varchar(8)"`
}

type Commit struct {
	Id           uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Project_Id   int          `gorm:"type:int"`
	Branch_Id    uuid.UUID    `gorm:"type:uuid"`
	User_Id      uuid.UUID    `gorm:"type:uuid"`
	User_name    string       `gorm:"type:varchar(512)"`
	Message      string       `gorm:"type:varchar(512)"`
	Line_added   int          `gorm:"type:int"`
	Line_removed int          `gorm:"type:int"`
	CommitDate   *time.Time   `gorm:"type:timestamp with time zone"`
	CommitID     string       `gorm:"type:varchar(64)"`
	FileCommits  []FileCommit `gorm:"foreignkey:Commit_Id;references:Id"`
}

type Branch struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name       string    `gorm:"varchar(512)"`
	Project_Id int       `gorm:"type:int"`
	Commits    []Commit  `gorm:"foreignkey:Branch_Id;references:Id"`
}

type Project struct {
	Id      int      `gorm:"type:int;primaryKey;"`
	Name    string   `gorm:"type:varchar(512)"`
	Commits []Commit `gorm:"foreignkey:Project_Id;references:Id"`
	Branchs []Branch `gorm:"foreignkey:Project_Id;references:Id"`
}

type User struct {
	Id     uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name   string    `gorm:"varchar(512)"`
	Email  string    `gorm:"varchar(512)"`
	Commit Commit    `gorm:"foreignkey:User_Id;references:Id" `
}
