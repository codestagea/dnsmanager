<script setup lang="ts"></script>

<template>
  <section class="app-main">
    <router-view v-slot="{ Component, route }">
      <transition name="router-fade" mode="out-in">
        <keep-alive>
          <component :is="Component" :key="route.fullPath" />
        </keep-alive>
      </transition>
    </router-view>
  </section>
</template>

<style lang="scss" scoped>
.app-main {
  /* 50= navbar  50  */
  min-height: calc(100vh - 84px);
  width: 100%;
  position: relative;
  overflow: auto;
  background-color: var(--el-bg-color-page);
}

.fixed-header + .app-main {
  padding-top: 50px;
  min-height: 100vh;
}
</style>
