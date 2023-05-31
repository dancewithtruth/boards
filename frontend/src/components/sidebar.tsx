'use client';
import { useUser } from '@/providers/user';
import Avatar from './avatar';
import { useBoard } from '@/providers/board';

interface SidebarProps {
  width: number;
}

const Sidebar = ({ width }: SidebarProps) => {
  const {
    state: { user },
  } = useUser();
  const {
    state: { board },
    dispatch,
  } = useBoard();

  const onlineCount = 3;
  const totalAccessCount = 5;

  const onlineFraction = `${onlineCount}/${totalAccessCount}`;

  return (
    <div
      className="fixed top-h-16 left-0 bg-base-100 shadow-md z-40"
      style={{ height: `calc(100vh - 4rem)`, width: `${width}px` }}
    >
      <div className="flex flex-col items-center justify-between h-full py-8">
        <div className="flex flex-col items-center w-full">
          <p className="text-gray-700 text-md font-bold">Members</p>
          <div
            className="overflow-y-auto max-h-[500px] w-full p-6"
            style={{ background: 'linear-gradient(to bottom, rgba(255, 255, 255, 0), rgba(255, 255, 255, 1))' }}
          >
            <div className="flex flex-col items-start space-y-3">
              {board?.members.map(({ id, name }) => {
                return (
                  <div key={id} className="flex space-x-2 items-center">
                    <Avatar id={id} />
                    <span className="text-sm text-gray-700">{name}</span>
                  </div>
                );
              })}
            </div>
          </div>
          <div className="divider" />
        </div>

        <div className="flex flex-col justify-center mt-4 space-y-2">
          <div className="badge badge-primary rounded-md px-2 py-1">{onlineFraction}</div>
          <span className="text-xs text-gray-500">Online</span>
        </div>
      </div>
    </div>
  );
};

export default Sidebar;
