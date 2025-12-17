import { useState, useEffect, useRef } from 'react';
import { Wifi, ArrowUp, ArrowDown, Shield, ShieldAlert } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import StatCard from '../components/dashboard/StatCard';
import TrafficChart from '../components/dashboard/TrafficChart';
import TopUsageList from '../components/dashboard/TopUsageList';
import { GetStats, GetSystemConnections, GetChartData, EnableProtection, GetProtectionStatus } from '../../wailsjs/go/main/App';
import { core, system } from '../../wailsjs/go/models';

const generateMockChartData = () => {
    // Keep mock chart data for now as GetStats doesn't return historical data series yet
    // Future improvement: Backend should provide historical stats
    const data = [];
    for (let i = 0; i < 20; i++) {
        data.push({
            name: i.toString(),
            upload: Math.floor(Math.random() * 500) + 100,
            download: Math.floor(Math.random() * 1000) + 200,
        });
    }
    return data;
};

export default function Dashboard() {
    const { t } = useTranslation();
    const [chartData, setChartData] = useState<any[]>([]);

    // Real State
    const [stats, setStats] = useState<core.Stats>(new core.Stats());
    const [connections, setConnections] = useState<system.ConnectionInfo[]>([]);
    const [protectionEnabled, setProtectionEnabled] = useState(false);

    // Refs for rate calculation
    const prevStatsRef = useRef<core.Stats | null>(null);

    const fetchData = async () => {
        try {
            const [fetchedStats, fetchedConns] = await Promise.all([
                GetStats(),
                GetSystemConnections()
            ]);

            const currentStats = fetchedStats || new core.Stats();
            setStats(currentStats);
            setConnections(fetchedConns || []);

            // Calculate live rate for the chart (append to history)
            setChartData(prevData => {
                const now = new Date();
                const timeLabel = now.getHours() + ':' + now.getMinutes().toString().padStart(2, '0') + ':' + now.getSeconds().toString().padStart(2, '0');
                // Note: Chart history from backend is per minute (%H:%M). Live update is per 2s.
                // This mismatch might look weird (minute bars vs second bars).
                // Ideally backend history should match frontend granularity OR frontend should aggregate.
                // For now, let's just append.

                let uploadRate = 0;
                let downloadRate = 0;

                if (prevStatsRef.current) {
                    uploadRate = currentStats.total_upload - prevStatsRef.current.total_upload;
                    downloadRate = currentStats.total_download - prevStatsRef.current.total_download;
                }

                // Prevent negative spikes if stats reset
                if (uploadRate < 0) uploadRate = 0;
                if (downloadRate < 0) downloadRate = 0;

                const newData = [...prevData.slice(1)];
                newData.push({
                    name: timeLabel,
                    upload: uploadRate,
                    download: downloadRate,
                });
                return newData;
            });

            prevStatsRef.current = currentStats;

        } catch (error) {
            console.error("Failed to fetch dashboard data:", error);
        }
    };

    useEffect(() => {
        // Initial fetch with History
        const init = async () => {
            const history = await GetChartData();
            if (history && history.length > 0) {
                setChartData(history);
            } else {
                // Fill with empty if no history
                setChartData(Array(20).fill({ name: '', upload: 0, download: 0 }));
            }

            // Fetch initial protection status
            const enabled = await GetProtectionStatus();
            setProtectionEnabled(enabled);

            await fetchData();
        };
        init();

        const interval = setInterval(() => {
            fetchData();
        }, 2000); // Poll every 2 seconds

        return () => clearInterval(interval);
    }, []);

    const toggleProtection = async () => {
        try {
            await EnableProtection(!protectionEnabled);
            setProtectionEnabled(!protectionEnabled);
        } catch (error) {
            console.error("Failed to toggle protection:", error);
        }
    };

    const formatBytes = (bytes: number) => {
        if (!bytes || bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    };

    const getRootDomain = (domain: string) => {
        if (!domain) return '-';
        // Return IP as is
        if (/^(\d{1,3}\.){3}\d{1,3}$/.test(domain)) return domain;

        const parts = domain.split('.');
        // Simple heuristic: if parts > 2, take last 2 (e.g. mail.google.com -> google.com)
        // Note: This is imperfect for ccTLDs (google.co.uk -> co.uk), but sufficient for now without heavy libraries.
        if (parts.length > 2) {
            return parts.slice(-2).join('.');
        }
        return domain;
    };

    // Transform Top Domains for UI
    const topUsageList = Object.entries(stats.top_domains || {})
        .sort(([, usageA], [, usageB]) => usageB - usageA) // Sort by usage desc (value)
        .map(([domain, usage]) => ({
            domain: getRootDomain(domain),
            usage: formatBytes(usage),
            value: usage, // Pass raw numeric value for chart
            process: '-' // Process info not yet aggregated in stats top_domains map
        }))
        .slice(0, 5);


    return (
        <div className="p-6 space-y-6 bg-background min-h-screen text-foreground">
            <div className="flex justify-between items-center mb-4">
                <h1 className="text-2xl font-bold">{t('dashboard.title')}</h1>
                <button
                    onClick={toggleProtection}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-colors ${protectionEnabled
                        ? 'bg-green-500/20 text-green-500 border border-green-500/50 hover:bg-green-500/30'
                        : 'bg-red-500/20 text-red-500 border border-red-500/50 hover:bg-red-500/30'
                        } `}
                >
                    {protectionEnabled ? <Shield size={20} /> : <ShieldAlert size={20} />}
                    {protectionEnabled ? t('dashboard.disableProtection') : t('dashboard.enableProtection')}
                </button>
            </div>

            {/* Status Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <StatCard
                    title={t('dashboard.activeConnections')}
                    value={connections.length}
                    icon={Wifi}
                    iconColorClass="text-blue-400"
                    bgClass="bg-blue-500/20"
                />
                <StatCard
                    title={t('dashboard.totalUpload')}
                    value={formatBytes(stats.total_upload)}
                    icon={ArrowUp}
                    iconColorClass="text-green-400"
                    bgClass="bg-green-500/20"
                />
                <StatCard
                    title={t('dashboard.totalDownload')}
                    value={formatBytes(stats.total_download)}
                    icon={ArrowDown}
                    iconColorClass="text-purple-400"
                    bgClass="bg-purple-500/20"
                />
            </div>

            {/* Charts Section */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <TrafficChart data={chartData} />
                <TopUsageList items={topUsageList} />
            </div>
        </div>
    );
}
