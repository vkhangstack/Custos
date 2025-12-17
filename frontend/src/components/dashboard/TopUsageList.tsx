import { useTranslation } from 'react-i18next';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts';

interface UsageItem {
    domain: string;
    usage: string;
    process: string;
    // We need numeric value for chart, but 'usage' is pre-formatted string.
    // Ideally parent should pass number.
    // But let's assume valid formatted string for display, and we need to pass numeric for chart.
    // Wait, Dashboard passes formatted string.
    // I need value.
    // Let's refactor Dashboard to pass numeric usage and TopUsageList handles formatting?
    // OR: Parse the string roughly?
    // Ideally Dashboard should pass raw value.
    // Let's change this in the next step. For now, let's setup the component.
    value?: number; // Add optional numeric value
}

interface TopUsageListProps {
    items: (UsageItem & { value?: number })[]; // Extend type
}

const TopUsageList = ({ items }: TopUsageListProps) => {
    const { t } = useTranslation();

    // Sort items by value just in case
    // Colors for bars
    const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8'];

    return (
        <div className="bg-card p-6 rounded-xl shadow-lg border border-border h-[400px] flex flex-col">
            <h3 className="text-lg font-semibold mb-4">{t('dashboard.topUsage')}</h3>
            <div className="flex-1 w-full min-h-0">
                <ResponsiveContainer width="100%" height="100%">
                    <BarChart
                        layout="vertical"
                        data={items}
                        margin={{
                            top: 5,
                            right: 30,
                            left: 20,
                            bottom: 5,
                        }}
                    >
                        <XAxis type="number" hide />
                        <YAxis
                            dataKey="domain"
                            type="category"
                            width={100}
                            tick={{ fill: 'hsl(var(--muted-foreground))', fontSize: 12 }}
                            tickFormatter={(val) => val.length > 15 ? val.substring(0, 15) + '...' : val}
                        />
                        <Tooltip
                            cursor={{ fill: 'transparent' }}
                            contentStyle={{ backgroundColor: 'hsl(var(--popover))', borderColor: 'hsl(var(--border))', color: 'hsl(var(--popover-foreground))' }}
                            formatter={(value: any, name: any, props: any) => [props.payload.usage, 'Usage']}
                        />
                        <Bar dataKey="value" radius={[0, 4, 4, 0]} barSize={20}>
                            {items.map((entry, index) => (
                                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                            ))}
                        </Bar>
                    </BarChart>
                </ResponsiveContainer>
            </div>
            {/* Legend / List below chart */}
            <div className="mt-4 space-y-2 max-h-[100px] overflow-y-auto">
                {items.slice(0, 3).map((item, index) => (
                    <div key={index} className="flex items-center justify-between text-sm">
                        <div className="flex items-center gap-2">
                            <div className="w-2 h-2 rounded-full" style={{ backgroundColor: COLORS[index % COLORS.length] }}></div>
                            <span className="text-muted-foreground truncate max-w-[150px]" title={item.domain}>{item.domain}</span>
                        </div>
                        <span className="font-medium">{item.usage}</span>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default TopUsageList;
