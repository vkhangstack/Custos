import { LucideIcon } from 'lucide-react';
import React from 'react';

interface PageHeaderProps {
    title: string;
    icon: LucideIcon;
    iconColorClass?: string;
    description?: string;
    actions?: React.ReactNode;
}

const PageHeader = ({ title, icon: Icon, iconColorClass = "text-blue-400", description, actions }: PageHeaderProps) => {
    return (
        <div className="flex justify-between items-center mb-6">
            <div>
                <h1 className="text-2xl font-bold flex items-center gap-2">
                    <Icon className={iconColorClass} />
                    {title}
                </h1>
                {description && <p className="text-gray-400 text-sm mt-1">{description}</p>}
            </div>
            {actions && (
                <div>
                    {actions}
                </div>
            )}
        </div>
    );
};

export default PageHeader;
