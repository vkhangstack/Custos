import { Home, FileText, Settings, Network as NetworkIcon, ChevronLeft, ChevronRight, Globe } from 'lucide-react';
import { NavLink } from 'react-router-dom';
import { useState } from 'react';
import '../App.css'; 
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
        { icon: Home, label: t('sidebar.home'), path: '/' },
        { icon: FileText, label: t('sidebar.log'), path: '/log' },
        { icon: Settings, label: t('sidebar.config'), path: '/config' },
        { icon: NetworkIcon, label: t('sidebar.network'), path: '/network' },
    ];

    return (
        <div className={`sidebar ${isCollapsed ? 'collapsed' : ''}`}>
            <div className="sidebar-header">
                {!isCollapsed && <h2>{APP_CONFIG.appName.split(' ')[0]}</h2>}
                <button 
                    className="collapse-btn" 
                    onClick={() => setIsCollapsed(!isCollapsed)}
                    title={isCollapsed ? t('sidebar.expand') as string : t('sidebar.collapse') as string}
                >
                    {isCollapsed ? <ChevronRight size={20} /> : <ChevronLeft size={20} />}
                </button>
            </div>
            <nav className="sidebar-nav">
                {navItems.map((item) => (
                    <NavLink 
                        key={item.path} 
                        to={item.path}
                        className={({ isActive }) => 
                            `nav-item ${isActive ? 'active' : ''}`
                        }
                        title={isCollapsed ? item.label : ''}
                    >
                        <item.icon size={20} />
                        {!isCollapsed && <span>{item.label}</span>}
                    </NavLink>
                ))}
            </nav>
            <div className="sidebar-footer">
                <button className="lang-btn" onClick={toggleLanguage} title={i18n.language === 'en' ? 'Switch to Vietnamese' : 'Chuyển sang Tiếng Anh'}>
                    <Globe size={14} />
                    {!isCollapsed && <span>{i18n.language === 'en' ? 'EN' : 'VI'}</span>}
                </button>

                {!isCollapsed ? (
                    <>
                        <div className="footer-info">
                            <span className="app-name">{APP_CONFIG.appName}</span>
                            <span className="version">{APP_CONFIG.appVersion}</span>
                        </div>
                        <div className="author">{t('sidebar.author')}: {APP_CONFIG.appAuthor}</div>
                    </>
                ) : (
                    <div className="footer-info-collapsed">
                        <span className="version">{APP_CONFIG.appShortVersion}</span>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Sidebar;
