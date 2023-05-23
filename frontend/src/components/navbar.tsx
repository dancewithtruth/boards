import Link from 'next/link';

const Navbar = () => {
  return (
    <div className="navbar bg-base-100 shadow-md">
      <div className="flex-1">
        <a className="btn btn-ghost normal-case text-xl">Boards</a>
      </div>
      <div className="flex-none space-x-2">
        <button className="btn btn-secondary btn-outline">Sign in</button>
        <button className="btn btn-primary">Sign up</button>
      </div>
    </div>
  );
};

export default Navbar;
