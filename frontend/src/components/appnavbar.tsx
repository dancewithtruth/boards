'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useUser } from '@/providers/user';
import avatar from 'gradient-avatar';
import { LOCAL_STORAGE_AUTH_TOKEN } from '../../constants';

const AppNavbar = () => {
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
    router.replace('/');
  };

  return (
    <nav className="fixed top-0 left-0 w-full bg-base-100 shadow-md z-50">
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
                <div className="relative dropdown dropdown-end">
                  <button className="btn btn-ghost btn-circle avatar">
                    <div className="w-10 rounded-full">
                      <img src={dataUri} alt="SVG Image" />
                    </div>
                  </button>
                  <ul className="absolute right-0 mt-3 p-2 shadow menu menu-compact dropdown-content bg-base-100 rounded-box w-52">
                    <li>
                      <a className="flex items-center justify-between">
                        Profile
                        <span className="badge ml-2">New</span>
                      </a>
                    </li>
                    <li>
                      <a>Settings</a>
                    </li>
                    <li>
                      <a onClick={handleLogout}>Logout</a>
                    </li>
                  </ul>
                </div>
              </>
            ) : (
              <>
                <button className="btn btn-secondary btn-outline">Sign in</button>
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

export default AppNavbar;
