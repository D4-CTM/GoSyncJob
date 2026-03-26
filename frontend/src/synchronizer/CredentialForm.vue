<script setup lang="ts">
import EyeIcon from '@/assets/icons/EyeIcon.vue';
import EyeSlashIcon from '@/assets/icons/EyeSlashIcon.vue';
import { DatabaseType, DbCredential } from '@/models/Credentials';
import { ref } from 'vue';

const props = defineProps<{
    credential: DbCredential
}>()

props.credential.DbType = DatabaseType.ORACLE;
let showPass = ref(false);
</script>

<template>
    <form>
        <fieldset>
            <label>
                Select database provider...
                <select name="provider">
                    <option selected @onSelect="DatabaseType.ORACLE">Oracle</option>
                    <option @onSelect="DatabaseType.POSTGRES">Postgres</option>
                </select>
            </label>
            <label>
                Database
                <input name="database"
                       placeholder="public"
                       autocomplete="database"
                       type="text"
                       v-model.trim="credential.Database" />
            </label>
            <div class="grid">
                <label>
                    Server
                    <input name="server"
                           placeholder="localhost"
                           autocomplete="server"
                           type="text"
                           v-model.trim="credential.Server"/>
                </label>
                <label>
                    Port
                    <input name="port"
                           placeholder="5412"
                           autocomplete="port"
                           type="number"
                           v-model.trim="credential.Port" />
                </label>
            </div>
            <label>
                User
                <input name="user"
                       placeholder="user"
                       autocomplete="user"
                       type="text"
                       v-model.trim="credential.User"/>
            </label>
            <label>
                Password
            <fieldset role="group">
                <input name="password"
                       placeholder="password"
                       autocomplete="password"
                       :type="showPass ? 'text' : 'password'"
                       v-model.trim="credential.Password"/>
                <button type="button" class="img-btn"
                        @click="showPass = !showPass">
                    <EyeIcon v-if="!showPass"/>
                    <EyeSlashIcon v-else/>
                </button>
            </fieldset>
            </label>
        </fieldset>
    </form>
</template>
