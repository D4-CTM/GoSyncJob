export enum DatabaseType {
    ORACLE,
    POSTGRES
}

export interface DbCredential {
    Database: string
    Server: string
    Port: number

    User: string
    Password: string

    DbType: DatabaseType
}

export interface PutCredentialResult {
    OldName: string
    NewName: string
}
