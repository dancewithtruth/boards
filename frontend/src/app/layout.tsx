'use client';

import 'react-toastify/dist/ReactToastify.css';
import Navbar from '../components/navbar';
import { usePathname } from 'next/navigation';

import './globals.css';
import { UserProvider } from '@/providers/user';
import ConfiguredToastContainer from '@/components/toastcontainer';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const isOnBoardPage = pathname.startsWith('/boards/');

  return (
    <html lang="en" data-theme="lofi">
      <body>
        <ConfiguredToastContainer />
        <UserProvider>
          {isOnBoardPage ? null : <Navbar />}
          {children}
        </UserProvider>
      </body>
    </html>
  );
}
