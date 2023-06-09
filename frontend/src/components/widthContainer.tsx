import { ReactNode } from 'react';

export default function WidthContainer({ children }: { children: ReactNode }) {
  return <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 w-full">{children}</div>;
}
