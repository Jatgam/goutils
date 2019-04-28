package goutils

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

type (
	updateFileModifiedTimeTest struct {
		filename string
		time     time.Time
		create   bool
	}
	dirSizeMBTest struct {
		files []filesToCreate
	}
	deleteOldestFileUntilDirUnderSizeTest struct {
		files []filesToCreate
		limit int64
	}
	filesToCreate struct {
		filename string
		size     int64
		time     time.Time
	}
)

var (
	testDirPath = "build/tests/files"

	updateFileModifiedTimeTestData = []updateFileModifiedTimeTest{
		{"f1.tmp", time.Now().AddDate(-3, 0, 0), true},
		{"f12.tmp", time.Now().AddDate(-17, -23, -154), true},
		{"ha.tmp", time.Now().AddDate(-17, -23, -154), false},
	}

	dirSizeMBTestData = []dirSizeMBTest{
		{[]filesToCreate{
			{filename: "f1.tmp", size: int64(1e7)},
			{filename: "f2.tmp", size: int64(2e7)},
			{filename: "f3.tmp", size: int64(3e7)},
		}},
		{[]filesToCreate{
			{filename: "f1.tmp", size: int64(1e7)},
			{filename: "f2.tmp", size: int64(2e7)},
			{filename: "f3.tmp", size: int64(3e7)},
			{filename: "f4.tmp", size: int64(3e7)},
			{filename: "f5.tmp", size: int64(3e7)},
			{filename: "f6.tmp", size: int64(3e7)},
			{filename: "f7.tmp", size: int64(3e7)},
		}},
	}

	deleteOldestFileUntilDirUnderSizeTestData = []deleteOldestFileUntilDirUnderSizeTest{
		{files: []filesToCreate{
			{filename: "f1.tmp", size: int64(1e7), time: time.Now()},
			{filename: "f2.tmp", size: int64(2e7), time: time.Now().AddDate(-1, 0, 0)},
			{filename: "f3.tmp", size: int64(3e7), time: time.Now().AddDate(-2, 0, 0)},
		}, limit: 10},
		{files: []filesToCreate{
			{filename: "f1.tmp", size: int64(1e7), time: time.Now()},
			{filename: "f2.tmp", size: int64(2e7), time: time.Now().AddDate(-1, 0, 0)},
			{filename: "f3.tmp", size: int64(3e7), time: time.Now().AddDate(-2, 0, 0)},
			{filename: "f4.tmp", size: int64(1e7), time: time.Now()},
			{filename: "f5.tmp", size: int64(2e7), time: time.Now().AddDate(-1, -1, -1)},
			{filename: "f6.tmp", size: int64(3e7), time: time.Now().AddDate(-2, -1, -1)},
		}, limit: 30},
	}
)

func TestUpdateFileModifiedTime(t *testing.T) {
	thisTestDir := path.Join(testDirPath, "UpdateFileModifiedTime")
	err := os.MkdirAll(thisTestDir, 0755)
	if err != nil {
		t.Errorf("Failed to create test directory")
	}
	for _, testData := range updateFileModifiedTimeTestData {
		fileForTest := path.Join(thisTestDir, testData.filename)
		if testData.create {
			f, err := os.Create(fileForTest)
			if err != nil {
				t.Errorf("Failed to create file for test: %s", fileForTest)
			}
			f.Close()
			UpdateFileModifiedTime(path.Join(thisTestDir, testData.filename), testData.time)
			createdFileInfo, err := os.Stat(path.Join(thisTestDir, testData.filename))
			if err != nil {
				t.Errorf("Failed to read created file in UpdateFileModifiedTime test: %s", path.Join(thisTestDir, testData.filename))
			}
			if createdFileInfo.ModTime() != testData.time {
				t.Error(
					"Created File: ", testData.filename,
					"FileInfo Time: ", createdFileInfo.ModTime(),
					"Time Supposedly Set: ", testData.time,
				)
			}
		} else {
			testErr := UpdateFileModifiedTime(path.Join(thisTestDir, testData.filename), testData.time)
			if testErr.Error() != fmt.Sprintf("File doesn't exist: %s", path.Join(thisTestDir, testData.filename)) {
				t.Errorf("File doesn't exist, yet modified time worked: %s", path.Join(thisTestDir, testData.filename))
			}
		}
	}
}

func TestDirSizeMB(t *testing.T) {
	thisTestDir := path.Join(testDirPath, "TestDirSizeMB")
	err := os.MkdirAll(thisTestDir, 0755)
	if err != nil {
		t.Errorf("Failed to create test directory")
	}
	for _, testIteration := range dirSizeMBTestData {
		var totalBytes int64
		for _, testData := range testIteration.files {
			fileForTest := path.Join(thisTestDir, testData.filename)
			f, err := os.Create(fileForTest)
			defer f.Close()
			if err != nil {
				t.Errorf("Failed to create file for test: %s", fileForTest)
			}
			truncErr := f.Truncate(testData.size)
			if truncErr != nil {
				t.Errorf("Failed to create test file %s with size %v", testData.filename, testData.size)
			}
			totalBytes += testData.size
		}
		floatDirMB, err := DirSizeMB(thisTestDir)
		if err != nil {
			t.Errorf("Failed to determine test dir size: %s", err)
		}
		totalMB := float64(totalBytes) / 1024.0 / 1024.0
		if floatDirMB != totalMB {
			t.Error(
				"Got MB: ", floatDirMB,
				"Expected MB: ", totalMB,
			)
		}
	}
}

func TestDeleteOldestFileUntilDirUnderSize(t *testing.T) {
	thisTestDir := path.Join(testDirPath, "DeleteOldestFileUntilDirUnderSize")
	err := os.MkdirAll(thisTestDir, 0755)
	if err != nil {
		t.Errorf("Failed to create test directory")
	}
	for _, testIteration := range deleteOldestFileUntilDirUnderSizeTestData {
		for _, testData := range testIteration.files {
			fileForTest := path.Join(thisTestDir, testData.filename)
			f, err := os.Create(fileForTest)
			if err != nil {
				t.Errorf("Failed to create file for test: %s", fileForTest)
			}
			truncErr := f.Truncate(testData.size)
			if truncErr != nil {
				t.Errorf("Failed to create test file %s with size %v", testData.filename, testData.size)
			}
			f.Close()
			UpdateFileModifiedTime(path.Join(thisTestDir, testData.filename), testData.time)
			_, err = os.Stat(path.Join(thisTestDir, testData.filename))
			if err != nil {
				t.Errorf("Failed to read created file in TestDeleteOldestFileUntilDirUnderSize: %s", path.Join(thisTestDir, testData.filename))
			}
		}
		beforeDeleteDirSize, err := DirSizeMB(thisTestDir)
		if err != nil {
			t.Errorf("Failed to determine dir size before delete test: %s", err)
		}
		delErr := DeleteOldestFileUntilDirUnderSize(thisTestDir, float64(testIteration.limit))
		if delErr != nil {
			t.Errorf("Failed to delete any files: %s", delErr)
		}
		afterDeleteDirSize, err := DirSizeMB(thisTestDir)
		if err != nil {
			t.Errorf("Failed to determine dir size after delete test: %s", err)
		}
		if beforeDeleteDirSize <= afterDeleteDirSize {
			t.Error(
				"Before Delete: ", beforeDeleteDirSize,
				"After Delete: ", afterDeleteDirSize,
			)
		}
	}
}

func TestMain(m *testing.M) {
	setupTests()
	retCode := m.Run()
	cleanupTests()
	os.Exit(retCode)
}

func setupTests() {
	err := os.MkdirAll(testDirPath, 0755)
	if err != nil {
		fmt.Println("Failed to setup tests")
		os.Exit(1)
	}
}

func cleanupTests() {
	os.RemoveAll(testDirPath)
}
