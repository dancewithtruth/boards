'use client';
import { useState } from 'react';
import { useUser } from '@/providers/user';
import avatar from 'gradient-avatar';

const Sidebar = () => {
  const { state } = useUser();

  const boardMembers = [
    { id: 'member1', name: 'Victor Djokovich' },
    { id: 'member2', name: 'Alice' },
    { id: 'member3', name: 'Bob' },
    { id: 'member4', name: 'Charlie' },
    { id: 'member5', name: 'David' },
    { id: 'member6', name: 'Eva' },
    { id: 'member7', name: 'Frank' },
    { id: 'member8', name: 'Grace' },
    { id: 'member9', name: 'Henry' },
    { id: 'member10', name: 'Isabel' },
    { id: 'member11', name: 'Jack' },
    { id: 'member12', name: 'Karen' },
    { id: 'member13', name: 'Liam' },
    { id: 'member14', name: 'Mia' },
    { id: 'member15', name: 'Nora' },
  ];

  const onlineCount = 3;
  const totalAccessCount = 5;

  const onlineFraction = `${onlineCount}/${totalAccessCount}`;

  return (
    <div className="fixed top-h-16 left-0 bg-base-100 shadow-md w-48 z-40" style={{ height: `calc(100vh - 4rem)` }}>
      <div className="flex flex-col items-center justify-between h-full py-8">
        <div className="flex flex-col items-center w-full">
          <p className="text-gray-700 text-md font-bold">Collaborators</p>
          <div
            className="overflow-y-auto max-h-[500px] w-full p-6"
            style={{ background: 'linear-gradient(to bottom, rgba(255, 255, 255, 0), rgba(255, 255, 255, 1))' }}
          >
            <div className="flex flex-col items-start space-y-3">
              {boardMembers.map(({ id, name }) => {
                const avatarSVG = avatar(id);
                const dataUri = `data:image/svg+xml,${encodeURIComponent(avatarSVG)}`;
                return (
                  <div key={id} className="flex space-x-2 items-center">
                    <div className="w-8 h-8">
                      <img className="w-full h-full rounded-full" src={dataUri} alt="Avatar" />
                    </div>
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
