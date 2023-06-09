import { cookies } from 'next/headers';
import { getBoards } from '@/api/board';
import Board from '@/components/board';
import WidthContainer from '@/components/widthContainer';
import { COOKIE_NAME_JWT_TOKEN, FOOTER_HEIGHT, NAVBAR_HEIGHT } from '@/constants';

export const metadata = {
  title: 'Dashboard',
  description: 'Collaborate with your team',
};

async function getDashboardData() {
  const cookieStore = cookies();
  const jwtToken = cookieStore.get(COOKIE_NAME_JWT_TOKEN);
  if (jwtToken) {
    const data = await getBoards(jwtToken.value);
    return data;
  } else {
    throw new Error('Please log in.');
  }
}

export default async function DashboardPage() {
  const data = await getDashboardData();
  return (
    <div className="min-h-screen" style={{ minHeight: `calc(100vh - ${NAVBAR_HEIGHT} - ${FOOTER_HEIGHT})` }}>
      <WidthContainer>
        <h1 className="text-4xl font-bold mt-10 mb-10">Dashboard</h1>
        <div className="space-y-8 my-8">
          <div>
            <h2 className="text-2xl font-bold">My Boards</h2>
            <div className="divider" />
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 place-items-center">
              {data.owned.map((board) => {
                return <Board key={board.id} board={board} />;
              })}
            </div>
          </div>
          <div>
            <h2 className="text-2xl font-bold">Shared Boards</h2>
            <div className="divider" />
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 place-items-center">
              {data.shared.map((board) => {
                return <Board key={board.id} board={board} />;
              })}
            </div>
          </div>
        </div>
      </WidthContainer>
    </div>
  );
}
