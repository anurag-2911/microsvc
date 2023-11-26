package handlers

import (
	"net/http"
	"postdbsvc/db"
	"postdbsvc/models"

	"github.com/gin-gonic/gin"
)

// GetPosts fetches all posts
func GetPosts(c *gin.Context) {
    rows, err := db.DB.Query("SELECT id, title, content FROM posts")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    posts := make([]models.Post, 0)
    for rows.Next() {
        var p models.Post
        if err := rows.Scan(&p.ID, &p.Title, &p.Content); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        posts = append(posts, p)
    }

    if err = rows.Err(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, posts)
}

// CreatePost creates a new post
func CreatePost(c *gin.Context) {
    var newPost models.Post
    if err := c.ShouldBindJSON(&newPost); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    sqlStatement := `INSERT INTO posts (title, content) VALUES ($1, $2) RETURNING id`
    id := 0
    err := db.DB.QueryRow(sqlStatement, newPost.Title, newPost.Content).Scan(&id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    newPost.ID = id
    c.JSON(http.StatusCreated, newPost)
}

// UpdatePost updates an existing post
// Add similar logic as CreatePost, but with SQL UPDATE statement

// DeletePost deletes a post
// Add similar logic as CreatePost, but with SQL DELETE statement
