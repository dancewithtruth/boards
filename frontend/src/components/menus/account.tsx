'use client';
import { FaChevronDown } from 'react-icons/fa';

import { COOKIE_NAME_JWT_TOKEN } from '@/constants';
import { User } from '@/api';
import Cookies from 'universal-cookie';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function AccountMenu({ user, avatar }: { user: User; avatar: React.ReactNode }) {
  const cookies = new Cookies();
  const router = useRouter();

  const handleLogout = () => {
    cookies.remove(COOKIE_NAME_JWT_TOKEN);
    router.refresh();
    router.push('/');
  };

  return (
    <div className="dropdown dropdown-end">
      <div tabIndex={0} className="btn btn-ghost normal-case">
        <div className="w-10">{avatar}</div>
        <span>{user?.name}</span>
        <FaChevronDown />
      </div>
      <div className="right-0 mt-3 p-2 shadow menu menu-compact dropdown-content bg-base-100 rounded-box w-52">
        <ul className="menu menu-compact gap-1 p-3">
          <li>
            <Link href="/dashboard" className="flex items-center justify-between">
              Dashboard
              <span className="badge ml-2">New</span>
            </Link>
          </li>
          <li>
            <button onClick={handleLogout}>Logout</button>
          </li>
        </ul>
      </div>
    </div>
  );
}
