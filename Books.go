package main

type Book struct
{
	ID     		int    `json:"id"`
	Name  		string `json:"name"`
	Price 		float64 `json:"price"`
	Category  	string `json:"category"`
	Author 		string `json:"author"`
}