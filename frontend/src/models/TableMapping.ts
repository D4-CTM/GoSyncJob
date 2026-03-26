export enum TableOwner {
    MASTER,
    SLAVE
}

export interface ColumnMapping {
    SlaveName: string
    MasterName: string
}

export class TableMapping {
    Owner: TableOwner
    MasterTableName: string
    SlaveTableName: string
    LastSync: Date
    ColumnsMapped: ColumnMapping[]

    fmtName = () =>
        this.Owner === TableOwner.SLAVE
            ? this.SlaveTableName
            : this.MasterTableName
}

export interface Mapping {
    Tables: TableMapping[]
}
