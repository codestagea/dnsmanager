<script lang="ts">
export default {
  name: 'ZoneForm'
};
</script>
<script setup lang="ts">
import { ElForm, ElMessage } from 'element-plus';
import { QuestionFilled } from '@element-plus/icons-vue';
import { reactive, toRefs, ref } from 'vue';
import { cloneDeep } from 'lodash-es';
import { t } from '@/utils/i18n';
import { apiAddZone, apiUpdateZone } from '@/api/zone';
const state = reactive({
  dialog: {
    title: '',
    visible: false
  },
  formData: {
    id: undefined,
    zone: '',
    refresh: 600,
    retry: 600,
    expire: 600,
    minimum: 10800,
    hostMaster: '',
    primaryNs: '',
    remark: '',
    state: ''
  }
});

const rules = {
  zone: [{ required: true, message: '域名不能为空', trigger: 'blur' }],
  refresh: [{ required: true, message: '请输入', trigger: 'blur' }],
  retry: [{ required: true, message: '请输入', trigger: 'blur' }],
  expire: [{ required: true, message: '请输入', trigger: 'blur' }],
  minimum: [{ required: true, message: '请输入', trigger: 'blur' }],
  hostMaster: [{ required: true, message: '请输入', trigger: 'blur' }],
  primaryNs: [{ required: true, message: '请输入', trigger: 'blur' }]
};
const emit = defineEmits(['update-zone']);

const dataFormRef = ref(ElForm);

const showZoneForm = record => {
  if (record) {
    state.formData = cloneDeep(record);
    state.dialog.title = t('common.edit') + `-${record.zone}`;
  } else {
    state.formData = {
      id: undefined,
      zone: '',
      refresh: 600,
      retry: 600,
      expire: 600,
      minimum: 10800,
      hostMaster: '',
      primaryNs: '',
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
        apiUpdateZone(id, state.formData).then(() => {
          ElMessage.success('修改域名成功');
          closeDialog();
          emit('update-zone');
        });
      } else {
        apiAddZone(state.formData).then(() => {
          ElMessage.success('新增域名成功');
          closeDialog();
          emit('update-zone');
        });
      }
    }
  });
}

const { dialog, formData } = toRefs(state);

defineExpose({
  showZoneForm
});
</script>
<template>
  <div>
    <el-dialog :title="dialog.title" v-model="dialog.visible" width="600px" append-to-body @close="closeDialog">
      <el-form ref="dataFormRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="域名" prop="zone">
          <el-input :disabled="!!formData.id" v-model="formData.zone" placeholder="请输入域名" maxlength="100" />
        </el-form-item>
        <el-form-item label="描述" prop="remark">
          <el-input v-model="formData.remark" placeholder="请输入描述" maxlength="100" />
        </el-form-item>
        <el-form-item label="刷新时间(秒)" prop="refresh">
          <el-input-number v-model="formData.refresh" :step="60" :min="60" controls-position="right" />
          <el-tooltip content="如果zone有修改，辅名称服务器向主名称服务器请求同步的时间间隔" placement="bottom-start">
            <el-icon class="mx-1"> <QuestionFilled /> </el-icon>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="重试时间(秒)" prop="retry">
          <el-input-number v-model="formData.retry" :step="60" :min="60" controls-position="right" />
          <el-tooltip
            content="如果主名称服务器没有响应，辅名称服务器根据此值确定等待多久后提交刷新请求"
            placement="bottom-start"
          >
            <el-icon class="mx-1"> <QuestionFilled /> </el-icon>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="刷新时间(秒)" prop="expire">
          <el-input-number v-model="formData.expire" :step="60" controls-position="right" />
          <el-tooltip
            content="如果在该时间到期前，主名称服务器没有响应刷新请求，辅名称服务器将停止该namespace的权威响应"
            placement="bottom-start"
          >
            <el-icon class="mx-1"> <QuestionFilled /> </el-icon>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="最小TTL(秒)" prop="minimum">
          <el-input-number v-model="formData.minimum" :step="60" controls-position="right" />
          <el-tooltip
            content="在 BIND 9 中，它定义缓存negative answers（否定回答）的时间。negative answers缓存最长可设定为 3 小时（即 3H）"
            placement="bottom-start"
          >
            <el-icon class="mx-1"> <QuestionFilled /> </el-icon>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="邮件联系人" prop="hostMaster">
          <el-input v-model="formData.hostMaster" maxlength="64" class="form-item-without-icon" />
          <el-tooltip content="该namespace联系人邮件" placement="bottom-start">
            <el-icon class="mx-1"> <QuestionFilled /> </el-icon>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="权威DNS" prop="primaryNs">
          <el-input v-model="formData.primaryNs" maxlength="64" class="form-item-without-icon" />
          <el-tooltip content="该域权威DNS的主名称服务器的主机名" placement="bottom-start">
            <el-icon class="mx-1"> <QuestionFilled /> </el-icon>
          </el-tooltip>
        </el-form-item>
        <el-form-item label="状态" prop="state">
          <el-switch v-model="formData.state" inactive-value="stopped" active-value="running" />
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
