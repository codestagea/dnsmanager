<script lang="ts">
export default {
  name: 'DnsView'
};
</script>

<script setup lang="ts">
import { reactive, ref, onMounted, getCurrentInstance, toRefs } from 'vue';
import { apiListViews } from '@/api/dns-view';
import { useRoute, useRouter } from 'vue-router';
import DnsViewForm from './view-form.vue';
const state = reactive({
  loading: false,
  views: []
});

const { params } = useRoute();

function handleQuery() {
  state.loading = true;
  apiListViews()
    .then(data => {
      state.loading = false;
      state.views = data;
    })
    .catch(() => (state.loading = false));
}

function resetQuery() {
  handleQuery();
}

const router = useRouter();

function showHistory(record) {
  router.push({ name: 'target_operation_log', params: { targetType: 'view', targetId: record.id } });
}

const viewFormRef = ref<InstanceType<typeof DnsViewForm>>();
const showViewForm = record => viewFormRef.value.showViewForm(record);
onMounted(() => {
  handleQuery();
});
</script>

<template>
  <div class="app-container">
    <el-card shadow="never">
      <template #header>
        <div class="flex justify-between">
          <div>
            <el-button type="success" :icon="Plus" @click="showViewForm(undefined)">新增</el-button>
          </div>
        </div>
      </template>

      <el-table v-loading="state.loading" :data="state.views" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column key="id" label="ID" align="center" prop="id" width="100" />
        <el-table-column key="name" label="名称" align="center" prop="name" />
        <el-table-column label="描述" prop="description" />
        <el-table-column label="创建时间" align="center" prop="createdAt" width="180"></el-table-column>
        <el-table-column label="操作" align="left" width="200">
          <template #default="scope">
            <el-button type="success" link @click="showHistory(scope.row)">历史</el-button>
            <el-button type="warning" link @click="showViewForm(scope.row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <dns-view-form ref="viewFormRef" @update-dns="handleQuery" />
  </div>
</template>
