import { Info, Heart, Shield } from 'lucide-react';
import PageHeader from '../components/common/PageHeader';
import { GetAppInfo } from '../../wailsjs/go/main/App';
import { useEffect, useState } from 'react';
import { main } from '../../wailsjs/go/models';
import { useTranslation } from 'react-i18next';

export default function About() {
    const { t } = useTranslation();
    const [appInfo, setAppInfo] = useState<main.AppInfo | null>(null);

    useEffect(() => {
        GetAppInfo().then(setAppInfo);
    }, []);

    // Helper to safely display author info if it's a string or obj (depending on Go model, but TS err says string)
    const authorName = typeof appInfo?.author === 'object' ? (appInfo?.author as any).name : appInfo?.author;
    const authorContact = typeof appInfo?.contact === 'object' ? (appInfo?.contact as any).name : appInfo?.contact;

    return (
        <div className="p-6 bg-background min-h-screen text-foreground">
            <PageHeader
                title={t('about.title') as string}
                icon={Info}
                iconColorClass="text-blue-500"
            />

            <div className="max-w-2xl mx-auto">
                <div className="bg-card border border-border rounded-xl p-8 text-center shadow-lg relative overflow-hidden">
                    {/* Decorative Background Blur */}
                    <div className="absolute top-0 left-0 w-full h-32 bg-gradient-to-b from-primary/10 to-transparent pointer-events-none" />

                    <div className="relative z-10 flex flex-col items-center">
                        <div className="w-24 h-24 bg-background border-4 border-primary/20 rounded-3xl flex items-center justify-center mb-6 shadow-xl">
                            <Shield size={48} className="text-primary" />
                        </div>

                        <h1 className="text-3xl font-bold mb-2">{appInfo?.name || 'Custos'}</h1>
                        <div className="px-3 py-1 bg-secondary text-secondary-foreground rounded-full text-sm font-mono mb-6">
                            {t('about.version')} {appInfo?.version || '0.0.0'}
                        </div>

                        <p className="text-lg text-muted-foreground mb-8 max-w-md mx-auto leading-relaxed">
                            {appInfo?.description || t('about.description')}
                        </p>

                        <div className="flex justify-center gap-4 w-full max-w-sm mb-8">
                            <div className="bg-muted/50 p-4 rounded-lg">
                                <p className="text-xs text-muted-foreground uppercase font-bold mb-1">{t('about.author')}</p>
                                <p className="font-medium truncate">{authorName || 'Unknown'}</p>
                            </div>
                            <div className="bg-muted/50 p-4 rounded-lg">
                                <p className="text-xs text-muted-foreground uppercase font-bold mb-1">{t('about.contact')}</p>
                                <a href={authorContact} className="font-medium truncate">{authorContact}</a>
                            </div>
                        </div>

                        <div className="pt-6 border-t border-border w-full">
                            <p className="text-sm text-muted-foreground flex items-center justify-center gap-1">
                                {t('about.madeWith')} <Heart size={14} className="text-red-500 fill-red-500" /> {t('about.by')} {authorName}
                            </p>
                            <p className="text-xs text-muted-foreground mt-2">
                                &copy; {new Date().getFullYear()} {t('about.rights')}
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
