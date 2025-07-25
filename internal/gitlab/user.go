package gitlab

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

type UserService struct {
	client *Client
}

type User struct {
	ID           int        `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	State        string     `json:"state"`
	Note         string     `json:"note"`
	IsBot        bool       `json:"bot"`
	LastActivity *Date      `json:"last_activity_on"`
	LastLogin    *time.Time `json:"last_sign_in_at"`
	CreatedAt    *time.Time `json:"created_at"`
}

type ListUsersOptions struct {
	Page    int  `url:"page,omitempty"`
	PerPage int  `url:"per_page,omitempty"`
	Active  bool `url:"active,omitempty"`
}

func (s *UserService) ListUsers(opts ListUsersOptions) ([]User, error) {
	v, err := query.Values(opts)
	if err != nil {
		return nil, err
	}

	path := "/api/v4/users"
	if len(v) > 0 {
		path += "?" + v.Encode()
	}

	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("fail to create list users request: %w", err)
	}

	var users []User
	_, err = s.client.do(req, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) ListAllActiveUsers() ([]User, error) {
	var allUsers []User

	opts := ListUsersOptions{
		Page:    1,
		PerPage: 100,
		Active:  true,
	}

	for {
		users, err := s.ListUsers(opts)
		if err != nil {
			return nil, fmt.Errorf("fail to list users: %w", err)
		}

		if len(users) == 0 {
			break
		}

		allUsers = append(allUsers, users...)
		opts.Page++
	}

	return allUsers, nil
}

func (s *UserService) BlockUser(id int) error {
	path := fmt.Sprintf("/api/v4/users/%d/block", id)

	req, err := s.client.newRequest(http.MethodPost, path, nil)
	if err != nil {
		return fmt.Errorf("fail to create block user request: %w", err)
	}

	_, err = s.client.do(req, nil)
	if err != nil {
		return fmt.Errorf("fail to block user %d: %w", id, err)
	}

	return nil
}

type UpdateUserOptions struct {
	Note string `json:"note,omitempty"`
}

func (s *UserService) UpdateUser(id int, opts UpdateUserOptions) error {
	path := fmt.Sprintf("/api/v4/users/%d", id)

	req, err := s.client.newRequest(http.MethodPut, path, opts)
	if err != nil {
		return fmt.Errorf("fail to create update user request: %w", err)
	}

	_, err = s.client.do(req, nil)
	if err != nil {
		return fmt.Errorf("fail to update user %d: %w", id, err)
	}

	return nil
}
