<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessageBox } from 'element-plus';

import Hamburger from '@/components/Hamburger/index.vue';
import Breadcrumb from '@/components/Breadcrumb/index.vue';
// import SizeSelect from '@/components/SizeSelect/index.vue';
import LangSelect from '@/components/LangSelect/index.vue';
import ThemeSelect from '@/components/ThemeSelect/index.vue';
import MixNav from './Sidebar/MixNav.vue';
import { CaretBottom } from '@element-plus/icons-vue';

import { useAppStore } from '@/store/app';
import { useUserStore } from '@/store/user';

const appStore = useAppStore();
const userStore = useUserStore();

const route = useRoute();
const router = useRouter();

const device = computed(() => appStore.device);

function toggleSideBar() {
  appStore.toggleSidebar(true);
}

function logout() {
  ElMessageBox.confirm('确定注销并退出系统吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    userStore
      .logout()
      .then(() => {})
      .then(() => {
        router.push(`/login?redirect=${route.fullPath}`);
      });
  });
}
</script>

<template>
  <div class="navbar">
    <div class="flex justify-start" v-if="device === 'mobile' || appStore.layout === 'left'">
      <hamburger :is-active="appStore.sidebar.opened" @toggleClick="toggleSideBar" />
      <!-- 面包屑导航栏 -->
      <breadcrumb />
    </div>

    <mix-nav v-if="device !== 'mobile' && appStore.layout === 'mix'" />

    <div v-if="device === 'mobile' || appStore.layout === 'left'" class="flex justify-start">
      <div v-if="device !== 'mobile'" class="flex justify-center items-center pr-[5px]">
        <!-- 布局大小 -->
        <!-- <el-tooltip content="布局大小" effect="dark" placement="bottom">
          <size-select />
        </el-tooltip> -->
        <theme-select style="margin-right: 5px" />
        <!--语言选择-->
        <lang-select />
      </div>

      <el-dropdown trigger="click">
        <div class="flex justify-center items-center pr-[20px]">
          {{ userStore.userInfo.nickname || userStore.userInfo.username }}
          <CaretBottom class="w-3 h-3" />
        </div>

        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item divided @click="logout">
              {{ $t('navbar.logout') }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 0px 2px rgba(0, 0, 0, 0.2);
}
</style>
