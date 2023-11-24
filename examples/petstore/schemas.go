package main

type Pet struct {
	ID   int64   `json:"id,omitempty"`
	Name string  `json:"name,omitempty"`
	Tag  *string `json:"tag,omitempty"`
}

type NewPet struct {
	Name string  `json:"name,omitempty"`
	Tag  *string `json:"tag,omitempty"`
}

type Error struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}
