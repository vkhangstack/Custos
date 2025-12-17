import { useState, useEffect } from 'react';
import { Wifi, ArrowUp, ArrowDown } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import StatCard from '../components/dashboard/StatCard';
import TrafficChart from '../components/dashboard/TrafficChart';
import TopUsageList from '../components/dashboard/TopUsageList';

const generateData = () => {
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

const topUsageData = [
    { domain: 'google.com', usage: '1.2 GB', process: 'chrome.exe' },
    { domain: 'youtube.com', usage: '850 MB', process: 'chrome.exe' },
    { domain: 'api.github.com', usage: '120 MB', process: 'Code.exe' },
    { domain: 'slack.com', usage: '95 MB', process: 'slack.exe' },
    { domain: 'zoom.us', usage: '450 MB', process: 'zoom.exe' },
];

export default function Dashboard() {
    const { t } = useTranslation();
    const [data, setData] = useState(generateData());
    const [activeConnections, setActiveConnections] = useState(124);

    useEffect(() => {
        const interval = setInterval(() => {
            setData(prevData => {
                const newData = [...prevData.slice(1)];
                newData.push({
                    name: Math.floor(Math.random() * 100).toString(),
                    upload: Math.floor(Math.random() * 500) + 100,
                    download: Math.floor(Math.random() * 1000) + 200,
                });
                return newData;
            });
            setActiveConnections(prev => prev + Math.floor(Math.random() * 5) - 2);
        }, 1000);

        return () => clearInterval(interval);
    }, []);

    return (
        <div className="p-6 space-y-6 bg-background min-h-screen text-foreground">
            <h1 className="text-2xl font-bold mb-4">{t('dashboard.title')}</h1>

            {/* Status Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <StatCard
                    title={t('dashboard.activeConnections')}
                    value={activeConnections}
                    icon={Wifi}
                    iconColorClass="text-blue-400"
                    bgClass="bg-blue-500/20"
                />
                <StatCard
                    title={t('dashboard.totalUpload')}
                    value="4.2 GB"
                    icon={ArrowUp}
                    iconColorClass="text-green-400"
                    bgClass="bg-green-500/20"
                />
                <StatCard
                    title={t('dashboard.totalDownload')}
                    value="12.5 GB"
                    icon={ArrowDown}
                    iconColorClass="text-purple-400"
                    bgClass="bg-purple-500/20"
                />
            </div>

            {/* Charts Section */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <TrafficChart data={data} />
                <TopUsageList items={topUsageData} />
            </div>
        </div>
    );
}
