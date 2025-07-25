package main

import (
	"log"
	"strings"
	"time"

	"github.com/eynsfordcq/go-gitlab-user-prune/internal/config"
	"github.com/eynsfordcq/go-gitlab-user-prune/internal/gitlab"
)

func getUsersToBlock(users []gitlab.User, cfg *config.Config) []gitlab.User {
	now := time.Now()
	inactivityThreshold := now.AddDate(0, 0, -cfg.InactivityDays)
	var candidates []gitlab.User

	for _, user := range users {
		// check if user is in whitelist
		if _, ok := cfg.Whitelist[user.Email]; ok {
			log.Printf("user %s is in whitelist, skipping", user.Email)
			continue
		}

		// check notes
		if strings.Contains(user.Note, cfg.WhiteListText) {
			log.Printf("user %s is whitelisted by admin note, skipping", user.Email)
			continue
		}

		// check if is bot
		if user.IsBot {
			log.Printf("user %s is bot, skipping", user.Email)
		}

		// find the lastest activity whether it's just created
		// or just logged in
		// or recently has activities
		latestActivity := *user.CreatedAt
		if user.LastLogin != nil && user.LastLogin.After(latestActivity) {
			latestActivity = *user.LastLogin
		}

		if user.LastActivity != nil && user.LastActivity.Time.After(latestActivity) {
			latestActivity = user.LastActivity.Time
		}

		if latestActivity.After(inactivityThreshold) {
			continue
		}

		log.Printf("user %s set for pruning. last recent activity: %s",
			user.Email,
			latestActivity,
		)
		candidates = append(candidates, user)
	}

	return candidates
}

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

	candidates := getUsersToBlock(users, config)
	if len(candidates) == 0 {
		log.Println("no users found to prune")
		return
	}

	log.Printf("found %d users to prune", len(candidates))
	for _, user := range candidates {
		err = client.Users.BlockUser(user.ID)
		if err != nil {
			log.Printf("fail to block user: %s", user.Email)
		} else {
			log.Printf("block user %s success", user.Email)
		}
	}

	log.Printf("completed")
}
