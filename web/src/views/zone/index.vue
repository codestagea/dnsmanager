<script lang="ts">
export default {
  name: 'zone'
};
</script>

<script setup lang="ts">
import { ElForm } from 'element-plus';
import { reactive, ref, onMounted, getCurrentInstance, toRefs } from 'vue';
import { apiListZones } from '@/api/zone';
import { useRouter } from 'vue-router';
import ZoneForm from './zone-form.vue';
const state = reactive({
  loading: false,
  dnsList: [],
  total: 0,
  selectedDns: [],
  queryParams: {
    search: '',
    pageNum: 1,
    pageSize: 10
  }
});

const queryFormRef = ref(ElForm); // 查询表单

function handleQuery() {
  state.loading = true;
  apiListZones(state.queryParams).then(data => {
    state.loading = false;
    state.dnsList = data.records;
    state.total = data.total;
  });
}

function resetQuery() {
  queryFormRef.value.resetFields();
  handleQuery();
}

const router = useRouter();

function showRecord(record) {
  router.push({ name: 'dns_record', params: { zoneId: record.id } });
}

function showHistory(record) {
  router.push({ name: 'target_operation_log', params: { targetType: 'zone', targetId: record.id } });
}

const zoneFormRef = ref<InstanceType<typeof ZoneForm>>();
const showZoneForm = record => zoneFormRef.value.showZoneForm(record);

onMounted(() => {
  // 初始化用户列表数据
  handleQuery();
});
</script>

<template>
  <div class="app-container">
    <div class="search">
      <el-form ref="queryFormRef" :model="queryParams" :inline="true">
        <el-form-item label="关键字" prop="search">
          <el-input
            v-model="state.queryParams.search"
            placeholder="域名"
            clearable
            style="width: 200px"
            @keyup.enter="handleQuery"
          />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select v-model="state.queryParams.status" placeholder="全部" clearable style="width: 200px">
            <el-option label="启用" value="1" />
            <el-option label="禁用" value="0" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
          <el-button :icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-card shadow="never">
      <template #header>
        <div class="flex justify-between">
          <div>
            <el-button type="success" :icon="Plus" @click="showZoneForm(undefined)">新增</el-button>
          </div>
        </div>
      </template>

      <el-table
        v-loading="state.loading"
        :data="state.dnsList"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column key="id" label="ID" align="center" prop="id" width="50" />
        <el-table-column key="zone" label="域名" align="center" prop="zone" />
        <el-table-column label="刷新时间(秒)" prop="refresh" width="100" />
        <el-table-column label="重试时间(秒)" prop="retry" width="100" />
        <el-table-column label="刷新时间(秒)" prop="expire" width="100" />
        <el-table-column label="最小TTL(秒)" prop="minimum" width="100" />
        <el-table-column label="邮件联系人" prop="hostMaster" show-overflow-tooltip />
        <el-table-column label="权威DNS" prop="primaryNs" show-overflow-tooltip />
        <el-table-column label="描述" prop="remark" show-overflow-tooltip />
        <el-table-column label="状态" align="center" prop="state" fixed="right" width="80">
          <template #default="scope">
            <el-tag v-if="scope.row.state === 'running'" type="success" class="mx-1" effect="dark">解析中</el-tag>
            <el-tag v-else-if="scope.row.state === 'stopped'" type="danger" class="mx-1" effect="dark">停用</el-tag>
            <el-tag v-else>{{ scope.row.state }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" align="left" width="140" fixed="right">
          <template #default="scope">
            <el-button type="success" link @click="showHistory(scope.row)">历史</el-button>
            <el-button type="primary" link @click="showRecord(scope.row)">记录</el-button>
            <el-button type="danger" link @click="showZoneForm(scope.row)">修改</el-button>
          </template>
        </el-table-column>
      </el-table>

      <pagination
        v-if="state.total > 0"
        :total="state.total"
        v-model:page="state.queryParams.pageNum"
        v-model:limit="state.queryParams.pageSize"
        @pagination="handleQuery"
      />
    </el-card>
    <zone-form ref="zoneFormRef" @update-zone="handleQuery" />
  </div>
</template>

<style></style>
