package server

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ohmyc/messaggio/internal/app/dal"
	"github.com/ohmyc/messaggio/pkg/processing/request"
)

func process(p *request.Producer, d *dal.Dal) gin.HandlerFunc {
	return func(context *gin.Context) {
		req := &struct {
			Message string `json:"message"`
		}{}
		if context.BindJSON(&req) != nil {
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
	}
}

func processed(d *dal.Dal) gin.HandlerFunc {
	return func(context *gin.Context) {
		id := context.Query("id")
		if id == "" {
			context.String(400, "Missing \"id\" query parameter")
			return
		}
		message, err := d.GetProcessedMessage(id)
		if err != nil {
			if errors.Is(err, dal.ErrNotProcessed) {
				context.String(404, "Not processed message")
				return
			} else {
				fmt.Println("Error getting message:", err)
				context.String(500, "Internal server error")
				return
			}
		}
		context.JSON(200, message)
	}
}

func processedAll(d *dal.Dal) gin.HandlerFunc {
	return func(context *gin.Context) {
		messages, err := d.GetAll()
		if err != nil {
			fmt.Println("Error getting stats:", err)
			context.String(500, "Internal server error")
			return
		}
		context.JSON(200, messages)
	}
}

func stats(d *dal.Dal) gin.HandlerFunc {
	return func(context *gin.Context) {
		stats, err := d.GetStats()
		if err != nil {
			fmt.Println("Error getting stats:", err)
			context.String(500, "Internal server error")
			return
		}
		context.JSON(200, stats)
	}
}

func NewServer(p *request.Producer, d *dal.Dal) *gin.Engine {
	g := gin.Default()
	g.POST("/process", process(p, d))
	g.GET("/processed", processed(d))
	g.GET("/processed-all", processedAll(d))
	g.GET("/stats", stats(d))
	return g
}
