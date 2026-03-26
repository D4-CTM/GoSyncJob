package handler

import (
	"fmt"
	"net/http"
	database "syncjob/Database"
	"syncjob/Logger"

	"github.com/gin-gonic/gin"
)

const (
	SLAVE string = "Slave"
	MASTER string = "Master"
)

func GetSlaveMasterPairs(c *gin.Context) {
	pairs := make([]string, len(database.SMMap))
	idx := 0
	for key := range database.SMMap {
		pairs[idx] = key
		idx++
	}

	c.JSON(http.StatusOK, pairs)
}

func GetSlaveMasterPair(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		logger.BadRequestStr(c, "Please provide a key name")
		return
	}

	pair := database.SMMap[key]

	c.JSON(http.StatusOK, pair)
}

func DeleteSlaveMasterPair(c *gin.Context) {
	pairName := c.Param("key")
	pair := database.SMMap[pairName]
	if err := pair.Close(); err != nil {
		logger.InternalServerError(c, err)
		return
	}

	delete(database.SMMap, pairName)
	database.SaveSMPM()
	c.Status(http.StatusNoContent)
}

func PostSlaveMasterPair(c *gin.Context) {
	pair, err := database.CreatePairFromGin(c)
	if err != nil {
		logger.InternalServerError(c, err)
		return
	}

	pairName := pair.Name
	mapKey := pairName
	{
		idx := 1
		for name := range database.SMMap {
			if name == mapKey {
				mapKey = fmt.Sprintf("%s_%d", pairName, idx)
				idx++
			}	
		}
	}

	database.SMMap[mapKey] = *pair
	database.SaveSMPM()
}

func PutSlaveMasterPair(c *gin.Context) {
	oldName := c.Param("key")

	pair, _ := database.CreatePairFromGin(c)
	if err := pair.Ping(); err != nil {
		logger.InternalServerError(c, err)
		return
	}

	newName := pair.Name
	if newName == "" {
		logger.InternalServerErrorStr(c, "New pair name cannot be empty!",)
		return 
	}

	if newName != oldName {
		for name := range database.SMMap {
			if name == newName && name != oldName {
				logger.InternalServerErrorStr(c, "Connection name already exists")
				return
			}
		}
	}

	if newName != oldName {
		database.SMMap[newName] = *pair
		delete(database.SMMap, oldName)
	}
	database.SaveSMPM()

	c.JSON(http.StatusOK, gin.H{
		"OldName": oldName,
		"NewName": newName,
	})
}

func PostCredentialsPing(c *gin.Context) {
	_, err := database.CreateCredFromGin(c)
	if err != nil {
		logger.InternalServerError(c, err)
		return
	}

	c.Status(http.StatusOK)
}
