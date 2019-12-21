package helpers

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMarkdown(t *testing.T) {
	testString := "# Hello"
	total := string(Markdown(testString))
	output := "<h1>Hello</h1>\n"

	if total != output {
		t.Errorf("Markdown function doesn't output the correct data")
	}
}

func TestGenerateRandomStringURLSafe(t *testing.T) {
	testString, err := GenerateRandomStringURLSafe(5)
	if err != nil {
		t.Errorf("GenerateRandomStringURLSafe failed, error: %s", err.Error())
	}
	if len(testString) != 8 {
		t.Errorf("GenerateRandomStringURLSafe not correct length, outputted length %d wanted %d", len(testString), 8)
	}

	anotherTestString, err := GenerateRandomStringURLSafe(5)
	if err != nil {
		t.Errorf("GenerateRandomStringURLSafe failed, error: %s", err.Error())
	}
	if testString == anotherTestString {
		t.Errorf("GenerateRandomStringURLSafe not correct output")
	}
}

func TestCreateDBHandler(t *testing.T) {
	db, err := CreateDBHandler()
	if err != nil {
		t.Errorf("CreateDBHandler failed. Error message: %s", err.Error())
	}
	_, err = db.Exec("CREATE DATABASE testDB")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully created database..")
	}

	_, err = db.Exec("DROP DATABASE testDB")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Successfully dropped database..")
	}
	db.Close()
}
