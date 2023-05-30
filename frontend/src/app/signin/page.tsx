'use client';

import SigninPanel from '../../components/forms/signin';

const Signin = () => {
  return (
    <div className="flex items-center justify-center min-h-screen mt-[-5%]">
      <div className="w-full flex items-center justify-center flex-col space-y-4">
        <h1 className="text-5xl font-bold">Sign in</h1>
        <SigninPanel />
      </div>
    </div>
  );
};

export default Signin;
