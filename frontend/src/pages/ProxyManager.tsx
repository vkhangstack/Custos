import { useState } from 'react';
import { Globe, Plus } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import PageHeader from '../components/common/PageHeader';
import ProxyCard from '../components/proxy/ProxyCard';

interface Proxy {
    id: string;
    name: string;
    type: 'SOCKS5' | 'HTTP';
    address: string;
    port: number;
    latency: number;
    status: 'active' | 'dead' | 'testing';
    enabled: boolean;
    country: string;
}

const initialProxies: Proxy[] = [
    { id: '1', name: 'US Premium 1', type: 'SOCKS5', address: '192.168.1.101', port: 1080, latency: 45, status: 'active', enabled: true, country: 'US' },
    { id: '2', name: 'SG Low Latency', type: 'HTTP', address: '10.0.0.5', port: 8080, latency: 12, status: 'active', enabled: true, country: 'SG' },
    { id: '3', name: 'JP Privacy', type: 'SOCKS5', address: '172.16.0.22', port: 1085, latency: 89, status: 'testing', enabled: false, country: 'JP' },
    { id: '4', name: 'Dead Proxy', type: 'HTTP', address: '192.168.0.2', port: 3128, latency: -1, status: 'dead', enabled: false, country: 'VN' },
];

export default function ProxyManager() {
    const { t } = useTranslation();
    const [proxies, setProxies] = useState<Proxy[]>(initialProxies);

    const toggleProxy = (id: string) => {
        setProxies(proxies.map(p => p.id === id ? { ...p, enabled: !p.enabled } : p));
    };

    const deleteProxy = (id: string) => {
        setProxies(proxies.filter(p => p.id !== id));
    };

    const actions = (
        <button className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors font-medium">
            <Plus size={18} />
            {t('proxy.addProxy')}
        </button>
    );

    return (
        <div className="p-6 bg-background min-h-screen text-foreground">
            <PageHeader
                title={t('proxy.title')}
                icon={Globe}
                iconColorClass="text-purple-400"
                actions={actions}
            />

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                {proxies.map((proxy) => (
                    <ProxyCard
                        key={proxy.id}
                        proxy={proxy}
                        onToggle={toggleProxy}
                        onDelete={deleteProxy}
                    />
                ))}

                {/* Add New Placeholder Card */}
                <button className="bg-muted/30 rounded-xl border border-dashed border-border p-5 flex flex-col items-center justify-center gap-3 text-muted-foreground hover:text-primary hover:border-primary/50 hover:bg-card transition-all group min-h-[250px]">
                    <div className="p-4 rounded-full bg-muted group-hover:bg-primary/10 transition-colors">
                        <Plus size={32} />
                    </div>
                    <span className="font-medium">{t('proxy.addNewProxy')}</span>
                </button>
            </div>
        </div>
    );
}
