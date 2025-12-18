import { useState, useEffect } from 'react';
import { ToggleLeft, ToggleRight, Save, Settings as SettingsIcon, Shield, Sliders } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import PageHeader from '../components/common/PageHeader';
import Select from '../components/common/Select';
import { useTheme } from '../context/ThemeContext';
import { useToast } from '../context/ToastContext';
import { SetRunOnStartup, GetAppSettings, SaveAppSettings } from '../../wailsjs/go/main/App';

export default function Settings() {
    const { t, i18n } = useTranslation();
    const { theme, setTheme } = useTheme();
    const { showToast } = useToast();
    const [autoStart, setAutoStart] = useState(false);
    const [notifications, setNotifications] = useState(false);
    const [port, setPort] = useState(1080);
    const [language, setLanguage] = useState(i18n.language);

    // Sync language state with i18n on mount/change
    useEffect(() => {
        setLanguage(i18n.language);
    }, [i18n.language]);

    // Initialize Settings
    useEffect(() => {
        GetAppSettings().then(settings => {
            setPort(settings.port);
            setNotifications(settings.notifications);
            setAutoStart(settings.auto_start);
        });
    }, []);

    const handleAutoStartToggle = async () => {
        const newState = !autoStart;
        try {
            await SetRunOnStartup(newState);
            setAutoStart(newState);
            showToast(
                t('settings.general.autoStart') + (newState ? ' Enabled' : ' Disabled'),
                'success'
            );
        } catch (error) {
            console.error("Failed to toggle auto-start:", error);
            showToast("Failed to toggle auto-start", 'error');
        }
    };

    const handleLanguageChange = (val: string) => {
        setLanguage(val);
        i18n.changeLanguage(val);
    };

    const toggleTheme = () => {
        setTheme(theme === 'dark' ? 'light' : 'dark');
    };

    const isDark = theme === 'dark';

    const handleSave = async () => {
        try {
            await SaveAppSettings({
                port: port,
                notifications: notifications,
                auto_start: autoStart,
                adblock_enabled: true
            });
            showToast(t('settings.save') + ' Success', 'success');
        } catch (error) {
            console.error("Failed to save settings:", error);
            showToast("Failed to save settings", 'error');
        }
    };

    const actions = (
        <button
            onClick={handleSave}
            className="flex items-center gap-2 bg-primary hover:bg-primary/90 text-primary-foreground px-4 py-2 rounded-lg transition-colors font-medium shadow-lg shadow-primary/20"
        >
            <Save size={18} />
            {t('settings.save')}
        </button>
    );

    return (
        <div className="p-6 bg-background min-h-screen text-foreground">
            <PageHeader
                title={t('settings.title')}
                icon={SettingsIcon}
                iconColorClass="text-muted-foreground"
                actions={null}
            />

            <div className="space-y-6 max-w-4xl mx-auto">
                {/* General Settings */}
                <div className="bg-card rounded-xl border border-border p-6 shadow-lg">
                    <h2 className="text-lg font-bold mb-4 flex items-center gap-2 text-card-foreground">
                        <Sliders size={20} className="text-blue-500" />
                        {t('settings.general.title')}
                    </h2>
                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <span className="text-muted-foreground font-medium">{t('settings.general.autoStart')}</span>
                            <button onClick={handleAutoStartToggle} className="text-primary hover:text-primary/80 transition-colors">
                                {autoStart ? <ToggleRight size={32} /> : <ToggleLeft size={32} className="text-muted-foreground" />}
                            </button>
                        </div>
                        <div className="flex items-center justify-between">
                            <span className="text-muted-foreground font-medium">{t('settings.general.notifications')}</span>
                            <button disabled={true} onClick={() => setNotifications(!notifications)} className="text-primary hover:text-primary/80 transition-colors">
                                {notifications ? <ToggleRight size={32} /> : <ToggleLeft size={32} className="text-muted-foreground" />}
                            </button>
                        </div>
                        <div className="flex items-center justify-between">
                            <span className="text-muted-foreground font-medium">{t('settings.general.language')}</span>
                            <Select
                                value={language}
                                onChange={(val) => {
                                    handleLanguageChange(val);
                                }}
                                options={[
                                    { label: 'English', value: 'en' },
                                    { label: 'Tiếng Việt', value: 'vi' },
                                ]}
                                className="min-w-[150px]"
                            />
                        </div>
                    </div>
                </div>

                {/* Network & Security */}
                <div className="bg-card rounded-xl border border-border p-6 shadow-lg">
                    <h2 className="text-lg font-bold mb-4 flex items-center gap-2 text-card-foreground">
                        <Shield size={20} className="text-green-500" />
                        {t('settings.network.title')}
                    </h2>
                    <div className="space-y-6">
                        <div>
                            <label className="block mb-2 text-sm font-medium text-muted-foreground">{t('settings.network.port')}</label>
                            <input
                                type="number"
                                value={port}
                                disabled={true}
                                onChange={(e) => setPort(parseInt(e.target.value))}
                                className="bg-input border border-border text-foreground text-sm rounded-lg focus:ring-primary focus:border-primary block w-full p-2.5 font-mono"
                            />
                        </div>
                        <Select
                            disabled={true}
                            label={t('settings.network.protocol') as string}
                            value="socks5"
                            onChange={() => { }} // Placeholder for actual logic if needed
                            options={[
                                { label: 'SOCKS5', value: 'socks5' },
                                { label: 'HTTP', value: 'http' },
                                { label: 'HTTPS', value: 'https' },
                            ]}
                        />
                    </div>
                </div>

                {/* Appearance */}
                <div className="bg-card rounded-xl border border-border p-6 shadow-lg">
                    <h2 className="text-lg font-bold mb-4 flex items-center gap-2 text-card-foreground">
                        <SettingsIcon size={20} className="text-purple-500" />
                        {t('settings.appearance.title')}
                    </h2>
                    <div className="flex items-center justify-between">
                        <span className="text-muted-foreground font-medium">{t('settings.general.darkMode')}</span>
                        <button onClick={toggleTheme} className="text-primary hover:text-primary/80 transition-colors">
                            {isDark ? <ToggleRight size={32} /> : <ToggleLeft size={32} className="text-muted-foreground" />}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
