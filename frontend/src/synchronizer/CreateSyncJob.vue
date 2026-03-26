<script setup lang="ts">
import Modal from '@/components/Modal.vue';
import { DbCredential } from '@/models/Credentials';
import { ref } from 'vue';
import CredentialForm from './CredentialForm.vue';

let slave = ref({} as DbCredential)
let master = ref({} as DbCredential)

let currentStep = ref(0)

const steps = [
    'Create Slave',
    'Create Master',
    'Create Synchronizer',
]

const verifyCredentials = (cred: DbCredential) =>
    cred.Database != '' &&
    cred.Server != '' &&
    cred.Port > 0 &&
    cred.Password != '' &&
    cred.User != '';

function confirm(closeFn: Function) {
    switch (currentStep.value) {
        case 0: {        
            if (verifyCredentials(slave.value))
                currentStep.value++
            else alert('Invalid slave credentials')
        } break;

        case 1: {
            if (verifyCredentials(master.value))
                currentStep.value++
            else alert('Invalid master credentials')
        } break;

        case 2: {
            currentStep.value = 0;
            closeFn();
        } break;
    }
}

function cancel(closeFn: Function) {
    currentStep.value = 0;
    closeFn();
}
</script>

<template>
    <Modal Title="Create Synchronizer"
           BtnText="Create job"
           :ConfirmText="steps[currentStep]"
           @confirm="confirm"
           @cancel="cancel">
        <CredentialForm v-if="currentStep === 0"
                        :credential="slave" />
        <CredentialForm v-if="currentStep === 1"
                        :credential="master" />
    </Modal>
</template>

<style lang="css" scoped>
</style>
