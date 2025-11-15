package domain

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	id        uuid.UUID
	title     Title
	body      Body
	authorID  uuid.UUID
	createdAt time.Time
	updatedAt time.Time
}

func NewPost(authorID uuid.UUID, title Title, body Body) (*Post, error) {
	now := time.Now().UTC()
	return &Post{
		id:        uuid.New(),
		title:     title,
		body:      body,
		authorID:  authorID,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func RehydratePost(
	id uuid.UUID,
	title Title,
	body Body,
	authorID uuid.UUID,
	createdAt, updatedAt time.Time,
) (*Post, error) {
	return &Post{
		id:        id,
		title:     title,
		body:      body,
		authorID:  authorID,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (p *Post) ID() uuid.UUID        { return p.id }
func (p *Post) Title() Title         { return p.title }
func (p *Post) Body() Body           { return p.body }
func (p *Post) AuthorID() uuid.UUID  { return p.authorID }
func (p *Post) CreatedAt() time.Time { return p.createdAt }
func (p *Post) UpdatedAt() time.Time { return p.updatedAt }
