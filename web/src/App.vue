<script setup lang="ts">
import { ElConfigProvider } from 'element-plus';
import { useAppStore } from '@/store/app';

const appStore = useAppStore();
appStore.setDark(appStore.isDark);
</script>

<template>
  <el-config-provider :locale="appStore.locale" :size="appStore.size">
    <router-view />
  </el-config-provider>
</template>
