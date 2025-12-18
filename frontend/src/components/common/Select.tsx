import { LucideIcon, ChevronDown } from 'lucide-react';
import React from 'react';

interface SelectOption {
    label: string;
    value: string | number;
}

interface SelectProps extends Omit<React.SelectHTMLAttributes<HTMLSelectElement>, 'onChange'> {
    label?: string;
    options: SelectOption[];
    value: string | number;
    onChange: (value: string) => void;
    icon?: LucideIcon;
    className?: string;
    description?: string;
}

const Select: React.FC<SelectProps> = ({
    label,
    options,
    value,
    onChange,
    icon: Icon,
    className = '',
    description,
    disabled,
    ...props
}) => {
    return (
        <div className={className}>
            {label && (
                <label className="block mb-2 text-sm font-medium text-muted-foreground">
                    {label}
                </label>
            )}
            <div className="relative">
                {Icon && (
                    <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-muted-foreground">
                        <Icon size={18} />
                    </div>
                )}
                <select
                    value={value}
                    onChange={(e) => onChange(e.target.value)}
                    disabled={disabled}
                    className={`
                        bg-input border border-border text-foreground text-sm rounded-lg 
                        focus:ring-primary focus:border-primary block w-full p-2.5 
                        appearance-none pr-10
                        ${Icon ? 'pl-10' : ''}
                        ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
                        transition-colors
                    `}
                    {...props}
                >
                    {options.map((option) => (
                        <option key={option.value} value={option.value}>
                            {option.label}
                        </option>
                    ))}
                </select>
                <div className="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none text-muted-foreground/50">
                    <ChevronDown size={16} />
                </div>
            </div>
            {description && (
                <p className="mt-1 text-xs text-muted-foreground">{description}</p>
            )}
        </div>
    );
};

export default Select;
