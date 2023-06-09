import Centered from '@/components/centered';
import { NAVBAR_HEIGHT } from '@/constants';

export const metadata = {
  title: 'Boards',
  description: 'Collaborate with your team',
};

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div style={{ height: `calc(100vh - ${NAVBAR_HEIGHT})` }}>
      <Centered>{children}</Centered>
    </div>
  );
}
