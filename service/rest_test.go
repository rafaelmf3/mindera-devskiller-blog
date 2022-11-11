package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bitbucket.org/mindera/go-rest-blog/model"
	"bitbucket.org/mindera/go-rest-blog/repository"
)

func TestRestApiService_handleAddPost(t *testing.T) {
	tests := []struct {
		testName           string
		commentRepository  repository.CommentRepository
		postRepository     repository.PostRepository
		post               interface{}
		expectedHttpStatus int
		expectedHeader     string
		expectedResponse   interface{}
	}{
		{
			testName:           "testSuccessfullyAddPost",
			post:               model.Post{Id: 256, Title: "title", Content: "cntnt", CreationDate: time.Now()},
			commentRepository:  repository.CustomCommentRepository(make([]model.Comment, 0)),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			expectedHttpStatus: 200,
			expectedHeader:     "application/json",
			expectedResponse:   AckJsonResponse{Message: "post id: 256 successfully added", Status: http.StatusOK},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// GIVEN
			data, _ := json.Marshal(tc.post)
			req := httptest.NewRequest(http.MethodPost, postsPath, bytes.NewReader(data))
			w := httptest.NewRecorder()
			router := mux.NewRouter()
			svc := RestApiService{&tc.postRepository, &tc.commentRepository}

			// WHEN
			router.HandleFunc(postsPath, svc.handleAddPost)
			router.ServeHTTP(w, req)
			response := w.Result()
			body, _ := io.ReadAll(response.Body)
			var ackResponse AckJsonResponse
			err := json.Unmarshal(body, &ackResponse)
			require.NoError(t, err)

			// THEN
			assert.Equal(t, tc.expectedHttpStatus, response.StatusCode)
			assert.Equal(t, tc.expectedHeader, response.Header.Get("Content-Type"))
			assert.Equal(t, tc.expectedResponse, ackResponse)
		})
	}
}

func TestRestApiService_handleGetPostByPostId(t *testing.T) {
	var testDate = time.Date(2018, time.September, 16, 12, 0, 0, 0, time.UTC)
	var validPost = model.Post{Id: 34, Title: "happy post", Content: "test content", CreationDate: testDate}
	var badID = "badID"

	tests := []struct {
		testName           string
		commentRepository  repository.CommentRepository
		postRepository     repository.PostRepository
		postId             string
		expectedHttpStatus int
		expectedHeader     string
		expectedResponse   interface{}
		verifyResponseFunc func(t *testing.T, expectedResponse interface{}, body []byte)
	}{
		{
			testName:           "testSuccessfullyGetPost",
			commentRepository:  repository.CustomCommentRepository(make([]model.Comment, 0)),
			postRepository:     repository.CustomPostRepository([]model.Post{validPost}),
			postId:             strconv.FormatUint(validPost.Id, 10),
			expectedHttpStatus: 200,
			expectedHeader:     "application/json",
			expectedResponse:   validPost,
			verifyResponseFunc: func(t *testing.T, expectedResponse interface{}, body []byte) {
				t.Helper()
				var post model.Post
				err := json.Unmarshal(body, &post)
				require.NoError(t, err)
				assert.Equal(t, expectedResponse, post)
			},
		},
		{
			testName:           "testPostNotFound",
			commentRepository:  repository.CustomCommentRepository(make([]model.Comment, 0)),
			postRepository:     repository.CustomPostRepository([]model.Post{}),
			postId:             "111",
			expectedHttpStatus: 404,
			expectedHeader:     "application/json",
			expectedResponse:   AckJsonResponse{Message: "Post with id: 111 does not exist", Status: 404},
			verifyResponseFunc: func(t *testing.T, expectedResponse interface{}, body []byte) {
				t.Helper()
				var resp AckJsonResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.Equal(t, expectedResponse, resp)
			},
		},
		{
			testName:           "testPostBadRequest",
			commentRepository:  repository.CustomCommentRepository(make([]model.Comment, 0)),
			postRepository:     repository.CustomPostRepository([]model.Post{}),
			postId:             badID,
			expectedHttpStatus: 400,
			expectedHeader:     "application/json",
			expectedResponse:   AckJsonResponse{Message: "wrong id path variable: " + badID, Status: 400},
			verifyResponseFunc: func(t *testing.T, expectedResponse interface{}, body []byte) {
				t.Helper()
				var resp AckJsonResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.Equal(t, expectedResponse, resp)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// GIVEN
			svc := RestApiService{commentRepository: &tc.commentRepository,
				postRepository: &tc.postRepository}

			path := strings.Replace(getPostPath, "{id}", tc.postId, 1)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()
			router := mux.NewRouter()

			// WHEN
			router.HandleFunc(getPostPath, svc.handleGetPostByPostId)
			router.ServeHTTP(w, req)
			response := w.Result()
			body, _ := io.ReadAll(response.Body)

			// THEN
			assert.Equal(t, tc.expectedHttpStatus, response.StatusCode)
			assert.Equal(t, tc.expectedHeader, response.Header.Get("Content-Type"))
			tc.verifyResponseFunc(t, tc.expectedResponse, body)
		})
	}
}

func TestRestApiService_handleGetCommentsByPostId(t *testing.T) {
	var testDate = time.Date(2018, time.September, 16, 12, 0, 0, 0, time.UTC)
	var validComments = []model.Comment{
		{Id: 123, PostId: 3, Comment: "abc", Author: "cool author", CreationDate: testDate},
		{Id: 321, PostId: 3, Comment: "def", Author: "cool author2", CreationDate: testDate},
		{Id: 543, PostId: 3, Comment: "ghi", Author: "cool author3", CreationDate: testDate},
	}
	var badID = "badID"

	tests := []struct {
		testName           string
		commentRepository  repository.CommentRepository
		postRepository     repository.PostRepository
		postId             string
		expectedHttpStatus int
		expectedHeader     string
		expectedResponse   interface{}
		verifyResponseFunc func(t *testing.T, expectedResponse interface{}, body []byte)
	}{
		{
			testName:           "testSuccessfullyGetComments",
			commentRepository:  repository.CustomCommentRepository(validComments),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			postId:             "3",
			expectedHttpStatus: 200,
			expectedHeader:     "application/json",
			expectedResponse:   validComments,
			verifyResponseFunc: func(t *testing.T, expectedResponse interface{}, body []byte) {
				t.Helper()
				var commentsList []model.Comment
				err := json.Unmarshal(body, &commentsList)
				require.NoError(t, err)
				assert.ElementsMatch(t, expectedResponse, commentsList)
			},
		},
		{
			testName:           "testEmptyComments",
			commentRepository:  repository.CustomCommentRepository(validComments),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			postId:             "1",
			expectedHttpStatus: 200,
			expectedHeader:     "application/json",
			expectedResponse:   nil,
			verifyResponseFunc: func(t *testing.T, expectedResponse interface{}, body []byte) {
				t.Helper()
				var commentsList []model.Comment
				err := json.Unmarshal(body, &commentsList)
				require.NoError(t, err)
				assert.ElementsMatch(t, expectedResponse, commentsList)
			},
		},
		{
			testName:           "testBadRequest",
			commentRepository:  repository.CustomCommentRepository(validComments),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			postId:             badID,
			expectedHttpStatus: 400,
			expectedHeader:     "application/json",
			expectedResponse:   AckJsonResponse{Message: "wrong id path variable: " + badID, Status: 400},
			verifyResponseFunc: func(t *testing.T, expectedResponse interface{}, body []byte) {
				t.Helper()
				var resp AckJsonResponse
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.Equal(t, expectedResponse, resp)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// GIVEN
			svc := RestApiService{commentRepository: &tc.commentRepository,
				postRepository: &tc.postRepository}
			path := strings.Replace(getCommentPath, "{id}", tc.postId, 1)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			w := httptest.NewRecorder()
			router := mux.NewRouter()

			// WHEN
			router.HandleFunc(getCommentPath, svc.handleGetCommentsByPostId)
			router.ServeHTTP(w, req)
			response := w.Result()
			body, _ := io.ReadAll(response.Body)

			// THEN
			assert.Equal(t, tc.expectedHttpStatus, response.StatusCode)
			assert.Equal(t, tc.expectedHeader, response.Header.Get("Content-Type"))
			tc.verifyResponseFunc(t, tc.expectedResponse, body)
		})
	}
}

func TestRestApiService_handleAddComment(t *testing.T) {
	var validComment = model.Comment{Id: 123, PostId: 3, Comment: "cool cmnt", Author: "cool auth", CreationDate: time.Now()}
	var validReqBody, _ = json.Marshal(&validComment)

	tests := []struct {
		testName           string
		commentRepository  repository.CommentRepository
		postRepository     repository.PostRepository
		reqBody            []byte
		expectedHttpStatus int
		expectedHeader     string
		expectedResponse   interface{}
	}{
		{
			testName:           "testSuccessfullyAddComment",
			commentRepository:  repository.CustomCommentRepository(make([]model.Comment, 0)),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			reqBody:            validReqBody,
			expectedHttpStatus: 200,
			expectedHeader:     "application/json",
			expectedResponse:   AckJsonResponse{Message: "comment id: 123 successfully added", Status: 200},
		},
		{
			testName:           "testIncompleteData",
			commentRepository:  repository.CustomCommentRepository(make([]model.Comment, 0)),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			reqBody:            []byte("{}"),
			expectedHttpStatus: 400,
			expectedHeader:     "application/json",
			expectedResponse:   AckJsonResponse{Message: "could not deserialize comment json payload", Status: 400},
		},
		{
			testName:           "testBadPayload",
			commentRepository:  repository.CustomCommentRepository(make([]model.Comment, 0)),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			reqBody:            []byte("invalidJson"),
			expectedHttpStatus: 400,
			expectedHeader:     "application/json",
			expectedResponse:   AckJsonResponse{Message: "could not deserialize comment json payload", Status: 400},
		},
		{
			testName:           "testAlreadyExists",
			commentRepository:  repository.CustomCommentRepository([]model.Comment{validComment}),
			postRepository:     repository.CustomPostRepository(make([]model.Post, 0)),
			reqBody:            validReqBody,
			expectedHttpStatus: 400,
			expectedHeader:     "application/json",
			expectedResponse:   AckJsonResponse{Message: "Comment with id: 123 already exists in the database", Status: 400},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			// GIVEN
			svc := RestApiService{commentRepository: &tc.commentRepository,
				postRepository: &tc.postRepository}

			req := httptest.NewRequest(http.MethodPost, commentsPath, bytes.NewReader(tc.reqBody))
			w := httptest.NewRecorder()
			router := mux.NewRouter()

			// WHEN
			router.HandleFunc(commentsPath, svc.handleAddComment)
			router.ServeHTTP(w, req)
			response := w.Result()
			body, _ := io.ReadAll(response.Body)
			var resp AckJsonResponse
			err := json.Unmarshal(body, &resp)
			require.NoError(t, err)

			// THEN
			assert.Equal(t, tc.expectedHttpStatus, response.StatusCode)
			assert.Equal(t, tc.expectedHeader, response.Header.Get("Content-Type"))
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}
