import Link from 'next/link';
import { FaChevronDown } from 'react-icons/fa';

import Logo from './logo';
import WidthContainer from './widthContainer';
import { COOKIE_NAME_JWT_TOKEN, NAVBAR_HEIGHT } from '@/constants';
import { User } from '@/api';
import Avatar from './avatar';
import Cookies from 'universal-cookie';
import { useRouter } from 'next/navigation';
import AccountMenu from './menus/account';

export default function Navbar({ user }: { user: User | null }) {
  return (
    <div
      className="navbar fixed top-0 left-0 w-full bg-white shadow-md"
      style={{ height: NAVBAR_HEIGHT, zIndex: 10002 }}
    >
      <WidthContainer>
        <div className="flex justify-between items-center w-full">
          <Logo className="font-bold text-xl" />
          {user ? (
            <div className="flex items-center space-x-4">
              <Link href="/dashboard" className="btn btn-primary btn-sm">
                Dashboard
              </Link>
              <AccountMenu user={user} avatar={<Avatar id={user.id} />} />
            </div>
          ) : (
            <AuthNav />
          )}
        </div>
      </WidthContainer>
    </div>
  );
}

function AuthNav() {
  return (
    <div className="space-x-2">
      <Link href="/auth/signin" className="btn btn-secondary btn-outline">
        Sign in
      </Link>
      <Link href="/auth/signup" className="btn btn-primary">
        Sign up
      </Link>
    </div>
  );
}
