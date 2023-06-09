import { cookies } from 'next/headers';
import Navbar from '@/components/navbar';
import ToastProvider from '@/providers/ToastProvider';

import './globals.css';
import 'react-toastify/dist/ReactToastify.css';
import { getUserByJwt } from '@/api';
import { COOKIE_NAME_JWT_TOKEN } from '@/constants';
import { Tooltip } from 'react-tooltip';

export const metadata = {
  title: 'Boards',
  description: 'Collaborate with your team',
};

async function getUser() {
  const cookieStore = cookies();
  const jwtToken = cookieStore.get(COOKIE_NAME_JWT_TOKEN);
  if (jwtToken) {
    try {
      const user = await getUserByJwt(jwtToken.value);
      return user;
    } catch (e) {
      return null;
    }
  }
  return null;
}

export default async function RootLayout({ children }: { children: React.ReactNode }) {
  const user = await getUser();
  return (
    <html lang="en" data-theme="lofi">
      <body>
        <div className="pt-16">
          <Navbar user={user} />
          {children}
          <ToastProvider />
          <Tooltip id="my-tooltip" style={{ zIndex: 10001 }} />
        </div>
      </body>
    </html>
  );
}
