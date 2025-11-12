package main

import (
	"fmt"
	"time"
)

type Post struct {
	ID        string    `json:"id"`
	Author    *User     `json:"author"`
	Content   string    `json:"content"`
	MediaURL  string    `json:"mediaURL,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	Likes     int       `json:"likes"`
	Comments  []string  `json:"comments"`
}

type PostBuilder interface {
	SetID(id string) PostBuilder
	SetAuthor(user *User) PostBuilder
	SetContent(content string) PostBuilder
	SetMedia(url string) PostBuilder
	Build() *Post
}

type ConcretePostBuilder struct {
	post *Post
}

func NewPostBuilder() PostBuilder {
	return &ConcretePostBuilder{
		post: &Post{
			CreatedAt: time.Now(),
			Likes:     0,
			Comments:  make([]string, 0),
		},
	}
}

func (c *ConcretePostBuilder) SetID(id string) PostBuilder {
	c.post.ID = id
	return c
}

func (c *ConcretePostBuilder) SetAuthor(user *User) PostBuilder {
	c.post.Author = user
	return c
}

func (c *ConcretePostBuilder) SetContent(content string) PostBuilder {
	c.post.Content = content
	return c
}

func (c *ConcretePostBuilder) SetMedia(url string) PostBuilder {
	c.post.MediaURL = url
	return c
}

func (c *ConcretePostBuilder) Build() *Post {
	if c.post.ID == "" {
		c.post.ID = fmt.Sprintf("post-%d", time.Now().UnixNano())
	}
	return c.post
}
