<script setup lang='ts'>
import { ref,onMounted } from 'vue'
import {useUserStore} from "@/store"
import {updateUserPassword} from "@/api/user"
import { useI18n } from 'vue-i18n'
// 创建 i18n 实例
const { t } = useI18n()
const userinfo = ref()
const userStore = useUserStore()
// 获取用户信息
onMounted( async() => {
  userinfo.value = await userStore.getUserInfo()
})
const username:Ref<string> = ref('')
const password:Ref<string> = ref('')


/** 重置密码 */
function resetPassword(row: { [key: string]: any }) {
  if (!username.value || !password.value) {
    ElMessage.error(t('userset.message.xx1'))
    return
  }
  if ((password.value.length < 6)) {
    ElMessage.error(t('userset.message.xx2'))
    return
  }
  ElMessageBox.confirm(
    t('userset.message.xx3'),
    t('userset.message.title'),
    {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  )
    .then(() => {
      updateUserPassword({
          username:username.value.trim(),
          password:password.value.trim()

        }
      ).then(() => {
        ElMessage.success(t('userset.message.xx4') + password.value);
        window.location.reload();
      });
    })
}
</script>

<template>
  <div class="page-container userset-page">
    <el-card class="userset-card" shadow="never">
      <el-row :gutter="16" justify="center">
        <el-col :xs="22" :sm="18" :md="14" :lg="12">
          <h2 class="userset-title">{{ $t('userset.title') }}</h2>
        </el-col>
        <el-col :xs="22" :sm="18" :md="14" :lg="12" v-if="userinfo" class="userset-avatar-wrap">
          <el-badge :value="userinfo.username" class="item">
            <el-image :src="userinfo.avatar" class="userset-avatar" />
          </el-badge>
        </el-col>
        <el-col :xs="22" :sm="18" :md="14" :lg="12">
          <el-input
            v-model="username"
            :placeholder="$t('userset.newUsername')"
            size="large"
            clearable
          />
        </el-col>
        <el-col :xs="22" :sm="18" :md="14" :lg="12">
          <el-input
            v-model="password"
            type="password"
            :placeholder="$t('userset.newPassword')"
            size="large"
            show-password
            clearable
          />
        </el-col>
        <el-col :xs="22" :sm="18" :md="14" :lg="12" class="userset-actions">
          <el-button type="primary" size="large" @click="resetPassword">修改</el-button>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<style scoped>
.userset-page {
  max-width: 520px;
  margin: 0 auto;
}

.userset-card {
  text-align: center;
}

.userset-title {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
  color: var(--sl-text);
}

.userset-avatar-wrap {
  margin-bottom: 8px;
}

.userset-avatar {
  width: 96px;
  height: 96px;
  border-radius: 12px;
}

.userset-actions {
  margin-top: 8px;
}

.el-col {
  margin-bottom: 12px;
}
</style>
