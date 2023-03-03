<script setup lang="ts">
import { Sunny, Moon } from '@element-plus/icons-vue';
import { useAppStore } from '@/store/app';
import { computed } from 'vue';

const appStore = useAppStore();

const darkMode = computed(() => appStore.isDark);

function toggleTheme() {
  appStore.setDark(darkMode.value);
}
</script>

<template>
  <div @click.stop>
    <el-switch
      size="default"
      v-model="appStore.isDark"
      inline-prompt
      @change="toggleTheme"
      :active-icon="Sunny"
      :inactive-icon="Moon"
    />
  </div>
</template>
