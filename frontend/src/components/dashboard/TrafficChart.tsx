import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { Activity } from 'lucide-react';
import { useTranslation } from 'react-i18next';

interface TrafficData {
    name: string;
    upload: number;
    download: number;
}

interface TrafficChartProps {
    data: TrafficData[];
}

const TrafficChart = ({ data }: TrafficChartProps) => {
    const { t } = useTranslation();
    return (
        <div className="lg:col-span-2 bg-card p-6 rounded-xl shadow-lg border border-border">
            <div className="flex items-center justify-between mb-6">
                <h3 className="text-lg font-semibold">{t('dashboard.networkTraffic')}</h3>
                <Activity size={20} className="text-muted-foreground" />
            </div>
            <div className="h-[300px] w-full">
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
                        <XAxis dataKey="name" hide />
                        <YAxis stroke="hsl(var(--muted-foreground))" />
                        <Tooltip
                            contentStyle={{ backgroundColor: 'hsl(var(--popover))', borderColor: 'hsl(var(--border))', color: 'hsl(var(--popover-foreground))' }}
                            itemStyle={{ color: 'hsl(var(--popover-foreground))' }}
                        />
                        <Area type="monotone" dataKey="download" stroke="#8884d8" fillOpacity={1} fill="url(#colorDownload)" name={`${t('dashboard.download')} (KB/s)`} />
                        <Area type="monotone" dataKey="upload" stroke="#82ca9d" fillOpacity={1} fill="url(#colorUpload)" name={`${t('dashboard.upload')} (KB/s)`} />
                    </AreaChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
};

export default TrafficChart;
