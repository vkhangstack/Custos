import { LucideIcon } from 'lucide-react';

interface StatCardProps {
    title: string;
    value: string | number;
    icon: LucideIcon;
    iconColorClass: string;
    bgClass: string;
}

const StatCard = ({ title, value, icon: Icon, iconColorClass, bgClass }: StatCardProps) => {
    return (
        <div className="bg-card p-4 rounded-xl shadow-lg border border-border">
            <div className="flex items-center justify-between">
                <div>
                    <p className="text-muted-foreground text-sm">{title}</p>
                    <h3 className="text-2xl font-bold mt-1">{value}</h3>
                </div>
                <div className={`p-3 rounded-lg ${bgClass}`}>
                    <Icon className={iconColorClass} size={24} />
                </div>
            </div>
        </div>
    );
};

export default StatCard;
