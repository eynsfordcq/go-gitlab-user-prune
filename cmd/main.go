package main

import (
	"log"

	"github.com/eynsfordcq/go-gitlab-user-prune/internal/config"
	"github.com/eynsfordcq/go-gitlab-user-prune/internal/gitlab"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	client := gitlab.NewClient(config.APIBaseUrl, config.APIToken)
	users, err := client.Users.ListAllActiveUsers()
	if err != nil {
		log.Fatalf("fail to fetch users: %v", err)
	}

	log.Printf("fetched %d users", len(users))
	for _, user := range users {
		log.Print(user)
	}
}
