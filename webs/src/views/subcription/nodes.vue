<script setup lang='ts'>
import { ref, onMounted, nextTick, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search, Refresh, Plus, Download, Delete, Edit, CopyDocument } from '@element-plus/icons-vue'
import { getNodes, AddNodes, DelNode, UpdateNode } from "@/api/subcription/node"
import { getSubSchedulers, addSubScheduler, updateSubScheduler, deleteSubScheduler, type SubScheduler, type SubSchedulerRequest } from "@/api/subcription/scheduler"
import { ElMessage, ElMessageBox } from 'element-plus'

interface Node {
  ID: number;
  Name: string;
  Link: string;
  DialerProxyName: string;
  CreateDate: string;
}

const route = useRoute()
const router = useRouter()

const tableData = ref<Node[]>([])
const loading = ref(false)
const Nodelink = ref('')
const NodeOldlink = ref('')
const Nodename = ref('')
const NodeOldname = ref('')
const DialerProxyName = ref('')
const dialogVisible = ref(false)
const table = ref()
const NodeTitle = ref('')
const radio1 = ref('1')

const subSchedulerData = ref<SubScheduler[]>([])
const subSchedulerDialogVisible = ref(false)
const subSchedulerFormVisible = ref(false)
const subSchedulerForm = ref<SubSchedulerRequest>({
  name: '',
  url: '',
  cron_expr: '',
  enabled: true
})
const subSchedulerFormTitle = ref('')
const subSchedulerTable = ref()
const subSchedulerSelection = ref<SubScheduler[]>([])

const cronValidationStatus = ref<{ isValid: boolean, message: string }>({
  isValid: true,
  message: ''
})

const subCurrentPage = ref(1)
const subPageSize = ref(10)

const currentSubSchedulerData = computed(() => {
  const start = (subCurrentPage.value - 1) * subPageSize.value;
  const end = start + subPageSize.value;
  return subSchedulerData.value.slice(start, end);
})

async function getnodes() {
  loading.value = true;
  try {
    const { data } = await getNodes();
    tableData.value = data
  } catch (error) {
    console.error('获取节点列表失败:', error);
  } finally {
    loading.value = false;
  }
}

onMounted(async () => { getnodes() })

const handleAddNode = () => {
  NodeTitle.value = '添加节点'
  Nodelink.value = ''
  Nodename.value = ''
  radio1.value = '1'
  DialerProxyName.value = ''
  dialogVisible.value = true
}

const addnodes = async () => {
  let nodelinks = Nodelink.value.split(/[\r\n,]/)
    .map((item) => item.trim())
    .filter((item) => item !== '');

  if (NodeTitle.value == '添加节点') {
    if (radio1.value === '1') {
      if (Nodename.value.trim() === '') {
        ElMessage.error('备注不能为空')
        return
      }
      if (nodelinks.length > 0) {
        const processedLink = nodelinks.join(',');
        await AddNodes({
          link: processedLink,
          name: Nodename.value.trim(),
          dialerProxyName: DialerProxyName.value.trim(),
        })
      }
    } else {
      for (let i = 0; i < nodelinks.length; i++) {
        await AddNodes({
          link: nodelinks[i],
          name: "",
          dialerProxyName: DialerProxyName.value.trim(),
        })
      }
    }
    ElMessage.success("添加成功");
  } else {
    const processedLink = nodelinks.join(',');
    await UpdateNode({
      oldname: NodeOldname.value.trim(),
      oldlink: NodeOldlink.value.trim(),
      link: processedLink,
      name: Nodename.value.trim(),
      dialerProxyName: DialerProxyName.value.trim(),
    })
    ElMessage.success("更新成功");
  }
  getnodes()
  Nodelink.value = ''
  Nodename.value = ''
  dialogVisible.value = false;
}

const multipleSelection = ref<Node[]>([])
const handleSelectionChange = (val: Node[]) => { multipleSelection.value = val }

const searchQuery = ref('')
const handleSearch = () => { }

const filteredTableData = computed(() => {
  if (!searchQuery.value) return tableData.value;
  const query = searchQuery.value.toLowerCase();
  return tableData.value.filter(node =>
    node.Name.toLowerCase().includes(query) ||
    node.Link.toLowerCase().includes(query)
  );
});

const selectAll = () => {
  nextTick(() => {
    tableData.value.forEach(row => {
      table.value.toggleRowSelection(row, true)
    })
  })
}

const handleEdit = (row: any) => {
  radio1.value = '1'
  for (let i = 0; i < tableData.value.length; i++) {
    if (tableData.value[i].ID === row.ID) {
      NodeTitle.value = '编辑节点'
      Nodename.value = tableData.value[i].Name
      NodeOldname.value = Nodename.value
      Nodelink.value = tableData.value[i].Link
      NodeOldlink.value = Nodelink.value
      DialerProxyName.value = tableData.value[i].DialerProxyName
      dialogVisible.value = true
    }
  }
}

const toggleSelection = () => { table.value.clearSelection() }

const handleDel = (row: any) => {
  ElMessageBox.confirm(`确定要删除节点「${row.Name}」吗？`, '确认删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    await DelNode({ id: row.ID })
    ElMessage.success('删除成功')
    getnodes()
  })
}

const selectDel = () => {
  if (multipleSelection.value.length === 0) return
  ElMessageBox.confirm(`确定要删除选中的 ${multipleSelection.value.length} 个节点吗？`, '确认批量删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    for (let i = 0; i < multipleSelection.value.length; i++) {
      DelNode({ id: multipleSelection.value[i].ID })
      tableData.value = tableData.value.filter((item) => item.ID !== multipleSelection.value[i].ID)
    }
    ElMessage.success('删除成功')
  })
}

const currentPage = ref(1);
const pageSize = ref(10);
const handleSizeChange = (val: number) => { pageSize.value = val; }
const handleCurrentChange = (val: number) => { currentPage.value = val; }

const currentTableData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return filteredTableData.value.slice(start, end);
});

const copyUrl = async (url: string) => {
  try {
    await navigator.clipboard?.writeText(url);
    ElMessage.success('复制成功');
  } catch {
    const textarea = document.createElement('textarea');
    textarea.value = url;
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
    ElMessage.success('复制成功');
  }
};
const copyInfo = (row: any) => { copyUrl(row.Link) }

const getSubSchedulerList = async () => {
  try {
    const response = await getSubSchedulers()
    if (response) {
      subSchedulerData.value = response.data || []
    }
  } catch {
    ElMessage.error('获取订阅列表失败')
  }
}

const handleImportSubscription = () => {
  subSchedulerDialogVisible.value = true
  getSubSchedulerList()
}

const SCHEDULER_QUERY_ON = new Set(['1', 'true'])
function openSchedulerIfQuery() {
  if (route.path !== '/subcription/nodes') return
  const raw = route.query.scheduler
  const v = Array.isArray(raw) ? raw[0] : raw
  if (!SCHEDULER_QUERY_ON.has(String(v ?? ''))) return
  handleImportSubscription()
  router.replace({ path: '/subcription/nodes' })
}

watch(() => route.fullPath, openSchedulerIfQuery, { immediate: true })

const handleAddSubScheduler = () => {
  subSchedulerFormTitle.value = '添加订阅'
  subSchedulerForm.value = { name: '', url: '', cron_expr: '', enabled: true }
  subSchedulerFormVisible.value = true
}

const handleEditSubScheduler = (row: SubScheduler) => {
  subSchedulerFormTitle.value = '编辑订阅'
  subSchedulerForm.value = {
    id: row.ID, name: row.Name, url: row.URL,
    cron_expr: row.CronExpr, enabled: row.Enabled
  }
  subSchedulerFormVisible.value = true
}

const handleDeleteSubScheduler = (row: SubScheduler) => {
  ElMessageBox.confirm(`确定要删除订阅「${row.Name}」吗？`, '确认删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    await deleteSubScheduler(row.ID)
    ElMessage.success('删除成功')
    await getSubSchedulerList()
  }).catch(() => { })
}

const handleSubSchedulerSelectionChange = (val: SubScheduler[]) => {
  subSchedulerSelection.value = val
}

const handleBatchDeleteSubScheduler = () => {
  if (subSchedulerSelection.value.length === 0) {
    ElMessage.warning('请选择要删除的项目')
    return
  }
  ElMessageBox.confirm(`确定要删除选中的 ${subSchedulerSelection.value.length} 个订阅吗？`, '确认批量删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    await Promise.all(subSchedulerSelection.value.map(item => deleteSubScheduler(item.ID)))
    ElMessage.success('批量删除成功')
    await getSubSchedulerList()
  }).catch(() => { })
}

const validateCronExpression = (cron: string): boolean => {
  cron = cron.trim()
  const parts = cron.split(/\s+/)
  if (parts.length !== 5) return false
  const ranges = [59, 23, 31, 12, 6]
  for (let i = 0; i < parts.length; i++) {
    const part = parts[i]
    const maxVal = ranges[i]
    if (part === '*' || part === '?') continue
    if (part.includes('-')) {
      const [start, end] = part.split('-').map(Number)
      if (isNaN(start) || isNaN(end) || start < 0 || end > maxVal || start > end) return false
      continue
    }
    if (part.includes('/')) {
      const [base, step] = part.split('/')
      if (isNaN(Number(step)) || Number(step) <= 0) return false
      if (base === '*') continue
      if (base.includes('-')) {
        const [start, end] = base.split('-').map(Number)
        if (isNaN(start) || isNaN(end) || start < 0 || end > maxVal || start > end) return false
      } else {
        const num = Number(base)
        if (isNaN(num) || num < 0 || num > maxVal) return false
      }
      continue
    }
    if (part.includes(',')) {
      const values = part.split(',').map(Number)
      for (const val of values) {
        if (isNaN(val) || val < 0 || val > maxVal) return false
      }
      continue
    }
    const num = Number(part)
    if (isNaN(num) || num < 0 || num > maxVal) return false
  }
  return true
}

watch(
  () => subSchedulerForm.value.cron_expr,
  (newCron) => {
    if (!newCron || newCron.trim() === '') {
      cronValidationStatus.value = { isValid: true, message: '' }
      return
    }
    const isValid = validateCronExpression(newCron.trim())
    if (isValid) {
      cronValidationStatus.value = { isValid: true, message: '格式正确' }
    } else {
      const parts = newCron.trim().split(/\s+/)
      let errorMsg = 'Cron 表达式格式不正确'
      if (parts.length !== 5) {
        errorMsg = `需要5个字段，当前有 ${parts.length} 个`
      }
      cronValidationStatus.value = { isValid: false, message: errorMsg }
    }
  }
)

const submitSubSchedulerForm = async () => {
  if (!subSchedulerForm.value.name.trim()) {
    ElMessage.warning('请输入名称')
    return
  }
  if (!subSchedulerForm.value.url.trim()) {
    ElMessage.warning('请输入URL')
    return
  }
  if (!subSchedulerForm.value.cron_expr.trim()) {
    ElMessage.warning('请输入Cron表达式')
    return
  }
  if (!validateCronExpression(subSchedulerForm.value.cron_expr.trim())) {
    ElMessage.error('请输入正确的5字段Cron表达式（分 时 日 月 周）')
    return
  }
  try {
    if (subSchedulerFormTitle.value === '添加订阅') {
      const response = await addSubScheduler(subSchedulerForm.value)
      if (response) {
        ElMessage.success('添加成功')
        subSchedulerFormVisible.value = false
        await getSubSchedulerList()
      }
    } else {
      const response = await updateSubScheduler(subSchedulerForm.value)
      if (response) {
        ElMessage.success('更新成功')
        subSchedulerFormVisible.value = false
        await getSubSchedulerList()
      }
    }
  } catch {
    ElMessage.error('操作失败')
  }
}

const handleSubSizeChange = (val: number) => { subPageSize.value = val }
const handleSubCurrentChange = (val: number) => { subCurrentPage.value = val }

const formatDateTime = (dateTimeString: string) => {
  if (!dateTimeString) return '-'
  try {
    const date = new Date(dateTimeString)
    if (isNaN(date.getTime())) return '-'
    return date.toLocaleString('zh-CN', {
      year: 'numeric', month: '2-digit', day: '2-digit',
      hour: '2-digit', minute: '2-digit', second: '2-digit'
    })
  } catch { return '-' }
}
</script>

<template>
  <div class="page-container">
    <!-- Add/Edit Node Dialog -->
    <el-dialog v-model="dialogVisible" :title="NodeTitle" width="640px">
      <el-form label-position="top">
        <el-form-item label="节点链接">
          <el-input v-model="Nodelink"
            placeholder="输入节点链接，多行使用回车或逗号分隔，支持 base64 格式的 URL 订阅"
            type="textarea" :autosize="{ minRows: 3, maxRows: 10 }" />
        </el-form-item>
        <el-form-item v-if="NodeTitle == '添加节点'" label="导入模式">
          <el-radio-group v-model="radio1">
            <el-radio-button value="1">合并为一个节点</el-radio-button>
            <el-radio-button value="2">分开为多个节点</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="节点备注" v-if="radio1 != '2'">
          <el-input v-model="Nodename" placeholder="请输入节点备注名称" clearable />
        </el-form-item>
        <el-form-item label="前置代理">
          <el-input v-model="DialerProxyName"
            placeholder="前置代理节点名称或策略组名称（仅 Clash-Meta 内核可用）" clearable />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="addnodes">确定</el-button>
      </template>
    </el-dialog>

    <!-- Main Card -->
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <div class="header-left">
            <span class="card-title">节点管理</span>
            <el-input v-model="searchQuery" placeholder="搜索节点..." class="header-search-input" clearable
              @input="handleSearch" size="default">
              <template #prefix>
                <el-icon><Search /></el-icon>
              </template>
            </el-input>
          </div>
          <div class="header-right">
            <el-button type="primary" @click="handleAddNode">
              <el-icon class="mr-1"><Plus /></el-icon>添加节点
            </el-button>
            <el-button type="success" @click="handleImportSubscription">
              <el-icon class="mr-1"><Download /></el-icon>导入订阅
            </el-button>
          </div>
        </div>
      </template>

      <el-table ref="table" v-loading="loading" :data="currentTableData" style="width: 100%" stripe
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" fixed width="50" />
        <el-table-column prop="Name" label="备注" min-width="160">
          <template #default="scope">
            <el-tag type="success" effect="light" size="small">{{ scope.row.Name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="Link" label="节点链接" sortable :show-overflow-tooltip="true" min-width="250" />
        <el-table-column prop="CreateDate" label="创建时间" sortable width="180" />
        <el-table-column fixed="right" label="操作" width="160">
          <template #default="scope">
            <el-button link type="primary" size="small" @click="copyInfo(scope.row)">
              <el-icon><CopyDocument /></el-icon>复制
            </el-button>
            <el-button link type="primary" size="small" @click="handleEdit(scope.row)">
              <el-icon><Edit /></el-icon>编辑
            </el-button>
            <el-button link type="danger" size="small" @click="handleDel(scope.row)">
              <el-icon><Delete /></el-icon>删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="!loading && (!filteredTableData || filteredTableData.length === 0)"
        description="暂无节点数据">
        <el-button type="primary" @click="getnodes">
          <el-icon><Refresh /></el-icon> 重新加载
        </el-button>
      </el-empty>

      <div class="table-footer">
        <div class="batch-actions">
          <el-button size="small" @click="selectAll()">全选</el-button>
          <el-button size="small" @click="toggleSelection()">取消</el-button>
          <el-button size="small" type="danger" @click="selectDel"
            :disabled="multipleSelection.length === 0">
            批量删除 ({{ multipleSelection.length }})
          </el-button>
        </div>
        <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange"
          :current-page="currentPage" :page-size="pageSize" layout="total, sizes, prev, pager, next"
          :page-sizes="[10, 20, 30, 50]" :total="filteredTableData.length" small />
      </div>
    </el-card>

    <!-- Subscription Scheduler Dialog -->
    <el-dialog v-model="subSchedulerDialogVisible" title="订阅管理" width="900px" :close-on-click-modal="false">
      <div class="scheduler-header">
        <el-button type="primary" size="small" @click="handleAddSubScheduler">
          <el-icon class="mr-1"><Plus /></el-icon>添加订阅
        </el-button>
        <el-button type="danger" size="small" @click="handleBatchDeleteSubScheduler"
          :disabled="subSchedulerSelection.length === 0">
          批量删除 ({{ subSchedulerSelection.length }})
        </el-button>
      </div>

      <el-table ref="subSchedulerTable" :data="currentSubSchedulerData" style="width: 100%" stripe
        @selection-change="handleSubSchedulerSelectionChange">
        <el-table-column type="selection" width="45" />
        <el-table-column prop="Name" label="名称" min-width="120">
          <template #default="scope">
            <el-tag type="primary" effect="light" size="small">{{ scope.row.Name }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="URL" label="订阅地址" min-width="200" :show-overflow-tooltip="true" />
        <el-table-column prop="CronExpr" label="Cron" width="110" />
        <el-table-column prop="SuccessCount" label="节点数" width="80" />
        <el-table-column prop="LastRunTime" label="上次运行" width="160">
          <template #default="scope">
            <span v-if="scope.row.LastRunTime">{{ formatDateTime(scope.row.LastRunTime) }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="Enabled" label="状态" width="80">
          <template #default="scope">
            <el-tag :type="scope.row.Enabled ? 'success' : 'danger'" size="small" effect="light">
              {{ scope.row.Enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button link type="primary" size="small" @click="handleEditSubScheduler(scope.row)">编辑</el-button>
            <el-button link type="danger" size="small" @click="handleDeleteSubScheduler(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination class="mt-4" @size-change="handleSubSizeChange" @current-change="handleSubCurrentChange"
        :current-page="subCurrentPage" :page-size="subPageSize" layout="total, sizes, prev, pager, next"
        :page-sizes="[10, 20, 30]" :total="subSchedulerData.length" small />
    </el-dialog>

    <!-- Add/Edit Scheduler Form -->
    <el-dialog v-model="subSchedulerFormVisible" :title="subSchedulerFormTitle" width="520px"
      :close-on-click-modal="false">
      <el-form :model="subSchedulerForm" label-position="top">
        <el-form-item label="名称" required>
          <el-input v-model="subSchedulerForm.name" placeholder="订阅名称" clearable />
        </el-form-item>
        <el-form-item label="订阅地址" required>
          <el-input v-model="subSchedulerForm.url" placeholder="订阅 URL 地址" clearable />
        </el-form-item>
        <el-form-item label="Cron 表达式" required>
          <el-input v-model="subSchedulerForm.cron_expr" placeholder="例如: 0 */6 * * *" clearable />
          <div class="cron-help">
            <div v-if="subSchedulerForm.cron_expr.trim() && cronValidationStatus.isValid"
              class="cron-valid">{{ cronValidationStatus.message }}</div>
            <div v-if="subSchedulerForm.cron_expr.trim() && !cronValidationStatus.isValid"
              class="cron-invalid">{{ cronValidationStatus.message }}</div>
            <div class="cron-examples">
              <span>常用: </span>
              <code>0 */6 * * *</code> 每6小时 |
              <code>0 0 * * *</code> 每天0点 |
              <code>0 */2 * * *</code> 每2小时
            </div>
          </div>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="subSchedulerForm.enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="subSchedulerFormVisible = false">取消</el-button>
        <el-button type="primary" @click="submitSubSchedulerForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  min-width: 0;
}

.header-search-input {
  width: 220px;
  max-width: 100%;
}

.header-right {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

@media (max-width: 768px) {
  .header-left {
    flex-direction: column;
    align-items: stretch;
    width: 100%;
  }

  .header-search-input {
    width: 100%;
  }

  .header-right .el-button {
    flex: 1 1 calc(50% - 4px);
    min-width: 0;
  }
}

.scheduler-header {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.text-muted {
  color: var(--sl-text-muted);
}

.cron-help {
  font-size: 12px;
  margin-top: 6px;
  color: var(--sl-text-secondary);
}

.cron-valid {
  color: var(--sl-success);
  font-weight: 500;
  margin-bottom: 4px;
}

.cron-invalid {
  color: var(--sl-danger);
  font-weight: 500;
  margin-bottom: 4px;
}

.cron-examples {
  margin-top: 4px;
  line-height: 1.8;

  code {
    background: var(--sl-border-light);
    padding: 1px 6px;
    border-radius: 4px;
    font-size: 12px;
  }
}

.mr-1 {
  margin-right: 4px;
}

.mt-4 {
  margin-top: 16px;
}
</style>
