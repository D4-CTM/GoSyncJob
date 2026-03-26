package dtos

import "syncjob/Database"

type TriggerSyncDto struct {
	Owner database.TableOwner
	All   bool
	Table string
}
