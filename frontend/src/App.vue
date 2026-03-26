<script setup lang="ts">
import { ref } from 'vue';
import { Get, Post } from './helper/HttpCaller';
import { TableOwner } from './models/TableMapping';
import { TriggerSyncDto } from './models/TriggerSyncDto';
import SyncJobItem from './synchronizer/SyncJobItem.vue';

let pairNames = ref([] as string[])

async function getPairs() {
    try {
        let result = await Get<string[]>(`/api/pairs`);
        pairNames.value = result;
    } catch (ex) {
        alert(ex);
    }
}
getPairs();

async function syncIn() {
    const pairName = pairNames.value[0];
    const dto: TriggerSyncDto = {
        owner: TableOwner.MASTER,
        all: true
    };
    try {
        const result = await Post<string, TriggerSyncDto>(`/api/pairs/${pairName}/sync`, dto);
        alert(result);
    } catch (ex) {
        alert(ex);
    }
}
</script>

<template>
    <main>
        <header>
            <h1>Synchronizer</h1>
            <div class="sync-container">
                <button class="secondary">Sync-OUT</button>
                <button @click="syncIn">Sync-IN</button>
            </div>
        </header>
        <div>
            <SyncJobItem v-for="pair in pairNames" :PairName="pair"/>           
        </div>
    </main>
</template>

<style lang="css" scoped>
main {
    padding: 10px;
    margin: 10px;
    max-width: 100vw;
}

header {
    display: flex;
    justify-content: space-between;
    padding-bottom: 15px;
}

.sync-container button {
    margin-left: 10px;
    margin-right: 10px;
}
</style>
