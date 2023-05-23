import Navbar from '../components/navbar';
import './globals.css';
import { Inter } from 'next/font/google';

const inter = Inter({ subsets: ['latin'] });

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" data-theme="lofi">
      <body className={inter.className}>
        <Navbar />
        {children}
      </body>
    </html>
  );
}
