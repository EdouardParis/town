package resources

import (
	"fmt"
	"time"

	"github.com/EdouardParis/town/constants"
	"github.com/EdouardParis/town/models"
)

type Article struct {
	Title           string  `json:"title"`
	Subtitle        *string `json:"subtitle"`
	Body            string  `json:"body"`
	Lang            string  `json:"lang"`
	AmountCollected int64   `json:"amount_collected"`
	AmountReceived  int64   `json:"amount_received"`
	Status          string  `json:"status"`
	Slug            string  `json:"slug"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	PublishedAt *time.Time `json:"published_at"`

	Address *Address `json:"address,omitempty"`
	Node    *Node    `json:"node,omitempty"`

	Reactions Reactions `json:"reactions"`
}

func (a Article) AbsoluteURL() string {
	return fmt.Sprintf("articles/%s", a.Slug)
}

func NewArticle(article *models.Article) *Article {
	if article == nil {
		return nil
	}

	resource := &Article{
		Title:           article.Title,
		Body:            article.Body,
		Lang:            constants.LangIntToStr[article.Lang],
		AmountCollected: article.AmountCollected,
		AmountReceived:  article.AmountReceived,
		Status:          constants.ArticleStatusIntToStr[article.Status],
		Slug:            article.Slug,
		CreatedAt:       article.CreatedAt,
	}

	if article.Subtitle.Valid {
		resource.Subtitle = &article.Subtitle.String
	}

	if article.UpdatedAt.Valid {
		resource.UpdatedAt = &article.UpdatedAt.Time
	}

	if article.PublishedAt.Valid {
		resource.PublishedAt = &article.PublishedAt.Time
	}

	if article.Node != nil {
		resource.Node = NewNode(article.Node)
	}

	if article.Address != nil {
		resource.Address = NewAddress(article.Address)
	}

	if len(article.Reactions) > 0 {
		resource.Reactions = NewReactions(article.Reactions)
	}

	return resource
}

func NewArticleList(articles []models.Article) []Article {
	resources := make([]Article, len(articles))
	for i := range articles {
		resource := NewArticle(&articles[i])
		resources[i] = *resource
	}
	return resources
}
