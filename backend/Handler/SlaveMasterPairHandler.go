package handler

import (
	"net/http"
	database "syncjob/Database"
	"syncjob/Handler/Dtos"
	jobs "syncjob/Handler/Jobs"
	logger "syncjob/Logger"

	"github.com/gin-gonic/gin"
)

func PostSlaveMasterPairSync(c *gin.Context) {
	key := c.Param("key")

	var dto dtos.TriggerSyncDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid JSON body")
		return
	}

	pair := database.SMMap[key]
	owner := dto.Owner

	if dto.All {
		for _, tn := range pair.Mappings.Tables {
			if tn.Owner != owner {
				continue
			}
			jobs.Queue <- jobs.SyncQueue{
				Key:   key,
				Owner: owner,
				Table: tn,
			}
		}
		c.JSON(http.StatusOK, "Sync job queued")
		return
	}

	tn := pair.Mappings.ContainsTable(dto.Table, owner)
	if tn == nil {
		logger.InternalServerErrorStr(c, "Specified table does not exists")
		return
	}

	jobs.Queue <- jobs.SyncQueue{
		Key:   key,
		Owner: owner,
		Table: *tn,
	}

	c.JSON(http.StatusOK, "Sync job queued")
}
