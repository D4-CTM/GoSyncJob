<script setup lang="ts">
import { Get } from '@/helper/HttpCaller';
import { DatabaseType } from '@/models/Credentials';
import { SlaveMasterPair } from '@/models/SlaveMasterPair';
import { TableOwner } from '@/models/TableMapping';
import { ref } from 'vue';

const props = defineProps({
    PairName: {
        type: String,
        required: true
    }
});

let pair = ref({
    Master: {
        DbType: DatabaseType.POSTGRES,
        Database: 'Pagila'
    },
    Slave: {
        DbType: DatabaseType.ORACLE,
        Database: 'FREEPDB1'
    },
    Mappings: {
        Tables: [
            {
                Owner: TableOwner.MASTER,
                MasterTableName: 'rental',
                SlaveTableName: 'RENTAL',
                LastSync: new Date(),
            },
            {
                Owner: TableOwner.SLAVE,
                MasterTableName: 'Payment',
                SlaveTableName: 'PAYMENT',
                LastSync: new Date(),
            }
        ]
    }
} as SlaveMasterPair);

async function getPair() {
    const pairName = props.PairName;
    try {
        let result = await Get<SlaveMasterPair>(`/api/pairs/${pairName}`);
        pair.value = result;
    } catch (ex) {
        alert(ex);
    }
}
// getPair();
</script>

<template>
    <div class="grid">
        <div>
            <h5>Master</h5>
            <p>Database: {{ pair.Master.Database }} <b>[{{ DatabaseType[pair.Master.DbType] }}]</b></p>
        </div>
        <div>
            <h5>Slave</h5>
            <p>Database: {{ pair.Slave.Database }} <b>[{{ DatabaseType[pair.Slave.DbType] }}]</b></p>
        </div>
    </div>

    <div>
        <section v-for="table in pair.Mappings.Tables">
            <fieldset>
                <div class="flex">
                    <p>{{ table.Owner == TableOwner.MASTER ? table.MasterTableName : table.SlaveTableName }} | Owner: {{ TableOwner[table.Owner]}}</p>
                    <p data-tooltip="Last Sync">{{ new Date(table.LastSync).toLocaleString() }}</p>
                    <p>{{ table.Owner == TableOwner.MASTER ? table.MasterTableName : table.SlaveTableName }} | Owner: {{ TableOwner[table.Owner]}}</p>
                </div>
                <progress />
            </fieldset>
        </section>
    </div>
</template>

<style lang="css" scoped>
.grid div {
    justify-items: center;
    border-bottom: solid 3px darkgray;
}

.grid {
    padding-bottom: 15px;
}

fieldset {
    justify-items: left;
}

section {
    border-bottom: solid 3px gray;
}

section:last-of-type {
    border-bottom: none;
}

.flex {
    display: flex;
    justify-content: space-between;
    width: 100%;
}
</style>
