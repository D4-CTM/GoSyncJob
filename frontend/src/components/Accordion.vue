<script setup lang="ts">
const props = defineProps({
    Name: {
        type: String,
        required: true
    },
    Summary: {
        type: String,
        required: true
    }
});
</script>

<template>
    <details role="button" 
             class="outline contrast"
             :name="Name" >
        <summary>
            {{ Summary }}
        </summary>
        <slot />
    </details>
</template>
