'use client';

import { InviteWithBoardAndSender, UpdateInviteParams, updateInvite } from '@/api';
import { COOKIE_NAME_JWT_TOKEN, INVITE_STATUS } from '@/constants';
import { useRouter } from 'next/navigation';
import { FC, useState } from 'react';
import { FaUserPlus } from 'react-icons/fa';
import Cookies from 'universal-cookie';

type InviteAlertsProps = {
  invites: InviteWithBoardAndSender[];
};

const InviteAlerts: FC<InviteAlertsProps> = ({ invites }) => {
  const [pendingInvites, setPendingInvites] = useState<InviteWithBoardAndSender[]>(invites);
  const latestInvite = pendingInvites[0];
  const cookies = new Cookies();
  const token = cookies.get(COOKIE_NAME_JWT_TOKEN);
  const router = useRouter();

  const handleAccept = async () => {
    const params: UpdateInviteParams = { status: INVITE_STATUS.ACCEPTED };
    await updateInvite(latestInvite.id, params, token);
    router.push(`/boards/${latestInvite.board.id}`);
  };

  const handleIgnore = async () => {
    const params: UpdateInviteParams = { status: INVITE_STATUS.IGNORED };
    await updateInvite(latestInvite.id, params, token);
    const newPendingInvites = pendingInvites.filter(({ id }) => id != latestInvite.id);
    setPendingInvites(newPendingInvites);
  };

  return latestInvite ? (
    <div className="alert">
      <FaUserPlus />
      <span className="text-sm">{`${latestInvite.sender.name} (${latestInvite.sender.email}) invited you to a board.`}</span>
      <div className="space-x-2">
        <button className="btn btn-sm btn-primary" onClick={handleAccept}>
          Accept
        </button>
        <button className="btn btn-sm" onClick={handleIgnore}>
          Ignore
        </button>
      </div>
    </div>
  ) : null;
};

export default InviteAlerts;
