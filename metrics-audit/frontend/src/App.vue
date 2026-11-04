<template>
    <div class="padding-20 audit-page">
        <a-alert v-if="status.message" class="mb-20" :type="status.enabled && status.installed ? 'success' : 'warning'">
            <template #title>{{ statusTitle }}</template>
            <div>{{ status.message }}</div>
        </a-alert>

        <div class="bg-white padding-20">
            <div class="df ai-c jc-b audit-header">
                <a-tabs v-model:active-key="activeTab" hide-content @change="changeTab">
                    <a-tab-pane key="operation" title="操作日志"></a-tab-pane>
                    <a-tab-pane key="login" title="登录日志"></a-tab-pane>
                </a-tabs>
                <a-button @click="getList">
                    <template #icon><icon-refresh /></template>
                </a-button>
            </div>

            <div class="audit-filter mt-10">
                <a-input v-model="filter.username" placeholder="用户名" allow-clear @press-enter="search" />
                <a-select v-model="filter.success" placeholder="状态" allow-clear>
                    <a-option value="true">成功</a-option>
                    <a-option value="false">失败</a-option>
                </a-select>
                <a-select v-if="activeTab === 'operation'" v-model="filter.method" placeholder="请求方法" allow-clear>
                    <a-option value="GET">GET</a-option>
                    <a-option value="POST">POST</a-option>
                    <a-option value="PUT">PUT</a-option>
                    <a-option value="PATCH">PATCH</a-option>
                    <a-option value="DELETE">DELETE</a-option>
                </a-select>
                <a-input v-if="activeTab === 'operation'" v-model="filter.path" placeholder="请求路径" allow-clear @press-enter="search" />
                <a-range-picker
                    v-model:model-value="filter.timeRange"
                    show-time
                    style="width: 360px;"
                    @ok="search"
                    @clear="search"
                />
                <a-button type="primary" @click="search">
                    <template #icon><icon-search /></template>
                    查询
                </a-button>
                <a-button @click="resetFilter">重置</a-button>
            </div>

            <a-table
                class="cptable mt-20"
                :data="list"
                :pagination="false"
                :bordered="false"
                :scroll="{ x: activeTab === 'operation' ? 1240 : 710 }"
                row-key="rowKey"
            >
                <template #columns>
                    <a-table-column title="时间" :width="180" fixed="left">
                        <template #cell="{ record }">{{ formatTime(record.time || record._time) }}</template>
                    </a-table-column>
                    <a-table-column title="用户名" :width="140" data-index="username"></a-table-column>
                    <template v-if="activeTab === 'login'">
                        <a-table-column title="登录方式" :width="120" data-index="login_method"></a-table-column>
                        <a-table-column title="状态" :width="90">
                            <template #cell="{ record }">
                                <span :class="record.success ? 'c-green' : 'c-red'">{{ record.success ? '成功' : '失败' }}</span>
                            </template>
                        </a-table-column>
                        <a-table-column title="失败原因" :width="220">
                            <template #cell="{ record }">{{ record.reason || '-' }}</template>
                        </a-table-column>
                    </template>
                    <template v-else>
                        <a-table-column title="操作记录" :width="200">
                            <template #cell="{ record }">
                                <span class="ellipsis" :title="operationDescription(record)">{{ operationDescription(record) }}</span>
                            </template>
                        </a-table-column>
                        <a-table-column title="方法" :width="90" data-index="method"></a-table-column>
                        <a-table-column title="路径" :width="300">
                            <template #cell="{ record }">
                                <span class="mono ellipsis" :title="record.path">{{ record.path }}</span>
                            </template>
                        </a-table-column>
                        <a-table-column title="路由" :width="260">
                            <template #cell="{ record }">
                                <span class="mono ellipsis" :title="record.route">{{ record.route || '-' }}</span>
                            </template>
                        </a-table-column>
                        <a-table-column title="状态码" :width="100" data-index="status_code"></a-table-column>
                        <a-table-column title="耗时" :width="100">
                            <template #cell="{ record }">{{ record.duration_ms }}ms</template>
                        </a-table-column>
                        <a-table-column title="状态" :width="90">
                            <template #cell="{ record }">
                                <span :class="record.success ? 'c-green' : 'c-red'">{{ record.success ? '成功' : '失败' }}</span>
                            </template>
                        </a-table-column>
                    </template>
                    <a-table-column title="IP" :width="140" data-index="ip"></a-table-column>
                    <a-table-column title="操作" :width="90" fixed="right">
                        <template #cell="{ record }">
                            <a-tooltip content="详情">
                                <i class="opt-icon" @click="openDetail(record)"><icon-eye /></i>
                            </a-tooltip>
                        </template>
                    </a-table-column>
                </template>
            </a-table>

            <div class="df jc-e ai-c mt-10">
                <a-pagination
                    v-model:current="pagination.current"
                    v-model:page-size="pagination.pageSize"
                    :total="pagination.total"
                    :page-size-options="[20, 50, 100, 200]"
                    show-page-size
                    @change="pageChange"
                    @page-size-change="pageSizeChange"
                />
            </div>
        </div>

        <a-drawer :width="720" :visible="detail.show" @cancel="detail.show=false" :footer="false">
            <template #title>日志详情</template>
            <table class="com-table audit-detail-table">
                <tbody>
                    <tr v-for="item in detailRows" :key="item.key">
                        <td style="width:140px;">{{ item.label }}</td>
                        <td>
                            <pre v-if="item.pre">{{ item.value }}</pre>
                            <span v-else>{{ item.value }}</span>
                        </td>
                    </tr>
                </tbody>
            </table>
        </a-drawer>
    </div>
</template>

<script>
import { auditGet } from './api';
import { Message } from '@arco-design/web-vue';
import {
    getLoadedApiRouteDescriptions,
    loadApiRouteDescriptions,
    resolveApiRouteDescription,
} from './api-route-description';

const defaultFilter = () => ({
    username: '',
    success: '',
    method: '',
    path: '',
    timeRange: [],
});

export default {
    data() {
        return {
            activeTab: 'operation',
            status: {
                enabled: false,
                installed: false,
                baseUrl: '',
                message: '',
            },
            filter: defaultFilter(),
            list: [],
            pagination: {
                current: 1,
                pageSize: 20,
                total: 0,
            },
            detail: {
                show: false,
                data: {},
            },
            apiRouteDescriptions: getLoadedApiRouteDescriptions(),
        };
    },
    computed: {
        statusTitle() {
            if (this.status.enabled && this.status.installed) { return '日志服务正常'; }
            if (!this.status.enabled) { return '日志服务未启用'; }
            return '日志服务不可用';
        },
        detailRows() {
            const labels = {
                time: '时间',
                _time: '时间',
                audit_type: '类型',
                username: '用户名',
                operation_description: '操作记录',
                route_description: '操作记录',
                login_method: '登录方式',
                success: '成功',
                reason: '失败原因',
                method: '请求方法',
                path: '请求路径',
                route: '路由',
                params: '参数',
                status_code: '状态码',
                duration_ms: '耗时',
                ip: 'IP',
                user_agent: 'User-Agent',
                message: '消息',
            };
            return Object.keys(this.detail.data || {}).filter((key) => !['tenant', 'user_mode', 'route_description', 'message'].includes(key)).map((key) => {
                let value = this.detail.data[key];
                if (key === 'time' || key === '_time') { value = this.formatTime(value); }
                if (key === 'success') { value = this.normalizeSuccess(value) ? '成功' : '失败'; }
                if (key === 'duration_ms') { value = `${value}ms`; }
                if (value && typeof value === 'object') { value = JSON.stringify(value, null, 2); }
                return {
                    key,
                    label: labels[key] || key,
                    value: value === '' || value === undefined || value === null ? '-' : value,
                    pre: key === 'params' || key === 'user_agent' || (typeof value === 'string' && value.length > 120),
                };
            });
        },
    },
    created() {
        this.loadRouteDescriptions();
        this.getStatus();
        this.getList();
    },
    methods: {
        loadRouteDescriptions() {
            loadApiRouteDescriptions().then((descriptions) => {
                this.apiRouteDescriptions = descriptions;
            });
        },
        normalizeData(res) {
            return res?.data?.data || res?.data || {};
        },
        getStatus() {
            auditGet('logs/status').then((res) => {
                this.status = {
                    ...this.status,
                    ...this.normalizeData(res),
                };
            }).catch((err) => {
                this.status = {
                    enabled: false,
                    installed: false,
                    baseUrl: '',
                    message: err?.response?.data?.message || err?.message || '获取日志服务状态失败',
                };
            });
        },
        buildParams() {
            const params = {
                page: this.pagination.current,
                pageSize: this.pagination.pageSize,
                username: this.filter.username,
                success: this.filter.success,
            };
            if (this.activeTab === 'operation') {
                params.method = this.filter.method;
                params.path = this.filter.path;
            }
            if (this.filter.timeRange && this.filter.timeRange.length === 2) {
                params.startTime = this.toQueryTime(this.filter.timeRange[0]);
                params.endTime = this.toQueryTime(this.filter.timeRange[1]);
            }
            Object.keys(params).forEach((key) => {
                if (params[key] === '' || params[key] === undefined || params[key] === null) {
                    delete params[key];
                }
            });
            return params;
        },
        getList() {
            const path = this.activeTab === 'login' ? '/audit/login-logs' : '/audit/operation-logs';
            const currentPage = this.pagination.current;
            const currentPageSize = this.pagination.pageSize;
            auditGet(path.replace('/audit/', ''), { params: this.buildParams() }).then((res) => {
                const data = this.normalizeData(res);
                const list = data?.list || [];
                this.list = list.map((item, index) => ({
                    rowKey: `${this.activeTab}-${currentPage}-${index}-${item.time || item._time || ''}`,
                    ...item,
                    success: this.normalizeSuccess(item?.success),
                    time: item.time || item._time || '',
                }));
                this.pagination.current = currentPage;
                this.pagination.pageSize = currentPageSize;
                this.pagination.total = Number(data?.total || 0);
            }).catch((err) => {
                this.list = [];
                this.pagination.total = 0;
                Message.error(err?.message || '读取审计日志失败');
            });
        },
        search() {
            this.pagination.current = 1;
            this.getList();
        },
        resetFilter() {
            this.filter = defaultFilter();
            this.search();
        },
        changeTab() {
            this.filter.method = '';
            this.filter.path = '';
            this.pagination.current = 1;
            this.getList();
        },
        pageChange(current) {
            this.pagination.current = current;
            this.getList();
        },
        pageSizeChange(pageSize) {
            this.pagination.pageSize = pageSize;
            this.pagination.current = 1;
            this.getList();
        },
        openDetail(record) {
            this.detail = {
                show: true,
                data: {
                    operation_description: this.operationDescription(record),
                    ...record,
                },
            };
        },
        operationDescription(record) {
            const hasStructuredRoute = record?.method && (record?.route || record?.path);
            if (hasStructuredRoute) {
                return resolveApiRouteDescription(this.apiRouteDescriptions, {
                    method: record.method,
                    route: record.route,
                    path: record.path,
                });
            }
            return record?.route_description || record?.message || '-';
        },
        normalizeSuccess(value) {
            if (typeof value === 'boolean') { return value; }
            if (typeof value === 'number') { return value !== 0; }
            if (typeof value === 'string') {
                const normalized = value.trim().toLowerCase();
                if (['false', '0', 'fail', 'failed', 'error'].includes(normalized)) { return false; }
                if (['true', '1', 'success', 'ok'].includes(normalized)) { return true; }
            }
            return Boolean(value);
        },
        formatTime(value) {
            if (!value) { return '-'; }
            const dateValue = this.normalizeTimeValue(value);
            if (Number.isNaN(new Date(dateValue).getTime())) { return value; }
            return window.formatDate ? window.formatDate(dateValue) : dateValue;
        },
        normalizeTimeValue(value) {
            if (typeof value !== 'string') { return value; }
            return value.replace(/\.(\d{3})\d+(Z|[+-]\d{2}:?\d{2})$/, '.$1$2');
        },
        toQueryTime(value) {
            if (!value) { return ''; }
            if (value instanceof Date) { return value.toISOString(); }
            const date = new Date(value);
            if (!Number.isNaN(date.getTime())) { return date.toISOString(); }
            return value;
        },
    },
};
</script>

<style scoped>
.audit-header {
    min-height: 32px;
}

.audit-filter {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 12px;
}

.audit-filter :deep(.arco-input-wrapper),
.audit-filter :deep(.arco-select) {
    width: 180px;
}

.ellipsis {
    display: inline-block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    vertical-align: middle;
    white-space: nowrap;
}

.mono {
    font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}

.audit-detail-table pre {
    margin: 0;
    white-space: pre-wrap;
    word-break: break-all;
}
</style>
