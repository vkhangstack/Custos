import { Power, Trash2, Edit2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';

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

interface ProxyCardProps {
    proxy: Proxy;
    onToggle: (id: string) => void;
    onDelete: (id: string) => void;
}

const ProxyCard = ({ proxy, onToggle, onDelete }: ProxyCardProps) => {
    const { t } = useTranslation();
    return (
        <div className={`flex flex-col justify-between h-full bg-card rounded-xl border p-4 shadow-lg transition-all ${proxy.enabled ? 'border-primary/50 shadow-primary/10' : 'border-border opacity-75'
            }`}>
            <div className="flex justify-between items-center mb-4 gap-2">
                <div className="flex items-center gap-3">
                    <span className={`w-3 h-3 rounded-full ${proxy.status === 'active' ? 'bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.6)]' :
                        proxy.status === 'testing' ? 'bg-yellow-500' : 'bg-red-500'
                        }`} />
                    <div>
                        <h3 className="font-bold text-base text-card-foreground">{proxy.name}</h3>
                        <p className="text-xs text-muted-foreground uppercase tracking-wider font-semibold bg-muted inline-block px-1.5 py-0.5 rounded mt-1">{proxy.type}</p>
                    </div>
                </div>
                <button
                    onClick={() => onToggle(proxy.id)}
                    className={`rounded-lg p-2 flex items-center justify-center transition-colors ${proxy.enabled ? 'bg-primary/20 text-primary hover:bg-primary/30' : 'bg-muted text-muted-foreground hover:text-foreground'
                        }`}
                >
                    <Power size={20} />
                </button>
            </div>

            <div className="space-y-2 mb-6">
                <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">{t('proxy.address')}</span>
                    <span className="font-mono text-foreground">{proxy.address}:{proxy.port}</span>
                </div>
                <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">{t('proxy.latency')}</span>
                    <span className={`font-mono font-medium ${proxy.latency < 50 && proxy.latency > 0 ? 'text-green-400' :
                        proxy.latency > 200 ? 'text-red-400' :
                            proxy.latency === -1 ? 'text-muted-foreground' : 'text-yellow-400'
                        }`}>
                        {proxy.latency > -1 ? `${proxy.latency}ms` : t('proxy.timeout')}
                    </span>
                </div>
                <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">{t('proxy.location')}</span>
                    <span className="text-foreground flex items-center gap-1">
                        {proxy.country}
                    </span>
                </div>
            </div>

            <div className="flex gap-2 justify-center">
                <button className="flex-1 py-3 rounded-lg bg-muted hover:bg-accent text-sm font-medium text-muted-foreground hover:text-accent-foreground transition-colors flex items-center justify-center gap-2">
                    <Edit2 size={16} /> {t('proxy.edit')}
                </button>
                <button
                    onClick={() => onDelete(proxy.id)}
                    className="flex-1 py-3 rounded-lg bg-muted hover:bg-destructive/20 text-sm font-medium text-destructive hover:text-destructive transition-colors flex items-center justify-center gap-2"
                >
                    <Trash2 size={16} /> {t('proxy.delete')}
                </button>
            </div>
        </div>
    );
};

export default ProxyCard;
