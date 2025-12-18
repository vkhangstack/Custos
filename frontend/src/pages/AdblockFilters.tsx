import { useState, useEffect } from 'react';
import { Shield, Plus, Search, Filter, X, RefreshCw, Trash2, ExternalLink } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import PageHeader from '../components/common/PageHeader';
import { GetAdblockFilters, AddAdblockFilter, DeleteAdblockFilter, ToggleAdblockFilter, RefreshAdblockFilters } from '../../wailsjs/go/main/App';
import { core } from '../../wailsjs/go/models';
import { useToast } from '../context/ToastContext';

export default function AdblockFilters() {
    const { t } = useTranslation();
    const { showToast } = useToast();
    const [filters, setFilters] = useState<core.AdblockFilter[]>([]);
    const [searchTerm, setSearchTerm] = useState('');
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isRefreshing, setIsRefreshing] = useState(false);

    // New Filter Form State
    const [newName, setNewName] = useState('');
    const [newURL, setNewURL] = useState('');

    const fetchFilters = async () => {
        try {
            const fetched = await GetAdblockFilters();
            setFilters(fetched || []);
        } catch (e) {
            console.error(e);
        }
    };

    useEffect(() => {
        fetchFilters();
    }, []);

    const handleToggleFilter = async (id: string, enabled: boolean) => {
        try {
            await ToggleAdblockFilter(id, enabled);
            fetchFilters();
        } catch (e) {
            console.error(e);
        }
    };

    const handleDeleteFilter = async (id: string) => {
        if (!confirm('Are you sure you want to delete this filter?')) return;
        try {
            await DeleteAdblockFilter(id);
            fetchFilters();
        } catch (e) {
            console.error(e);
        }
    };

    const handleAddFilter = async () => {
        if (!newName || !newURL) return;
        try {
            await AddAdblockFilter(newName, newURL);
            setNewName('');
            setNewURL('');
            setIsModalOpen(false);
            fetchFilters();
        } catch (e) {
            console.error(e);
        }
    };

    const handleRefresh = async () => {
        setIsRefreshing(true);
        try {
            await RefreshAdblockFilters();
            showToast("Filters refreshed and reloaded", "success");
            fetchFilters();
        } catch (e) {
            console.error(e);
            showToast("Failed to refresh filters", "error");
        } finally {
            setIsRefreshing(false);
        }
    };

    const filteredFilters = filters.filter(f =>
        f.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
        f.url.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const actions = (
        <div className="flex gap-2">
            <button
                onClick={handleRefresh}
                disabled={isRefreshing}
                className={`flex items-center gap-2 bg-secondary hover:bg-secondary/80 text-secondary-foreground px-4 py-2 rounded-lg transition-colors font-medium ${isRefreshing ? 'opacity-50 cursor-not-allowed' : ''}`}
            >
                <RefreshCw size={18} className={isRefreshing ? 'animate-spin' : ''} />
                Refresh
            </button>
            <button
                onClick={() => setIsModalOpen(true)}
                className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors font-medium shadow-lg shadow-blue-900/20"
            >
                <Plus size={18} />
                New List
            </button>
        </div>
    );

    return (
        <div className="p-6 bg-background min-h-screen text-foreground relative">
            <PageHeader
                title="Adblock Filters"
                icon={Shield}
                iconColorClass="text-blue-400"
                description="Manage adblock filter lists and rules."
                actions={actions}
            />

            {/* Search Bar */}
            <div className="flex gap-4 mb-6">
                <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" size={18} />
                    <input
                        type="text"
                        placeholder="Search filters..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full bg-input border border-border rounded-lg py-2.5 pl-10 pr-4 text-foreground focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary transition-all font-medium placeholder:text-muted-foreground"
                    />
                </div>
            </div>

            {/* Filters List */}
            <div className="bg-card rounded-xl border border-border overflow-hidden shadow-lg flex flex-col">
                <div className="divide-y divide-border flex-1">
                    {filteredFilters.map((filter) => (
                        <div key={filter.id} className="p-4 flex items-center justify-between hover:bg-muted/50 transition-colors">
                            <div className="flex items-center gap-4">
                                <div className={`w-10 h-10 rounded-lg flex items-center justify-center ${filter.enabled ? 'bg-blue-500/10 text-blue-500' : 'bg-muted text-muted-foreground'}`}>
                                    <Shield size={20} />
                                </div>
                                <div>
                                    <div className="font-bold flex items-center gap-2">
                                        {filter.name}
                                        {filter.url && (
                                            <a href={filter.url} target="_blank" rel="noopener noreferrer" className="text-muted-foreground hover:text-primary">
                                                <ExternalLink size={14} />
                                            </a>
                                        )}
                                    </div>
                                    <div className="text-xs text-muted-foreground truncate max-w-md">{filter.url}</div>
                                    <div className="text-[10px] text-muted-foreground mt-1">
                                        Last Updated: {filter.last_updated ? new Date(filter.last_updated).toLocaleString() : 'Never'}
                                    </div>
                                </div>
                            </div>
                            <div className="flex items-center gap-4">
                                {/* <div className="text-right mr-4">
                                    <div className="text-xs text-muted-foreground">Hits</div>
                                    <div className="font-mono font-bold text-blue-400">{filter.hits.toLocaleString()}</div>
                                </div> */}
                                <button
                                    onClick={() => ToggleAdblockFilter(filter.id, !filter.enabled).then(fetchFilters)}
                                    className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${filter.enabled ? 'bg-blue-600' : 'bg-muted'}`}
                                >
                                    <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${filter.enabled ? 'translate-x-6' : 'translate-x-1'}`} />
                                </button>
                                <button
                                    onClick={() => handleDeleteFilter(filter.id)}
                                    className="p-2 text-muted-foreground hover:text-red-500 hover:bg-red-500/10 rounded-lg transition-colors"
                                >
                                    <Trash2 size={18} />
                                </button>
                            </div>
                        </div>
                    ))}

                    {filteredFilters.length === 0 && (
                        <div className="p-12 text-center text-muted-foreground">
                            <Filter size={48} className="mx-auto mb-4 opacity-20" />
                            <p>No filters found.</p>
                        </div>
                    )}
                </div>
            </div>

            {/* Add Filter Modal */}
            {isModalOpen && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
                    <div className="bg-card w-full max-w-md rounded-xl shadow-2xl border border-border p-6">
                        <div className="flex justify-between items-center mb-6">
                            <h3 className="text-xl font-bold">New Filter List</h3>
                            <button onClick={() => setIsModalOpen(false)} className="text-muted-foreground hover:text-foreground">
                                <X size={24} />
                            </button>
                        </div>

                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium mb-1">Name</label>
                                <input
                                    autoFocus
                                    className="w-full bg-input border border-border rounded-lg p-2"
                                    placeholder="e.g. EasyList"
                                    value={newName}
                                    onChange={e => setNewName(e.target.value)}
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium mb-1">URL</label>
                                <input
                                    className="w-full bg-input border border-border rounded-lg p-2"
                                    placeholder="https://..."
                                    value={newURL}
                                    onChange={e => setNewURL(e.target.value)}
                                />
                            </div>

                            <div className="pt-4 flex justify-end gap-2">
                                <button onClick={() => setIsModalOpen(false)} className="px-4 py-2 rounded-lg text-muted-foreground hover:bg-accent">Cancel</button>
                                <button onClick={handleAddFilter} className="px-4 py-2 rounded-lg bg-primary text-primary-foreground font-medium">Add Filter</button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
