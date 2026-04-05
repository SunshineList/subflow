<script setup lang='ts'>
import { ref, onMounted, computed, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete, View, Edit, CopyDocument, Sort, Check, Close, Link as LinkIcon } from '@element-plus/icons-vue'
import { getSubs, AddSub, DelSub, UpdateSub, SortSub } from "@/api/subcription/subs"
import { getTemp } from "@/api/subcription/temp"
import { getNodes } from "@/api/subcription/node"
import QrcodeVue from 'qrcode.vue'
import md5 from 'md5'

interface Sub {
  ID: number;
  Name: string;
  CreateDate: string;
  Config: Config;
  Nodes: Node[];
  SubLogs: SubLogs[];
}
interface Node {
  ID: number;
  Name: string;
  Link: string;
  CreateDate: string;
  Sort?: number;
}
interface Config {
  clash: string;
  surge: string;
  udp: string;
  cert: string;
}
interface SubLogs {
  ID: number;
  IP: string;
  Date: string;
  Addr: string;
  Count: number;
  SubcriptionID: number;
}
interface Temp {
  file: string;
  text: string;
  CreateDate: string;
}

const tableData = ref<Sub[]>([])
const Clash = ref('')
const Surge = ref('')
const SubTitle = ref('')
const Subname = ref('')
const oldSubname = ref('')
const dialogVisible = ref(false)
const table = ref()
const NodesList = ref<Node[]>([])
const value1 = ref<string[]>([])
const checkList = ref<string[]>([])
const iplogsdialog = ref(false)
const IplogsList = ref<SubLogs[]>([])
const qrcode = ref('')
const templist = ref<Temp[]>([])

async function getsubs() {
  const { data } = await getSubs();
  tableData.value = data;
  processTableData();
}

async function gettemps() {
  const { data } = await getTemp();
  templist.value = data
}

onMounted(() => {
  getsubs()
  gettemps()
})

onMounted(async () => {
  const { data } = await getNodes();
  NodesList.value = data
})

const addSubs = async () => {
  const config = JSON.stringify({
    "clash": Clash.value.trim(),
    "surge": Surge.value.trim(),
    "udp": checkList.value.includes('udp') ? true : false,
    "cert": checkList.value.includes('cert') ? true : false
  })
  if (SubTitle.value === '添加订阅') {
    await AddSub({
      config: config,
      name: Subname.value.trim(),
      nodes: value1.value.join(',')
    })
    getsubs()
    ElMessage.success("添加成功");
  } else {
    await UpdateSub({
      config: config,
      name: Subname.value.trim(),
      nodes: value1.value.join(','),
      oldname: oldSubname.value
    })
    getsubs()
    ElMessage.success("更新成功");
  }
  dialogVisible.value = false;
}

const multipleSelection = ref<Sub[]>([])
const handleSelectionChange = (val: Sub[]) => {
  multipleSelection.value = val
}

const selectAll = () => {
  tableData.value.forEach(row => {
    table.value.toggleRowSelection(row, true)
  })
}

const handleIplogs = (row: any) => {
  iplogsdialog.value = true
  nextTick(() => {
    tableData.value.forEach((item) => {
      if (item.ID === row.ID) {
        IplogsList.value = item.SubLogs
      }
    })
  })
}

const getRowKey = function (row: any): string {
  if (row.Nodes) {
    return row.ID;
  } else {
    return 'node_' + row.ID;
  }
}

const processTableData = () => {
  tableData.value.forEach(subscription => {
    if (subscription.Nodes) {
      subscription.Nodes.forEach((node, index) => {
        (node as any).parentId = subscription.ID;
        if (node.Sort === undefined || node.Sort === null) {
          node.Sort = index;
        }
      });
      if (subscription.Nodes.length > 0 && subscription.Nodes[0].Sort !== undefined) {
        subscription.Nodes.sort((a, b) => {
          return (a.Sort || 0) - (b.Sort || 0);
        });
      }
    }
  });
}

const toggleSelection = () => {
  table.value.clearSelection()
}

const handleAddSub = () => {
  SubTitle.value = '添加订阅'
  Subname.value = ''
  oldSubname.value = ''
  checkList.value = []
  Clash.value = './template/clash.yaml'
  Surge.value = './template/surge.conf'
  dialogVisible.value = true
  value1.value = []
}

const handleEdit = (row: any) => {
  for (let i = 0; i < tableData.value.length; i++) {
    if (tableData.value[i].ID === row.ID) {
      function toConfig(value: string | Config): Config {
        if (typeof value === 'string') {
          return JSON.parse(value) as Config;
        } else {
          return value as Config;
        }
      }
      const config = toConfig(tableData.value[i].Config);
      SubTitle.value = '编辑订阅'
      Subname.value = tableData.value[i].Name
      oldSubname.value = Subname.value
      if (config.udp) {
        checkList.value.push('udp')
      }
      if (config.cert) {
        checkList.value.push('cert')
      }
      Clash.value = config.clash
      Surge.value = config.surge
      dialogVisible.value = true
      value1.value = tableData.value[i].Nodes.map((item) => item.Name)
    }
  }
}

const handleDel = (row: any) => {
  ElMessageBox.confirm(
    `确定要删除订阅「${row.Name}」吗？`,
    '确认删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    await DelSub({ id: row.ID })
    getsubs()
    ElMessage.success('删除成功')
  })
}

const selectDel = () => {
  if (multipleSelection.value.length === 0) return
  ElMessageBox.confirm(
    `确定要删除选中的 ${multipleSelection.value.length} 个订阅吗？`,
    '确认批量删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(() => {
    for (let i = 0; i < multipleSelection.value.length; i++) {
      if (!multipleSelection.value[i].Nodes) continue
      DelSub({ id: multipleSelection.value[i].ID })
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
  let data: Sub[] = JSON.parse(JSON.stringify(tableData.value));
  return data.slice(start, end);
});

const copyUrl = (url: string) => {
  navigator.clipboard?.writeText(url).then(() => {
    ElMessage.success('复制成功');
  }).catch(() => {
    const textarea = document.createElement('textarea');
    textarea.value = url;
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
    ElMessage.success('复制成功');
  });
};

const copyInfo = (row: any) => { copyUrl(row.Link) }

const ClientDiaLog = ref(false)
const ClientList = ['v2ray', 'clash', 'surge']
const ClientUrls = ref<Record<string, string>>({})
const ClientUrl = ref('')

const handleClient = (name: string) => {
  let serverAddress = location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '');
  ClientDiaLog.value = true
  ClientUrl.value = `${serverAddress}/c/?token=${md5(name)}`
  ClientList.forEach((item: string) => {
    ClientUrls.value[item] = `${serverAddress}/c/?token=${md5(name)}`
  })
}

const Qrdialog = ref(false)
const QrTitle = ref('')
const handleQrcode = (url: string, title: string) => {
  Qrdialog.value = true
  qrcode.value = url
  QrTitle.value = title
}
const OpenUrl = (url: string) => { window.open(url) }
const clientradio = ref('1')

// 排序相关
const sortingSubscriptionId = ref<number | null>(null)
const tempNodeSort = ref<{ ID: number, Sort: number }[]>([])
const originalNodesOrder = ref<Node[]>([])
const dragSource = ref<number | null>(null)
const dragTarget = ref<number | null>(null)

const handleDragStart = (e: DragEvent, nodeId: number) => {
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
    e.dataTransfer.setData('text/plain', nodeId.toString())
    dragSource.value = nodeId
  }
}

const handleDragOver = (e: DragEvent, nodeId: number) => {
  if (e.preventDefault) e.preventDefault()
  if (e.dataTransfer) e.dataTransfer.dropEffect = 'move'
  dragTarget.value = nodeId
  return false
}

const handleDrop = (e: DragEvent, targetNodeId: number, subscriptionId: number) => {
  e.stopPropagation()
  if (sortingSubscriptionId.value !== subscriptionId) return
  const sourceNodeId = parseInt(e.dataTransfer?.getData('text/plain') || '0')
  if (sourceNodeId === targetNodeId) return
  const subscription = tableData.value.find(sub => sub.ID === subscriptionId)
  if (!subscription || !subscription.Nodes) return
  const sourceIndex = subscription.Nodes.findIndex(node => node.ID === sourceNodeId)
  const targetIndex = subscription.Nodes.findIndex(node => node.ID === targetNodeId)
  if (sourceIndex > -1 && targetIndex > -1) {
    const [movedNode] = subscription.Nodes.splice(sourceIndex, 1)
    subscription.Nodes.splice(targetIndex, 0, movedNode)
    subscription.Nodes.forEach((node, index) => {
      node.Sort = index + 1
      const sortItem = tempNodeSort.value.find(item => item.ID === node.ID)
      if (sortItem) {
        sortItem.Sort = index + 1
      } else {
        tempNodeSort.value.push({ ID: node.ID, Sort: index + 1 })
      }
    })
  }
  dragSource.value = null
  dragTarget.value = null
  return false
}

const handleDragEnter = (e: DragEvent, nodeId: number) => { dragTarget.value = nodeId }
const handleDragLeave = () => { dragTarget.value = null }

const handleStartSort = (row: any) => {
  sortingSubscriptionId.value = row.ID
  originalNodesOrder.value = JSON.parse(JSON.stringify(row.Nodes))
  tempNodeSort.value = row.Nodes.map((node: any, index: number) => ({
    ID: node.ID,
    Sort: node.Sort !== undefined ? node.Sort : (index + 1)
  }))
  ElMessage.info('已进入排序模式，拖动节点调整顺序')
}

const handleConfirmSort = async (row: any) => {
  row.Nodes.forEach((node: Node, index: number) => {
    const nodeSort = tempNodeSort.value.find(item => item.ID === node.ID)
    if (nodeSort) {
      nodeSort.Sort = index + 1
    } else {
      tempNodeSort.value.push({ ID: node.ID, Sort: index + 1 })
    }
  })
  try {
    await SortSub({ ID: row.ID, NodeSort: tempNodeSort.value })
    ElMessage.success('排序已保存')
    sortingSubscriptionId.value = null
    tempNodeSort.value = []
    originalNodesOrder.value = []
    await getsubs()
  } catch (error) {
    ElMessage.error('排序保存失败')
  }
}

const handleCancelSort = () => {
  if (sortingSubscriptionId.value !== null) {
    for (let i = 0; i < tableData.value.length; i++) {
      if (tableData.value[i].ID === sortingSubscriptionId.value) {
        tableData.value[i].Nodes = JSON.parse(JSON.stringify(originalNodesOrder.value))
        break
      }
    }
  }
  ElMessage.info('已取消排序')
  sortingSubscriptionId.value = null
  tempNodeSort.value = []
  originalNodesOrder.value = []
}
</script>

<template>
  <div class="page-container">
    <!-- QR Code Dialog -->
    <el-dialog v-model="Qrdialog" width="340px" :title="QrTitle" class="qr-dialog">
      <div class="qr-content">
        <qrcode-vue :value="qrcode" :size="200" level="H" />
        <el-input v-model="qrcode" readonly class="qr-url" />
        <div class="qr-actions">
          <el-button type="primary" @click="copyUrl(qrcode)">
            <el-icon class="mr-1"><CopyDocument /></el-icon>复制
          </el-button>
          <el-button @click="OpenUrl(qrcode)">
            <el-icon class="mr-1"><LinkIcon /></el-icon>打开
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- Client Dialog -->
    <el-dialog v-model="ClientDiaLog" title="客户端订阅链接" width="560px">
      <div class="client-list">
        <div class="client-item" @click="handleQrcode(ClientUrl, '自动识别')">
          <div class="client-left">
            <el-tag type="success" effect="plain" round>AUTO</el-tag>
            <span class="client-name">自动识别客户端</span>
          </div>
          <el-button type="primary" link>二维码</el-button>
        </div>
        <div v-for="(item, index) in ClientUrls" :key="index" class="client-item"
          @click="handleQrcode(`${item}&client=${index}`, String(index))">
          <div class="client-left">
            <el-tag effect="plain" round>{{ index }}</el-tag>
            <span class="client-name">{{ index }} 客户端</span>
          </div>
          <el-button type="primary" link>二维码</el-button>
        </div>
      </div>
    </el-dialog>

    <!-- IP Logs Dialog -->
    <el-dialog v-model="iplogsdialog" title="访问记录" width="80%" draggable>
      <el-table :data="IplogsList" border stripe>
        <el-table-column prop="IP" label="IP 地址" />
        <el-table-column prop="Count" label="访问次数" width="100" />
        <el-table-column prop="Addr" label="来源地区" />
        <el-table-column prop="Date" label="最近访问" />
      </el-table>
    </el-dialog>

    <!-- Add/Edit Subscription Dialog -->
    <el-dialog v-model="dialogVisible" :title="SubTitle" width="640px">
      <el-form label-position="top">
        <el-form-item label="订阅名称">
          <el-input v-model="Subname" placeholder="请输入订阅名称" clearable />
        </el-form-item>

        <el-form-item label="Clash 模板">
          <div class="template-selector">
            <el-radio-group v-model="clientradio" size="small">
              <el-radio-button value="1">本地文件</el-radio-button>
              <el-radio-button value="2">URL 链接</el-radio-button>
            </el-radio-group>
            <el-select v-model="Clash" placeholder="选择 Clash 模板" v-if="clientradio === '1'" class="template-input">
              <el-option v-for="template in templist" :key="template.file" :label="template.file"
                :value="'./template/' + template.file" />
            </el-select>
            <el-input v-model="Clash" placeholder="输入模板 URL" v-else class="template-input" />
          </div>
        </el-form-item>

        <el-form-item label="Surge 模板">
          <div class="template-selector">
            <el-select v-model="Surge" placeholder="选择 Surge 模板" v-if="clientradio === '1'" class="template-input">
              <el-option v-for="template in templist" :key="template.file" :label="template.file"
                :value="'./template/' + template.file" />
            </el-select>
            <el-input v-model="Surge" placeholder="输入模板 URL" v-else class="template-input" />
          </div>
        </el-form-item>

        <el-form-item label="强制选项">
          <el-checkbox-group v-model="checkList">
            <el-checkbox :value="'udp'">强制开启 UDP</el-checkbox>
            <el-checkbox :value="'cert'">跳过证书验证</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="选择节点">
          <el-select v-model="value1" multiple placeholder="选择要包含的节点" style="width: 100%"
            filterable collapse-tags collapse-tags-tooltip :max-collapse-tags="3">
            <el-option v-for="item in NodesList" :key="item.Name" :label="item.Name" :value="item.Name" />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="addSubs">确定</el-button>
      </template>
    </el-dialog>

    <!-- Main Content -->
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span class="card-title">订阅管理</span>
          <el-button type="primary" @click="handleAddSub">
            <el-icon class="mr-1"><Plus /></el-icon>添加订阅
          </el-button>
        </div>
      </template>

      <el-table ref="table" :data="currentTableData" style="width: 100%" stripe
        @selection-change="handleSelectionChange" :row-key="getRowKey" :tree-props="{ children: 'Nodes' }">
        <el-table-column type="selection" fixed width="50" />
        <el-table-column prop="Name" label="名称 / 节点" min-width="200">
          <template #default="{ row }">
            <el-tag v-if="row.Nodes" type="primary" effect="plain">
              {{ row.Name }}
              <span v-if="sortingSubscriptionId === row.ID" class="sorting-badge">排序中</span>
            </el-tag>
            <div v-else
              :draggable="sortingSubscriptionId !== null && row.parentId === sortingSubscriptionId"
              @dragstart="(e) => handleDragStart(e, row.ID)"
              @dragover="(e) => handleDragOver(e, row.ID)"
              @drop="(e) => handleDrop(e, row.ID, row.parentId)"
              @dragenter="(e) => handleDragEnter(e, row.ID)"
              @dragleave="handleDragLeave"
              :class="{
                'node-dragging': dragSource === row.ID,
                'node-drag-over': dragTarget === row.ID,
                'node-draggable': sortingSubscriptionId !== null && row.parentId === sortingSubscriptionId
              }">
              <el-tag type="success" effect="light" size="small">{{ row.Name }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="Link" label="链接" min-width="120">
          <template #default="{ row }">
            <el-button v-if="row.Nodes" type="primary" link size="small" @click="handleClient(row.Name)">
              <el-icon class="mr-1"><LinkIcon /></el-icon>获取链接
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="CreateDate" label="创建时间" sortable width="180" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="scope">
            <div v-if="scope.row.Nodes" class="action-btns">
              <el-button link type="primary" size="small" @click="handleIplogs(scope.row)">
                <el-icon><View /></el-icon>记录
              </el-button>
              <el-button link type="primary" size="small" @click="handleEdit(scope.row)">
                <el-icon><Edit /></el-icon>编辑
              </el-button>
              <el-button link type="danger" size="small" @click="handleDel(scope.row)">
                <el-icon><Delete /></el-icon>删除
              </el-button>
              <el-button v-if="sortingSubscriptionId !== scope.row.ID" link type="warning" size="small"
                @click="handleStartSort(scope.row)">
                <el-icon><Sort /></el-icon>排序
              </el-button>
              <template v-if="sortingSubscriptionId === scope.row.ID">
                <el-button link type="success" size="small" @click="handleConfirmSort(scope.row)">
                  <el-icon><Check /></el-icon>保存
                </el-button>
                <el-button link type="info" size="small" @click="handleCancelSort()">
                  <el-icon><Close /></el-icon>取消
                </el-button>
              </template>
            </div>
            <div v-else>
              <el-button link type="primary" size="small" @click="copyInfo(scope.row)">
                <el-icon><CopyDocument /></el-icon>复制
              </el-button>
            </div>
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

.template-selector {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.template-input {
  width: 100%;
}

.action-btns {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
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

.sorting-badge {
  margin-left: 6px;
  font-size: 11px;
  color: var(--sl-accent);
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.6; }
  50% { opacity: 1; }
}

.node-draggable {
  padding: 8px 12px;
  margin: 2px 0;
  border: 1px dashed var(--sl-border);
  border-radius: var(--sl-radius);
  background: var(--sl-border-light);
  cursor: move;
  transition: all 200ms ease-out;
}

.node-draggable:hover {
  border-color: var(--sl-primary);
  background: rgba(37, 99, 235, 0.04);
}

.node-dragging {
  opacity: 0.5;
  border-color: var(--sl-primary) !important;
  background: rgba(37, 99, 235, 0.08) !important;
}

.node-drag-over {
  border-color: var(--sl-primary) !important;
  background: rgba(37, 99, 235, 0.1) !important;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.15);
}

/* Client list */
.client-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.client-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-radius: var(--sl-radius);
  border: 1px solid var(--sl-border);
  cursor: pointer;
  transition: all 200ms ease-out;
}

.client-item:hover {
  border-color: var(--sl-primary);
  background: rgba(37, 99, 235, 0.04);
}

.client-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.client-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--sl-text);
}

/* QR Code dialog */
.qr-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.qr-url {
  width: 100%;
}

.qr-actions {
  display: flex;
  gap: 8px;
}

.mr-1 {
  margin-right: 4px;
}
</style>
