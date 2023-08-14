import { ReactNode } from 'react';

export default function WidthContainer({ children, className }: { children: ReactNode; className?: string }) {
  return <div className={`max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 w-full ${className}`}>{children}</div>;
}
