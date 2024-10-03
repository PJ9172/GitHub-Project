package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
)

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Svn_url     string `json:"svn_url"`
}

type Info struct {
	Username   string `json:"login"`
	Name       string `json:"name"`
	Bio        string `json:"bio"`
	Followers  int    `json:"followers"`
	Following  int    `json:"following"`
	PublicRepo int    `json:"public_repos"`
	Avatar_url string `json:"avatar_url"`
	Repos_url  string `json:"repos_url"` 
	Repolen    int    `json:"repolen"`
	Repos      []Repo `json:"repos"`
}


func main() {

	http.HandleFunc("/", rootUrlHandler)
	http.HandleFunc("/submit", submitUrlHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server starting on port 3000!!!!")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}
}

func rootUrlHandler(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "index.html")
}

func submitUrlHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method == "POST" {
		req.ParseForm()
		user := req.FormValue("user")
		url := "https://api.github.com/users/" + user
		// fmt.Println(url)

		// First API call to get user information
		response, err := http.Get(url)
		if err != nil {
			fmt.Fprint(res, "Error fetching user data: ", err)
			return
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Fprint(res, "Error reading response body: ", err)
			return
		}

		var data Info
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Fprint(res, "Error unmarshalling JSON: ", err)
			return
		}

		// Check if the user exists
		if data.Username == "" {
			fmt.Fprint(res, "Invalid Username! Please enter a valid GitHub username.")
			return
		}

		// Second API call to get user repositories
		response, err = http.Get(data.Repos_url)
		if err != nil {
			fmt.Fprint(res, "Error fetching repositories: ", err)
			return
		}

		defer response.Body.Close()
		body, err = io.ReadAll(response.Body)
		if err != nil {
			fmt.Fprint(res, "Error reading repositories data: ", err)
			return
		}

		// Unmarshal the repositories data
		err = json.Unmarshal(body, &data.Repos)
		if err != nil {
			fmt.Fprint(res, "Error unmarshalling repository JSON: ", err)
			return
		}

		data.Repolen = len(data.Repos) // Set the length of repos

		// Marshal the entire Info struct into JSON (for debugging or other purposes)
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Fprint(res, "Error marshalling data to JSON: ", err)
			return
		}

		// Load the HTML template and pass the data
		temp, err := template.ParseFiles("users.html")
		if err != nil {
			fmt.Fprint(res, "Error loading template: ", err)
			return
		}

		// Render the template with the data
		temp.Execute(res, template.JS(jsonData))
	} else {
		fmt.Fprint(res, "Invalid request method: ", http.StatusBadRequest)
	}
}

