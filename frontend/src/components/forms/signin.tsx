'use client';

import { login } from '@/api';
import { COOKIE_NAME_JWT_TOKEN } from '@/constants';
import { useRouter } from 'next/navigation';
import { ChangeEvent, FormEvent, useState } from 'react';
import { toast } from 'react-toastify';
import Cookies from 'universal-cookie';

export default function SignInForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const cookies = new Cookies();
  const router = useRouter();

  const handleEmailChange = (e: ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
  };

  const handlePasswordChange = (e: ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const { token } = await login({ email, password });
      toast.success('Successfully signed in');
      const expirationDate = new Date();
      expirationDate.setDate(expirationDate.getDate() + 30);
      cookies.set(COOKIE_NAME_JWT_TOKEN, token, { path: '/', expires: expirationDate });
      router.refresh();
      router.push('/dashboard');
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
          <span className="label-text text-stone-600">Email</span>
        </label>
        <input
          type="email"
          id="email"
          className="input input-bordered w-full max-w-xs"
          value={email}
          onChange={handleEmailChange}
          required
        />
      </div>
      <div className="form-control">
        <label className="label">
          <span className="label-text  text-stone-600">Password</span>
        </label>
        <input
          type="password"
          id="password"
          className="input input-bordered w-full max-w-x"
          value={password}
          onChange={handlePasswordChange}
          required
        />
      </div>
      <div className="form-control mt-6">
        <button type="submit" className="btn btn-primary" disabled={isLoading}>
          {isLoading ? 'Signing in...' : 'Sign in'}
        </button>
      </div>
    </form>
  );
}
