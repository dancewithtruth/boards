import 'react-toastify/dist/ReactToastify.min.css';
import Navbar from '../components/navbar';

import './globals.css';
import { UserProvider } from '@/providers/user';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" data-theme="lofi">
      <body>
        <UserProvider>
          <Navbar />
          {children}
        </UserProvider>
      </body>
    </html>
  );
}
