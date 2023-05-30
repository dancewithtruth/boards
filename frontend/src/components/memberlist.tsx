import React from 'react';
import { Member } from '../../helpers/api/boards';
import Avatar from './avatar';

interface MemberListProps {
  members: Member[];
}

const MemberList: React.FC<MemberListProps> = ({ members }) => {
  return (
    <div className="flex">
      {members.map(({ id, name, membership }) => (
        <div data-tooltip-id="my-tooltip" data-tooltip-html={name}>
          <div key={id} className="relative flex-shrink-0 w-8 h-8 mr-2 bg-gray-200 rounded-full overflow-hidden">
            <Avatar id={id} />
          </div>
        </div>
      ))}
    </div>
  );
};

export default MemberList;
