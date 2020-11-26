package models

import (
	"log"

	"abidhmuhsin.com/gowebapp/mvc/config"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Book struct {
	gorm.Model // for ID,CreatedAT,UpdatedAt,DeletedAt
	//Id          string `json:"id"`
	Name        string `gorm:"" json:"name"`
	Author      string `gorm:"size:255" json:"author"`
	Publication string `gorm:"type:varchar(100)" json:"publ"`
}

func init() {
	log.Println("init func in mvc/book.go - connecting db")
	config.Connect()
	db = config.GetDB()
	db.DropTableIfExists(&Book{}) //Drops the table if already exists
	db.AutoMigrate(&Book{})

	// Add some initial data -- Internally it will create the query like
	// INSERT INTO `book` (`name`,`address`) VALUES ('John','New York')

	var books []Book = []Book{
		{Name: "Ricky", Author: "Sydney", Publication: "RickysBook"},
		{Name: "Adam", Author: "Brisbane", Publication: "AdamsBook"},
		{Name: "Justin", Author: "California", Publication: "JustinsBook"},
		{Name: "Beiber", Author: "California", Publication: "BeibersBook"},
		{Name: "Amy", Author: "Boston", Publication: "AmysBook"},
	}

	for _, book := range books {
		// db.Debug().Create(&book) // use  db.Debug().* to see queries that are run in sql by gorm
		db.Create(&book)
	}
	// ------------------------------
}

func (b *Book) CreateBook() *Book {
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func GetAllBooks() []Book {
	var Books []Book
	db.Debug().Find(&Books)
	return Books
}

func GetBookById(Id int64) (*Book, *gorm.DB) {
	var getBook Book
	db := db.Where("ID = ?", Id).Find(&getBook)
	return &getBook, db
}

func DeleteBook(ID int64) Book {
	var book Book
	db.Where("ID = ?", ID).Delete(book)
	return book
}
