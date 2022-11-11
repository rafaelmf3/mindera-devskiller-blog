## Introduction

You were hired as a consultant for BestBlogs<sup>TM</sup> company. The company 
needs your help with finishing the implementation of a REST api for its newest blog product.
Your client has already written tests and implemented some parts of the golang service. You were tasked to finish 
the implementation and make all tests pass. Follow this README to complete the code in a correct order. Good luck!


### Persistence layer

Firstly, you need to implement empty methods in `CommentRepository` and `PostRepository`. These structs represent an in-memory storage
for our entities. You will find stubs of the methods in `repository/repository.go` file.
* `CommentRepository.Insert` - inserts a given comment to the database
* `CommentRepository.GetById` - find the comment by the id
* `CommentRepository.GetAllByPostId` - finds all comments by post id
* `PostRepository.Insert` - inserts a given post to the database
* `PostRepository.GetById` - finds a given post by id
  
Detailed descriptions of desirable implementation may be found as comments in the methods' body.


### Web layer
Secondly, you need to implement a web layer for the blog API. This step involves implementing REST API handlers in `service/rest.go`.
Each empty handler function you need to implement contains comments that form a specification of what needs to be done.
As a reference, you can use the implementation of `/api/posts` rest endpoint in `service.rest.go`
The REST endpoints to implement are:
* `/api/posts/comments` - deserializes JSON request payload as `model.Comment` and persists it into `CommentRepository`. Otherwise, appropriate error message and status code are returned.
* `/api/posts/[POST_ID]` -  looks for a post with given id in the database and returns it. Otherwise, appropriate error message and status code are returned.
* `/api/posts/comments/[POST_ID]` - looks for all comments with given post id in the database and returns them. Otherwise, appropriate error message and status code are returned.

## Building and testing
#### Prerequisites: 
1. `make` is installed on your system

#### Building
To build the binary run `make` command in the root directory of this repository.


#### Testing
To run all unit tests issue `make test` command in the root directory of this repository. 


## Good luck!
