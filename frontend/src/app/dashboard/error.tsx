'use client';

import Centered from '@/components/centered';
import { FOOTER_HEIGHT_PX, NAVBAR_HEIGHT_PX } from '@/constants';
import Link from 'next/link';

export default function Error({ error, reset }: { error: Error; reset: () => void }) {
  return (
    <div className="h-screen" style={{ height: `calc(100vh - ${NAVBAR_HEIGHT_PX} - ${FOOTER_HEIGHT_PX})` }}>
      <Centered>
        <div className="text-center">
          <h1 className="text-4xl font-bold mb-4 text-primary">Oops! Something went wrong.</h1>
          <p className="mb-8 text-secondary">"Error {error.message}"</p>
          <Link href="/" className="font-bold text-gray-700">
            Go back to homepage
          </Link>
        </div>
      </Centered>
    </div>
  );
}
