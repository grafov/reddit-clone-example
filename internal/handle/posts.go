package handle

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handlePostList(c *gin.Context) {
	r := bytes.NewReader([]byte(xxxList))
	c.DataFromReader(http.StatusOK, int64(len(xxxList)), "application/json", r, nil)
}

func handlePost(c *gin.Context) {
	c.JSON(200, gin.H{"token": "zzz"})
}

const xxxList = `
[
    {
	"author": {
	    "id": "5ed3fba46fe531000740facd",
	    "username": "Jeanvaljean"
	},
	"category": "news",
	"comments": [],
	"created": "2020-07-21T22:43:24.584Z",
	"id": "5f176f8ce64be600073e73a1",
	"score": -2,
	"text": "The Les Amis de l'abc caught him and I freed him. But he killed himself at the end. RIP Javert! I will not forget his name. I will not forget him. I am Jean Valjean.",
	"title": "Javert committed suicide!",
	"type": "text",
	"upvotePercentage": 0,
	"views": 6,
	"votes": [
	    {
		"user": "5ed3fba46fe531000740facd",
		"vote": -1
	    },
	    {
		"user": "5f69013457b29900078c0938",
		"vote": -1
	    }
	]
    },
    {
	"author": {
	    "id": "5c6a4ced0bab4adf9b1df352",
	    "username": "Hello"
	},
	"category": "news",
	"comments": [
	    {
		"author": {
		    "id": "5c6a4ced0bab4adf9b1df352",
		    "username": "Hello"
		},
		"body": "June, Mary july/July. Pink room. What did the dad see?",
		"created": "2019-02-18T06:16:03.943Z",
		"id": "5c6a4da30bab4a6fe01df355"
	    }
	],
	"created": "2019-02-18T06:15:23.740Z",
	"id": "5c6a4d7b0bab4a5d9b1df354",
	"score": -5,
	"title": "True detective season 3 episode 7 discussion",
	"type": "link",
	"upvotePercentage": 14,
	"url": "http://www.hbonow.com",
	"views": 67,
	"votes": [
	    {
		"user": "5c6a4ced0bab4adf9b1df352",
		"vote": 1
	    },
	    {
		"user": "5c6a74cde72901539d8a296c",
		"vote": -1
	    },
	    {
		"user": "5cc38daca3060a148aaa4210",
		"vote": -1
	    },
	    {
		"user": "5d45a5c8e8a04502f8cb5aec",
		"vote": -1
	    },
	    {
		"user": "5d4bb4685e1e66707ac194e2",
		"vote": -1
	    },
	    {
		"user": "5df5577268840d4391c70979",
		"vote": -1
	    },
	    {
		"user": "5e570492f34b110007c14aaa",
		"vote": -1
	    }
	]
    }
]
`
