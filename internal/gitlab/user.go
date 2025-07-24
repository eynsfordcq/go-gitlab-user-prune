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
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	State     string    `json:"state"`
	Note      string    `json:"note"`
	LastLogin time.Time `json:"last_sign_in_at"`
	CreatedAt time.Time `json:"created_at"`
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

	req, err := s.client.newRequest(http.MethodGet, path)
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

//	{
//	        "id": 131,
//	        "username": "test",
//	        "name": "Test",
//	        "state": "active",
//	        "locked": false,
//	        "avatar_url": "https://secure.gravatar.com/avatar/fce99064b7fd70aec498a5c44f61e3703cd9a1f9663144a559ff599dc34b81f1?s=80&d=identicon",
//	        "web_url": "https://gitlab.globeoss.com/test",
//	        "created_at": "2025-07-24T12:53:38.785Z",
//	        "bio": "",
//	        "location": "",
//	        "public_email": null,
//	        "skype": "",
//	        "linkedin": "",
//	        "twitter": "",
//	        "discord": "",
//	        "website_url": "",
//	        "organization": "",
//	        "job_title": "",
//	        "pronouns": null,
//	        "bot": false,
//	        "work_information": null,
//	        "followers": 0,
//	        "following": 0,
//	        "is_followed": false,
//	        "local_time": null,
//	        "last_sign_in_at": null,
//	        "confirmed_at": "2025-07-24T12:53:38.557Z",
//	        "last_activity_on": null,
//	        "email": "test@globeoss.com",
//	        "theme_id": 3,
//	        "color_scheme_id": 1,
//	        "projects_limit": 100000,
//	        "current_sign_in_at": null,
//	        "identities": [],
//	        "can_create_group": true,
//	        "can_create_project": true,
//	        "two_factor_enabled": false,
//	        "external": false,
//	        "private_profile": false,
//	        "commit_email": "test@globeoss.com",
//	        "is_admin": false,
//	        "note": "User created to test script; PIC: ZJ.",
//	        "namespace_id": 747,
//	        "created_by": {
//	            "id": 8,
//	            "username": "zhengjie",
//	            "name": "Chia Zheng Jie",
//	            "state": "active",
//	            "locked": false,
//	            "avatar_url": "https://gitlab.globeoss.com/uploads/-/system/user/avatar/8/avatar.png",
//	            "web_url": "https://gitlab.globeoss.com/zhengjie"
//	        },
//	        "email_reset_offered_at": null
//	    },
