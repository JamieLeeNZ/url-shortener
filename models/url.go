package models

type URLMapping struct {
	Key       string `json:"key"`
	Original  string `json:"original_url"`
	CreatedAt string `json:"created_at"`
}

type URLShortenRequest struct {
	Original string `json:"original_url"`
}

type URLShortenResponse struct {
	Key string `json:"key"`
}
