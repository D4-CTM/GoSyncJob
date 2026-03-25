import { DbCredential } from "./Credentials"
import { Mapping } from "./TableMapping"

export interface SlaveMasterPair {
    Name: string
    Slave: DbCredential
    Master: DbCredential
    Mappings: Mapping
}
