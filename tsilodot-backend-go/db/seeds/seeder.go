package main

import (
	"fmt"
	"math/rand"
	"time"
	"tsilodot/db"
	"tsilodot/helpers"
	"tsilodot/model"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	helpers.InitLogger()
	if err := godotenv.Load(); err != nil {
		log.Info().Msg("No .env file found, using environment variables from OS")
	} else {
		log.Info().Msg("Loading .env success")
	}
}

func main() {
	db.InitDBConnection()
	defer db.StopDBConnection()

	fmt.Println("Cleaning up database...")
	db.DB.Exec("TRUNCATE TABLE tasks RESTART IDENTITY CASCADE")
	db.DB.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")

	rand.Seed(time.Now().UnixNano())

	fmt.Println("Seeding database...")
	seed()
	fmt.Println("Seeding completed!")
}

func seed() {
	password, _ := helpers.HashPassword("password")

	for i := 1; i <= 5; i++ {
		user := model.User{
			Name:     fmt.Sprintf("User %d", i),
			Email:    fmt.Sprintf("user%d@example.com", i),
			Password: password,
		}

		if err := db.DB.Create(&user).Error; err != nil {
			log.Error().Err(err).Msgf("Failed to create user %d", i)
			continue
		}

		// Each user has 4-15 tasks
		numTasks := rand.Intn(12) + 4 // [0, 11] + 4 = [4, 15]
		for j := 1; j <= numTasks; j++ {
			dueDate := time.Now().AddDate(0, 0, rand.Intn(30))
			status := "pending"
			if rand.Intn(2) == 0 {
				status = "completed"
			}

			task := model.Task{
				UserID:      user.ID,
				Title:       fmt.Sprintf("Task %d for %s", j, user.Name),
				Description: fmt.Sprintf("Description for task %d of %s", j, user.Name),
				Status:      status,
				DueDate:     &dueDate,
			}

			if err := db.DB.Create(&task).Error; err != nil {
				log.Error().Err(err).Msgf("Failed to create task %d for user %d", j, i)
			}
		}
		fmt.Printf("Seeded user %d with %d tasks\n", i, numTasks)
	}
}
