package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"postsvc/models"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// in memory store

var (
	posts = make(map[int]*models.Post)
	mu    sync.Mutex
	id    = 1
)

func init() {
	//populate posts from db or some dummy values
}
func GetPosts(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	resp, err := http.Get("http://localhost:9090/posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var fetchedPosts []*models.Post

	if err := json.NewDecoder(resp.Body).Decode(&fetchedPosts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, post := range fetchedPosts {
		posts[post.ID] = post
	}

	c.JSON(http.StatusOK, fetchedPosts)

}

func GetPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Post Id"})
		return
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:9090/posts/%d", postID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var fetchedPost models.Post
	if err := json.NewDecoder(resp.Body).Decode(&fetchedPost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mu.Lock()
	posts[fetchedPost.ID] = &fetchedPost
	mu.Unlock()

	c.JSON(http.StatusOK, &fetchedPost)
}

func CreatePost(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	var newPost models.Post

	if err := c.ShouldBindJSON(&newPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newPost.ID = id
	id++
	dbpost, err := createPostInDB(&newPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	posts[newPost.ID] = dbpost

	c.JSON(http.StatusCreated, newPost)
}

func UpdatePost(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	// get the post id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post, exists := posts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var updatedPost models.Post

	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.Title = updatedPost.Title
	post.Content = updatedPost.Content

	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	mu.Lock()

	defer mu.Unlock()

	//get the id
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	_, exists := posts[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "post with this id is not present"})
		return
	}
	delete(posts, id)

	c.JSON(http.StatusOK, gin.H{"message ": "Post deleted"})

}

func createPostInDB(post *models.Post) (*models.Post, error) {

	jsonvalue, err := json.Marshal(post)
	if err != nil {
		return nil, fmt.Errorf("error in json data %v\n", err)
	}

	resp, err := http.Post("http://localhost:9090/posts", "application/json", bytes.NewBuffer(jsonvalue))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respdata, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var newpost *models.Post
	json.Unmarshal(respdata, &newpost)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("received non 200 status code %d", resp.StatusCode)
	}
	return newpost, nil

}
