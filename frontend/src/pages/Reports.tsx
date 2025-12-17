import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, PieChart, Pie, Cell } from 'recharts';
import { ArrowDown, ArrowUp, Clock, FileText } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import PageHeader from '../components/common/PageHeader';

const dailyData = [
    { day: 'Mon', upload: 400, download: 2400 },
    { day: 'Tue', upload: 300, download: 1398 },
    { day: 'Wed', upload: 200, download: 9800 },
    { day: 'Thu', upload: 278, download: 3908 },
    { day: 'Fri', upload: 189, download: 4800 },
    { day: 'Sat', upload: 239, download: 3800 },
    { day: 'Sun', upload: 349, download: 4300 },
];

const protocolData = [
    { name: 'HTTPS', value: 400 },
    { name: 'HTTP', value: 300 },
    { name: 'DNS', value: 300 },
    { name: 'SOCKS5', value: 200 },
];

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

export default function Reports() {
    const { t } = useTranslation();

    return (
        <div className="p-6 bg-background min-h-screen text-foreground">
            <PageHeader
                title={t('reports.title')}
                description="Network usage analytics and insights"
                icon={FileText}
                iconColorClass="text-orange-400"
            />

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mt-6">
                {/* Traffic History Chart */}
                <div className="bg-card rounded-xl border border-border p-6 shadow-lg">
                    <h2 className="text-lg font-bold mb-6 text-card-foreground">{t('reports.dailyTraffic')}</h2>
                    <div className="h-80">
                        <ResponsiveContainer width="100%" height="100%">
                            <BarChart data={dailyData}>
                                <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--muted-foreground)/0.2)" />
                                <XAxis dataKey="day" stroke="hsl(var(--muted-foreground))" />
                                <YAxis stroke="hsl(var(--muted-foreground))" />
                                <Tooltip
                                    contentStyle={{ backgroundColor: 'hsl(var(--popover))', borderColor: 'hsl(var(--border))', color: 'hsl(var(--popover-foreground))' }}
                                    cursor={{ fill: 'hsl(var(--muted)/0.2)' }}
                                    itemStyle={{ color: 'hsl(var(--popover-foreground))' }}
                                />
                                <Bar dataKey="download" name={t('dashboard.download') || ''} fill="#8884d8" radius={[4, 4, 0, 0]} />
                                <Bar dataKey="upload" name={t('dashboard.upload') || ''} fill="#82ca9d" radius={[4, 4, 0, 0]} />
                            </BarChart>
                        </ResponsiveContainer>
                    </div>
                </div>

                {/* Protocol Distribution Chart */}
                <div className="bg-card rounded-xl border border-border p-6 shadow-lg">
                    <h2 className="text-lg font-bold mb-6 text-card-foreground">{t('reports.protocolDistribution')}</h2>
                    <div className="h-80 flex items-center justify-center">
                        <ResponsiveContainer width="100%" height="100%">
                            <PieChart>
                                <Pie
                                    data={protocolData}
                                    cx="50%"
                                    cy="50%"
                                    innerRadius={60}
                                    outerRadius={100}
                                    fill="#8884d8"
                                    paddingAngle={5}
                                    dataKey="value"
                                >
                                    {protocolData.map((entry, index) => (
                                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                    ))}
                                </Pie>
                                <Tooltip
                                    contentStyle={{ backgroundColor: 'hsl(var(--popover))', borderColor: 'hsl(var(--border))', color: 'hsl(var(--popover-foreground))' }}
                                    itemStyle={{ color: 'hsl(var(--popover-foreground))' }}
                                />
                            </PieChart>
                        </ResponsiveContainer>
                    </div>
                </div>
            </div>

            {/* Summary Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-6">
                <div className="bg-card rounded-xl border border-border p-6 shadow-lg">
                    <div className="flex items-center gap-4">
                        <div className="p-3 rounded-full bg-blue-500/20 text-blue-400">
                            <ArrowDown size={24} />
                        </div>
                        <div>
                            <p className="text-muted-foreground text-sm font-medium">{t('reports.totalDownload')}</p>
                            <h3 className="text-2xl font-bold text-foreground">45.2 GB</h3>
                        </div>
                    </div>
                </div>
                <div className="bg-card rounded-xl border border-border p-6 shadow-lg">
                    <div className="flex items-center gap-4">
                        <div className="p-3 rounded-full bg-green-500/20 text-green-400">
                            <ArrowUp size={24} />
                        </div>
                        <div>
                            <p className="text-muted-foreground text-sm font-medium">{t('reports.totalUpload')}</p>
                            <h3 className="text-2xl font-bold text-foreground">12.8 GB</h3>
                        </div>
                    </div>
                </div>
                <div className="bg-card rounded-xl border border-border p-6 shadow-lg">
                    <div className="flex items-center gap-4">
                        <div className="p-3 rounded-full bg-purple-500/20 text-purple-400">
                            <Clock size={24} />
                        </div>
                        <div>
                            <p className="text-muted-foreground text-sm font-medium">{t('reports.activeTime')}</p>
                            <h3 className="text-2xl font-bold text-foreground">142h 30m</h3>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
