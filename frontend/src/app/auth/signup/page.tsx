import SignUpForm from '@/components/forms/signup';
import Redirect from '@/components/redirect';

export default function SignUpPage() {
  return (
    <div className="card card-bordered shadow-xl w-96 bg-white">
      <div className="card-body space-y-3">
        <div className="card-title text-2xl">Sign up</div>
        <SignUpForm />
        <Redirect url="/auth/signin" redirectText="Already have an account?" buttonText="Sign in" />
      </div>
    </div>
  );
}
