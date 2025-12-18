import { useState, useEffect, useRef } from 'react';
import { Wifi, ArrowUp, ArrowDown, Shield, ShieldAlert, Loader2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import StatCard from '../components/dashboard/StatCard';
import TrafficChart from '../components/dashboard/TrafficChart';
import TopUsageList from '../components/dashboard/TopUsageList';
import { GetStats, GetSystemConnections, GetChartData, EnableProtection, GetProtectionStatus, GetStartupStatus } from '../../wailsjs/go/main/App';
import { core, system } from '../../wailsjs/go/models';
import { formatBytes } from '../utils/formatting';
import { useToast } from '../context/ToastContext';


export default function Dashboard() {
    const { t } = useTranslation();
    const [chartData, setChartData] = useState<any[]>([]);
    const [timeRange, setTimeRange] = useState<string>('1h');
    const { showToast } = useToast();

    // Real State
    const [stats, setStats] = useState<core.Stats>(new core.Stats());
    const [connections, setConnections] = useState<system.ConnectionInfo[]>([]);
    const [protectionEnabled, setProtectionEnabled] = useState<boolean | null>(null);
    const [isProtectionLoading, setIsProtectionLoading] = useState(false);

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
                const timeLabel = `${now.getHours()}:${now.getMinutes().toString().padStart(2, '0')}:${now.getSeconds().toString().padStart(2, '0')}`;

                let uploadRate = 0;
                let downloadRate = 0;

                if (prevStatsRef.current) {
                    const up = currentStats.total_upload || 0;
                    const down = currentStats.total_download || 0;
                    const prevUp = prevStatsRef.current.total_upload || 0;
                    const prevDown = prevStatsRef.current.total_download || 0;
                    
                    uploadRate = up - prevUp;
                    downloadRate = down - prevDown;
                }

                // Prevent negative spikes if stats reset
                if (uploadRate < 0) uploadRate = 0;
                if (downloadRate < 0) downloadRate = 0;

                const newData = [...prevData];
                newData.push({
                    name: timeLabel,
                    timestamp: now.toISOString(),
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
            const history = await GetChartData(timeRange);
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
    }, [timeRange]); // Re-run when timeRange changes

    const toggleProtection = async () => {
        if (isProtectionLoading || protectionEnabled === null) return;
        setIsProtectionLoading(true);
        try {
            // Artificial delay for UX
            await new Promise(resolve => setTimeout(resolve, 500));
            await EnableProtection(!protectionEnabled);
            setProtectionEnabled(!protectionEnabled);
        } catch (error) {
            console.error("Failed to toggle protection:", error);
        } finally {
            setIsProtectionLoading(false);
            showToast((!protectionEnabled ? t('dashboard.enableProtection') : t('dashboard.disableProtection')), 'success');
        }
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
                    disabled={isProtectionLoading || protectionEnabled === null}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg font-medium transition-colors ${
                        protectionEnabled === null
                            ? 'bg-muted text-muted-foreground border border-muted-foreground/20'
                            : protectionEnabled
                            ? 'bg-green-500/20 text-green-500 border border-green-500/50 hover:bg-green-500/30'
                            : 'bg-red-500/20 text-red-500 border border-red-500/50 hover:bg-red-500/30'
                    } ${isProtectionLoading || protectionEnabled === null ? 'opacity-50 cursor-not-allowed' : ''} min-w-[180px] justify-center`}
                >
                    {isProtectionLoading || protectionEnabled === null ? (
                        <Loader2 size={20} className="animate-spin" />
                    ) : (
                        protectionEnabled ? <Shield size={20} /> : <ShieldAlert size={20} />
                    )}
                    {protectionEnabled === null
                        ? t('Loading...')
                        : protectionEnabled
                        ? t('dashboard.disableProtection')
                        : t('dashboard.enableProtection')}
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
                <div className="lg:col-span-2 space-y-4">
                    <div className="flex justify-end space-x-2">
                        {['1h', '3h', '24h'].map((range) => (
                            <button
                                key={range}
                                onClick={() => setTimeRange(range)}
                                className={`px-3 py-1 text-sm rounded-md transition-colors ${
                                    timeRange === range
                                        ? 'bg-primary text-primary-foreground'
                                        : 'bg-muted text-muted-foreground hover:bg-muted/80'
                                }`}
                            >
                                {range}
                            </button>
                        ))}
                    </div>
                    <TrafficChart data={chartData} timeRange={timeRange} />
                </div>
                <div className="lg:col-span-1">
                     <TopUsageList items={topUsageList} />
                </div>
            </div>
        </div>
    );
}
