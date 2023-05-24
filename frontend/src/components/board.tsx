import React from 'react';
import { FaEllipsisV } from 'react-icons/fa';
import TimeAgo from './timeago';

interface BoardroomCardProps {
  title: string;
  description: string;
  createdAt: string;
}

const BoardroomCard: React.FC<BoardroomCardProps> = ({ title, description, createdAt }) => {
  return (
    <div className="bg-base-200 shadow-sm rounded-md p-4 h-[200px] w-[300px] flex flex-col justify-evenly">
      <div className="mb-2 flex justify-between">
        <div>
          <p className="text-xs text-gray-400">Board Name</p>
          <h3 className="text-lg font-bold">{title}</h3>
        </div>
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
      </div>
      <div className="mb-4">
        <p className="text-xs text-gray-400">Description</p>
        <p className="text-sm text-gray-600">{description}</p>
      </div>
      <div className="flex justify-between items-center">
        <div className="flex space-x-2">
          <span className="bg-gray-200 text-xs text-gray-600 p-2 rounded">
            Created <TimeAgo timestamp={createdAt} />
          </span>
        </div>
        <button className="btn btn-secondary btn-sm btn-outline">View</button>
      </div>
    </div>
  );
};

export default BoardroomCard;
