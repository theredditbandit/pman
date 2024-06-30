package utils_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/charmbracelet/glamour"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	c "github.com/theredditbandit/pman/constants"
	"github.com/theredditbandit/pman/pkg/db"
	"github.com/theredditbandit/pman/pkg/utils"
)

const (
	dbname      = db.DBTestName
	expectedMsg = "Something went wrong"
	projectName = "project_name"
	projectPath = "./" + projectName
)

func Test_TitleCase(t *testing.T) {
	t.Run("Test TitleCase", func(t *testing.T) {
		s := "hello world"
		expected := "Hello World"

		actual := utils.TitleCase(s)

		assert.Equal(t, expected, actual)
	})
}

func Test_FilterByStatus(t *testing.T) {
	t.Run("Test FilterByStatus", func(t *testing.T) {
		data := map[string]string{
			"key1": "status1",
			"key2": "status2",
			"key3": "status1",
		}
		status := []string{"status1"}

		expectedData := map[string]string{
			"key1": "status1",
			"key3": "status1",
		}

		actualData := utils.FilterByStatuses(data, status)

		assert.Equal(t, expectedData, actualData)
	})

	t.Run("Test FilterByStatus with empty data", func(t *testing.T) {
		data := map[string]string{}
		status := []string{"status1"}

		expectedData := map[string]string{}

		actualData := utils.FilterByStatuses(data, status)

		assert.Equal(t, expectedData, actualData)
	})
}

func Test_BeautifyMD(t *testing.T) {
	t.Run("Test BeautifyMD under normal conditions", func(t *testing.T) {
		data := []byte("# i am a test")
		r, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(120),
			glamour.WithAutoStyle(),
		)
		require.NoError(t, err)

		expected, err := r.Render(string(data))
		require.NoError(t, err)

		actual, err := utils.BeautifyMD(data)

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Test BeautifyMD with empty data", func(t *testing.T) {
		data := []byte("")
		r, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(120),
			glamour.WithAutoStyle(),
		)
		require.NoError(t, err)

		expected, err := r.Render(string(data))
		require.NoError(t, err)

		actual, err := utils.BeautifyMD(data)

		require.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func Test_ReadREADME(t *testing.T) {
	t.Run("Test ReadREADME under normal conditions", func(t *testing.T) {
		expected := []byte{}

		t.Cleanup(func() {
			err := db.DeleteDb(dbname)
			require.NoError(t, err)
			_ = os.RemoveAll(projectPath)
		})

		err := os.Mkdir(projectPath, 0755)
		require.NoError(t, err)
		f, err := os.Create(projectPath + "/README.md")
		require.NoError(t, err)
		f.Close()

		err = db.WriteToDB(dbname, map[string]string{projectName: projectPath}, c.ProjectPaths)
		require.NoError(t, err)

		actual, err := utils.ReadREADME(dbname, projectName)

		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})
	t.Run("Test ReadREADME with well alias and project", func(t *testing.T) {
		alias := "project_alias"
		expected := []byte{}

		t.Cleanup(func() {
			err := db.DeleteDb(dbname)
			require.NoError(t, err)
			_ = os.RemoveAll(projectPath)
		})

		err := os.Mkdir(projectPath, 0755)
		require.NoError(t, err)
		f, err := os.Create(projectPath + "/README.md")
		require.NoError(t, err)
		f.Close()

		err = db.WriteToDB(dbname, map[string]string{alias: projectName}, c.ProjectAliasBucket)
		require.NoError(t, err)
		err = db.WriteToDB(dbname, map[string]string{projectName: projectPath}, c.ProjectPaths)
		require.NoError(t, err)

		actual, err := utils.ReadREADME(dbname, alias)

		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})
	t.Run("Test ReadREADME with well alias and bad project", func(t *testing.T) {
		alias := "project_alias"
		expected := []byte(nil)

		t.Cleanup(func() {
			err := db.DeleteDb(dbname)
			require.NoError(t, err)
			_ = os.RemoveAll(projectPath)
		})

		err := db.WriteToDB(dbname, map[string]string{alias: projectName}, c.ProjectAliasBucket)
		require.NoError(t, err)

		actual, err := utils.ReadREADME(dbname, alias)

		assert.Equal(t, expected, actual)
		require.ErrorIs(t, err, utils.ErrGetProject)
	})

	t.Run("Test ReadREADME with empty project name", func(t *testing.T) {
		projectName := ""
		expected := []byte(nil)

		actual, err := utils.ReadREADME(dbname, projectName)

		assert.Equal(t, expected, actual)
		require.ErrorIs(t, err, utils.ErrGetAlias)
	})

	t.Run("Test ReadREADME with invalid README file name", func(t *testing.T) {
		expected := []byte(fmt.Sprintf("# README does not exist for %s", projectName))

		t.Cleanup(func() {
			err := db.DeleteDb(dbname)
			require.NoError(t, err)
			_ = os.RemoveAll(projectPath)
		})

		err := os.Mkdir(projectPath, 0755)
		require.NoError(t, err)
		f, err := os.Create(projectPath + "/README.txt")
		require.NoError(t, err)
		f.Close()

		err = db.WriteToDB(dbname, map[string]string{projectName: projectPath}, c.ProjectPaths)
		require.NoError(t, err)

		actual, err := utils.ReadREADME(dbname, projectName)

		assert.Equal(t, expected, actual)
		require.NoError(t, err)
	})
}

func Test_GetLastModifiedTime(t *testing.T) {
	t.Run("Test GetLastModifiedTime under normal conditions: case Today modification", func(t *testing.T) {
		t.Cleanup(func() {
			err := db.DeleteDb(dbname)
			require.NoError(t, err)
			_ = os.RemoveAll(projectPath)
		})

		err := os.Mkdir(projectPath, 0755)
		require.NoError(t, err)
		f, err := os.Create(projectPath + "/README.md")
		require.NoError(t, err)
		fCreateTime := time.Now()
		correctModTime := fCreateTime.Format("02 Jan 06 15:04")
		err = os.Chtimes(projectPath+"/README.md", fCreateTime, fCreateTime)
		require.NoError(t, err)
		f.Close()

		err = db.WriteToDB(dbname, map[string]string{projectName: projectPath}, c.ProjectPaths)
		require.NoError(t, err)

		actual := utils.GetLastModifiedTime(dbname, projectName)

		assert.NotEqual(t, expectedMsg, actual)
		assert.Equal(t, correctModTime, actual)
		assert.NotEmpty(t, actual)
	})

	t.Run("Test GetLastModifiedTime under normal conditions: case Yesterday modification", func(t *testing.T) {
		t.Cleanup(func() {
			err := db.DeleteDb(dbname)
			require.NoError(t, err)
			_ = os.RemoveAll(projectPath)
		})

		err := os.Mkdir(projectPath, 0755)
		require.NoError(t, err)
		f, err := os.Create(projectPath + "/README.md")
		require.NoError(t, err)
		fCreateTime := time.Now().AddDate(0, 0, -1)
		correctModTime := fCreateTime.Format("02 Jan 06 15:04")
		err = os.Chtimes(projectPath+"/README.md", fCreateTime, fCreateTime)
		require.NoError(t, err)
		f.Close()

		err = db.WriteToDB(dbname, map[string]string{projectName: projectPath}, c.ProjectPaths)
		require.NoError(t, err)

		actual := utils.GetLastModifiedTime(dbname, projectName)

		assert.NotEqual(t, expectedMsg, actual)
		assert.Equal(t, correctModTime, actual)
		assert.NotEmpty(t, actual)
	})
	t.Run("Test GetLastModifiedTime under normal conditions: case old modification", func(t *testing.T) {
		t.Cleanup(func() {
			err := db.DeleteDb(dbname)
			require.NoError(t, err)
			_ = os.RemoveAll(projectPath)
		})

		err := os.Mkdir(projectPath, 0755)
		require.NoError(t, err)
		f, err := os.Create(projectPath + "/README.md")
		require.NoError(t, err)
		err = os.Chtimes(projectPath+"/README.md", time.Now().AddDate(0, 0, -5), time.Now().AddDate(0, 0, -5))
		require.NoError(t, err)
		f.Close()

		err = db.WriteToDB(dbname, map[string]string{projectName: projectPath}, c.ProjectPaths)
		require.NoError(t, err)

		actual := utils.GetLastModifiedTime(dbname, projectName)

		assert.NotEqual(t, expectedMsg, actual)
		assert.NotEmpty(t, actual)
	})
	t.Run("Test GetLastModifiedTime with invalid project", func(t *testing.T) {
		projectPath := "./invalid_project"

		actual := utils.GetLastModifiedTime(dbname, projectPath)

		assert.Equal(t, expectedMsg, actual)
	})
}
