<script lang="ts">
export default {
  name: 'DnsRecord'
};
</script>

<script setup lang="ts">
import { ElForm } from 'element-plus';
import { reactive, ref, onMounted, getCurrentInstance, toRefs } from 'vue';
import { apiListRecords } from '@/api/dns-records';
import { useRoute, useRouter } from 'vue-router';
import RecordForm from './record-form.vue';
const state = reactive({
  loading: false,
  zoneId: 0,
  records: [],
  total: 0,
  selected: [],
  queryParams: {
    search: '',
    pageNum: 1,
    pageSize: 10
  }
});

const { params } = useRoute();

const queryFormRef = ref(ElForm); // 查询表单

function handleQuery() {
  state.loading = true;
  apiListRecords(state.zoneId, state.queryParams).then(data => {
    state.loading = false;
    state.records = data.records;
    state.total = data.total;
  });
}

function resetQuery() {
  queryFormRef.value.resetFields();
  handleQuery();
}

const router = useRouter();

function showHistory(record) {
  router.push({ name: 'target_operation_log', params: { targetType: 'record', targetId: record.id } });
}

const recordFormRef = ref<InstanceType<typeof RecordForm>>();
const showRecordForm = record => recordFormRef.value.showRecordForm(state.zoneId, record);

onMounted(() => {
  // 初始化用户列表数据
  state.zoneId = params.zoneId;
  handleQuery();
});
</script>

<template>
  <div class="app-container">
    <div class="search">
      <el-form ref="queryFormRef" :model="state.queryParams" :inline="true">
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
            <el-button type="success" :icon="Plus" @click="showRecordForm(undefined)">新增</el-button>
          </div>
        </div>
      </template>

      <el-table v-loading="state.loading" :data="state.records" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column key="id" label="ID" align="center" prop="id" width="100" />
        <el-table-column key="host" label="主机记录" align="center" prop="host" />
        <el-table-column label="记录类型" prop="type" />
        <el-table-column label="记录值" prop="data" />
        <el-table-column label="状态" align="center" prop="state">
          <template #default="scope">
            <el-tag v-if="scope.row.state === 'running'" type="success" class="mx-1" effect="dark">解析中</el-tag>
            <el-tag v-else-if="scope.row.state === 'stopped'" type="danger" class="mx-1" effect="dark">停用</el-tag>
            <el-tag v-else>{{ scope.row.state }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="权重" prop="mx" width="80">
          <template #default="scope">
            {{ scope.row.type == 'mx' ? scope.row.mx : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="TTL" prop="ttl" width="80" />
        <el-table-column label="记录值" prop="data" />
        <el-table-column label="描述" prop="remark" />
        <el-table-column label="创建时间" align="center" prop="createdAt" width="180"></el-table-column>
        <el-table-column label="操作" align="left" width="200">
          <template #default="scope">
            <el-button type="success" link @click="showHistory(scope.row)">历史</el-button>
            <el-button type="warning" link @click="showRecordForm(scope.row)">编辑</el-button>
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
      <record-form ref="recordFormRef" @update-record="handleQuery" />
    </el-card>
  </div>
</template>
