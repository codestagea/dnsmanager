<script lang="ts">
export default {
  name: 'RecordForm'
};
</script>
<script setup lang="ts">
import { ElForm, ElMessage } from 'element-plus';
import { QuestionFilled } from '@element-plus/icons-vue';
import { reactive, toRefs, ref, onMounted } from 'vue';
import { cloneDeep } from 'lodash-es';
import { t } from '@/utils/i18n';
import { apiAddRecord, apiUpdateRecord } from '@/api/dns-records';
import { apiListViews } from '@/api/dns-view';
const state = reactive({
  dialog: {
    title: '',
    visible: false
  },
  zoneId: 0,
  formData: {
    id: undefined,
    host: '',
    type: '',
    data: '',
    ttl: 600,
    mx: 0,
    view: '',
    remark: '',
    state: ''
  },
  viewOpts: []
});

const rules = {
  host: [
    {
      validator: (rule: any, value: any, callback: any) => {
        if (state.formData.type !== 'NS' && !state.formData.host) {
          callback(new Error('非NS记录需要输入host'));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ],
  type: [{ required: true, message: '请选择', trigger: 'blur' }],
  data: [
    {
      validator: (rule: any, value: any, callback: any) => {
        if (
          state.formData.type === 'A' &&
          !/^(?!0)(?!.*\.$)((1?\d?\d|25[0-5]|2[0-4]\d)(\.|$)){4}$/.test(state.formData.data)
        ) {
          callback(new Error('A记录记录值必须为ip'));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ],
  ttl: [{ required: true, message: '请输入', trigger: 'blur' }],
  mx: [
    {
      validator: (rule: any, value: any, callback: any) => {
        if (state.formData.type == 'MX' && state.formData.mx <= 0) {
          callback(new Error('MX记录优先级必须大于0'));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ],
  view: [{ required: true, message: '请选择', trigger: 'blur' }]
};

const typeOpts = ['A', 'MX', 'CNAME', 'NS', 'PTR', 'TXT', 'AAAA', 'SVR', 'URL'];
const emit = defineEmits(['update-record']);

const dataFormRef = ref(ElForm);

const showRecordForm = (zoneId, record) => {
  state.zoneId = zoneId;
  if (record) {
    state.formData = cloneDeep(record);
    state.dialog.title = t('common.edit') + `-${record.host}`;
  } else {
    state.formData = {
      id: undefined,
      host: '',
      type: '',
      data: '',
      ttl: 600,
      mx: 0,
      view: '',
      remark: '',
      state: 'running'
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
        apiUpdateRecord(state.zoneId, id, state.formData).then(() => {
          ElMessage.success('修改域名成功');
          closeDialog();
          emit('update-record');
        });
      } else {
        apiAddRecord(state.zoneId, state.formData).then(() => {
          ElMessage.success('新增域名成功');
          closeDialog();
          emit('update-record');
        });
      }
    }
  });
}

const { dialog, formData, viewOpts } = toRefs(state);

onMounted(() => {
  apiListViews().then(res => {
    state.viewOpts = res;
  });
});

defineExpose({
  showRecordForm
});
</script>
<template>
  <div>
    <el-dialog :title="dialog.title" v-model="dialog.visible" width="600px" append-to-body @close="closeDialog">
      {{ viewOpts }}
      <el-form ref="dataFormRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="主机记录" prop="host">
          <el-input v-model="formData.host" placeholder="请输入主机记录" maxlength="255" />
        </el-form-item>
        <el-form-item label="记录类型" prop="type">
          <el-select v-model="formData.type" placeholder="Select">
            <el-option v-for="item in typeOpts" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="记录值" prop="data">
          <el-input v-model="formData.data" maxlength="255" />
        </el-form-item>
        <el-form-item label="TTL(秒)" prop="ttl">
          <el-input-number v-model="formData.ttl" :step="60" controls-position="right" />
        </el-form-item>
        <el-form-item label="优先级" prop="mx">
          <el-input-number
            v-model="formData.mx"
            :step="1"
            :min="1"
            controls-position="right"
            :disabled="formData.type !== 'MX'"
          />
        </el-form-item>
        <el-form-item label="视图" prop="view">
          <el-select v-model="formData.view" placeholder="Select">
            <el-option v-for="item in viewOpts" :key="item.name" :label="item.name" :value="item.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="state">
          <el-switch v-model="formData.state" inactive-value="stopped" active-value="running" />
        </el-form-item>
        <el-form-item label="描述" prop="remark">
          <el-input v-model="formData.remark" placeholder="请输入描述" maxlength="100" />
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

<style lang="scss">
.form-item-without-icon {
  width: -webkit-calc(100% - 20px) !important;
  width: -moz-calc(100% - 20px) !important;
  width: calc(100% - 20px) !important;
}
</style>
