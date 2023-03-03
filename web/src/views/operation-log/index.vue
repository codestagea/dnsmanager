<script lang="ts">
export default {
  name: 'OperationLog'
};
</script>
<script setup lang="ts">
import { ref } from 'vue';
import { ElForm } from 'element-plus';
import LogTable from './log-table.vue';
const queryParams = ref({ search: '', targetType: '', pageNum: 1, pageSize: 10 });

const logTableRef = ref<InstanceType<typeof LogTable>>();

const handleQuery = () => {
  logTableRef.value.handleQuery();
};

const queryFormRef = ref(ElForm); // 查询表单
const resetQuery = () => {
  queryFormRef.value.resetFields();
  handleQuery();
};
</script>
<template>
  <div class="app-container">
    <div class="search">
      <el-form ref="queryFormRef" :model="queryParams" :inline="true">
        <el-form-item label="关键字" prop="search">
          <el-input
            v-model="queryParams.search"
            placeholder="域名"
            clearable
            style="width: 200px"
            @keyup.enter="handleQuery"
          />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-select v-model="queryParams.targetType" placeholder="全部" clearable style="width: 200px">
            <el-option label="域名" value="domain" />
            <el-option label="记录" value="record" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
          <el-button :icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <log-table :queryParams="queryParams" ref="logTableRef"></log-table>
  </div>
</template>
