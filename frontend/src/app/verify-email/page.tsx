'use client';

import { verifyEmail } from '@/api/board';
import WidthContainer from '@/components/widthContainer';
import { COOKIE_NAME_JWT_TOKEN, FOOTER_HEIGHT_PX, INVITE_STATUS, NAVBAR_HEIGHT_PX } from '@/constants';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { toast } from 'react-toastify';
import Cookies from 'universal-cookie';

export default async function VerifyEmailPage({
  searchParams,
}: {
  searchParams: { [key: string]: string | undefined };
}) {
  const router = useRouter();
  const cookies = new Cookies();
  const token = cookies.get(COOKIE_NAME_JWT_TOKEN);
  const code = searchParams['code'];

  useEffect(() => {
    (async () => {
      if (code && !token) {
        router.push(`/auth/signin?verify-email=true&code=${code}`);
      }

      if (!code) {
        toast.error('Please request a new verification code.');
      }

      if (code && token) {
        try {
          await verifyEmail(code, token);
          toast.success('Successfully verified email.');
          router.push('/dashboard');
        } catch (e) {
          toast.error('Issue verifying email--please request a new verification code.');
        }
      }
    })();
  }, []);

  return (
    <div className="min-h-screen" style={{ minHeight: `calc(100vh - ${NAVBAR_HEIGHT_PX} - ${FOOTER_HEIGHT_PX})` }}>
      <WidthContainer>
        <h1 className="text-4xl font-bold my-5">Verify Email</h1>
      </WidthContainer>
    </div>
  );
}
