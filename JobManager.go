package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Job struct {
	CompanyName string
	Rating      string
	Notes       string
	HasAnswered bool
}

var client *mongo.Client

//go:embed ascii.txt
var ascii string

func dotenv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func init() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dotenv("MONGO_DB_URI")).SetServerAPIOptions(serverAPI)
	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	fmt.Println("Successfully connected to MongoDB")
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func mainScreen() {
	fmt.Println(ascii)
}

func InsertJob(CompanyName string, Rating string, Notes string, OfferLink string, ReviewLink string, HasAnswered bool) {
	coll := client.Database(dotenv("MONGO_DB_DATABASE")).Collection(dotenv("MONGO_DB_COLLECTION"))
	docs := []interface{}{
		Job{CompanyName: CompanyName, Rating: Rating, Notes: Notes, HasAnswered: HasAnswered},
	}
	coll.InsertMany(context.TODO(), docs)
}

func ListJobs() {
	coll := client.Database("Jobs").Collection("Jobs")
	opts := options.Find().SetSort(bson.D{{"hasanswered", -1}})
	cursor, err := coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Fatalf("Failed to retrieve jobs: %v", err)
	}
	defer cursor.Close(context.TODO())

	var jobs []Job
	if err := cursor.All(context.TODO(), &jobs); err != nil {
		log.Fatalf("Failed to decode all jobs: %v", err)
	}

	printJobsTable(jobs)
}

func UpdateJob(companyName string) {
	coll := client.Database("Jobs").Collection("Jobs")

	filter := bson.D{{"companyname", companyName}}
	update := bson.D{{"$set", bson.D{{"hasanswered", true}}}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatalf("Failed to update job: %v", err)
	}
}

func printJobsTable(jobs []Job) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Company Name", "Rating", "Notes", "Has Answered"})

	for _, job := range jobs {
		hasAnsweredStr := "false"
		if job.HasAnswered {
			hasAnsweredStr = "true"
		}
		table.Append([]string{job.CompanyName, job.Rating, job.Notes, hasAnsweredStr})
	}

	table.Render()
}

func main() {
LOOP:
	clearScreen()
	mainScreen()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("JobManager$> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		input = strings.TrimSpace(input)

		if strings.HasPrefix(input, "jm -a") {
			cmdParts := strings.Split(input, "-")

			var companyName, rating, notes, offerLink, reviewLink string
			var hasAnswered bool

			for i := 1; i < len(cmdParts); i++ {
				part := strings.TrimSpace(cmdParts[i])
				switch {
				case strings.HasPrefix(part, "cn "):
					companyName = strings.TrimSpace(strings.TrimPrefix(part, "cn "))
				case strings.HasPrefix(part, "r "):
					rating = strings.TrimSpace(strings.TrimPrefix(part, "r "))
				case strings.HasPrefix(part, "n "):
					notes = strings.TrimSpace(strings.TrimPrefix(part, "n "))
				case strings.HasPrefix(part, "ha "):
					hasAnswered = strings.TrimSpace(strings.TrimPrefix(part, "ha ")) == "true"
				}
			}

			InsertJob(companyName, rating, notes, offerLink, reviewLink, hasAnswered)
			fmt.Println("Inserted job application successfully!")
		} else if strings.HasPrefix(input, "jm -ls") {
			args := strings.Fields(input)
			if len(args) != 2 {
				fmt.Println("Usage: jm -l")
				continue
			}
			ListJobs()
		} else if strings.HasPrefix(input, "jm -u") {
			ListJobs()
			fmt.Println("Which application would you like to update: ")
			companyNameToUpdate, _ := reader.ReadString('\n')
			companyNameToUpdate = strings.TrimSpace(companyNameToUpdate)

			fmt.Println("Updating application:", companyNameToUpdate)

			UpdateJob(companyNameToUpdate)
			fmt.Println("Updated job application successfully!")
		} else if strings.HasPrefix(input, "clear") {
			goto LOOP
		} else {
			fmt.Println("Unknown command. Usage: jm [-a|-ls|-u] [-cn|-r|-n|-ha] [true|false]")
		}
	}
}
