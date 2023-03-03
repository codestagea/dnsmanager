<script lang="ts">
export default {
  name: 'DnsViewForm'
};
</script>
<script setup lang="ts">
import { ElForm, ElMessage } from 'element-plus';
import { reactive, toRefs, ref } from 'vue';
import { cloneDeep } from 'lodash-es';
import { t } from '@/utils/i18n';
import { apiAddView, apiUpdateView } from '@/api/dns-view';
const state = reactive({
  dialog: {
    title: '',
    visible: false
  },
  formData: {
    id: undefined,
    name: '',
    description: ''
  }
});

const emit = defineEmits(['update-dns']);

const dataFormRef = ref(ElForm);

const showViewForm = record => {
  console.log('record', record);
  if (record) {
    state.formData = cloneDeep(record);
    state.dialog.title = t('common.edit') + `-${record.name}`;
  } else {
    state.formData = {
      id: undefined,
      name: '',
      description: ''
    };
    state.dialog.title = t('common.add');
  }
  state.dialog.visible = true;
};

const closeDialog = () => {
  state.dialog.visible = false;
  dataFormRef.value.resetFields();
  dataFormRef.value.clearValidate();
  state.formData.id = undefined;
};

/**
 * 表单提交
 */
function submitForm() {
  dataFormRef.value.validate((valid: any) => {
    if (valid) {
      const id = state.formData.id;
      if (id) {
        apiUpdateView(id, state.formData).then(() => {
          ElMessage.success('修改视图成功');
          closeDialog();
          emit('update-dns');
        });
      } else {
        apiAddView(state.formData).then(() => {
          ElMessage.success('新增视图成功');
          closeDialog();
          emit('update-dns');
        });
      }
    }
  });
}

const { dialog, formData } = toRefs(state);

defineExpose({
  showViewForm
});
</script>
<template>
  <div>
    <el-dialog :title="dialog.title" v-model="dialog.visible" width="600px" append-to-body @close="closeDialog">
      <el-form ref="dataFormRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="名字" prop="name">
          <el-input :disabled="!!formData.id" v-model="formData.name" placeholder="请输入名字" />
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="submitForm">确 定</el-button>
          <el-button @click="closeDialog">取 消</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>
