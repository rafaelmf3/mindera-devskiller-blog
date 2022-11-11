package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"bitbucket.org/mindera/go-rest-blog/model"
	"bitbucket.org/mindera/go-rest-blog/repository"
)

type RestApiService struct {
	postRepository    *repository.PostRepository
	commentRepository *repository.CommentRepository
}

type AckJsonResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func NewRestApiService() RestApiService {
	return RestApiService{
		postRepository:    repository.NewPostRepository(),
		commentRepository: repository.NewCommentRepository(),
	}
}

func (svc *RestApiService) ServeContent(port int) error {
	portString := ":" + strconv.Itoa(port)
	svc.initializeHandlers()
	return http.ListenAndServe(portString, nil)
}

const (
	postsPath      = "/api/posts"
	getPostPath    = postsPath + "/{id}"
	commentsPath   = "/api/posts/comments"
	getCommentPath = commentsPath + "/{id}"
)

func (svc *RestApiService) initializeHandlers() {
	r := mux.NewRouter()

	r.HandleFunc(postsPath, svc.handleAddPost).Methods(http.MethodPost)
	r.HandleFunc(getPostPath, svc.handleGetPostByPostId).Methods(http.MethodGet)
	r.HandleFunc(getCommentPath, svc.handleGetCommentsByPostId).Methods(http.MethodGet)
	r.HandleFunc(commentsPath, svc.handleAddComment).Methods(http.MethodPost)
	http.Handle("/", r)
}

func (svc *RestApiService) handleAddPost(w http.ResponseWriter, r *http.Request) {
	var post model.Post

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return
	}
	if err := svc.postRepository.Insert(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(&AckJsonResponse{Message: fmt.Sprintf("post id: %d successfully added", post.Id), Status: http.StatusOK})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (svc *RestApiService) handleGetPostByPostId(w http.ResponseWriter, r *http.Request) {
	// TODO example valid api call: GET /api/posts/42
	//  Every response should have Content-Type=application/json header set
	//  should respond with a json response in a format of `AckJsonResponse` with Status 400 and Message "wrong id path variable: PATH_VARIABLE" when invalid ID given,
	//  note that also the HTTP response code should be set to 400!
	//  e.g. GET /api/posts/abc --> AckJsonResponse{Message: "wrong id path variable: abc", Status: 400}
	//  should respond with a json response in a format of `AckJsonResponse` with Status 404 and Message "Post with id: [POST_ID] does not exist"
	//  note that also the HTTP response code should be set to 404!
	//  when given postID does not exist
	//  e.g. GET /api/posts/35 --> '{"Message": "post with id: 35 does not exist", Status: 404}'
	//  should respond with valid post entity when post with given id exists:
	//  e.g. GET /api/posts/2 --> {"Id": 2, "Title": "test title", "Content": "this is a post content", "CreationDate": "1970-01-01T03:46:40+01:00"}

	// Given that this project uses gorilla/mux as a router you can access the path params with following code:
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(AckJsonResponse{
			Message: fmt.Sprintf("wrong id path variable: %s", vars["id"]),
			Status:  http.StatusBadRequest,
		})
		w.Write(response)
		return
	}
	res, err := svc.postRepository.GetById(uint64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal(AckJsonResponse{
			Message: fmt.Sprintf("Post with id: %d does not exist", id),
			Status:  http.StatusNotFound,
		})
		w.Write(response)
		return
	}
	response, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func (svc *RestApiService) handleGetCommentsByPostId(w http.ResponseWriter, r *http.Request) {
	// TODO example valid api call: GET /api/posts/comments/4
	//  Every response should have Content-Type=application/json header set
	//  should respond with a json response in a format of `AckJsonResponse` with Status 400 and Message "wrong id path variable: PATH_VARIABLE" when invalid ID given, e.g. GET /api/posts/comments/abc
	//  should respond with a valid json response with a list of comments for given postId. If there are no comments for a given postId, should return an empty list
	//  e.g. example valid api call: GET /api/posts/comments/101 -->
	//  '[
	//	 	{"Id": 1, "PostId": 101, "Comment": "comment1", "Author": "author5", "CreationDate" :"1970-01-01T03:46:40+01:05"},
	//	 	{"Id": 3, "PostId": 101, "Comment": "comment2", "Author": "author4", "CreationDate" :"1970-01-01T03:46:40+01:10"},
	//	 	{"Id": 5, "PostId": 101, "Comment": "comment3", "Author": "author13", "CreationDate" :"1970-01-01T03:46:40+01:15"}
	//	 ]'

	// Given that this project uses gorilla/mux as a router you can access the path params with following code:
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(AckJsonResponse{
			Message: fmt.Sprintf("wrong id path variable: %s", vars["id"]),
			Status:  http.StatusBadRequest,
		})
		w.Write(response)
		return
	}
	res := svc.commentRepository.GetAllByPostId(uint64(id))
	response, _ := json.Marshal(res)
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (svc *RestApiService) handleAddComment(w http.ResponseWriter, r *http.Request) {
	// TODO: example valid api call: POST /api/posts/comments '{"Id": 1, "PostId": 101, "Comment": "comment1", "Author": "author1", "CreationDate" :"1970-01-01T03:46:40+01:00"}'
	//  Every response should have Content-Type=application/json header set
	//  should respond with a json response in a format of `AckJsonResponse` with Status code 400 and Message "could not deserialize comment json payload"
	//  when invalid or incomplete data posted. Data is considered
	//  incomplete when payload misses any member property of the model.
	//  Note that HTTP response code also should be 400
	//  e.g. POST /api/posts/comments '{"weird_payload": "weird value"}' --> '{"Message": "could not deserialize comment json payload", Status: 400}'
	//  should respond with a json response in a format of `AckJsonResponse` with Status code 400 and json payload Message
	//  "Comment with id: COMMENT_ID already exists in the database"
	//  when comment with given id already exists in the database.
	//  e.g. POST /api/posts/comments '{"Id": 30, "PostId": 23123, "Comment": "comment1", "Author": "author1", "CreationDate" :"1970-01-01T03:46:40+01:00"}'
	//  --> '{"Message": "Comment with id: 30 already exists in the database", Status: 400}'
	//  should respond with a json response in a format of `AckJsonResponse` with status code 200 and message 'comment id: COMMENT_ID successfully added' when data was posted successfully.
	//  e.g. POST /api/posts/comments '{"Id": 123, "PostId": 663, "Comment": "this is a comment", "Author": "blogger", "CreationDate" :"1970-01-01T03:46:40+01:00"}' -->
	//  '{"Message": "comment id: 123 successfully added", Status: 200}'

	body := model.Comment{}

	decoder := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err := decoder.Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(AckJsonResponse{
			Message: "could not deserialize comment json payload",
			Status:  http.StatusBadRequest,
		})
		w.Write(response)
		return
	}

	if body.Id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(AckJsonResponse{
			Message: "could not deserialize comment json payload",
			Status:  http.StatusBadRequest,
		})
		w.Write(response)
		return
	}

	err := svc.commentRepository.Insert(body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(AckJsonResponse{
			Message: fmt.Sprintf("Comment with id: %d already exists in the database", body.Id),
			Status:  http.StatusBadRequest,
		})
		w.Write(response)
		return
	}

	response, _ := json.Marshal(AckJsonResponse{
		Message: fmt.Sprintf("comment id: %d successfully added", body.Id),
		Status:  http.StatusOK,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
