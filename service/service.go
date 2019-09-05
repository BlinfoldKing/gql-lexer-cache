package service

import (
	"gql-lexer-cache/repository"
	"gql-lexer-cache/repository/model"

	"github.com/sirupsen/logrus"
)

var TotalStory int32 = 0

type StoryService struct {
	repo  *repository.Repository
	cache *repository.Cache
}

func New() *StoryService {
	db, err := repository.InitCockroach(repository.DB_URL)
	if err != nil {
		logrus.Panic(err)
	}

	logrus.Println("DB Connected")

	repo := repository.NewRepo(db)
	cache := repository.Cache{}
	cache.Flush()

	return &StoryService{
		repo:  repo,
		cache: &cache,
	}
}

func (s StoryService) CreateStory(title string) (model.Story, error) {
	story := model.Story{
		ID:    TotalStory + 1,
		Title: title,
	}
	err := s.repo.Save(story)
	if err != nil {
		return model.Story{}, err
	}
	s.cache.Flush()
	return story, nil
}

func (s StoryService) GetQuery(query string) ([]model.Story, error) {
	return s.cache.Get(query)
}

func (s StoryService) SetQuery(query string, value []model.Story) {
	s.cache.Index(query, value)
}

func (s StoryService) GetAll() ([]model.Story, error) {
	return s.repo.GetAll()
}

func (s StoryService) GetOddStories() ([]model.Story, error) {
	stories, err := s.repo.GetAll()
	if err != nil {
		return []model.Story{}, nil
	}

	var result []model.Story
	for _, v := range stories {
		if v.ID%2 == 1 {
			result = append(result, v)
		}
	}

	return result, nil
}

func (s StoryService) GetEvenStories() ([]model.Story, error) {
	stories, err := s.repo.GetAll()
	if err != nil {
		return []model.Story{}, nil
	}

	var result []model.Story
	for _, v := range stories {
		if v.ID%2 == 0 {
			result = append(result, v)
		}
	}

	return result, nil
}
