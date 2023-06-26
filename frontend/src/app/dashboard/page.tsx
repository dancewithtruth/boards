import { cookies } from 'next/headers';
import { listBoards, listInvitesByReceiver } from '@/api/board';
import Board from '@/components/board';
import WidthContainer from '@/components/widthContainer';
import { COOKIE_NAME_JWT_TOKEN, FOOTER_HEIGHT, INVITE_STATUS, NAVBAR_HEIGHT } from '@/constants';
import NewBoardModal from '@/components/modals/newBoard';
import InviteAlerts from '@/components/inviteAlerts';

export const metadata = {
  title: 'Dashboard',
  description: 'Collaborate with your team',
};

async function listBoardsData() {
  const cookieStore = cookies();
  const jwtToken = cookieStore.get(COOKIE_NAME_JWT_TOKEN);
  if (jwtToken) {
    const data = await listBoards(jwtToken.value);
    return data;
  } else {
    throw new Error('Please log in.');
  }
}

async function getPendingInvites() {
  const cookieStore = cookies();
  const jwtToken = cookieStore.get(COOKIE_NAME_JWT_TOKEN);
  if (jwtToken) {
    const response = await listInvitesByReceiver(jwtToken.value, INVITE_STATUS.PENDING);
    return response.result;
  } else {
    throw new Error('Please log in.');
  }
}

export default async function DashboardPage() {
  const boardsData = listBoardsData();
  const invitesData = getPendingInvites();

  const [boards, pendingInvites] = await Promise.all([boardsData, invitesData]);

  return (
    <div className="min-h-screen" style={{ minHeight: `calc(100vh - ${NAVBAR_HEIGHT} - ${FOOTER_HEIGHT})` }}>
      <WidthContainer>
        <h1 className="text-4xl font-bold my-5">Dashboard</h1>
        <InviteAlerts invites={pendingInvites} />
        <div className="space-y-8 my-8">
          <div>
            <div className="flex justify-between">
              <h2 className="text-2xl font-bold">My Boards</h2>
              <NewBoardModal />
            </div>
            <div className="divider" />
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 place-items-top">
              {boards.owned.map((board) => {
                return <Board key={board.id} board={board} />;
              })}
            </div>
          </div>
          <div>
            <h2 className="text-2xl font-bold">Shared Boards</h2>
            <div className="divider" />
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 place-items-center">
              {boards.shared.map((board) => {
                return <Board key={board.id} board={board} />;
              })}
            </div>
          </div>
        </div>
      </WidthContainer>
    </div>
  );
}
