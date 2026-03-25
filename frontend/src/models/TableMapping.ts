export enum TableOwner {
    MASTER,
    SLAVE
}

export interface ColumnMapping {
    SlaveName: string
    MasterName: string
}

export interface TableMapping {
    Owner: TableOwner
    MasterTableName: string,
    SlaveTableName: string
    ColumnsMapped: ColumnMapping[]
}

export interface Mapping {
    Tables: TableMapping[]
	LastMasterSync: Date
	LastSlaveSync: Date
}
