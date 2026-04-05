<script setup lang='ts'>
import { ref, onMounted, nextTick, computed } from 'vue'
import { Plus, Delete, Edit } from '@element-plus/icons-vue'
import { getTemp, AddTemp, UpdateTemp, DelTemp } from "@/api/subcription/temp"
import { ElMessage, ElMessageBox } from 'element-plus'

interface Temp {
  file: string;
  text: string;
  CreateDate: string;
}

const tableData = ref<Temp[]>([])
const Tempoldname = ref('')
const Tempname = ref('')
const TempText = ref('')
const dialogVisible = ref(false)
const table = ref()
const TempTitle = ref('')

async function gettemps() {
  const { data } = await getTemp();
  tableData.value = data
}

onMounted(async () => { gettemps() })

const handleAddTemp = () => {
  TempTitle.value = '添加模板'
  Tempname.value = ''
  TempText.value = ''
  dialogVisible.value = true
}

const addtemp = async () => {
  if (TempTitle.value == '添加模板') {
    await AddTemp({
      filename: Tempname.value.trim(),
      text: TempText.value.trim(),
    })
    ElMessage.success("添加成功");
  } else {
    await UpdateTemp({
      filename: Tempname.value.trim(),
      oldname: Tempoldname.value.trim(),
      text: TempText.value.trim(),
    })
    ElMessage.success("更新成功");
  }
  gettemps()
  Tempname.value = ''
  TempText.value = ''
  dialogVisible.value = false;
}

const multipleSelection = ref<Temp[]>([])
const handleSelectionChange = (val: Temp[]) => { multipleSelection.value = val }

const selectAll = () => {
  nextTick(() => {
    tableData.value.forEach(row => {
      table.value.toggleRowSelection(row, true)
    })
  })
}

const handleEdit = (row: any) => {
  for (let i = 0; i < tableData.value.length; i++) {
    if (tableData.value[i].file === row.file) {
      TempTitle.value = '编辑模板'
      Tempname.value = tableData.value[i].file
      Tempoldname.value = Tempname.value
      TempText.value = tableData.value[i].text
      dialogVisible.value = true
    }
  }
}

const toggleSelection = () => { table.value.clearSelection() }

const handleDel = (row: any) => {
  ElMessageBox.confirm(`确定要删除模板「${row.file}」吗？`, '确认删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    await DelTemp({ filename: row.file, type: row.type })
    ElMessage.success('删除成功')
    gettemps()
  })
}

const selectDel = () => {
  if (multipleSelection.value.length === 0) return
  ElMessageBox.confirm(`确定要删除选中的 ${multipleSelection.value.length} 个模板吗？`, '确认批量删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    for (let i = 0; i < multipleSelection.value.length; i++) {
      DelTemp({ filename: multipleSelection.value[i].file })
      tableData.value = tableData.value.filter((item) => item.file !== multipleSelection.value[i].file)
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
  return tableData.value.slice(start, end);
});
</script>

<template>
  <div class="page-container">
    <el-dialog v-model="dialogVisible" :title="TempTitle" width="80%">
      <el-form label-position="top">
        <el-form-item label="模板内容">
          <el-input v-model="TempText" placeholder="输入模板配置内容 (YAML/CONF)" :rows="16"
            type="textarea" class="template-editor" />
        </el-form-item>
        <el-form-item label="文件名">
          <el-input v-model="Tempname" placeholder="例如: clash.yaml 或 surge.conf" clearable />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="addtemp">确定</el-button>
      </template>
    </el-dialog>

    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">模板管理</span>
          <el-button type="primary" @click="handleAddTemp">
            <el-icon class="mr-1"><Plus /></el-icon>添加模板
          </el-button>
        </div>
      </template>

      <el-table ref="table" :data="currentTableData" style="width: 100%" stripe
        @selection-change="handleSelectionChange">
        <el-table-column type="selection" fixed width="50" />
        <el-table-column prop="file" label="模板文件名" min-width="200">
          <template #default="scope">
            <el-tag type="success" effect="light" size="small">{{ scope.row.file }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="create_date" label="创建时间" sortable width="180" />
        <el-table-column fixed="right" label="操作" width="140">
          <template #default="scope">
            <el-button link type="primary" size="small" @click="handleEdit(scope.row)">
              <el-icon><Edit /></el-icon>编辑
            </el-button>
            <el-button link type="danger" size="small" @click="handleDel(scope.row)">
              <el-icon><Delete /></el-icon>删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

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
          :page-sizes="[10, 20, 30, 50]" :total="tableData.length" small />
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.page-container {
  padding: 16px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--sl-text);
}

.table-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.batch-actions {
  display: flex;
  gap: 8px;
}

.template-editor :deep(textarea) {
  font-family: 'SF Mono', 'Fira Code', 'Cascadia Code', 'Consolas', monospace;
  font-size: 13px;
  line-height: 1.6;
  tab-size: 2;
}

.mr-1 {
  margin-right: 4px;
}
</style>
