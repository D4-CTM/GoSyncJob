<script setup lang="ts">
import { ref } from 'vue';
import Accordion from './components/Accordion.vue';
import CreateSyncJob from './synchronizer/CreateSyncJob.vue';
import { Get } from './helper/HttpCaller';
import SyncJobItem from './synchronizer/SyncJobItem.vue';

let pairNames = ref([] as string[])

async function getPairs() {
    try {
        // let result = await Get<string[]>(`/api/pairs`);
        // pairNames.value = result;
        pairNames.value = [
            'Test Mapper'
        ]
    } catch (ex) {
        alert(ex);
    }
}
getPairs();
</script>

<template>
    <main>
        <header>
            <h1>Synchronizer</h1>
            <!-- <CreateSyncJob/> -->
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
</style>
