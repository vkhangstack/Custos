import { useState, useEffect } from 'react';
import { Activity, ArrowDown, ArrowUp, Filter, Search, RefreshCw, XCircle, CheckCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { GetLogs, GetStats, GetSystemConnections } from '../../wailsjs/go/main/App';
import { core, system } from '../../wailsjs/go/models';
import PageHeader from '../components/common/PageHeader';
import Select from '../components/common/Select';

export default function Traffic() {
    const { t } = useTranslation();
    const [searchQuery, setSearchQuery] = useState('');
    const [logs, setLogs] = useState<core.LogEntry[]>([]);
    const [stats, setStats] = useState<core.Stats>(new core.Stats());
    const [connections, setConnections] = useState<system.ConnectionInfo[]>([]);
    const [activeTab, setActiveTab] = useState<'logs' | 'connections'>('logs');
    const [statusFilter, setStatusFilter] = useState('all');
    const [typeFilter, setTypeFilter] = useState('all');
    const [isRefreshing, setIsRefreshing] = useState(false);

    const fetchData = async () => {
        try {
            const [fetchedLogs, fetchedStats, fetchedConns] = await Promise.all([
                GetLogs(),
                GetStats(),
                GetSystemConnections()
            ]);

            setLogs(fetchedLogs || []);
            setStats(fetchedStats || new core.Stats());
            setConnections(fetchedConns || []);
        } catch (error) {
            console.error("Failed to fetch traffic data:", error);
        }
    };

    const handleRefresh = async () => {
        setIsRefreshing(true);
        await fetchData();
        setTimeout(() => setIsRefreshing(false), 500);
    };

    useEffect(() => {
        fetchData();
        // Poll every 2 seconds for fresh data
        const interval = setInterval(fetchData, 2000);
        return () => clearInterval(interval);
    }, []);

    const formatBytes = (bytes: number) => {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    const getStatusColor = (status: string) => {
        switch (status) {
            case 'blocked': return 'text-red-500 bg-red-500/10';
            case 'allowed': return 'text-green-500 bg-green-500/10';
            default: return 'text-yellow-500 bg-yellow-500/10';
        }
    };

    const filteredLogs = logs.filter(log => {
        const matchesSearch = log.domain?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            log.process_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            log.dst_ip?.includes(searchQuery);
        const matchesStatus = statusFilter === 'all' || log.status === statusFilter;
        const matchesType = typeFilter === 'all' || log.type === typeFilter;
        return matchesSearch && matchesStatus && matchesType;
    });

    const filteredConnections = connections.filter(conn =>
        conn.process_name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
        conn.remote_addr?.includes(searchQuery)
    );

    const actions = (
        <div className="flex gap-2 items-center">
            <Select
                value={statusFilter}
                onChange={setStatusFilter}
                options={[
                    { label: 'All Status', value: 'all' },
                    { label: 'Allowed', value: 'allowed' },
                    { label: 'Blocked', value: 'blocked' },
                ]}
                className="w-32"
            />
            <Select
                value={typeFilter}
                onChange={setTypeFilter}
                options={[
                    { label: 'All Types', value: 'all' },
                    { label: 'DNS', value: 'dns' },
                    { label: 'Proxy', value: 'proxy' },
                ]}
                className="w-28"
            />
            <div className="relative">
                <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground" size={16} />
                <input
                    type="text"
                    placeholder={t('traffic.searchPlaceholder') as string}
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="max-w-32 pl-9 pr-4 py-2 bg-muted border border-border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/50 text-foreground w-64"
                />
            </div>
            <button
                onClick={handleRefresh}
                className="p-2 hover:bg-muted rounded-lg border border-border text-muted-foreground hover:text-foreground transition-colors"
                title="Refresh"
            >
                <RefreshCw size={18} className={isRefreshing ? 'animate-spin' : ''} />
            </button>
        </div>
    );

    return (
        <div className="p-6 bg-background min-h-screen text-foreground">
            <PageHeader
                title={t('traffic.title')}
                icon={Activity}
                iconColorClass="text-blue-400"
                actions={actions}
            />

            {/* Stats Overview */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                <div className="bg-card border border-border rounded-xl p-4 flex items-center justify-between shadow-sm">
                    <div>
                        <p className="text-sm text-muted-foreground font-medium mb-1">{t('traffic.upload')}</p>
                        <h3 className="text-2xl font-bold text-foreground">{formatBytes(stats.total_upload || 0)}</h3>
                    </div>
                    <div className="p-3 bg-blue-500/10 text-blue-500 rounded-full">
                        <ArrowUp size={24} />
                    </div>
                </div>
                <div className="bg-card border border-border rounded-xl p-4 flex items-center justify-between shadow-sm">
                    <div>
                        <p className="text-sm text-muted-foreground font-medium mb-1">{t('traffic.download')}</p>
                        <h3 className="text-2xl font-bold text-foreground">{formatBytes(stats.total_download || 0)}</h3>
                    </div>
                    <div className="p-3 bg-green-500/10 text-green-500 rounded-full">
                        <ArrowDown size={24} />
                    </div>
                </div>
                <div className="bg-card border border-border rounded-xl p-4 flex items-center justify-between shadow-sm">
                    <div>
                        <p className="text-sm text-muted-foreground font-medium mb-1">{t('traffic.activeConnections')}</p>
                        <h3 className="text-2xl font-bold text-foreground">{stats.active_connections || connections.length}</h3>
                    </div>
                    <div className="p-3 bg-purple-500/10 text-purple-500 rounded-full">
                        <Activity size={24} />
                    </div>
                </div>
            </div>

            {/* Tabs */}
            <div className="flex gap-4 mb-4 border-b border-border">
                <button
                    onClick={() => setActiveTab('logs')}
                    className={`pb-2 px-4 text-sm font-medium transition-colors border-b-2 ${activeTab === 'logs' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'}`}
                >
                    {t('traffic.logs')}
                </button>
                <button
                    onClick={() => setActiveTab('connections')}
                    className={`pb-2 px-4 text-sm font-medium transition-colors border-b-2 ${activeTab === 'connections' ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'}`}
                >
                    System Connections
                </button>
            </div>

            {/* Content Table */}
            <div className="bg-card rounded-xl border border-border shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className={`w-full text-left text-sm`}>
                        <thead className="bg-muted/50 border-b border-border">
                            <tr>
                                {activeTab === 'logs' ? (
                                    <>
                                        <th className="p-4 font-medium text-muted-foreground">Time</th>
                                        <th className="p-4 font-medium text-muted-foreground">Process</th>
                                        <th className="p-4 font-medium text-muted-foreground">Domain / IP</th>
                                        <th className="p-4 font-medium text-muted-foreground">Type</th>
                                        <th className="p-4 font-medium text-muted-foreground">Protocol</th>
                                        <th className="p-4 font-medium text-muted-foreground">Status</th>
                                        <th className="p-4 font-medium text-muted-foreground text-right">Size</th>
                                    </>
                                ) : (
                                    <>
                                        <th className="p-4 font-medium text-muted-foreground w-[100px]">PID</th>
                                        <th className="p-4 font-medium text-muted-foreground">Process</th>
                                        <th className="p-4 font-medium text-muted-foreground w-[200px]">Local Address</th>
                                        <th className="p-4 font-medium text-muted-foreground w-[200px]">Remote Address</th>
                                        <th className="p-4 font-medium text-muted-foreground">Protocol</th>
                                        <th className="p-4 font-medium text-muted-foreground">Status</th>
                                    </>
                                )}
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-border">
                            {activeTab === 'logs' ? (
                                filteredLogs.length > 0 ? (
                                    filteredLogs.slice().reverse().map((log) => (
                                        <tr key={log.id} className="hover:bg-muted/30 transition-colors">
                                            <td className="p-4 text-muted-foreground whitespace-nowrap">
                                                {new Date(log.timestamp).toLocaleTimeString()}
                                            </td>
                                            <td className="p-4 font-medium text-foreground">
                                                {log.process_name || '-'}
                                                {log.process_id ? <span className="text-xs text-muted-foreground ml-1">({log.process_id})</span> : ''}
                                            </td>
                                            <td className="p-4 text-foreground font-mono text-xs min-w-[200px]">
                                                <div className="font-medium text-sm mb-0.5">{log.domain || '-'}</div>
                                                <div className="text-muted-foreground">{log.dst_ip}:{log.dst_port}</div>
                                            </td>
                                            <td className="p-4">
                                                <span className={`px-2 py-0.5 rounded text-xs font-medium border ${log.type === 'dns' ? 'bg-blue-500/10 text-blue-500 border-blue-500/20' : 'bg-purple-500/10 text-purple-500 border-purple-500/20'}`}>
                                                    {log.type.toUpperCase()}
                                                </span>
                                            </td>
                                            <td className="p-4 text-muted-foreground uppercase text-xs font-semibold">{log.protocol}</td>
                                            <td className="p-4">
                                                <span className={`flex items-center gap-1.5 px-2 py-0.5 rounded text-xs font-medium w-fit ${getStatusColor(log.status)}`}>
                                                    {log.status === 'blocked' ? <XCircle size={12} /> : <CheckCircle size={12} />}
                                                    {log.status}
                                                </span>
                                            </td>
                                            <td className="p-4 text-right text-muted-foreground font-mono">
                                                {formatBytes((log.bytes_sent || 0) + (log.bytes_recv || 0))}
                                            </td>
                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={7} className="p-8 text-center text-muted-foreground">
                                            No traffic logs found.
                                        </td>
                                    </tr>
                                )
                            ) : (
                                filteredConnections.length > 0 ? (
                                    filteredConnections.map((conn, idx) => (
                                        <tr key={`${conn.pid}-${idx}`} className="hover:bg-muted/30 transition-colors">
                                            <td className="p-4 text-muted-foreground font-mono">{conn.pid}</td>
                                            <td className="p-4 font-medium text-foreground">{conn.process_name}</td>
                                            <td className="p-4 text-muted-foreground font-mono text-xs max-w-[200px] truncate" title={conn.local_addr}>{conn.local_addr}</td>
                                            <td className="p-4 text-muted-foreground font-mono text-xs max-w-[200px] truncate" title={conn.remote_addr}>{conn.remote_addr}</td>
                                            <td className="p-4 text-muted-foreground uppercase text-xs font-semibold">{conn.protocol}</td>
                                            <td className="p-4">
                                                <span className="px-2 py-0.5 rounded text-xs font-medium bg-secondary text-secondary-foreground">
                                                    {conn.status}
                                                </span>
                                            </td>
                                        </tr>
                                    ))
                                ) : (
                                    <tr>
                                        <td colSpan={6} className="p-8 text-center text-muted-foreground">
                                            No active connections found.
                                        </td>
                                    </tr>
                                )
                            )}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}
