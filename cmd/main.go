package main

import (
	"fmt"
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
			continue
		}

		// check notes
		if strings.Contains(user.Note, cfg.WhiteListText) {
			continue
		}

		// check if is bot
		if user.IsBot {
			continue
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

func blockAndSetNotes(c *gitlab.Client, user gitlab.User) {
	// block user
	err := c.Users.BlockUser(user.ID)
	if err != nil {
		log.Printf("fail to block user: %s", user.Email)
	} else {
		log.Printf("block user %s success", user.Email)
	}

	// update user note
	note := fmt.Sprintf("blocked by pruning script due to inactivity. %s",
		time.Now(),
	)

	err = c.Users.UpdateUser(user.ID, gitlab.UpdateUserOptions{
		Note: note,
	})
	if err != nil {
		log.Printf("fail to update user: %s", user.Email)
	} else {
		log.Printf("update user %s success", user.Email)
	}
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

	candidates := getUsersToBlock(users, config)
	if len(candidates) == 0 {
		log.Println("no users found to prune")
		return
	}

	log.Printf("found %d users to prune", len(candidates))
	log.Printf("dry run flag: %v", config.DryRun)
	for _, user := range candidates {
		log.Printf("block user: %s with id: %d ", user.Email, user.ID)
		if config.DryRun {
			continue
		}
		blockAndSetNotes(client, user)
	}

	log.Printf("completed")
}
