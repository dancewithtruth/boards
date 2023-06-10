'use client';

import { createUser, login } from '@/api';
import { COOKIE_NAME_JWT_TOKEN } from '@/constants';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { ChangeEvent, FormEvent, useState } from 'react';
import { toast } from 'react-toastify';
import Cookies from 'universal-cookie';

export default function SignUpGuestForm() {
  const [name, setName] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const cookies = new Cookies();
  const router = useRouter();

  const handleNameChange = (e: ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const { jwt_token: token } = await createUser({ name, isGuest: true });
      toast.success('Successfully signed up as guest.');
      const expirationDate = new Date();
      expirationDate.setDate(expirationDate.getDate() + 30);
      cookies.set(COOKIE_NAME_JWT_TOKEN, token, { path: '/', expires: expirationDate });
      router.push('/dashboard');
      router.refresh();
    } catch (error) {
      toast.error(String(error));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="form-control">
        <label className="label">
          <span className="label-text text-stone-600">Name</span>
        </label>
        <input
          type="name"
          id="name"
          className="input input-bordered w-full max-w-xs"
          value={name}
          onChange={handleNameChange}
          required
        />
      </div>
      <div className="form-control mt-6">
        <button type="submit" className="btn btn-secondary" disabled={isLoading}>
          {isLoading ? 'Creating account...' : 'Sign up'}
        </button>
      </div>
    </form>
  );
}
