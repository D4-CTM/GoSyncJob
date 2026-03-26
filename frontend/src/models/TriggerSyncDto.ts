import { TableOwner } from './TableMapping';

export interface TriggerSyncDto {
    Owner: TableOwner
    All: boolean
    Table: string
}
