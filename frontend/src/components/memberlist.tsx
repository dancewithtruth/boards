import React from 'react';
import { Member } from '../api/boards';
import Avatar from './avatar';

interface MemberListProps {
  members: Member[];
}

const MemberList: React.FC<MemberListProps> = ({ members }) => {
  return (
    <div className="flex">
      {members.map(({ id, name }) => (
        <div key={id} data-tooltip-id="my-tooltip" data-tooltip-content={name}>
          <div className="relative flex-shrink-0 w-8 h-8 mr-2 bg-gray-200 rounded-full overflow-hidden">
            <Avatar id={id} />
          </div>
        </div>
      ))}
    </div>
  );
};

export default MemberList;
