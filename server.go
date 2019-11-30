package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	cPort = 3000

	cPathPing = "/ping"
	cPathSum  = "/sum"
)

type sumRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"value":  "pong",
	})
}

func handleSum(c *gin.Context) {
	if value, err := computeSum(c); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"result":  "error",
			"message": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result": "success",
			"value":  value,
		})
	}
}

func computeSum(c *gin.Context) (int, error) {
	defer c.Request.Body.Close()

	// read request body
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("[ERROR] unable to read request body ::", err)
		return 0, err
	}

	// unmarshal request body
	var request sumRequest
	if err := json.Unmarshal(data, &request); err != nil {
		log.Println("[INFO] unable to parse json request ::", err)
		return 0, err
	}

	return request.A + request.B, nil
}

// Serve starts the API server
func main() {
	router := gin.Default()
	router.GET(cPathPing, handlePing)
	router.POST(cPathSum, handleSum)

	server := &http.Server{}
	server.Addr = fmt.Sprintf(":%d", cPort)
	server.Handler = router

	log.Println("[INFO] running server on port:", cPort)
	server.ListenAndServe()
}
