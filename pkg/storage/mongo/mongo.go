package mongo

import (
	"GoNews/pkg/storage"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Хранилище данных.
type Store struct {
	db *mongo.Client
}

const (
	databaseName   = "news"  // имя БД
	collectionName = "posts" // имя коллекции в БД
)

// Конструктор объекта хранилища.
func New(constr string) (*Store, error) {
	mongoOpts := options.Client().ApplyURI(constr)
	client, err := mongo.Connect(context.Background(), mongoOpts)

	if err != nil {
		return nil, err
	}
	s := Store{
		db: client,
	}
	return &s, nil
}

func (s *Store) Posts() (posts []storage.Post, err error) {
	err, posts = nil, nil
	collection := s.db.Database(databaseName).Collection(collectionName)
	filter := bson.D{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {
		var l storage.Post
		err := cur.Decode(&l)
		if err != nil {
			return nil, err
		}
		posts = append(posts, l)
	}
	return posts, cur.Err()

}

func (s *Store) AddPost(p storage.Post) error {
	collection := s.db.Database(databaseName).Collection(collectionName)
	//insert post
	_, err := collection.InsertOne(context.Background(), p)
	return err
}

func (s *Store) UpdatePost(p storage.Post) error {
	collection := s.db.Database(databaseName).Collection(collectionName)
	filter := bson.M{"id": p.ID}
	update := bson.M{"$set": p}
	_, err := collection.UpdateOne(context.Background(), filter, update)

	return err
}
func (s *Store) DeletePost(p storage.Post) error {
	collection := s.db.Database(databaseName).Collection(collectionName)
	filter := bson.M{"id": p.ID}
	_, err := collection.DeleteOne(context.Background(), filter)

	return err
}
