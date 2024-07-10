package main

import (
	"bufio"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
	_ "modernc.org/sqlite"
)

type Job struct {
	CompanyName string
	Rating      string
	Notes       string
	HasAnswered bool
}

type ListJob struct {
	ID          int
	CompanyName string
	Rating      string
	Notes       string
	HasAnswered bool
}

var db *sql.DB

//go:embed ascii.txt
var ascii string

const YELLOW = "\033[33;1m"
const RESET = "\033[0m"

func dotenv(key string) string {
	err := godotenv.Load(".env")
	checkErr(err)
	return os.Getenv(key)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func mainScreen() {
	fmt.Println(YELLOW + ascii + RESET)
}

func InsertJob(db *sql.DB, CompanyName string, Rating string, Notes string, HasAnswered bool) {
	insertSQL := `INSERT INTO jobs (company_name, rating, notes, has_answered) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertSQL)
	if err != nil {
		log.Fatal(err)
	}
	defer statement.Close()

	_, err = statement.Exec(CompanyName, Rating, Notes, HasAnswered)
	if err != nil {
		log.Fatal(err)
	}
}

func ListJobs(db *sql.DB) {
	rows, err := db.Query("SELECT id, company_name, rating, notes, has_answered FROM jobs ORDER BY has_answered DESC")
	checkErr(err)
	defer rows.Close()

	var jobs []ListJob
	for rows.Next() {
		var job ListJob
		err := rows.Scan(&job.ID, &job.CompanyName, &job.Rating, &job.Notes, &job.HasAnswered)
		checkErr(err)
		jobs = append(jobs, job)
	}

	checkErr(rows.Err())
	printJobsTable(jobs)
}

func UpdateJob(db *sql.DB, ID string) {
	updateSQL := `UPDATE jobs SET has_answered = true WHERE ID = ?`
	statement, err := db.Prepare(updateSQL)
	checkErr(err)
	defer statement.Close()

	id, err := strconv.Atoi(ID)
	checkErr(err)
	_, err = statement.Exec(id)
	checkErr(err)

}

func printJobsTable(jobs []ListJob) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Company Name", "Rating", "Notes", "Has Answered"})

	for _, job := range jobs {
		hasAnsweredStr := "false"
		if job.HasAnswered {
			hasAnsweredStr = "true"
		}
		var jobID string = strconv.Itoa(job.ID)
		table.Append([]string{jobID, job.CompanyName, job.Rating, job.Notes, hasAnsweredStr})
	}
	table.Render()
}

func createTable(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS jobs(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		company_name TEXT,
		rating TEXT,
		notes TEXT,
		offer_link TEXT,
		review_link TEXT,
		has_answered BOOLEAN
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func printHelpMenu() {
  fmt.Println("-a: Add a new job application")
  fmt.Println("\t-cn: Company Name")
	fmt.Println("\t-r: Rating")
  fmt.Println("\t-n: Notes. Short description of the company")
  fmt.Println("\t-ha: Has Answered. Boolean value (true|false). Used to track if you have gotten a response from the company.")
  fmt.Println("Example: jm -a -cn Company Name -r 5/5 -n Great company -ha false")
  fmt.Println("-ls: List all job applications")
  fmt.Println("-u: Update a job application")
  fmt.Println("clear: Clear the screen")
  fmt.Println("--help: Display this help message")
}

func main() {
	command_prefix := "jm -"
	RATING_PATTERN := `^(?:[0-4](?:\.\d+)?|5(?:\.0+)?)/5$`
	NOTES_PATTERN := `^.{1,30}$`
	rating_pattern_compiled := regexp.MustCompile(RATING_PATTERN)
	notes_pattern_compiled := regexp.MustCompile(NOTES_PATTERN)

	db, err := sql.Open("sqlite", "jobs.db")
	checkErr(err)
	defer db.Close()

	createTable(db)
	clearScreen()
	mainScreen()

	reader := bufio.NewReader(os.Stdin)
mainLoop:
	for {
		fmt.Print(YELLOW + "JobManager$> " + RESET)
		input, err := reader.ReadString('\n')
		checkErr(err)
		input = strings.TrimSpace(input)

		if strings.HasPrefix(input, command_prefix+"a") {
			cmdParts := strings.Split(input, "-")

			var companyName, rating, notes string
			var hasAnswered bool

			for i := 1; i < len(cmdParts); i++ {
				part := strings.TrimSpace(cmdParts[i])
				switch {
				case strings.HasPrefix(part, "cn "):
					companyName = strings.TrimSpace(strings.TrimPrefix(part, "cn "))
				case strings.HasPrefix(part, "r "):
					rating = strings.TrimSpace(strings.TrimPrefix(part, "r "))
					if !rating_pattern_compiled.MatchString(rating) {
						fmt.Println("Invalid rating. Please use the format x/5 where x is a number between 0 and 5")
						continue mainLoop
					}
				case strings.HasPrefix(part, "n "):
					notes = strings.TrimSpace(strings.TrimPrefix(part, "n "))
					if !notes_pattern_compiled.MatchString(notes) {
						fmt.Println("Invalid notes. Please use a short description of the company (Max 30 characters)")
						continue mainLoop
					}
				case strings.HasPrefix(part, "ha "):
					hasAnswered = strings.TrimSpace(strings.TrimPrefix(part, "ha ")) == "true"
				}
			}

			InsertJob(db, companyName, rating, notes, hasAnswered)
			fmt.Println("Inserted job application successfully!")
		} else if strings.HasPrefix(input, command_prefix+"ls") {
			args := strings.Fields(input)
			if len(args) != 2 {
				fmt.Println("Usage: jm -l")
				continue
			}
			ListJobs(db)
		} else if strings.HasPrefix(input, command_prefix+"u") {
			ListJobs(db)
			fmt.Println("Which application would you like to update: ")
			companyID, _ := reader.ReadString('\n')
			companyID = strings.TrimSpace(companyID)

			fmt.Println("Updating application:", companyID)

			UpdateJob(db, companyID)
			fmt.Println("Updated job application successfully!")
		} else if strings.HasPrefix(input, command_prefix+"-help") {
      printHelpMenu()	
		} else if strings.HasPrefix(input, "clear") {
			clearScreen()
			mainScreen()
		} else {
			fmt.Println("Unknown command. Usage: jm [-a|-ls|-u] [-cn|-r|-n|-ha] [true|false]")
		}
	}
}


