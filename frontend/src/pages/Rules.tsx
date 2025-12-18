import { useState, useEffect } from 'react';
import { Shield, Plus, Search, Filter, X } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import PageHeader from '../components/common/PageHeader';
import RuleItem from '../components/rules/RuleItem';
import { AddRule, GetRules, GetRulesPaginated, DeleteRule, ToggleRule } from '../../wailsjs/go/main/App';
import { core } from '../../wailsjs/go/models';

// UI Rule interface matching backend core.Rule but with UI specifics if needed
// Actually we can just use core.Rule but RuleItem expects specific props.
// Let's adapt or change RuleItem.
// RuleItem expects: id, name, target, type, active, hits, category.
// core.Rule has: ID, Type, Pattern, Enabled.
// We miss: hits, category, name (Pattern covers target/name).
// We can mock hits/category for now or remove them from RuleItem.
// Let's modify RuleItem to be simpler/match backend, or adapt here.
interface DisplayRule extends core.Rule {
    hits: number; // Mocked
    category: 'custom' | 'ads' | 'security'; // Mocked
    name: string; // derived from Pattern
}

export default function Rules() {
    const { t } = useTranslation();
    const [rules, setRules] = useState<DisplayRule[]>([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [isModalOpen, setIsModalOpen] = useState(false);

    // New Rule Form State
    const [newPattern, setNewPattern] = useState('');
    const [newType, setNewType] = useState('BLOCK'); // core.RuleBlock

    // Pagination State
    const [currentPage, setCurrentPage] = useState(1);
    const [totalItems, setTotalItems] = useState(0);
    const pageSize = 50;

    const fetchRules = async () => {
        try {
            // Use server-side pagination
            const fetched = await GetRulesPaginated(currentPage, pageSize, searchTerm);
            if (fetched && fetched.rules) {
                const displayRules = fetched.rules.map((r: core.Rule) => ({
                    ...r,
                    active: r.enabled, // map enabled -> active for RuleItem
                    target: r.pattern, // map pattern -> target for RuleItem
                    name: r.pattern,   // use pattern as name
                    type: r.type === 'BLOCK' ? 'block' : 'allow', // map enum
                    hits: r.hit_count, // Use real hit count from backend
                    category: r.source === 'default' ? 'ads' : 'custom'
                } as any));
                setRules(displayRules);
                setTotalItems(fetched.total);
            } else {
                setRules([]);
                setTotalItems(0);
            }
        } catch (e) {
            console.error(e);
        }
    };

    useEffect(() => {
        fetchRules();

        // Implement 2s auto-refresh
        const interval = setInterval(() => {
            fetchRules();
        }, 2000);

        return () => clearInterval(interval);
    }, [currentPage, searchTerm]); // Refetch when page or search changes

    const handleToggleRule = async (id: string) => {
        const rule = rules.find(r => r.id === id);
        if (rule) {
            await ToggleRule(id, !rule.enabled);
            fetchRules();
        }
    };

    const handleDeleteRule = async (id: string) => {
        await DeleteRule(id);
        fetchRules();
    }

    const handleAddRule = async () => {
        if (!newPattern) return;
        await AddRule(newPattern, newType);
        setNewPattern('');
        setIsModalOpen(false);
        fetchRules();
    };

    // Server-side pagination handles filtering and slicing
    const filteredRules = rules;
    const totalPages = Math.ceil(totalItems / pageSize);
    const paginatedRules = rules; // rules are already sliced by server

    const actions = (
        <button
            onClick={() => setIsModalOpen(true)}
            className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors font-medium shadow-lg shadow-blue-900/20"
        >
            <Plus size={18} />
            {t('rules.newRule')}
        </button>
    );

    return (
        <div className="p-6 bg-background min-h-screen text-foreground relative">
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
            <div className="bg-card rounded-xl border border-border overflow-hidden shadow-lg flex flex-col">
                <div className="divide-y divide-border flex-1">
                    {paginatedRules.map((rule) => (
                        <RuleItem
                            key={rule.id}
                            rule={rule as any} // Pass adapted props
                            onToggle={handleToggleRule}
                            onDelete={() => handleDeleteRule(rule.id)} // Need to update RuleItem to accept onDelete
                        />
                    ))}

                    {filteredRules.length === 0 && (
                        <div className="p-12 text-center text-muted-foreground">
                            <Filter size={48} className="mx-auto mb-4 opacity-20" />
                            <p>{t('rules.noRulesFound')} "{searchTerm}"</p>
                        </div>
                    )}
                </div>

                {totalItems > pageSize && (
                    <div className="bg-muted/30 p-4 border-t border-border flex items-center justify-between">
                        <div className="text-sm text-muted-foreground">
                            Showing <span className="font-medium">{(currentPage - 1) * pageSize + 1}</span> to <span className="font-medium">{Math.min(currentPage * pageSize, totalItems)}</span> of <span className="font-medium">{totalItems}</span> rules
                        </div>
                        <div className="flex gap-2">
                            <button
                                onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                                disabled={currentPage === 1}
                                className="px-3 py-1.5 rounded-lg border border-border bg-background hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-sm font-medium"
                            >
                                Previous
                            </button>
                            <button
                                onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                                disabled={currentPage === totalPages}
                                className="px-3 py-1.5 rounded-lg border border-border bg-background hover:bg-muted disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-sm font-medium"
                            >
                                Next
                            </button>
                        </div>
                    </div>
                )}
            </div>

            {/* Simple Add Rule Modal */}
            {isModalOpen && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
                    <div className="bg-card w-full max-w-md rounded-xl shadow-2xl border border-border p-6">
                        <div className="flex justify-between items-center mb-6">
                            <h3 className="text-xl font-bold">{t('rules.newRule')}</h3>
                            <button onClick={() => setIsModalOpen(false)} className="text-muted-foreground hover:text-foreground">
                                <X size={24} />
                            </button>
                        </div>

                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium mb-1">Domain Pattern</label>
                                <input
                                    autoFocus
                                    className="w-full bg-input border border-border rounded-lg p-2"
                                    placeholder="*.example.com"
                                    value={newPattern}
                                    onChange={e => setNewPattern(e.target.value)}
                                />
                                <p className="text-xs text-muted-foreground mt-1">Use * for wildcards (e.g. *.ads.com)</p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium mb-1">Type</label>
                                <div className="flex gap-4">
                                    <button
                                        className={`flex-1 py-2 rounded-lg border ${newType === 'BLOCK' ? 'bg-red-500/20 border-red-500 text-red-500' : 'border-border opacity-50'}`}
                                        onClick={() => setNewType('BLOCK')}
                                    >
                                        Block
                                    </button>
                                    <button
                                        className={`flex-1 py-2 rounded-lg border ${newType === 'ALLOW' ? 'bg-green-500/20 border-green-500 text-green-500' : 'border-border opacity-50'}`}
                                        onClick={() => setNewType('ALLOW')}
                                    >
                                        Allow
                                    </button>
                                </div>
                            </div>

                            <div className="pt-4 flex justify-end gap-2">
                                <button onClick={() => setIsModalOpen(false)} className="px-4 py-2 rounded-lg text-muted-foreground hover:bg-accent">Cancel</button>
                                <button onClick={handleAddRule} className="px-4 py-2 rounded-lg bg-primary text-primary-foreground font-medium">Add Rule</button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
