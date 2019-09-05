package resolver

import (
	"context"
	"gql-lexer-cache/repository/model"
	"gql-lexer-cache/service"
)

var ServiceConnection *Resolver

type Resolver struct {
	Service *service.StoryService
}

func New() *Resolver {
	return &Resolver{service.New()}
}

func (r Resolver) CreateStory(ctx context.Context, args struct{ Title string }) (StoryResolver, error) {
	story, err := r.Service.CreateStory(args.Title)
	return StoryResolver{m: story}, err
}

func (r Resolver) GetAllStories() ([]StoryResolver, error) {

	result, err := r.Service.GetAll()
	if err != nil {
		return []StoryResolver{}, err
	}

	var resolves []StoryResolver
	for _, v := range result {
		resolves = append(resolves, StoryResolver{v})
	}

	r.Service.SetQuery("GetAllStories", result)
	return resolves, nil
}

func (r Resolver) GetOddStories() ([]StoryResolver, error) {
	result, err := r.Service.GetOddStories()
	if err != nil {
		return []StoryResolver{}, err
	}

	var resolves []StoryResolver
	for _, v := range result {
		resolves = append(resolves, StoryResolver{v})
	}

	r.Service.SetQuery("GetOddStories", result)
	return resolves, nil
}

func (r Resolver) GetEvenStories() ([]StoryResolver, error) {
	result, err := r.Service.GetEvenStories()
	if err != nil {
		return []StoryResolver{}, err
	}

	var resolves []StoryResolver
	for _, v := range result {
		resolves = append(resolves, StoryResolver{v})
	}

	r.Service.SetQuery("GetEvenStories", result)
	return resolves, nil
}

type StoryResolver struct {
	m model.Story
}

func (s StoryResolver) Id() int32 {
	return s.m.ID
}

func (s StoryResolver) Title() string {
	return s.m.Title
}
