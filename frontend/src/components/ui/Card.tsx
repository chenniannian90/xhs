import type { HTMLAttributes } from 'react';
import { clsx } from 'clsx';

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'bordered' | 'elevated';
}

export default function Card({ variant = 'default', className, children, ...props }: CardProps) {
  const baseStyles = 'rounded-lg p-6';

  const variantStyles = {
    default: 'bg-white dark:bg-gray-800',
    bordered: 'bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700',
    elevated: 'bg-white dark:bg-gray-800 shadow-md',
  };

  return (
    <div className={clsx(baseStyles, variantStyles[variant], className)} {...props}>
      {children}
    </div>
  );
}
