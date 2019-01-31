package datastore

import (
	"context"
	"errors"
	"fmt"
	"time"

	arxivlib "github.com/jacobkaufmann/arxivlib-papers"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

type papersStore struct {
	*Datastore
}

func (s *papersStore) Get(id primitive.ObjectID) (*arxivlib.Paper, error) {
	coll := s.db.Db.Collection("papers")
	var paper *arxivlib.Paper

	filter := bson.M{"_id": id}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := coll.FindOne(ctx, filter).Decode(&paper)
	if err != nil {
		return nil, err
	}

	return paper, nil
}

func (s *papersStore) List(opt *arxivlib.PaperListOptions) ([]*arxivlib.Paper, error) {
	if opt == nil {
		opt = &arxivlib.PaperListOptions{}
	}
	coll := s.db.Db.Collection("papers")
	var papers []*arxivlib.Paper

	filter := bson.M{
		"title": primitive.Regex{Pattern: opt.Title, Options: "i"},
		"updated": bson.M{
			"$gte": opt.Updated,
		},
		"abstract": primitive.Regex{Pattern: opt.Abstract, Options: "i"},
	}

	if len(opt.Categories) > 0 {
		cats := bson.A{}
		for i := 0; i < len(opt.Categories); i++ {
			cats = append(cats, fmt.Sprintf("%s", opt.Categories[i]))
		}
		filter["categories"] = bson.M{"$in": cats}
	}

	if opt.Author != "" {
		filter["authors"] = primitive.Regex{Pattern: opt.Author, Options: "i"}
	}

	opts := &options.FindOptions{}
	opts = opts.SetSort(bson.M{"updated": -1})
	if opt.MaxResults > 0 {
		opts.SetLimit(int64(opt.MaxResults))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		paper := &arxivlib.Paper{}
		if err := cursor.Decode(paper); err != nil {
			return nil, err
		}
		papers = append(papers, paper)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return papers, nil
}

func (s *papersStore) Update(paper *arxivlib.Paper) (updated bool, err error) {
	coll := s.db.Db.Collection("papers")

	filter := bson.M{"_id": paper.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := coll.UpdateOne(ctx, filter, paper)
	if err != nil {
		return false, err
	} else if result.MatchedCount != 1 && result.ModifiedCount != 1 {
		return false, nil
	}

	return true, nil
}

func (s *papersStore) Upload(papers []*arxivlib.Paper) (uploaded bool, err error) {
	coll := s.db.Db.Collection("papers")

	// Convert slice of type Paper to slice of type interface{} for insertion
	many := make([]interface{}, len(papers))
	for i, v := range papers {
		many[i] = v
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if len(many) > 1 {
		_, err = coll.InsertMany(ctx, many)
	} else {
		_, err = coll.InsertOne(ctx, many[0])
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *papersStore) Remove(id primitive.ObjectID) (removed bool, err error) {
	coll := s.db.Db.Collection("papers")
	var paper *arxivlib.Paper

	filter := bson.M{"_id": paper.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return false, err
	}
	if result.DeletedCount != 1 {
		return false, arxivlib.ErrPaperNotFound
	}

	return true, nil
}

func (s *papersStore) AddRating(id primitive.ObjectID, r *arxivlib.Rating) (added bool, err error) {
	coll := s.db.Db.Collection("papers")

	filter := bson.M{"_id": id}
	update := bson.M{"$push": bson.M{"ratings": r}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	} else if result.MatchedCount != 1 {
		return false, arxivlib.ErrPaperNotFound
	} else if result.ModifiedCount != 1 {
		return false, errors.New("paper not updated")
	}

	return true, nil
}
