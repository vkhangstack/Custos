import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Activity } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { formatBytes } from '../../utils/formatting';

interface TrafficData {
    name: string;
    timestamp: string | number; // Support both ISO string (time.Time) and Unix Milli
    upload: number;
    download: number;
}

interface TrafficChartProps {
    data: TrafficData[];
    timeRange: string;
}

const TrafficChart = ({ data, timeRange }: TrafficChartProps) => {
    const { t } = useTranslation();
    return (
        <div className="lg:col-span-2 bg-card p-6 rounded-xl shadow-lg border border-border">
            <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-semibold">{t('dashboard.networkTraffic')}</h3>
                <Activity size={20} className="text-muted-foreground" />
            </div>
            <div className="h-[300px] w-full pointer-events-none">
                <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={data}>
                        <defs>
                            <linearGradient id="colorDownload" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
                                <stop offset="95%" stopColor="#8884d8" stopOpacity={0} />
                            </linearGradient>
                            <linearGradient id="colorUpload" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#82ca9d" stopOpacity={0.8} />
                                <stop offset="95%" stopColor="#82ca9d" stopOpacity={0} />
                            </linearGradient>
                        </defs>
                        <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--muted-foreground)/0.2)" />
                        <XAxis 
                            dataKey="timestamp" 
                            hide={false}
                            stroke="hsl(var(--muted-foreground))"
                            tick={{ fontSize: 10 }}
                            tickFormatter={(str) => {
                                if (!str) return '';
                                const date = new Date(str);
                                return timeRange === '24h' 
                                    ? date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) // Just time is cleaner, or include date if needed
                                    : date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
                            }}
                            minTickGap={30}
                        />
                        <YAxis 
                            stroke="hsl(var(--muted-foreground))" 
                            tickFormatter={formatBytes} 
                            tick={{ fontSize: 10 }}
                        />
                        <Tooltip
                            labelFormatter={(label) => {
                                if (!label) return '';
                                const date = new Date(label);
                                return date.toLocaleString();
                            }}
                            contentStyle={{ backgroundColor: 'hsl(var(--popover))', borderColor: 'hsl(var(--border))', color: 'hsl(var(--popover-foreground))' }}
                            itemStyle={{ color: 'hsl(var(--popover-foreground))' }}
                            formatter={(value) => formatBytes(value as number)}
                        />
                        <Area type="monotone" dataKey="download" stroke="#8884d8" fillOpacity={1} fill="url(#colorDownload)" name={t('dashboard.download') as string} />
                        <Area type="monotone" dataKey="upload" stroke="#82ca9d" fillOpacity={1} fill="url(#colorUpload)" name={t('dashboard.upload') as string} />
                    </AreaChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
};

export default TrafficChart;
