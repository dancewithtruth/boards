import SignInForm from '@/components/forms/signin';
import Redirect from '@/components/redirect';

export default function SigninPage() {
  return (
    <div className="card card-bordered shadow-xl w-96 bg-white">
      <div className="card-body space-y-3">
        <div className="card-title text-2xl">Sign in</div>
        <SignInForm />
        <Redirect url="/auth/signup" redirectText="Don't have an account?" buttonText="Sign up" />
      </div>
    </div>
  );
}
