package app

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ohmyc/messaggio/pkg/processing/request"
)

func NewServer(p *request.Producer, d *Dal) *gin.Engine {
	g := gin.Default()
	g.POST("/process", func(context *gin.Context) {
		req := &struct {
			Message string `json:"message"`
		}{}
		if context.BindJSON(&req) != nil {
			context.String(400, "Invalid body: expected single field \"message\"")
			return
		}
		produce, err := p.Produce(req.Message)
		if err != nil {
			fmt.Println("Error producing message:", err)
			context.String(500, "Internal server error")
			return
		}
		if err = d.InsertNewMessage(produce, req.Message); err != nil {
			fmt.Println("Error inserting message:", err)
			context.String(500, "Internal server error")
			return
		}
		context.JSON(200, gin.H{"id": produce})
	})
	g.GET("/processed", func(context *gin.Context) {
		id := context.Query("id")
		if id == "" {
			context.String(400, "Missing \"id\" query parameter")
			return
		}
		message, err := d.GetProcessedMessage(id)
		if err != nil {
			if errors.Is(err, ErrNotProcessed) {
				context.String(404, "Not processed message")
				return
			} else {
				fmt.Println("Error getting message:", err)
				context.String(500, "Internal server error")
				return
			}
		}
		context.JSON(200, message)
	})
	g.GET("/processed-all", func(context *gin.Context) {
		messages, err := d.GetAll()
		if err != nil {
			fmt.Println("Error getting stats:", err)
			context.String(500, "Internal server error")
			return
		}
		context.JSON(200, messages)
	})
	g.GET("/stats", func(context *gin.Context) {
		stats, err := d.GetStats()
		if err != nil {
			fmt.Println("Error getting stats:", err)
			context.String(500, "Internal server error")
			return
		}
		context.JSON(200, stats)
	})
	return g
}
