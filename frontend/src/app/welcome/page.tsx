'use client';

import { useUser } from '@/providers/user';
import Link from 'next/link';

const Welcome = () => {
  const {
    state: { user },
  } = useUser();
  return (
    <div className="hero min-h-screen bg-base-200">
      <div className="hero-content text-center">
        <div className="max-w-md">
          <h1 className="text-5xl font-bold">{`Hi ${user?.name || 'there'}!`}</h1>
          <p className="py-6">
            {`Boards is a live collaboration tool focused on making team retrospectives fast and simple.`}
          </p>
          <div className="flex justify-center space-x-4">
            <a href="https://github.com/Wave-95/boards" target="_blank" className="btn btn-secondary btn-outline">
              View GitHub
            </a>
            <Link href="/dashboard" className="btn btn-primary">
              Explore App
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Welcome;
