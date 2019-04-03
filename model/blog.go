package model

import (
	"fmt"
	"log"
)

type Page struct {
	Id    int
	Title string
	Body  []byte
}

func (p Page) SelectBlog() (Page, error) {
	title := p.Title
	page := Page{}
	err := db.QueryRow("select * from blog where title=?", title).Scan(&page.Id, &page.Title, &page.Body)
	return page, err
}

func (_ Page) SelectAll() []Page {
	rows, _ := db.Query("select * from blog")
	page := Page{}
	var page_list []Page
	for rows.Next() {
		err := rows.Scan(&page.Id, &page.Title, &page.Body)
		if err != nil {
			log.Fatal(err)
		}
		page_list = append(page_list, page)
	}
	return page_list
}

func (p Page) InsertBlog() {
	title := p.Title
	body := p.Body
	stmt, _ := db.Prepare("insert into blog(title, body)values(?, ?)")
	res, _ := stmt.Exec(title, body)
	id, err := res.LastInsertId()
	fmt.Println(id, err)
}

func (p Page) UpdateBlog() {
	title := p.Title
	body := p.Body
	stmt, _ := db.Prepare("update blog set body=? where title=?")
	res, _ := stmt.Exec(body, title)
	id, err := res.RowsAffected()
	fmt.Println(id, err)
}
