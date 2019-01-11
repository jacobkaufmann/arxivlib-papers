package datastore

import (
	"context"
	"fmt"

	arxivlib "github.com/jacobkaufmann/arxivlib-papers"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type papersStore struct {
	*Datastore
}

func (s *papersStore) Get(id primitive.ObjectID) (*arxivlib.Paper, error) {
	coll := s.db.Collection("papers")
	var paper *arxivlib.Paper

	filter := bson.M{"_id": id}

	err := coll.FindOne(context.Background(), filter).Decode(&paper)
	if err != nil {
		return nil, err
	}

	return paper, nil
}

func (s *papersStore) List(opt *arxivlib.PaperListOptions) ([]*arxivlib.Paper, error) {
	if opt == nil {
		opt = &arxivlib.PaperListOptions{}
	}
	coll := s.db.Collection("papers")
	var papers []*arxivlib.Paper

	filter := bson.D{
		{"title", primitive.Regex{Pattern: opt.Title, Options: "i"}},
		{"updated", bson.D{
			{"$gte", opt.Updated},
		}},
		{"abstract", primitive.Regex{Pattern: opt.Abstract, Options: "i"}},
	}

	cats := bson.A{}
	if len(opt.Categories) > 0 {
		for i := 0; i < len(opt.Categories); i++ {
			cats = append(cats, fmt.Sprintf("/%s/", opt.Categories[i]))
		}
		filter = append(filter, bson.D{{"categories", bson.D{{"$in", cats}}}}...)
	}

	auths := bson.A{}
	if opt.Author != "" {
		auths = append(auths, fmt.Sprintf("/%s/", opt.Author))
		filter = append(filter, bson.D{{"authors", bson.D{{"$in", auths}}}}...)
	}

	cursor, err := coll.Find(
		context.Background(),
		filter,
	)
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

func (s *papersStore) Update(paper *arxivlib.Paper) (bool, error) {
	coll := s.db.Collection("papers")

	filter := bson.M{"_id": paper.ID}

	result, err := coll.UpdateOne(context.Background(), filter, paper)
	if err != nil {
		return false, err
	} else if result.MatchedCount != 1 && result.ModifiedCount != 1 {
		return false, nil
	}

	return true, nil
}

func (s *papersStore) Upload(paper *arxivlib.Paper) (bool, error) {
	coll := s.db.Collection("papers")

	_, err := coll.InsertOne(
		context.Background(),
		&paper,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}
