import { Home, FileText, Settings, Network as NetworkIcon, ChevronLeft, ChevronRight, Globe } from 'lucide-react';
import { NavLink } from 'react-router-dom';
import { useState } from 'react';
import { APP_CONFIG } from '../config';
import { useTranslation } from 'react-i18next';

const Sidebar = () => {
    const [isCollapsed, setIsCollapsed] = useState(false);
    const { t, i18n } = useTranslation();

    const toggleLanguage = () => {
        const newLang = i18n.language === 'en' ? 'vi' : 'en';
        i18n.changeLanguage(newLang);
    };

    const navItems = [
        { icon: Home, label: t('sidebar.dashboard'), path: '/' },
        { icon: NetworkIcon, label: t('sidebar.traffic'), path: '/traffic' },
        { icon: Globe, label: t('sidebar.proxy'), path: '/proxy' },
        { icon: FileText, label: t('sidebar.rules'), path: '/rules' },
        { icon: FileText, label: t('sidebar.reports'), path: '/reports' },
        { icon: Settings, label: t('sidebar.settings'), path: '/settings' },
    ];

    return (
        <div className={`flex flex-col h-screen bg-card border-r border-border text-card-foreground shadow-xl transition-[width] duration-300 ease-in-out shrink-0 ${isCollapsed ? 'w-20' : 'w-64'}`}>
            {/* Header */}
            <div className={`h-16 flex items-center justify-between px-6 border-b border-border ${isCollapsed ? 'justify-center p-0' : ''}`}>
                {!isCollapsed && <h2 className="text-xl font-semibold whitespace-nowrap overflow-hidden text-primary">{APP_CONFIG.appName.split(' ')[0]}</h2>}
                <button
                    className="p-1 rounded text-muted-foreground hover:text-foreground hover:bg-accent transition-colors"
                    onClick={() => setIsCollapsed(!isCollapsed)}
                    title={isCollapsed ? t('sidebar.expand') as string : t('sidebar.collapse') as string}
                >
                    {isCollapsed ? <ChevronRight size={20} /> : <ChevronLeft size={20} />}
                </button>
            </div>

            {/* Navigation */}
            <nav className="flex-1 py-6 overflow-y-auto">
                {navItems.map((item) => (
                    <NavLink
                        key={item.path}
                        to={item.path}
                        className={({ isActive }) =>
                            `flex items-center gap-3 px-6 py-3 text-muted-foreground transition-all duration-200 hover:bg-accent hover:text-accent-foreground whitespace-nowrap overflow-hidden ${isActive ? '!bg-primary !text-primary-foreground border-r-4 border-primary-foreground/30' : ''
                            } ${isCollapsed ? 'justify-center px-0' : ''}`
                        }
                        title={isCollapsed ? item.label : ''}
                    >
                        <item.icon size={20} className="shrink-0" />
                        {!isCollapsed && <span>{item.label}</span>}
                    </NavLink>
                ))}
            </nav>

            {/* Footer */}
            <div className={`p-4 border-t border-border text-sm text-muted-foreground bg-muted/20 ${isCollapsed ? 'text-center p-4' : ''}`}>
                <button
                    className={`flex items-center justify-center gap-2 w-20 p-1 mb-2 rounded border border-border text-muted-foreground hover:bg-accent hover:text-accent-foreground transition-colors text-xs ${isCollapsed ? 'w-auto mx-auto mb-3 border-none' : ''}`}
                    onClick={toggleLanguage}
                    title={i18n.language === 'en' ? 'Switch to Vietnamese' : 'Chuyển sang Tiếng Anh'}
                >
                    <Globe size={14} />
                    {!isCollapsed && <span>{i18n.language === 'en' ? 'EN' : 'VI'}</span>}
                </button>

                {!isCollapsed ? (
                    <>
                        <div className="flex justify-between items-center mb-1">
                            <span className="font-semibold text-foreground">{APP_CONFIG.appName}</span>
                            <span className="bg-accent px-1.5 py-0.5 rounded text-xs text-accent-foreground">{APP_CONFIG.appVersion}</span>
                        </div>
                        <div className="text-xs">{t('sidebar.author')}: {APP_CONFIG.appAuthor}</div>
                    </>
                ) : (
                    <div className="text-xs">
                        <span className="block mb-1">{APP_CONFIG.appShortVersion}</span>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Sidebar;
