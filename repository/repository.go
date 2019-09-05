package repository

import (
	"database/sql"
	"errors"
	"gql-lexer-cache/repository/model"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var DB_URL = "postgresql://root@localhost:26257/postgres?sslmode=disable"

func InitCockroach(url string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", url)
	return
}

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(story model.Story) error {
	_, err := r.db.Exec(
		`INSERT INTO "Story" VALUES($1, $2)`,
		story.ID,
		story.Title)

	return err
}

func parseStory(rows *sql.Rows) ([]model.Story, error) {
	var stories []model.Story
	for rows.Next() {
		var (
			id    int32
			title string
		)
		err := rows.Scan(&id, &title)
		if err != nil {
			return []model.Story{}, err
		}

		stories = append(stories, model.Story{
			ID:    id,
			Title: title,
		})
	}

	return stories, nil
}

func (r *Repository) GetAll() ([]model.Story, error) {
	rows, err := r.db.Query(
		`SELECT * FROM "Story"`)
	if err != nil {
		return []model.Story{}, err
	}
	result, err := parseStory(rows)

	return result, nil
}

type Cache struct {
	story map[string][]model.Story
}

func (c *Cache) Index(query string, data []model.Story) {
	c.story[query] = data
}

func (c *Cache) Get(key string) ([]model.Story, error) {
	data, ok := c.story[key]
	if ok {
		logrus.Println("get from cache")
		return data, nil
	}

	return []model.Story{}, errors.New("no story found")
}

func (c *Cache) Flush() {
	logrus.Println("flush cache")
	c.story = make(map[string][]model.Story)
}
