import Link from 'next/link';

import Logo from './logo';
import WidthContainer from './widthContainer';
import { NAVBAR_HEIGHT } from '@/constants';
import { User } from '@/api';
import Avatar from './avatar';
import AccountMenu from './menus/account';
import { FaRegBell } from 'react-icons/fa';

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
            <div className="flex items-center">
              <div className="btn btn-ghost btn-circle">
                <FaRegBell className="cursor-pointer" size={18} />
              </div>
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
