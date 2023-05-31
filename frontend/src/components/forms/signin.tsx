'use client';

import { useState, ChangeEvent, FormEvent } from 'react';
import { toast } from 'react-toastify';
import { useRouter } from 'next/navigation';
import { login } from '../../../api/auth';
import ConfiguredToastContainer from '../toastcontainer';
import { useUser } from '@/providers/user';
import { LOCAL_STORAGE_AUTH_TOKEN } from '../../../constants';
import Link from 'next/link';
import { getUserByJwt } from '../../../api/users';

const SignInPanel = (): JSX.Element => {
  const { dispatch } = useUser();
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [isLoading, setIsLoading] = useState(false);

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
      localStorage.setItem(LOCAL_STORAGE_AUTH_TOKEN, token);
      const user = await getUserByJwt(token);
      dispatch({ type: 'set_user', payload: user });
      dispatch({ type: 'set_is_authenticated', payload: true });
      router.push('/dashboard');
    } catch (error) {
      toast.error(String(error));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-base-100 border border-base-300">
      <ConfiguredToastContainer />
      <form onSubmit={handleSubmit}>
        <div className="card-body">
          <>
            <div className="form-control">
              <label className="label">
                <span className="label-text">Email</span>
              </label>
              <input
                type="email"
                id="email"
                className="input input-bordered w-full max-w-xs"
                placeholder="Email"
                value={email}
                onChange={handleEmailChange}
                required
              />
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">Password</span>
                <span className="label-text-alt text-xs text-gray-300">min. 8 char</span>
              </label>
              <input
                type="password"
                id="password"
                className="input input-bordered w-full max-w-xs"
                placeholder="Password"
                value={password}
                onChange={handlePasswordChange}
                required
              />
            </div>
          </>
          <div className="form-control mt-6">
            <div className="flex flex-col w-full border-opacity-50">
              <button type="submit" className="btn btn-primary">
                Sign in
              </button>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
};

export default SignInPanel;
