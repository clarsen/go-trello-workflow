// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package handle_graphql

import (
	"time"
)

type BoardList struct {
	Board string `json:"board"`
	List  string `json:"list"`
}

type BoardListInput struct {
	Board string `json:"board"`
	List  string `json:"list"`
}

type FinishResult struct {
	Message *string `json:"message"`
	Ok      bool    `json:"ok"`
}

type GenerateResult struct {
	Message *string `json:"message"`
	Ok      bool    `json:"ok"`
}

type Task struct {
	ID               string     `json:"id"`
	Title            string     `json:"title"`
	CreatedDate      *time.Time `json:"createdDate"`
	URL              *string    `json:"url"`
	Due              *time.Time `json:"due"`
	List             *BoardList `json:"list"`
	Period           *string    `json:"period"`
	DateLastActivity *time.Time `json:"dateLastActivity"`
	Desc             string     `json:"desc"`
	ChecklistItems   []string   `json:"checklistItems"`
}

type Timer struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type WeeklyGoal struct {
	IDCard      string  `json:"idCard"`
	IDCheckitem string  `json:"idCheckitem"`
	Title       string  `json:"title"`
	Tasks       []*Task `json:"tasks"`
	Year        *int    `json:"year"`
	Month       *int    `json:"month"`
	Week        *int    `json:"week"`
	Done        *bool   `json:"done"`
	Status      *string `json:"status"`
}
