import { useTranslation } from 'react-i18next';

interface UsageItem {
    domain: string;
    usage: string;
    process: string;
}

interface TopUsageListProps {
    items: UsageItem[];
}

const TopUsageList = ({ items }: TopUsageListProps) => {
    const { t } = useTranslation();
    return (
        <div className="bg-card p-6 rounded-xl shadow-lg border border-border">
            <h3 className="text-lg font-semibold mb-4">{t('dashboard.topUsage')}</h3>
            <div className="space-y-4">
                {items.map((item, index) => (
                    <div key={index} className="flex items-center justify-between p-3 bg-muted/50 rounded-lg">
                        <div className="flex flex-col">
                            <span className="font-medium text-foreground">{item.domain}</span>
                            <span className="text-xs text-muted-foreground">{item.process}</span>
                        </div>
                        <span className="text-sm font-bold text-primary">{item.usage}</span>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default TopUsageList;
