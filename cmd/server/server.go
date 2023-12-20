package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongo"
	"GoNews/pkg/storage/postgres"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	//authentication information
	authInfo := os.Getenv("crend")
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.
	//
	// БД в памяти.
	db := memdb.New()

	// Реляционная БД PostgreSQL.
	log.Println("connecting postgre")
	db2, err := postgres.New(fmt.Sprintf("postgresql://%s@localhost/db31a_3_1", authInfo))
	if err != nil {
		log.Fatal(err)
	}
	// posts, err := db2.Posts()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(posts)

	// Документная БД MongoDB.
	log.Println("connecting mongodb")
	db3, err := mongo.New("mongodb://localhost:27017/")
	if err != nil {
		log.Fatal(err)
	}
	// err=db3.UpdatePost(storage.Post{
	// 		ID:          0,
	// 	AuthorID:    0,
	// 	Title:       "title",
	// 	Content:     "content",
	// 	AuthorName:  "renamer",
	// 	CreatedAt:   0,
	// 	PublishedAt: 0,
	// })
	// if err != nil {
	// 	log.Println(err)
	// }
	// err = db3.AddPost(storage.Post{
	// 	ID:          11,
	// 	AuthorID:    0,
	// 	Title:       "art",
	// 	Content:     "monaco",
	// 	AuthorName:  "weter",
	// 	CreatedAt:   0,
	// 	PublishedAt: 0,
	// })
	// if err != nil {
	// 	log.Println(err)
	// }
	// posts, err = db3.Posts()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(posts)
	// err=db3.DeletePost(storage.Post{ID: 11})
	// if err != nil {
	// 	log.Println(err)
	// }
	// posts, err = db3.Posts()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println(posts)
	
	_, _ = db2, db3

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
