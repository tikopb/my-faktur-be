package model

type Pagination struct {
	Current_page int `json:"current_page"`
	Total_page   int `json:"total_page"`
	Per_page     int `json:"per_page"`
	Total_data   int `json:"total_data"`
}
