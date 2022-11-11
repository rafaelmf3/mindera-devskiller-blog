package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bitbucket.org/mindera/go-rest-blog/model"
)

func TestCommentRepository_Insert(t *testing.T) {
	var (
		comment1 = model.Comment{Id: 1, PostId: 101, Comment: "comment2", Author: "author2", CreationDate: time.Unix(10011, 0)}
	)

	t.Run("insert new Comment", func(t *testing.T) {
		c := NewCommentRepository()
		err := c.Insert(comment1)
		assert.NoError(t, err)
	})

	t.Run("comment already exists", func(t *testing.T) {
		c := NewCommentRepository()
		err := c.Insert(comment1)
		require.NoError(t, err)
		err = c.Insert(comment1)
		require.Error(t, err)
		assert.ErrorIs(t, err, CommentAlreadyExistsError{comment1.Id})
	})
}

func TestCommentRepository_GetAllByPostId(t *testing.T) {
	var (
		comment1          = model.Comment{Id: 1, PostId: 101, Comment: "comment2", Author: "author2", CreationDate: time.Unix(10011, 0)}
		comment2          = model.Comment{Id: 2, PostId: 101, Comment: "comment2", Author: "author2", CreationDate: time.Unix(10011, 0)}
		comment3          = model.Comment{Id: 3, PostId: 100, Comment: "comment3", Author: "author3", CreationDate: time.Unix(10022, 0)}
		NonExistentPostId = uint64(10101010)
	)

	t.Run("single comment", func(t *testing.T) {
		c := NewCommentRepository()
		err := c.Insert(comment1)
		require.NoError(t, err)
		assert.ElementsMatch(t, c.GetAllByPostId(comment1.PostId), []model.Comment{comment1})
	})

	t.Run("multiple comments", func(t *testing.T) {
		c := NewCommentRepository()
		err := c.Insert(comment1)
		require.NoError(t, err)
		err = c.Insert(comment2)
		require.NoError(t, err)
		err = c.Insert(comment3)
		require.NoError(t, err)

		expectedResult := []model.Comment{comment1, comment2}
		result := c.GetAllByPostId(comment1.PostId)
		assert.ElementsMatch(t, expectedResult, result)
	})

	t.Run("no comments", func(t *testing.T) {
		c := NewCommentRepository()
		err := c.Insert(comment1)
		require.NoError(t, err)
		err = c.Insert(comment2)
		require.NoError(t, err)
		err = c.Insert(comment3)
		require.NoError(t, err)

		assert.ElementsMatch(t, c.GetAllByPostId(NonExistentPostId), make([]*model.Comment, 0))
	})
}

func TestCommentRepository_GetById(t *testing.T) {
	var (
		comment1             = model.Comment{Id: 1, PostId: 101, Comment: "comment2", Author: "author2", CreationDate: time.Unix(10011, 0)}
		comment2             = model.Comment{Id: 2, PostId: 101, Comment: "comment2", Author: "author2", CreationDate: time.Unix(10011, 0)}
		comment3             = model.Comment{Id: 3, PostId: 100, Comment: "comment3", Author: "author3", CreationDate: time.Unix(10022, 0)}
		comments             = []model.Comment{comment1, comment2, comment3}
		NonExistentCommentID = uint64(20202020)
	)

	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Comment
		wantErr error
	}{
		{name: "comment exists", args: args{comment1.Id}, want: &comment1, wantErr: nil},
		{name: "comment not found", args: args{NonExistentCommentID}, want: nil, wantErr: CommentNotFoundError{NonExistentCommentID}},
	}
	c := CustomCommentRepository(comments)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comment, err := c.GetById(tt.args.id)
			if tt.wantErr != nil {
				assert.Nil(t, comment)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, comment)
			}

		})
	}
}

func TestPostRepository_Insert(t *testing.T) {
	var (
		post1 = model.Post{Id: 101, Title: "post1", Content: "content", CreationDate: time.Unix(10011, 0)}
	)

	t.Run("insert new Post", func(t *testing.T) {
		p := NewPostRepository()
		err := p.Insert(post1)
		assert.NoError(t, err)
	})

	t.Run("post already exists", func(t *testing.T) {
		p := NewPostRepository()
		err := p.Insert(post1)
		require.NoError(t, err)
		err = p.Insert(post1)
		require.Error(t, err)
		assert.ErrorIs(t, err, PostAlreadyExistsError{id: post1.Id})
	})
}

func TestPostRepository_GetById(t *testing.T) {
	var (
		NonExistentPostId = uint64(10101010)
		post1             = model.Post{Id: 101, Title: "post1", Content: "content", CreationDate: time.Unix(10011, 0)}
		post2             = model.Post{Id: 102, Title: "post2", Content: "content", CreationDate: time.Unix(10012, 0)}
		post3             = model.Post{Id: 103, Title: "post3", Content: "content", CreationDate: time.Unix(10013, 0)}
		posts             = []model.Post{post1, post2, post3}
	)

	type args struct {
		id uint64
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Post
		wantErr error
	}{
		{name: "post exists", args: args{post1.Id}, want: &post1, wantErr: nil},
		{name: "post not found", args: args{NonExistentPostId}, want: nil, wantErr: PostNotFoundError{NonExistentPostId}},
	}
	p := CustomPostRepository(posts)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := p.GetById(tt.args.id)
			if tt.wantErr != nil {
				assert.Nil(t, res)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, res)
			}
		})
	}
}
