import SignUpGuestForm from '@/components/forms/signupGuest';
import Redirect from '@/components/redirect';

export default function SignUpGuestPage() {
  return (
    <div className="card card-bordered shadow-xl w-96 bg-white">
      <div className="card-body space-y-3">
        <div className="card-title text-2xl">Guest sign up</div>
        <SignUpGuestForm />
        <Redirect url="/auth/signin" redirectText="Already have an account?" buttonText="Sign in" />
      </div>
    </div>
  );
}
