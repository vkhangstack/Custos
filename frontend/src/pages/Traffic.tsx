import { useState, useEffect } from 'react';
import { ArrowUpDown } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import LogTable from '../components/traffic/LogTable';
import PageHeader from '../components/common/PageHeader';

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

const processes = ['chrome.exe', 'code.exe', 'spotify.exe', 'slack.exe', 'system'];
const domains = ['google.com', 'a.ads-api.com', 'github.com', 'api.twitter.com', 'netflix.com', 'doubleclick.net'];
const methods = ['GET', 'POST', 'CONNECT', 'HEAD'];

const generateLog = (): RequestLog => {
    const isHttp = Math.random() > 0.3;
    const domain = domains[Math.floor(Math.random() * domains.length)];

    return {
        id: Math.random().toString(36).substr(2, 9),
        timestamp: new Date().toLocaleTimeString(),
        process: processes[Math.floor(Math.random() * processes.length)],
        domain: domain,
        protocol: isHttp ? 'HTTPS' : 'DNS',
        method: isHttp ? methods[Math.floor(Math.random() * methods.length)] : 'QUERY',
        status: isHttp ? (Math.random() > 0.9 ? 404 : 200) : (Math.random() > 0.95 ? 503 : 200),
        size: Math.floor(Math.random() * 5000) + ' B',
        latency: Math.floor(Math.random() * 200) + 'ms',
    };
};

export default function Traffic() {
    const { t } = useTranslation();
    const [logs, setLogs] = useState<RequestLog[]>([]);

    useEffect(() => {
        // Initial data
        const initialLogs = Array.from({ length: 10 }, generateLog);
        setLogs(initialLogs);

        const interval = setInterval(() => {
            setLogs(prev => {
                const newLog = generateLog();
                return [newLog, ...prev].slice(0, 50); // Keep last 50 logs
            });
        }, 2000);

        return () => clearInterval(interval);
    }, []);

    return (
        <div className="p-6 bg-background min-h-screen text-foreground">
            <PageHeader
                title={t('traffic.title')}
                icon={ArrowUpDown}
                iconColorClass="text-blue-400"
            />

            <LogTable logs={logs} />
        </div>
    );
}
