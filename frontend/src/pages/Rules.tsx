import { useState } from 'react';
import { Shield, Plus, Search, Filter } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import PageHeader from '../components/common/PageHeader';
import RuleItem from '../components/rules/RuleItem';

interface Rule {
    id: string;
    name: string;
    target: string;
    type: 'block' | 'allow' | 'redirect';
    active: boolean;
    hits: number;
    category: 'ads' | 'security' | 'custom' | 'privacy';
}

const initialRules: Rule[] = [
    { id: '1', name: 'Block Ads (General)', target: '*.doubleclick.net', type: 'block', active: true, hits: 1452, category: 'ads' },
    { id: '2', name: 'Allow Google APIs', target: '*.googleapis.com', type: 'allow', active: true, hits: 890, category: 'custom' },
    { id: '3', name: 'Block Facebook Tracking', target: 'connect.facebook.net', type: 'block', active: true, hits: 320, category: 'privacy' },
    { id: '4', name: 'Redirect Malware Domain', target: 'malicious-site.com', type: 'redirect', active: true, hits: 5, category: 'security' },
    { id: '5', name: 'Allow Local Network', target: '192.168.1.*', type: 'allow', active: false, hits: 0, category: 'custom' },
];

export default function Rules() {
    const { t } = useTranslation();
    const [rules, setRules] = useState<Rule[]>(initialRules);
    const [searchTerm, setSearchTerm] = useState('');

    const toggleRule = (id: string) => {
        setRules(rules.map(r => r.id === id ? { ...r, active: !r.active } : r));
    };

    const filteredRules = rules.filter(r =>
        r.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        r.target.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const actions = (
        <button className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors font-medium shadow-lg shadow-blue-900/20">
            <Plus size={18} />
            {t('rules.newRule')}
        </button>
    );

    return (
        <div className="p-6 bg-background min-h-screen text-foreground">
            <PageHeader
                title={t('rules.title')}
                icon={Shield}
                iconColorClass="text-green-400"
                description={t('rules.description') || ''}
                actions={actions}
            />

            {/* Search & Filter Bar */}
            <div className="flex gap-4 mb-6">
                <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={18} />
                    <input
                        type="text"
                        placeholder={t('rules.searchPlaceholder') || ''}
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full bg-input border border-border rounded-lg py-2.5 pl-10 pr-4 text-foreground focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary transition-all font-medium placeholder:text-muted-foreground"
                    />
                </div>
                <button className="bg-secondary border border-border text-secondary-foreground px-4 rounded-lg hover:bg-secondary/80 transition-colors flex items-center gap-2 font-medium">
                    <Filter size={18} />
                    {t('rules.filter')}
                </button>
            </div>

            {/* Rules List */}
            <div className="bg-card rounded-xl border border-border overflow-hidden shadow-lg">
                <div className="divide-y divide-border">
                    {filteredRules.map((rule) => (
                        <RuleItem
                            key={rule.id}
                            rule={rule}
                            onToggle={toggleRule}
                        />
                    ))}

                    {filteredRules.length === 0 && (
                        <div className="p-12 text-center text-muted-foreground">
                            <Filter size={48} className="mx-auto mb-4 opacity-20" />
                            <p>{t('rules.noRulesFound')} "{searchTerm}"</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
