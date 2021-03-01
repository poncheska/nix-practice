package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// docker run --name some-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root -e MYSQL_PASSWORD=password -e MYSQL_USER=user -e MYSQL_DATABASE=dbname -d mysql:latest
// docker exec -it some-mysql bash
// mysql -u user -p dbname

type Post struct {
	Id     int    `json:"id" db:"id"`
	UserId int    `json:"userId" db:"user_id"`
	Title  string `json:"title" db:"title"`
	Body   string `json:"body" db:"body"`
}

type Comment struct {
	Id     int    `json:"id" db:"id"`
	PostId int    `json:"postId" db:"post_id"`
	Name   string `json:"name" db:"name"`
	Email  string `json:"email" db:"email"`
	Body   string `json:"body" db:"body"`
}

type Store struct {
	DB *sqlx.DB
}

func main() {
	db, err := sqlx.Open("mysql", "user:password@tcp(localhost:3306)/dbname")
	if err != nil {
		log.Fatalln(err)
	}

	store := &Store{db}

	store.MustExecFile("beginner/task6/schema/init.sql")
	defer store.MustExecFile("beginner/task6/schema/drop.sql")

	userId := 7
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts?userId=" + strconv.Itoa(userId))
	if err != nil {
		log.Fatalln(err)
	}

	posts := new([]Post)
	err = json.NewDecoder(resp.Body).Decode(posts)
	if err != nil {
		log.Fatalln(err)
	}

	wg := new(sync.WaitGroup)
	for _, v := range *posts {
		wg.Add(1)
		go func(p Post) {
			err := store.WritePost(p)
			if err != nil {
				log.Fatalln(err)
			}

			resp, err := http.Get("https://jsonplaceholder.typicode.com/comments?postId=" + strconv.Itoa(p.Id))
			if err != nil {
				log.Fatalln(err)
			}

			comments := new([]Comment)
			err = json.NewDecoder(resp.Body).Decode(comments)
			if err != nil {
				log.Fatalln(err)
			}

			wg1 := new(sync.WaitGroup)
			for _, vv := range *comments {
				wg1.Add(1)
				go func(c Comment) {
					err := store.WriteComment(c)
					if err != nil {
						log.Fatalln(err)
					}
					wg1.Done()
				}(vv)
			}
			wg1.Wait()

			wg.Done()
		}(v)
	}
	wg.Wait()
	var input string
	fmt.Scanln(&input)
}

func (s *Store) MustExecFile(path string) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range strings.Split(string(bs), "\n\n") {
		fmt.Println(v)
		s.DB.MustExec(v)
	}
}

func (s *Store) WriteComment(c Comment) error {
	_, err := s.DB.Exec("INSERT INTO comment (id, post_id, name, email, body) VALUES (?, ?, ?, ?, ?)",
		strconv.Itoa(c.Id), strconv.Itoa(c.PostId), c.Name, c.Email, c.Body)
	return err
}

func (s *Store) WritePost(p Post) error {
	_, err := s.DB.Exec("INSERT INTO post (id, user_id, title, body) VALUES (?, ?, ?, ?)",
		strconv.Itoa(p.Id), strconv.Itoa(p.UserId), p.Title, p.Body)
	return err
}
