package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

// Конструктор объекта хранилища.
func New(constr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

func (s *Store) Posts() (posts []storage.Post, err error) {
	var rows pgx.Rows
	rows, err = s.db.Query(context.Background(), `
	SELECT 
		id, 
		author_id, 
		title, 
		content, 
		created_at
	FROM posts`)

	if err != nil {
		return nil, err
	}

	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var t storage.Post
		err = rows.Scan(
			&t.ID,
			&t.AuthorID,
			&t.Title,
			&t.Content,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, t)
	}

	// ВАЖНО не забыть проверить rows.Err()
	err = rows.Err()
	return
}

func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		INSERT INTO posts (author_id, title, content, created_at, published_at)
		VALUES ($1, $2, $3, COALESCE(NULLIF($4,0), EXTRACT(EPOCH FROM TIMESTAMP WITH TIME ZONE now())), COALESCE(NULLIF($5,0), EXTRACT(EPOCH FROM TIMESTAMP WITH TIME ZONE now()));
		`,
		p.AuthorID,		
		p.Title,
		p.Content,
		p.CreatedAt,
		p.PublishedAt,
	)
	return err
}
func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE posts
		SET
		author_id=$2, 
		title=$3, 
		content=$4, 
		published_at=COALESCE(NULLIF($5,0), EXTRACT(EPOCH FROM TIMESTAMP WITH TIME ZONE now()))
		WHERE id=$1
		`,
		p.ID,
		p.AuthorID,		
		p.Title,
		p.Content,		
		p.PublishedAt,
	)
	return err
}
func (s *Store) DeletePost(p storage.Post) error {
	_, err := s.db.Exec(context.Background(), `
		DELETE FROM posts		
		WHERE id=$1
		`,
		p.ID,		
	)
	return err
}
