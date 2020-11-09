package models

type User struct {
	Id    int
	Email string
	Name  string
	Login string
}

type Event struct {
	Type    string
	Actor   User
	Payload Payload
}

type Commit struct {
	Message   string
	Author    User
	Committer User
}

type Payload struct {
	Commits []Commit
}

type Repo struct {
	Id       int
	FullName string `json:"full_name"`
}

type RepoCommits struct {
	Commit Commit
}

type Activity struct {
	Id                int
	Number            int
	Title             string
	Url               string
	User              User
	State             string
	AuthorAssociation string `json:"author_association"`
}
