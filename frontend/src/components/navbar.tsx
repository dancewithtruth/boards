'use client';

import { useUser } from '@/providers/user';
import avatar from 'gradient-avatar';
import Link from 'next/link';
import { LOCAL_STORAGE_AUTH_TOKEN } from '../../constants';
import { useRouter } from 'next/navigation';
import { FaChevronDown } from 'react-icons/fa';

const Navbar = () => {
  const {
    state: { user, isAuthenticated },
    dispatch,
  } = useUser();
  const avatarSVG = avatar(user?.id || 'default');
  const dataUri = `data:image/svg+xml,${encodeURIComponent(avatarSVG)}`;
  const router = useRouter();

  const handleLogout = () => {
    dispatch({ type: 'set_is_authenticated', payload: false });
    dispatch({ type: 'set_user', payload: null });
    localStorage.removeItem(LOCAL_STORAGE_AUTH_TOKEN);
    router.push('/');
  };
  return (
    <nav className="bg-base-100 shadow-md">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          <div className="flex items-center">
            <Link href="/" className="font-bold text-xl">
              Boards
            </Link>
          </div>
          <div className="flex items-center space-x-4">
            {isAuthenticated ? (
              <>
                <Link href="/dashboard" className="btn btn-primary btn-sm">
                  Dashboard
                </Link>
                <div className="dropdown dropdown-end z-51">
                  <div tabIndex={0} className="btn btn-ghost normal-case space-x-2">
                    <button className="rounded-full avatar">
                      <div className="w-10 rounded-full">
                        <img src={dataUri} alt="SVG Image" />
                      </div>
                    </button>
                    <span>{user?.name}</span>
                    <FaChevronDown />
                  </div>
                  <div className="right-0 mt-3 p-2 shadow menu menu-compact dropdown-content bg-base-100 rounded-box w-52">
                    <ul className="menu menu-compact gap-1 p-3">
                      <li>
                        <button className="flex items-center justify-between">
                          Profile
                          <span className="badge ml-2">New</span>
                        </button>
                      </li>
                      <li>
                        <button>Settings</button>
                      </li>
                      <li>
                        <button onClick={handleLogout}>Logout</button>
                      </li>
                    </ul>
                  </div>
                </div>
              </>
            ) : (
              <>
                <Link href="/signin" className="btn btn-secondary btn-outline">
                  Sign in
                </Link>
                <Link href="/signup" className="btn btn-primary">
                  Sign up
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
