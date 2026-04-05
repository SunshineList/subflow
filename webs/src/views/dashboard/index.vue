<template>
  <div class="dashboard-container">
    <!-- 欢迎区 -->
    <div class="welcome-section">
      <div class="welcome-text">
        <h2 class="greeting">{{ greetings }}</h2>
        <p class="subtitle">SubFlow 订阅转换管理面板 &mdash; 当前时间 {{ currentTime }}</p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="12" :sm="6">
        <router-link to="/subcription/subs" class="stat-card stat-primary">
          <div class="stat-icon">
            <el-icon :size="24"><Connection /></el-icon>
          </div>
          <div class="stat-detail">
            <span class="stat-value">{{ subTotal }}</span>
            <span class="stat-label">订阅总数</span>
          </div>
        </router-link>
      </el-col>
      <el-col :xs="12" :sm="6">
        <router-link to="/subcription/nodes" class="stat-card stat-success">
          <div class="stat-icon">
            <el-icon :size="24"><Cpu /></el-icon>
          </div>
          <div class="stat-detail">
            <span class="stat-value">{{ nodeTotal }}</span>
            <span class="stat-label">节点总数</span>
          </div>
        </router-link>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card stat-warning" @click="handleGoScheduler" style="cursor:pointer">
          <div class="stat-icon">
            <el-icon :size="24"><Timer /></el-icon>
          </div>
          <div class="stat-detail">
            <span class="stat-value">{{ schedulerList.length }}</span>
            <span class="stat-label">定时任务</span>
          </div>
        </div>
      </el-col>
      <el-col :xs="12" :sm="6">
        <div class="stat-card stat-info-card">
          <div class="stat-icon">
            <el-icon :size="24"><InfoFilled /></el-icon>
          </div>
          <div class="stat-detail">
            <span class="stat-value">{{ version }}</span>
            <span class="stat-label">当前版本</span>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 第二行：订阅概览 + 协议分布 -->
    <el-row :gutter="16" class="content-row">
      <el-col :xs="24" :lg="16">
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header-flex">
              <span class="card-title">订阅概览</span>
              <router-link to="/subcription/subs" class="header-link">
                查看全部 <el-icon :size="12"><ArrowRight /></el-icon>
              </router-link>
            </div>
          </template>
          <el-table :data="subsList.slice(0, 5)" style="width: 100%" size="small" stripe
            v-loading="subsLoading" :show-header="true">
            <el-table-column prop="Name" label="订阅名称" min-width="140">
              <template #default="{ row }">
                <el-tag type="primary" effect="light" size="small">{{ row.Name }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="节点数" width="80" align="center">
              <template #default="{ row }">
                <span class="num-badge">{{ row.Nodes ? row.Nodes.length : 0 }}</span>
              </template>
            </el-table-column>
            <el-table-column label="访问次数" width="90" align="center">
              <template #default="{ row }">
                <span class="num-badge accent">{{ getAccessCount(row) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="CreateDate" label="创建时间" width="160" />
          </el-table>
          <el-empty v-if="!subsLoading && subsList.length === 0" description="暂无订阅" :image-size="60" />
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="8">
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header-flex">
              <span class="card-title">协议分布</span>
            </div>
          </template>
          <div class="protocol-list" v-if="protocolStats.length > 0">
            <div v-for="item in protocolStats" :key="item.name" class="protocol-item">
              <div class="protocol-header">
                <span class="protocol-name">{{ item.name }}</span>
                <span class="protocol-count">{{ item.count }}</span>
              </div>
              <el-progress :percentage="item.percent" :stroke-width="6"
                :color="item.color" :show-text="false" />
            </div>
          </div>
          <el-empty v-else description="暂无节点数据" :image-size="60" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 第三行：定时任务 + 最近访问 -->
    <el-row :gutter="16" class="content-row">
      <el-col :xs="24" :lg="12">
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header-flex">
              <span class="card-title">定时订阅任务</span>
              <el-tag :type="schedulerList.length > 0 ? 'success' : 'info'" size="small" effect="light">
                {{ schedulerList.filter(s => s.Enabled).length }} 个运行中
              </el-tag>
            </div>
          </template>
          <div v-if="schedulerList.length > 0" class="scheduler-list">
            <div v-for="item in schedulerList.slice(0, 5)" :key="item.ID" class="scheduler-item">
              <div class="scheduler-left">
                <el-tag :type="item.Enabled ? 'success' : 'danger'" size="small" effect="light" class="scheduler-status">
                  {{ item.Enabled ? '运行' : '停止' }}
                </el-tag>
                <div class="scheduler-info">
                  <span class="scheduler-name">{{ item.Name }}</span>
                  <span class="scheduler-cron">{{ item.CronExpr }}</span>
                </div>
              </div>
              <div class="scheduler-right">
                <span v-if="item.SuccessCount" class="scheduler-nodes">{{ item.SuccessCount }} 节点</span>
                <span v-if="item.LastRunTime" class="scheduler-time">{{ formatTime(item.LastRunTime) }}</span>
                <span v-else class="scheduler-time muted">未运行</span>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无定时任务" :image-size="60">
            <router-link to="/subcription/nodes">
              <el-button type="primary" size="small">前往添加</el-button>
            </router-link>
          </el-empty>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="12">
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header-flex">
              <span class="card-title">最近访问记录</span>
              <el-tag size="small" effect="light">
                共 {{ allLogs.length }} 条
              </el-tag>
            </div>
          </template>
          <div v-if="allLogs.length > 0" class="log-list">
            <div v-for="log in allLogs.slice(0, 8)" :key="log.ID" class="log-item">
              <div class="log-left">
                <span class="log-ip">{{ log.IP }}</span>
                <span class="log-addr" v-if="log.Addr">{{ log.Addr }}</span>
              </div>
              <div class="log-right">
                <span class="log-sub">{{ log.SubName }}</span>
                <span class="log-date">{{ log.Date }}</span>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无访问记录" :image-size="60" />
        </el-card>
      </el-col>
    </el-row>

    <!-- 第四行：快捷操作 + 项目信息 -->
    <el-row :gutter="16" class="content-row">
      <el-col :xs="24" :md="12">
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header-flex">
              <span class="card-title">快捷操作</span>
            </div>
          </template>
          <div class="quick-actions">
            <router-link to="/subcription/subs" class="action-item">
              <el-icon :size="20" color="#2563EB"><FolderOpened /></el-icon>
              <span>管理订阅</span>
            </router-link>
            <router-link to="/subcription/nodes" class="action-item">
              <el-icon :size="20" color="#10B981"><Cpu /></el-icon>
              <span>管理节点</span>
            </router-link>
            <router-link to="/subcription/template" class="action-item">
              <el-icon :size="20" color="#F59E0B"><Document /></el-icon>
              <span>管理模板</span>
            </router-link>
            <router-link to="/apikey" class="action-item">
              <el-icon :size="20" color="#6366F1"><Key /></el-icon>
              <span>API 密钥</span>
            </router-link>
            <router-link to="/plugin" class="action-item">
              <el-icon :size="20" color="#EC4899"><SetUp /></el-icon>
              <span>插件管理</span>
            </router-link>
            <div class="action-item" @click="refreshAll" style="cursor:pointer">
              <el-icon :size="20" color="#64748B"><Refresh /></el-icon>
              <span>刷新数据</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :md="12">
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header-flex">
              <span class="card-title">项目信息</span>
            </div>
          </template>
          <div class="project-info">
            <div class="info-row">
              <span class="info-label">项目名称</span>
              <span class="info-value">SubFlow</span>
            </div>
            <div class="info-row">
              <span class="info-label">项目地址</span>
              <a href="https://github.com/SunshineList/sublinkE" target="_blank" class="info-link">
                SunshineList/sublinkE <el-icon :size="12"><TopRight /></el-icon>
              </a>
            </div>
            <div class="info-row">
              <span class="info-label">技术栈</span>
              <span class="info-value">Go + Gin + Vue3 + Element Plus</span>
            </div>
            <div class="info-row">
              <span class="info-label">支持协议</span>
              <span class="info-value">SS/SSR/Trojan/VMess/VLESS/HY/HY2/TUIC/AnyTLS</span>
            </div>
            <div class="info-row">
              <span class="info-label">客户端格式</span>
              <span class="info-value">V2Ray / Clash / Surge</span>
            </div>
            <div class="info-row">
              <span class="info-label">原始项目</span>
              <a href="https://github.com/eun1e/sublinkE" target="_blank" class="info-link">
                eun1e/sublinkE <el-icon :size="12"><TopRight /></el-icon>
              </a>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
defineOptions({
  name: "Dashboard",
  inheritAttrs: false,
});

import {
  Connection, Cpu, Timer, Key, FolderOpened, Document, TopRight,
  InfoFilled, ArrowRight, SetUp, Refresh
} from "@element-plus/icons-vue";
import { useUserStore } from "@/store/modules/user";
import { getSubTotal, getNodeTotal } from "@/api/total";
import { GetVersion } from "@/api/auth";
import { getSubs } from "@/api/subcription/subs";
import { getNodes } from "@/api/subcription/node";
import { getSubSchedulers } from "@/api/subcription/scheduler";
import router from "@/router";

const userStore = useUserStore();
const subTotal = ref(0);
const nodeTotal = ref(0);
const version = ref('--');
const currentTime = ref('');
const subsLoading = ref(false);

interface SubLog {
  ID: number;
  IP: string;
  Date: string;
  Addr: string;
  Count: number;
  SubcriptionID: number;
  SubName: string;
}

const subsList = ref<any[]>([]);
const nodesList = ref<any[]>([]);
const schedulerList = ref<any[]>([]);
const allLogs = ref<SubLog[]>([]);

const protocolColors: Record<string, string> = {
  'VMess': '#2563EB',
  'VLESS': '#6366F1',
  'Trojan': '#EC4899',
  'SS': '#10B981',
  'SSR': '#14B8A6',
  'Hysteria': '#F59E0B',
  'Hysteria2': '#F97316',
  'TUIC': '#8B5CF6',
  'AnyTLS': '#EF4444',
  'Socks5': '#64748B',
  'Other': '#94A3B8',
};

const protocolStats = computed(() => {
  if (nodesList.value.length === 0) return [];
  const counts: Record<string, number> = {};
  for (const node of nodesList.value) {
    const link = (node.Link || '').toLowerCase();
    let proto = 'Other';
    if (link.startsWith('vmess://')) proto = 'VMess';
    else if (link.startsWith('vless://')) proto = 'VLESS';
    else if (link.startsWith('trojan://')) proto = 'Trojan';
    else if (link.startsWith('ss://')) proto = 'SS';
    else if (link.startsWith('ssr://')) proto = 'SSR';
    else if (link.startsWith('hysteria2://') || link.startsWith('hy2://')) proto = 'Hysteria2';
    else if (link.startsWith('hysteria://') || link.startsWith('hy://')) proto = 'Hysteria';
    else if (link.startsWith('tuic://')) proto = 'TUIC';
    else if (link.startsWith('anytls://')) proto = 'AnyTLS';
    else if (link.startsWith('socks5://')) proto = 'Socks5';
    counts[proto] = (counts[proto] || 0) + 1;
  }
  const total = nodesList.value.length;
  return Object.entries(counts)
    .sort((a, b) => b[1] - a[1])
    .map(([name, count]) => ({
      name,
      count,
      percent: Math.round((count / total) * 100),
      color: protocolColors[name] || '#94A3B8',
    }));
});

const getAccessCount = (row: any) => {
  if (!row.SubLogs || row.SubLogs.length === 0) return 0;
  return row.SubLogs.reduce((sum: number, log: any) => sum + (log.Count || 0), 0);
};

const formatTime = (dateStr: string) => {
  if (!dateStr) return '--';
  try {
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return '--';
    const now = new Date();
    const diff = now.getTime() - d.getTime();
    const mins = Math.floor(diff / 60000);
    if (mins < 1) return '刚刚';
    if (mins < 60) return `${mins} 分钟前`;
    const hours = Math.floor(mins / 60);
    if (hours < 24) return `${hours} 小时前`;
    const days = Math.floor(hours / 24);
    if (days < 7) return `${days} 天前`;
    return d.toLocaleDateString('zh-CN');
  } catch {
    return '--';
  }
};

const updateTime = () => {
  const now = new Date();
  currentTime.value = now.toLocaleString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  });
};

let timeInterval: ReturnType<typeof setInterval>;

const fetchAll = async () => {
  subsLoading.value = true;
  try {
    const [subTotalRes, nodeTotalRes, subsRes, nodesRes, schedulerRes] = await Promise.allSettled([
      getSubTotal(),
      getNodeTotal(),
      getSubs(),
      getNodes(),
      getSubSchedulers(),
    ]);

    if (subTotalRes.status === 'fulfilled') subTotal.value = subTotalRes.value.data;
    if (nodeTotalRes.status === 'fulfilled') nodeTotal.value = nodeTotalRes.value.data;
    if (subsRes.status === 'fulfilled') {
      subsList.value = subsRes.value.data || [];
      const logs: SubLog[] = [];
      for (const sub of subsList.value) {
        if (sub.SubLogs && sub.SubLogs.length > 0) {
          for (const log of sub.SubLogs) {
            logs.push({ ...log, SubName: sub.Name });
          }
        }
      }
      logs.sort((a, b) => new Date(b.Date).getTime() - new Date(a.Date).getTime());
      allLogs.value = logs;
    }
    if (nodesRes.status === 'fulfilled') nodesList.value = nodesRes.value.data || [];
    if (schedulerRes.status === 'fulfilled') schedulerList.value = schedulerRes.value.data || [];
  } catch (e) {
    /* ignore */
  } finally {
    subsLoading.value = false;
  }

  GetVersion().then((res) => { version.value = res.data; }).catch(() => {});
};

const refreshAll = () => {
  ElMessage.info('正在刷新数据...');
  fetchAll().then(() => ElMessage.success('刷新完成'));
};

const handleGoScheduler = () => {
  router.push('/subcription/nodes');
};

onMounted(() => {
  fetchAll();
  updateTime();
  timeInterval = setInterval(updateTime, 1000);
});

onUnmounted(() => {
  clearInterval(timeInterval);
});

const greetings = computed(() => {
  const hours = new Date().getHours();
  const name = userStore.user.nickname || 'Admin';
  if (hours >= 6 && hours < 12) return `早上好，${name}`;
  if (hours >= 12 && hours < 18) return `下午好，${name}`;
  if (hours >= 18 && hours < 22) return `晚上好，${name}`;
  return `夜深了，${name}`;
});
</script>

<style lang="scss" scoped>
.dashboard-container {
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
}

.welcome-section {
  margin-bottom: 24px;
}

.greeting {
  font-size: 28px;
  font-weight: 700;
  color: var(--sl-text);
  margin: 0 0 4px 0;
  letter-spacing: -0.5px;
}

.subtitle {
  font-size: 14px;
  color: var(--sl-text-secondary);
  margin: 0;
}

/* 统计卡片 */
.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  border-radius: var(--sl-radius-lg);
  background: var(--sl-bg-card);
  border: 1px solid var(--sl-border);
  box-shadow: var(--sl-shadow-sm);
  transition: all 200ms ease-out;
  cursor: default;
  text-decoration: none;

  &:hover {
    box-shadow: var(--sl-shadow-md);
    transform: translateY(-2px);
  }
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-primary .stat-icon { background: rgba(37, 99, 235, 0.1); color: #2563EB; }
.stat-success .stat-icon { background: rgba(16, 185, 129, 0.1); color: #10B981; }
.stat-warning .stat-icon { background: rgba(245, 158, 11, 0.1); color: #F59E0B; }
.stat-info-card .stat-icon { background: rgba(99, 102, 241, 0.1); color: #6366F1; }

.stat-detail {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: var(--sl-text);
  line-height: 1.2;
}

.stat-label {
  font-size: 13px;
  color: var(--sl-text-secondary);
  margin-top: 2px;
}

/* 通用内容行 */
.content-row {
  margin-bottom: 20px;

  .el-col {
    margin-bottom: 16px;
  }
}

.info-card {
  height: 100%;
}

.card-header-flex {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--sl-text);
}

.header-link {
  font-size: 13px;
  color: var(--sl-primary);
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 2px;
  transition: color 200ms;

  &:hover {
    color: var(--sl-primary-light);
  }
}

/* 数字徽章 */
.num-badge {
  display: inline-block;
  min-width: 28px;
  padding: 2px 8px;
  background: rgba(37, 99, 235, 0.08);
  color: var(--sl-primary);
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  text-align: center;

  &.accent {
    background: rgba(249, 115, 22, 0.08);
    color: #F97316;
  }
}

/* 协议分布 */
.protocol-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.protocol-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.protocol-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.protocol-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--sl-text);
}

.protocol-count {
  font-size: 12px;
  color: var(--sl-text-muted);
  font-weight: 500;
}

/* 定时任务 */
.scheduler-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.scheduler-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: var(--sl-radius);
  background: var(--sl-border-light);
  transition: background 200ms;

  &:hover {
    background: rgba(37, 99, 235, 0.04);
  }
}

.scheduler-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.scheduler-status {
  flex-shrink: 0;
}

.scheduler-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.scheduler-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--sl-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.scheduler-cron {
  font-size: 11px;
  color: var(--sl-text-muted);
  font-family: 'SF Mono', 'Cascadia Code', monospace;
}

.scheduler-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  flex-shrink: 0;
}

.scheduler-nodes {
  font-size: 12px;
  color: var(--sl-primary);
  font-weight: 500;
}

.scheduler-time {
  font-size: 11px;
  color: var(--sl-text-muted);

  &.muted {
    color: var(--sl-text-muted);
    font-style: italic;
  }
}

/* 访问日志 */
.log-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.log-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  border-radius: var(--sl-radius);
  transition: background 200ms;

  &:hover {
    background: var(--sl-border-light);
  }

  &:not(:last-child) {
    border-bottom: 1px solid var(--sl-border-light);
  }
}

.log-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.log-ip {
  font-size: 13px;
  font-weight: 600;
  color: var(--sl-text);
  font-family: 'SF Mono', 'Cascadia Code', monospace;
}

.log-addr {
  font-size: 12px;
  color: var(--sl-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 120px;
}

.log-right {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  flex-shrink: 0;
}

.log-sub {
  font-size: 12px;
  color: var(--sl-primary);
  font-weight: 500;
}

.log-date {
  font-size: 11px;
  color: var(--sl-text-muted);
}

/* 快捷操作 */
.quick-actions {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

@media (max-width: 640px) {
  .quick-actions {
    grid-template-columns: repeat(2, 1fr);
  }
}

.action-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-radius: var(--sl-radius);
  background: var(--sl-border-light);
  transition: all 200ms ease-out;
  cursor: pointer;
  text-decoration: none;
  color: var(--sl-text);
  font-size: 14px;
  font-weight: 500;

  &:hover {
    background: rgba(37, 99, 235, 0.06);
    color: var(--sl-primary);
    transform: translateX(4px);
  }
}

/* 项目信息 */
.project-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 0;
  border-bottom: 1px solid var(--sl-border-light);

  &:last-child {
    border-bottom: none;
  }
}

.info-label {
  font-size: 13px;
  color: var(--sl-text-secondary);
  font-weight: 500;
}

.info-value {
  font-size: 13px;
  color: var(--sl-text);
  font-weight: 500;
  text-align: right;
  max-width: 60%;
  word-break: break-all;
}

.info-link {
  font-size: 13px;
  color: var(--sl-primary);
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  text-decoration: none;
  transition: color 200ms;

  &:hover {
    color: var(--sl-primary-light);
  }
}
</style>
