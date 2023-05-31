'use client';

import { useState, ChangeEvent, FormEvent } from 'react';
import { toast } from 'react-toastify';
import { useRouter } from 'next/navigation';
import { createUser } from '@/api/users';
import ConfiguredToastContainer from '@/components/toastcontainer';
import { useUser } from '@/providers/user';
import { LOCAL_STORAGE_AUTH_TOKEN } from '@/constants';
import Link from 'next/link';

type SignUpPanelParams = {
  isGuest?: boolean;
};
const SignUpPanel = ({ isGuest = false }: SignUpPanelParams): JSX.Element => {
  const { dispatch } = useUser();
  const [name, setName] = useState('');
  const [email, setEmail] = useState<string | undefined>();
  const [password, setPassword] = useState<string | undefined>();
  const [isLoading, setIsLoading] = useState(false);

  const router = useRouter();

  const handleNameChange = (e: ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
  };

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
      const response = await createUser({ name, email, password, isGuest });
      toast.success('Account created!');
      localStorage.setItem(LOCAL_STORAGE_AUTH_TOKEN, response.jwt_token);
      dispatch({ type: 'set_user', payload: response.user });
      dispatch({ type: 'set_is_authenticated', payload: true });
      router.push('/welcome');
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
          <div className="form-control">
            <label className="label">
              <span className="label-text">Name</span>
              <span className="label-text-alt text-xs text-gray-300">2 to 12 char</span>
            </label>
            <input
              type="text"
              id="name"
              className="input input-bordered w-full max-w-xs"
              placeholder="Name"
              value={name}
              onChange={handleNameChange}
              required
            />
          </div>
          {isGuest ? null : (
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
          )}
          <div className="form-control mt-6">
            <div className="flex flex-col w-full border-opacity-50">
              {isGuest ? (
                <button type="submit" className="btn btn-secondary">
                  {isLoading ? 'Loading...' : 'Create guest account'}
                </button>
              ) : (
                <>
                  <button type="submit" className="btn btn-secondary btn-outline">
                    {isLoading ? 'Loading...' : 'Sign Up'}
                  </button>
                  <div className="divider">OR</div>
                  <Link href="/signup/guest" className="btn btn-primary">
                    Continue as guest
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </form>
    </div>
  );
};

export default SignUpPanel;
