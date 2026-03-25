<script setup lang="ts">
import { useTemplateRef } from 'vue';

const props = defineProps({
    Title: {
        type: String,
        required: true
    },
    BtnText: {
        type: String,
        required: true
    },
    ConfirmText: {
        type: String,
        required: true
    }
});

const dialogRef = useTemplateRef("dialogRef")

const emit = defineEmits(['confirm'])
function confirm() {
    emit('confirm', () => dialogRef.value.close())
}

dialogRef.value.addEventListener('cancel', e => {
  e.preventDefault()
})
</script>

<template>
    <button @click="dialogRef.showModal">
        {{ BtnText }}
    </button>

    <dialog ref="dialogRef">
        <article>
            <h2>{{ Title }}</h2>
            <slot/>
            <footer>
                <button class="secondary"
                        @click="dialogRef.close()">
                    Cancel
                </button>
                <button @click="confirm">
                    {{ ConfirmText }}
                </button>
            </footer>
        </article>
    </dialog>
</template>
