<script lang="ts">
export default {
  name: ''
};
</script>
<script setup lang="ts">
import { useRoute } from 'vue-router';
import { ref } from 'vue';
import LogTable from './log-table.vue';

const { params } = useRoute();

const queryParams = ref({ targetType: params.targetType, targetId: params.targetId, pageNum: 1, pageSize: 10 });

const logTableRef = ref<InstanceType<typeof LogTable>>();
</script>
<template>
  <div class="app-container"><log-table :queryParams="queryParams" ref="logTableRef" /></div>
</template>
