import { Clock, Globe } from 'lucide-react';
import { useTranslation } from 'react-i18next';

interface RequestLog {
    id: string;
    timestamp: string;
    process: string;
    domain: string;
    protocol: string;
    method?: string;
    status: number;
    size: string;
    latency: string;
}

interface LogTableProps {
    logs: RequestLog[];
}

const LogTable = ({ logs }: LogTableProps) => {
    const { t } = useTranslation();
    return (
        <div className="bg-card rounded-xl border border-border overflow-hidden shadow-lg">
            <div className="overflow-x-auto">
                <table className="w-full text-left text-sm text-muted-foreground">
                    <thead className="bg-muted/50 text-secondary-foreground uppercase font-medium">
                        <tr>
                            <th className="px-6 py-4">{t('traffic.table.time')}</th>
                            <th className="px-6 py-4">{t('traffic.table.process')}</th>
                            <th className="px-6 py-4">{t('traffic.table.domain')}</th>
                            <th className="px-6 py-4">{t('traffic.table.protocol')}</th>
                            <th className="px-6 py-4">{t('traffic.table.method')}</th>
                            <th className="px-6 py-4">{t('traffic.table.status')}</th>
                            <th className="px-6 py-4 text-right">{t('traffic.table.size')}</th>
                            <th className="px-6 py-4 text-right">{t('traffic.table.latency')}</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-border">
                        {logs.map((log) => (
                            <tr key={log.id} className="hover:bg-muted/50 transition-colors">
                                <td className="px-6 py-4 whitespace-nowrap font-mono text-xs text-muted-foreground">
                                    <div className="flex items-center gap-2">
                                        <Clock size={14} />
                                        {log.timestamp}
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-foreground font-medium">{log.process}</td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <div className="flex items-center gap-2">
                                        <Globe size={14} className={log.protocol === 'DNS' ? 'text-yellow-400' : 'text-blue-400'} />
                                        <span className={log.domain.includes('ads') ? 'text-red-400' : 'text-muted-foreground'}>
                                            {log.domain}
                                        </span>
                                    </div>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <span className={`px-2 py-1 rounded text-xs font-bold ${log.protocol === 'DNS' ? 'bg-yellow-500/10 text-yellow-400' : 'bg-blue-500/10 text-blue-400'
                                        }`}>
                                        {log.protocol}
                                    </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap font-mono text-xs">{log.method}</td>
                                <td className="px-6 py-4 whitespace-nowrap">
                                    <span className={`px-2 py-1 rounded text-xs font-bold ${log.status === 200 ? 'bg-green-500/10 text-green-400' : 'bg-red-500/10 text-red-400'
                                        }`}>
                                        {log.status}
                                    </span>
                                </td>
                                <td className="px-6 py-4 whitespace-nowrap text-right font-mono text-xs">{log.size}</td>
                                <td className="px-6 py-4 whitespace-nowrap text-right font-mono text-xs text-muted-foreground">{log.latency}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default LogTable;
