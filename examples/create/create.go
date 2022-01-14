package main

import (
	"log"
	"os"

	"github.com/Karitham/cronjoborg"
)

func main() {
	c := cronjoborg.New(os.Getenv("CRONJOB_API_KEY"))

	jobID, err := c.CreateJob(cronjoborg.DetailedJob{
		Job: cronjoborg.Job{
			URL:     "https://example.com",
			Title:   "Test Job",
			Enabled: true,
			Schedule: cronjoborg.Schedule{
				Minutes: []int{-1}, // Every minute
			},
			RequestTimeout: 10,
			SaveResponses:  true,
			RequestMethod:  cronjoborg.RequestMethodGet,
		},
	})
	if err != nil {
		log.Println("error creating job:", err)
		return
	}

	log.Println("job created:", jobID)
}
