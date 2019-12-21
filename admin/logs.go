package admin

import (
	"bufio"
	"database/sql"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
)

/*
InsertAccessLog ...
*/
func InsertAccessLog(db *sql.DB, accessTime time.Time, method, path, remoteAddress, userAgent string, speed int) error {
	_, err := db.Exec(`INSERT INTO access_logs (time, method, path, remote_address, user_agent, speed) VALUES (?, ?, ?, ?, ?, ?)`, accessTime, method, path, remoteAddress, userAgent, speed)
	if err != nil {
		return err
	}
	return nil
}

func parseLogString(logString string) {
	stringParts := strings.Split(logString, "*")
	db, err := helpers.CreateDBHandler()
	if err != nil {
		helpers.HandleError(err)
	}

	defer db.Close()

	accessTime, err := time.Parse(time.RFC1123, stringParts[0])
	if err != nil {
		helpers.HandleError(err)
	}
	performance, err := strconv.Atoi(stringParts[5])
	if err != nil {
		helpers.HandleError(err)
	}

	err = InsertAccessLog(db, accessTime, stringParts[1], stringParts[2], stringParts[3], stringParts[4], performance)
	if err != nil {
		helpers.HandleError(err)
	}
}

func parseLog() {
	random, err := helpers.GenerateRandomStringURLSafe(10)
	if err != nil {
		helpers.HandleError(err)
	}
	os.Rename("access.log", "access."+random+".log")
	file, err := os.Open("access." + random + ".log")
	defer file.Close()
	if err != nil {
		helpers.HandleError(err)
	} else {
		f, err := os.OpenFile("access.old.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			helpers.HandleError(err)
		} else {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				_, err = f.WriteString(scanner.Text() + "\n")
				if err != nil {
					helpers.HandleError(err)
				}
				parseLogString(scanner.Text())
			}
			f.Close()
			if err := scanner.Err(); err != nil {
				helpers.HandleError(err)
			}
			err = os.Remove("access." + random + ".log")
			if err != nil {
				helpers.HandleError(err)
			}
		}
	}
}
