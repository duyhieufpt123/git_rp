package database

import (
	"gitlabapi/models"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// setup
	configFile, err := os.Create("config.env")
	if err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}
	defer os.Remove(configFile.Name())
	defer configFile.Close()

	_, err = configFile.WriteString("DBSource=user=postgres dbname=database sslmode=disable")
	if err != nil {
		t.Fatalf("failed to write to config file: %v", err)
	}

	database := &Database{}
	database.Init()

	// test that tables were created
	db := database.Connect()
	defer db.Close()

	var count int
	err = db.Model(&models.Project{}).Count(&count).Error
	if err != nil {
		t.Fatalf("failed to count projects: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 projects, got %d", count)
	}

	err = db.Model(&models.User{}).Count(&count).Error
	if err != nil {
		t.Fatalf("failed to count users: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 users, got %d", count)
	}

	err = db.Model(&models.Branch{}).Count(&count).Error
	if err != nil {
		t.Fatalf("failed to count branches: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 branches, got %d", count)
	}

	err = db.Model(&models.Commit{}).Count(&count).Error
	if err != nil {
		t.Fatalf("failed to count commits: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 commits, got %d", count)
	}

	err = db.Model(&models.FileCommit{}).Count(&count).Error
	if err != nil {
		t.Fatalf("failed to count file commits: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 file commits, got %d", count)
	}
}
