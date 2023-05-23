import { createUser } from '../../helpers/api/users';
import { useState, ChangeEvent, FormEvent } from 'react';

const SignUpPanel = (): JSX.Element => {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isFormValid, setIsFormValid] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const handleNameChange = (e: ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
    checkFormValidity();
  };

  const handleEmailChange = (e: ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
    checkFormValidity();
  };

  const handlePasswordChange = (e: ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
    checkFormValidity();
  };

  const checkFormValidity = () => {
    setIsFormValid(name !== '' && email !== '' && password !== '');
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const userId = await createUser({ name, email, password });
      console.log(userId);
    } catch (error) {
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="card flex-shrink-0 w-full max-w-sm shadow-2xl bg-base-100 border border-base-300">
      <form onSubmit={handleSubmit}>
        <div className="card-body">
          <div className="form-control">
            <label className="label">
              <span className="label-text">Name</span>
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
          <div className="form-control mt-6">
            <div className="flex flex-col w-full border-opacity-50">
              <button type="submit" className="btn btn-secondary btn-outline" disabled={!isFormValid || isLoading}>
                {isLoading ? 'Loading...' : 'Sign Up'}
              </button>
              <div className="divider">OR</div>
              <button type="submit" className="btn btn-primary">
                Continue as guest
              </button>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
};

export default SignUpPanel;
