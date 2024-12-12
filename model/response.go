package model

type Success struct{
	Message string `json:"message"`
}

type Error struct{
	Error string `json:"message"`
}