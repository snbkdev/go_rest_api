package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	id int
	name string
	author string
}

func main() {
	db, err := sql.Open("sqlite3", "./books.db")
	log.Println(db)
	if err != nil {
		log.Println(err)
	}

	statement, err := db.Prepare("create table if not exists books(id integer primary key, isbn integer, author varchar(64), name varchar(64) null)")
	if err != nil {
		log.Println("Error in creating table")
	} else {
		log.Println("Successfully created table books!")
	}

	statement.Exec()

	statement, _ = db.Prepare("insert into books(name, author, isbn) values (?, ?, ?)")
	statement.Exec("A tale of two cities", "Charles Dickens", 140430547)
	log.Println("Inserted the book into database!")

	rows, _ := db.Query("select id, name, author from books")
	var tempBook Book
	for rows.Next() {
		rows.Scan(&tempBook.id, &tempBook.name, &tempBook.author)
		log.Printf("ID: %d, book: %s, author: %s\n", tempBook.id, tempBook.name, tempBook.author)
	}

	statement, _ = db.Prepare("update books set name = ? where id = ?")
	statement.Exec("The tale of Two Cities", 1)
	log.Println("Successfully updated the book in database!")

	statement, _ = db.Prepare("delete from books where id = ?")
	statement.Exec(1)
	log.Println("Successfully deleted the book in database!")
}