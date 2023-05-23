'use client';

import { useUser } from '@/providers/user';
import avatar from 'gradient-avatar';
import Link from 'next/link';
import { LOCAL_STORAGE_AUTH_TOKEN } from '../../constants';

const Navbar = () => {
  const {
    state: { user, isAuthenticated },
    dispatch,
  } = useUser();
  const avatarSVG = avatar(user?.id || 'default');
  const dataUri = `data:image/svg+xml,${encodeURIComponent(avatarSVG)}`;

  const handleLogout = () => {
    dispatch({ type: 'set_is_authenticated', payload: false });
    dispatch({ type: 'set_user', payload: null });
    localStorage.removeItem(LOCAL_STORAGE_AUTH_TOKEN);
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
          <div className="flex items-center space-x-2">
            {isAuthenticated ? (
              <div className="dropdown dropdown-end">
                <label tabIndex={0} className="btn btn-ghost btn-circle avatar">
                  <div className="w-10 rounded-full">
                    <img src={dataUri} alt="SVG Image" />
                  </div>
                </label>
                <ul
                  tabIndex={0}
                  className="mt-3 p-2 shadow menu menu-compact dropdown-content bg-base-100 rounded-box w-52"
                >
                  <li>
                    <a className="justify-between">
                      Profile
                      <span className="badge">New</span>
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

export default Navbar;
