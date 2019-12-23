package initiate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/OpenCoursePlatform/OpenCoursePlatform-Go/helpers"
)

/*
Tests should be used for testing functions that need a database.
It creates a database and adds tables.
*/
func Tests() (*sql.DB, error) {
	db, err := helpers.CreateDBHandlerWithDB("")
	if err != nil {
		fmt.Println("CreateDBHandler failed. Error message:", err.Error())
		return nil, err
	}
	_, err = db.Exec("CREATE DATABASE testDB")
	if err != nil {
		fmt.Println("Creation of database failed. Error message:", err.Error())
		return nil, err
	}
	fmt.Println("Successfully created database..")
	db.Close()

	db, err = helpers.CreateDBHandlerWithDB("testDB")
	if err != nil {
		fmt.Println("CreateDBHandler failed. Error message:", err.Error())
		return nil, err
	}

	err = CreateTables(db)
	if err != nil {
		fmt.Println("Creation of user table failed. Error message:", err.Error())
		return nil, err
	}
	fmt.Println("Successfully created user table..")
	return db, nil
}

/*
FinishTests should be used for testing functions that need a database.
It drops the database used in 'Tests' function.
*/
func FinishTests(db *sql.DB) error {
	_, err := db.Exec("DROP DATABASE testDB")
	if err != nil {
		fmt.Println("Dropping of database failed. Error message:", err.Error())
		return err
	}
	fmt.Println("Successfully dropped database..")

	defer db.Close()
	return nil
}

// CreateDatabase creates a database with the name sent in
func CreateDatabase(db *sql.DB, databaseName string) error {
	_, err := db.Exec(`CREATE DATABASE ` + databaseName)
	if err != nil {
		return err
	}
	return nil
}

// CreateAccessLogs creates access_logs table
func CreateAccessLogs(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE access_logs (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		time datetime NOT NULL,
		method varchar(255) NOT NULL,
		path text NOT NULL,
		remote_address varchar(255) NOT NULL,
		user_agent text NOT NULL,
		speed bigint(20) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateBlogPosts creates blog_posts table
func CreateBlogPosts(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE blog_posts (
		id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
		title text NOT NULL,
		text text NOT NULL,
		slug varchar(255) NOT NULL,
		published timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE KEY slug (slug)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateCourseCategories creates course_categories table
func CreateCourseCategories(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE course_categories (
		id int(11) NOT NULL AUTO_INCREMENT,
		name text NOT NULL,
		slug varchar(255) NOT NULL,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY slug (slug),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateCourses creates courses table
func CreateCourses(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE courses (
		id int(11) NOT NULL AUTO_INCREMENT,
		name text NOT NULL,
		description text NOT NULL,
		slug varchar(255) NOT NULL,
		category_id int(11) NOT NULL,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY slug (slug),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateEmailVerification creates email_verification table
func CreateEmailVerification(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE email_verification (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		user_id int(10) unsigned NOT NULL,
		token varchar(255) NOT NULL,
		verified_at datetime NULL,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE KEY token (token)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateFooter creates footer table
func CreateFooter(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE footer (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		page_id int(11) unsigned NOT NULL,
		footer_category_id int(10) unsigned NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateFooterCategories creates footer_categories table
func CreateFooterCategories(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE footer_categories (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		slug varchar(255) NOT NULL,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY slug (slug),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateForumPosts creates forum_posts table
func CreateForumPosts(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE forum_posts (
		id int(11) NOT NULL AUTO_INCREMENT,
		text text NOT NULL,
		topic_id int(11) NOT NULL,
		author_id int(11) NOT NULL,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateForumTopics creates forum_topics table
func CreateForumTopics(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE forum_topics (
		id int(11) NOT NULL AUTO_INCREMENT,
		title text NOT NULL,
		slug varchar(255) NOT NULL,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY slug (slug),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateGroupMembers creates group_members table
func CreateGroupMembers(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE group_members (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		user_id int(11) NOT NULL,
		group_id int(11) NOT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY user_id (user_id,group_id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateGroupPermissions creates group_permissions table
func CreateGroupPermissions(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE group_permissions (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		group_id int(11) NULL,
		permission_id int(11) NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateModule creates module table
func CreateModule(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE module (
		id int(11) NOT NULL AUTO_INCREMENT,
		name text NOT NULL,
		description text NOT NULL,
		image_link varchar(255) NOT NULL,
		slug varchar(255) NOT NULL,
		course_id int(11) NOT NULL,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY slug (slug),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreatePages creates pages table
func CreatePages(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE pages (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		title varchar(255) NOT NULL,
		content text CHARACTER SET utf8mb4,
		slug varchar(255),
		UNIQUE KEY slug (slug),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreatePermissions creates permission table
func CreatePermissions(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE permission (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreatePodcastPost creates podcast_posts table
func CreatePodcastPost(db *sql.DB) error {
	_, err := db.Exec(`
	  CREATE TABLE podcast_posts (
		id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
		title text NOT NULL,
		text text NOT NULL,
		file varchar(255) NOT NULL,
		slug varchar(255) NOT NULL,
		published timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id),
		UNIQUE KEY slug (slug),
		UNIQUE KEY file (file)
		) ENGINE=InnoDB CHARSET=utf8mb4;
		`)
	if err != nil {
		return err
	}
	return nil
}

// CreateSession creates session table
func CreateSession(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE session (
		id int(11) NOT NULL AUTO_INCREMENT,
		name text NOT NULL,
		slug varchar(255) NOT NULL,
		module_id int(11) NOT NULL,
		session_type int(11) NOT NULL,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY slug (slug),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateSessionText creates session_text table
func CreateSessionText(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE session_text (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		text text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
		session_id int(10) unsigned NOT NULL,
		created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
		) ENGINE=InnoDB CHARSET=utf8mb4;
		`)
	if err != nil {
		return err
	}
	return nil
}

// CreateSessionYoutube creates session_youtube table
func CreateSessionYoutube(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE session_youtube (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		text text NOT NULL,
		youtube_id varchar(255) NOT NULL,
		session_id int(10) unsigned NOT NULL,
		created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
		) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateSettings creates settings table and inserts basic settings
func CreateSettings(db *sql.DB) error {
	_, err := db.Exec(`
	  CREATE TABLE settings (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		option_name varchar(255) NOT NULL,
		option_value text NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	INSERT INTO settings (id, option_name, option_value)
	VALUES
		(1, 'SiteUrl', ''),
		(2, 'Name', ''),
		(3, 'Description', ''),
		(4, 'IndexTitle', ''),
		(5, 'IndexDescription', ''),
		(6, 'RedisEnabled', '0'),
		(7, 'MailSender', ''),
		(8, 'PostmarkEnabled', '0'),
		(9, 'PostmarkToken', ''),
		(10, 'RegistrationMailSubject', ''),
		(11, 'BlogActivated', '1'),
		(12, 'ForumActivated', '1'),
		(13, 'PodcastActivated', '1');
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateTicketResponses creates ticket_responses table
func CreateTicketResponses(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE ticket_responses (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		ticket_id int(10) unsigned NOT NULL,
		user_id int(10) unsigned NOT NULL,
		text text NOT NULL,
		created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateTickets creates tickets table
func CreateTickets(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE tickets (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		user_id int(11) unsigned NOT NULL,
		topic varchar(255) NOT NULL,
		slug varchar(255) NOT NULL,
		solved tinyint(1) NOT NULL DEFAULT 0,
		created datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY slug (slug),
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateToolbar creates toolbar table
func CreateToolbar(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE toolbar (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		page_id int(11) unsigned NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateUserGroups creates user_groups table
func CreateUserGroups(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE user_groups (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateUserPermissions creates user_permissions table
func CreateUserPermissions(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE user_permissions (
		id int(11) unsigned NOT NULL AUTO_INCREMENT,
		user_id int(11) NOT NULL,
		permission_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

// CreateUsers creates users table
func CreateUsers(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE users (
		id int(11) NOT NULL AUTO_INCREMENT,
		username text NOT NULL,
		email text NOT NULL,
		password varchar(255) NOT NULL,
		verified tinyint(1) NOT NULL DEFAULT 0,
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB CHARSET=utf8mb4;
	`)
	if err != nil {
		return err
	}
	return nil
}

/*
CreateUserRelatedTables creates user related tables
*/
func CreateUserRelatedTables(db *sql.DB) error {
	err := CreateGroupMembers(db)
	if err != nil {
		return err
	}
	err = CreateGroupPermissions(db)
	if err != nil {
		return err
	}
	err = CreateUserGroups(db)
	if err != nil {
		return err
	}
	err = CreateUserPermissions(db)
	if err != nil {
		return err
	}
	err = CreateUsers(db)
	if err != nil {
		return err
	}
	err = CreateEmailVerification(db)
	if err != nil {
		return err
	}
	return nil
}

/*
CreateCourseRelatedTables creates course related tables
*/
func CreateCourseRelatedTables(db *sql.DB) error {
	err := CreateCourses(db)
	if err != nil {
		return err
	}
	err = CreateCourseCategories(db)
	if err != nil {
		return err
	}
	err = CreateModule(db)
	if err != nil {
		return err
	}
	err = CreateSession(db)
	if err != nil {
		return err
	}
	err = CreateSessionText(db)
	if err != nil {
		return err
	}
	err = CreateSessionYoutube(db)
	if err != nil {
		return err
	}
	return nil
}

/*
CreateEssentials creates essential tables for running the application
*/
func CreateEssentials(db *sql.DB) error {
	err := CreateUserRelatedTables(db)
	if err != nil {
		return err
	}
	err = CreateCourseRelatedTables(db)
	if err != nil {
		return err
	}
	err = CreateAccessLogs(db)
	if err != nil {
		return err
	}
	err = CreateFooter(db)
	if err != nil {
		return err
	}
	err = CreateFooterCategories(db)
	if err != nil {
		return err
	}
	err = CreatePages(db)
	if err != nil {
		return err
	}
	err = CreatePermissions(db)
	if err != nil {
		return err
	}
	err = CreateSettings(db)
	if err != nil {
		return err
	}
	err = CreateToolbar(db)
	if err != nil {
		return err
	}
	return nil
}

/*
CreateTables creates all tables needed to run the application.
*/
func CreateTables(db *sql.DB) error {
	err := CreateEssentials(db)
	if err != nil {
		return err
	}
	err = CreateBlogPosts(db)
	if err != nil {
		return err
	}
	err = CreateForumPosts(db)
	if err != nil {
		return err
	}
	err = CreateForumTopics(db)
	if err != nil {
		return err
	}
	err = CreatePodcastPost(db)
	if err != nil {
		return err
	}
	err = CreateTicketResponses(db)
	if err != nil {
		return err
	}
	err = CreateTickets(db)
	if err != nil {
		return err
	}
	return nil
}

/*
GetSettingsName ...
*/
func GetSettingsName() (string, error) {
	var config helpers.Configuration
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadFile(path + "/" + helpers.SettingsFileName)
	if err != nil {
		data, err = ioutil.ReadFile(filepath.Dir(path) + "/" + helpers.SettingsFileName)
		if err != nil {
			return "", err
		}
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return "", err
	}
	return config.DatabaseName, nil
}

/*
UpdateDatabaseName updates the name of the database in the settings json file.
*/
func UpdateDatabaseName(input string) error {
	var config helpers.Configuration
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(path + "/" + helpers.SettingsFileName)
	if err != nil {
		data, err = ioutil.ReadFile(filepath.Dir(path) + "/" + helpers.SettingsFileName)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &config)
		if err != nil {
			return err
		}
		config.DatabaseName = input
		jsonData, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Dir(path)+"/"+helpers.SettingsFileName, jsonData, 0644)
		if err != nil {
			return err
		}
	} else {
		err = json.Unmarshal(data, &config)
		if err != nil {
			return err
		}
		config.DatabaseName = input
		jsonData, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path+"/"+helpers.SettingsFileName, jsonData, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
UpdateSettingToTestDB Updates the JSON Settings to the testDB.
*/
func UpdateSettingToTestDB() (string, error) {
	name, err := GetSettingsName()
	if err != nil {
		return "", err
	}
	err = UpdateDatabaseName("testDB")
	if err != nil {
		return "", err
	}
	return name, nil
}

/*
DeleteSettingsFile Deletes the JSON Settings file and returns it.
This function is mainly for testing purposes.
*/
func DeleteSettingsFile() (helpers.Configuration, error) {
	var config helpers.Configuration
	path, err := os.Getwd()
	if err != nil {
		return config, err
	}
	data, err := ioutil.ReadFile(path + "/" + helpers.SettingsFileName)
	if err != nil {
		data, err = ioutil.ReadFile(filepath.Dir(path) + "/" + helpers.SettingsFileName)
		if err != nil {
			return config, err
		}
		err = json.Unmarshal(data, &config)
		if err != nil {
			return config, err
		}
		err = ioutil.WriteFile(filepath.Dir(path)+"/"+helpers.SettingsFileName, []byte(""), 0644)
		if err != nil {
			return config, err
		}
	} else {
		err = json.Unmarshal(data, &config)
		if err != nil {
			return config, err
		}
		err = ioutil.WriteFile(path+"/"+helpers.SettingsFileName, []byte(""), 0644)
		if err != nil {
			return config, err
		}
	}
	return config, nil
}

/*
WriteSettingsFile writes the JSON Settings file.
*/
func WriteSettingsFile(config helpers.Configuration) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = ioutil.ReadFile(path + "/" + helpers.SettingsFileName)
	if err != nil {
		_, err = ioutil.ReadFile(filepath.Dir(path) + "/" + helpers.SettingsFileName)
		if err != nil {
			return err
		}
		jsonData, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Dir(path)+"/"+helpers.SettingsFileName, jsonData, 0644)
		if err != nil {
			return err
		}
	} else {
		jsonData, err := json.MarshalIndent(config, "", "\t")
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path+"/"+helpers.SettingsFileName, jsonData, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
