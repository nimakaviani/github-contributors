package scraper

type User struct {
	Login string
	Id    int
}

type GithubUser struct {
	Email string
	Name  string
	Login string
}

type GithubEvent struct {
	Type    string
	Actor   User
	Payload Payload
}

type Commit struct {
	Message   string
	Author    GithubUser
	Committer GithubUser
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

type Issue struct {
	Id                int
	Title             string
	User              User
	State             string
	AuthorAssociation string `json:"author_association"`
}
