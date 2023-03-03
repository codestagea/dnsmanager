<script lang="ts">
export default {
  name: 'LogTable'
};
</script>
<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { apiListOperationLog } from '@/api/operation-log';

const props = defineProps({
  queryParams: {
    required: true,
    type: Object,
    default: {}
  }
});

const operationLogs = ref([]);
const loading = ref<boolean>(false);
const total = ref<number>(0);

function handleQuery() {
  loading.value = true;
  apiListOperationLog(props.queryParams)
    .then(data => {
      total.value = data.total;
      operationLogs.value = data.records;
      loading.value = false;
    })
    .catch(() => (loading.value = false));
}
defineExpose({
  handleQuery
});

onMounted(() => {
  // 初始化用户列表数据
  handleQuery();
});
</script>
<template>
  <el-card shadow="never">
    <el-table v-loading="loading" :data="operationLogs">
      <el-table-column label="操作类型" align="center" prop="type" width="100">
        <template #default="scope">
          <span v-if="scope.row.type === 'add'">新增</span>
          <span v-else-if="scope.row.type === 'update'">修改</span>
          <span v-else>{{ scope.row.type }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作目标" align="center" prop="targetType" width="100">
        <template #default="scope">
          <span v-if="scope.row.targetType === 'zone'">区域</span>
          <span v-else-if="scope.row.targetType === 'domain'">域名</span>
          <span v-else-if="scope.row.targetType === 'record'">记录</span>
          <span v-else>{{ scope.row.targetType }}</span>
        </template>
      </el-table-column>
      <el-table-column label="时间" align="center" prop="createdAt" width="180"></el-table-column>
      <el-table-column label="操作者" align="center" prop="operator" width="180"></el-table-column>
      <el-table-column label="操作对象" width="200" prop="keyValue"> </el-table-column>
      <el-table-column label="内容" prop="diff"> </el-table-column>
    </el-table>

    <pagination
      v-if="total > 0"
      :total="total"
      v-model:page="queryParams.pageNum"
      v-model:limit="queryParams.pageSize"
      @pagination="handleQuery"
    />
  </el-card>
</template>
