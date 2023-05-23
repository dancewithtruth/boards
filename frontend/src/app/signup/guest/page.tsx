'use client';

import SignUpPanel from '../../../components/signup';

const Signup = () => {
  return (
    <div className="flex items-center justify-center min-h-screen mt-[-5%]">
      <div className="w-full flex items-center justify-center flex-col space-y-4">
        <h1 className="text-5xl font-bold">Sign up</h1>
        <SignUpPanel isGuest={true} />
      </div>
    </div>
  );
};

export default Signup;
