'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { FaEllipsisV } from 'react-icons/fa';

import { BoardWithMembers } from '@/api/board';
import TimeAgo from './timeago';
import MemberList from './memberlist';

const Board = ({ board }: { board: BoardWithMembers }) => {
  const router = useRouter();

  const handleClick = () => {
    router.push(`/boards/${board.id}`);
  };

  const BoardName = () => (
    <div>
      <p className="text-xs text-gray-400">Board Name</p>
      <h3 className="text-lg font-bold">{board.name}</h3>
    </div>
  );

  const BoardOptions = () => (
    <div className="dropdown dropdown-left">
      <label tabIndex={0} className="cursor-pointer">
        <FaEllipsisV />
      </label>
      <ul tabIndex={0} className="dropdown-content menu p-1 shadow bg-base-100 rounded-box w-25">
        <li>
          <a className="text-sm">Edit</a>
        </li>
        <li>
          <a className="text-sm">Archive</a>
        </li>
      </ul>
    </div>
  );

  const TopSection = () => (
    <div className="flex justify-between items-center">
      <BoardName />
      <BoardOptions />
    </div>
  );

  const MidSection = () => (
    <div>
      <p className="text-xs text-gray-400">Description</p>
      <p className="text-sm text-gray-600 max-h-[100px] overflow-auto">{board.description}</p>
    </div>
  );

  const BottomSection = () => (
    <>
      <div className="flex justify-between items-center">
        <div className="flex space-x-2">
          <span className="bg-gray-200 text-xs text-gray-600 p-2 rounded">
            Created <TimeAgo timestamp={board.created_at} />
          </span>
        </div>
        <button onClick={handleClick} className="btn btn-secondary btn-sm btn-outline">
          Open
        </button>
      </div>
      <div className="mt-2 overflow-x-auto">
        <MemberList members={board.members}/>
      </div>
    </>
  );

  return (
    <div className="card card-bordered bg-white w-[325px]">
      <div className="card-body">
        <TopSection />
        <MidSection />
        <BottomSection />
      </div>
    </div>
  );
};

export default Board;
