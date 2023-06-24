'use client';
import { NAVBAR_HEIGHT } from '@/constants';
import Avatar from './avatar';
import { User, BoardWithMembers } from '@/api';
import { isAdmin, mergeArrays } from '@/utils';
import InviteMemberModal from './modals/inviteMember';

interface SidebarProps {
  width: string;
  board: BoardWithMembers;
  user: User;
  connectedUsers: User[];
}

const Sidebar = ({ width, board, user, connectedUsers }: SidebarProps) => {
  const allUsers = mergeArrays('id', [user], connectedUsers, board.members);
  connectedUsers = mergeArrays('id', [user], connectedUsers);
  const onlineFraction = `${connectedUsers.length} / ${allUsers.length}`;

  return (
    <div
      className="fixed top-h-16 left-0 bg-base-100 shadow-md"
      style={{ height: `calc(100vh - ${NAVBAR_HEIGHT})`, width, zIndex: 10001 }}
    >
      <div className="flex flex-col items-center justify-between h-full py-8">
        <div className="flex flex-col items-center w-full">
          <p className="text-gray-700 text-md font-bold">Members</p>
          <div
            className="overflow-y-auto max-h-[500px] w-full p-6"
            style={{ background: 'linear-gradient(to bottom, rgba(255, 255, 255, 0), rgba(255, 255, 255, 1))' }}
          >
            <div className="flex flex-col items-start space-y-4">
              {allUsers.map((user) => (
                <SidebarMember key={user.id} user={user} isConnected={isConnected(user.id, connectedUsers)} />
              ))}
            </div>
          </div>
          {isAdmin(user.id, board) ? <InviteMemberModal board={board} user={user} /> : null}
          <div className="divider" />
        </div>
        <div className="flex flex-col justify-center items-center space-y-1">
          <div className="badge badge-primary rounded-md p-3">{onlineFraction}</div>
          <span className="text-sm text-gray-700">Online</span>
        </div>
      </div>
    </div>
  );
};

function SidebarMember({ user, isConnected }: { user: User; isConnected: boolean }) {
  return (
    <div key={user.id} className="flex space-x-2 items-center">
      <Avatar id={user.id} />
      <span className="text-sm" style={{ fontWeight: isConnected ? 700 : 300, color: isConnected ? 'black' : 'gray' }}>
        {user.name}
      </span>
    </div>
  );
}

function isConnected(userID: string, connectedUsers: User[]): boolean {
  let connected = false;
  connectedUsers.some((user) => {
    if (user.id == userID) {
      connected = true;
      return;
    }
  });
  return connected;
}

export default Sidebar;
