
import { Ban, CheckCircle, AlertTriangle, Trash2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';

interface Rule {
    id: string;
    name: string;
    target: string;
    type: 'block' | 'allow' | 'redirect';
    active: boolean;
    hits: number;
    category?: 'ads' | 'security' | 'custom' | 'privacy';
    source: 'default' | 'custom';
}

interface RuleItemProps {
    rule: Rule;
    onToggle: (id: string) => void;
    onDelete?: (id: string) => void;
}

const RuleItem = ({ rule, onToggle, onDelete }: RuleItemProps) => {
    console.log("rule", rule)
    const { t } = useTranslation();
    return (
        <div className="p-4 hover:bg-muted/50 transition-colors flex items-center justify-between group">
            <div className="flex items-center gap-4">
                <button
                    onClick={() => onToggle(rule.id)}
                    className={`w-12 h-6 rounded-full relative transition-colors ${rule.active ? 'bg-primary' : 'bg-input'
                        } `}
                >
                    <span className={`absolute top-1 w-4 h-4 rounded-full bg-background transition-all shadow-sm ${rule.active ? 'left-7' : 'left-1'
                        } `} />
                </button>

                <div className={`p-2 rounded-lg ${rule.type === 'block' ? 'bg-red-500/10 text-red-500' :
                    rule.type === 'allow' ? 'bg-green-500/10 text-green-500' :
                        'bg-yellow-500/10 text-yellow-500'
                    } `}>
                    {rule.type === 'block' ? <Ban size={20} /> :
                        rule.type === 'allow' ? <CheckCircle size={20} /> :
                            <AlertTriangle size={20} />}
                </div>

                <div>
                    <h3 className={`font-semibold ${rule.active ? 'text-foreground' : 'text-muted-foreground'} `}>
                        {rule.name}
                    </h3>
                    <p className="text-muted-foreground text-sm font-mono mt-0.5">{rule.target}</p>
                </div>
            </div>

            <div className="flex items-center gap-6">
                <div className="text-right">
                    <span className="block text-xs text-muted-foreground uppercase font-semibold">{t('rules.hits')}</span>
                    <span className="font-mono text-foreground">{rule.hits.toLocaleString()}</span>
                </div>
                <span className={`px-2 py-1 rounded text-xs font-bold uppercase ${rule.source === 'default' ? 'bg-muted text-muted-foreground' :
                    rule.source === 'custom' ? 'bg-primary text-primary-foreground' :
                        'bg-muted text-muted-foreground'
                    } `}>
                    {rule.source}
                </span>

                {onDelete && (
                    <button
                        onClick={() => onDelete(rule.id)}
                        className="p-2 text-muted-foreground hover:text-red-500 hover:bg-red-500/10 rounded transition-colors"
                        title="Delete Rule"
                    >
                        <Trash2 size={18} />
                    </button>
                )}
            </div>
        </div>
    );
};

export default RuleItem;
