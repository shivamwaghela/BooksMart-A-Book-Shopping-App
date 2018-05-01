/*
	Gumball API in Go
	Uses MySQL & Riak KV
*/

package main


type user struct {
	UserId    string 	`json:"UserId_register"`
	Name 	string 		`json:"Name_register"`
	Email string		`json:"Email_register"`
	
}


