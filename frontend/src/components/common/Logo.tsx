import React from 'react';
import logo from '../../../../build/appicon.png';
import { APP_CONFIG } from '../../config';

interface LogoProps {
    collapsed?: boolean;
    className?: string;
    appInfo: typeof APP_CONFIG
}

const Logo: React.FC<LogoProps> = ({ collapsed = false, className = '', appInfo }) => {
    return (
        <div className={`flex items-center gap-3 ${className} ${collapsed ? 'justify-center' : ''}`}>
            <div className={`
                flex items-center justify-center rounded-lg bg-primary shadow-lg shadow-primary/20
                ${!collapsed ? 'w-10 h-10' : 'w-8 h-8'}
                transition-all duration-300 ease-in-out
            `}>
                <img src={logo} className="w-full h-full object-contain rounded-lg" alt="" />
            </div>
            {!collapsed && (
                <span className="text-xl font-bold tracking-tight text-foreground whitespace-nowrap overflow-hidden transition-all duration-300">
                    <span className="text-primary">{appInfo.appName}</span>
                </span>
            )}
        </div>
    );
};

export default Logo;
