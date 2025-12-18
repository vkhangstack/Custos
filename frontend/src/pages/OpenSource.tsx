import { Code, ExternalLink, Github, Library, Heart } from 'lucide-react';
import PageHeader from '../components/common/PageHeader';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import { GetAppInfo } from '../../wailsjs/go/main/App';
import { useEffect, useState } from 'react';
import { main } from '../../wailsjs/go/models';

export default function OpenSource() {
    const [appInfo, setAppInfo] = useState<main.AppInfo | null>(null);

    useEffect(() => {
        GetAppInfo().then(setAppInfo);
    }, []);

    const libraries = [
        { name: 'Wails', license: 'MIT', url: 'https://wails.io' },
        { name: 'React', license: 'MIT', url: 'https://reactjs.org' },
        { name: 'Vite', license: 'MIT', url: 'https://vitejs.dev' },
        { name: 'Lucide React', license: 'ISC', url: 'https://lucide.dev' },
        { name: 'Recharts', license: 'MIT', url: 'https://recharts.org' },
        { name: 'Tailwind CSS', license: 'MIT', url: 'https://tailwindcss.com' },
        { name: 'i18next', license: 'MIT', url: 'https://www.i18next.com' },
        // Backend Libraries
        { name: 'go-socks5', license: 'MIT', url: 'https://github.com/armon/go-socks5' },
        { name: 'gopsutil', license: 'BSD', url: 'https://github.com/shirou/gopsutil' },
        { name: 'gorm', license: 'MIT', url: 'https://gorm.io' },
        { name: 'glebarez/sqlite', license: 'BSD', url: 'https://github.com/glebarez/sqlite' },
        { name: 'miekg/dns', license: 'BSD', url: 'https://github.com/miekg/dns' },
    ];

    const openUrl = (url: string) => {
        BrowserOpenURL(url);
    };

    return (
        <div className="p-6 bg-background min-h-screen text-foreground">
            <PageHeader
                title="Open Source"
                icon={Code}
                iconColorClass="text-purple-500"
            />

            <div className="max-w-4xl mx-auto space-y-6">
                {/* App Info Card */}
                <div className="bg-card border border-border rounded-xl p-8 text-center shadow-lg">
                    <div className="w-20 h-20 bg-primary/10 rounded-2xl flex items-center justify-center mx-auto mb-6">
                        <Heart size={40} className="text-primary animate-pulse" />
                    </div>
                    <h2 className="text-3xl font-bold mb-2">{appInfo?.name || 'Custos'}</h2>
                    <p className="text-muted-foreground mb-6">Version {appInfo?.version || '0.0.0'}</p>
                    
                    <div className="flex justify-center gap-4">
                        <button
                            onClick={() => openUrl('https://github.com/vkhangstack/Custos')}
                            className="flex items-center gap-2 px-6 py-2.5 bg-foreground text-background rounded-full font-medium hover:opacity-90 transition-opacity"
                        >
                            <Github size={20} />
                            GitHub Repository
                        </button>
                    </div>
                </div>

                {/* Libraries */}
                <div className="bg-card border border-border rounded-xl overflow-hidden shadow-lg">
                    <div className="p-6 border-b border-border bg-muted/30">
                        <h3 className="font-semibold text-lg flex items-center gap-2">
                            <Library size={20} className="text-blue-500" />
                            Open Source Libraries
                        </h3>
                    </div>
                    <div className="divide-y divide-border">
                        {libraries.map((lib) => (
                            <div key={lib.name} className="p-4 flex items-center justify-between hover:bg-muted/50 transition-colors">
                                <div>
                                    <h4 className="font-medium text-foreground">{lib.name}</h4>
                                    <span className="text-xs text-muted-foreground bg-secondary px-2 py-0.5 rounded-full mt-1 inline-block">
                                        {lib.license} License
                                    </span>
                                </div>
                                <button
                                    onClick={() => openUrl(lib.url)}
                                    className="p-2 text-muted-foreground hover:text-primary transition-colors"
                                    title="View Website"
                                >
                                    <ExternalLink size={18} />
                                </button>
                            </div>
                        ))}
                    </div>
                </div>

                <div className="text-center text-sm text-muted-foreground pt-8 pb-4">
                    Made with ❤️ by vkhangstack
                </div>
            </div>
        </div>
    );
}
