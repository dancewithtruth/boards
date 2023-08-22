'use client';

import Centered from '@/components/centered';
import { FOOTER_HEIGHT_PX, NAVBAR_HEIGHT_PX } from '@/constants';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { toast } from 'react-toastify';

export default function Error({ error, reset }: { error: Error; reset: () => void }) {
  const router = useRouter();
  if (error.message === 'Please log in.') {
    toast.error('You must be signed in to view the app.', { toastId: 'auth' });
    router.replace('/auth/signin');
  }
  
  return (
    <div className="h-screen" style={{ height: `calc(100vh - ${NAVBAR_HEIGHT_PX} - ${FOOTER_HEIGHT_PX})` }}>
      <Centered>
        <div className="text-center">
          <h1 className="text-4xl font-bold mb-4 text-primary">Oops! Something went wrong.</h1>
          <Link href="/" className="font-bold text-gray-700">
            Go back to homepage
          </Link>
        </div>
      </Centered>
    </div>
  );
}
