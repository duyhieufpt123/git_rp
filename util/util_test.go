package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_util_GetDate(t *testing.T) {
	expected := time.Now().Unix()
	res := GetDate(0, 0, 0).Unix()
	assert.Equal(t, expected, res)

}
func Test_util_ListBranchNameProject(t *testing.T) {
	expected := []string{"main"}
	res := ListBranchNameProject(176)
	assert.Equal(t, expected, res)

}

func Test_util_ListSubGroupsID(t *testing.T) {
	expected := []int{98, 99, 106, 88, 94, 111}
	res := ListSubGroupsID(83)
	assert.Equal(t, expected, res)

}

func Test_util_ListGroupProjectsID(t *testing.T) {
	expected := []int{154, 146, 145}
	res := ListGroupProjectsID(83)
	assert.Equal(t, expected, res)

}

func Test_util_LoadConfig1(t *testing.T) {
	_, err := LoadConfig(".")
	assert.NoError(t, err)
}
func Test_util_LoadConfig2(t *testing.T) {
	config, _ := LoadConfig(".")
	assert.Equal(t, "http://gitlab.itd.com.vn/api/v4", config.SERVER_GIT)
}
