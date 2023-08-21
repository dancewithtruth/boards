'use client';

import { sendVerificationEmail, verifyEmail } from '@/api/board';
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

  const handleResend = async () => {
    try {
      await sendVerificationEmail(token);
      toast.success('Email verification resent.');
    } catch (e) {
      toast.error('Issue resending email verification.');
    }
  };

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
        <p>
          By visiting this page, your account should be verified. If something went wrong and you'd like another
          verification email, please request one below.
        </p>
        <button className="btn btn-primary mt-4" onClick={handleResend}>
          Resend Email
        </button>
      </WidthContainer>
    </div>
  );
}
