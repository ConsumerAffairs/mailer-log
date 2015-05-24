# 1. Install process

## Pre install

1. brew install mongodb

## Backend development

### Install 

- Create go worspace runing `mkdir -p ~/go/src/github.com`
- Add to your `.zshrc` or `.bash_profile`

```
# Golang
export GOPATH=$HOME/go
export PATH=$PATH:/usr/local/opt/go/libexec/bin
```

- enter folder `~/go/src/github.com`
- clone this repo `git clone https://github.com/ConsumerAffairs/mailer-log.git`
- run `go get ../...`
	1. to build: `go build server.go`
	2. to run without build `go run server.go`


## Frontend development

### Install

- Enter in `frontend` folder
- run `npm install tsd@next -g`
- run `npm install gulp -g`
- run `npm install typescript -g`
- run `npm install`

### Watch build
- Just run `gulp` it will monitoring the changes and run typescript and webpack.


#2. API Doc

1. Document structure: [model.go](models/mail.go)
2. Endpoint list

| URL        | HTTP Verb | POST Body   |   URL Params   | Result        |
|------------|-----------|-------------|:--------------:|---------------|
| /mails     | GET       | empty       | page, per_page | mail list     |
| /mails     | POST      | JSON string |       N/A      | create mail   |
| /mails/:id | GET       | empty       |       N/A      | retrieve mail |
| /mails/:id | PUT       | JSON string |       N/A      | update mail   |
| /mails/:id | DELETE    | empty       |       N/A      | delete mail   |


# Todo

### General
-----------
1. Filter

### Frontend
------------
1. Tests
2. Gulp/webpack production build
3. Fix gulp task order issue

### Backend
-----------

1. Tests
2. Schema validation (do we need ?)

## References

1. JSON API: [http://jsonapi.org/recommendations/#filtering]()
