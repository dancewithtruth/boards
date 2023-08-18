'use client';

import {
  BoardWithMembers,
  InviteWithReceiver,
  UpdateInviteParams,
  User,
  createInvites,
  listInvitesByBoard,
  listUsersByFuzzyEmail,
  updateInvite,
} from '@/api';
import { ChangeEvent, useEffect, useState } from 'react';
import { FaPlus } from 'react-icons/fa';
import { FiX } from 'react-icons/fi';
import Avatar from '../avatar';
import Cookies from 'universal-cookie';
import { COOKIE_NAME_JWT_TOKEN, INVITE_STATUS } from '@/constants';
import { toast } from 'react-toastify';

export default function InviteMemberModal({ board, user }: { board: BoardWithMembers; user: User }) {
  const ID = 'modal_invite_member';
  const [selected, setSelected] = useState<User[]>([]);
  const [pendingInvites, setPendingInvites] = useState<InviteWithReceiver[]>([]);
  const cookies = new Cookies();
  const token = cookies.get(COOKIE_NAME_JWT_TOKEN);

  const handleOpen = async () => {
    (window as any)[ID].showModal();
    const response = await listInvitesByBoard(board.id, token, INVITE_STATUS.PENDING);
    const sentPendingInvites = response.result.filter((invite) => invite.sender_id == user.id);
    setPendingInvites(sentPendingInvites);
  };

  const handleClose = () => {
    setSelected([]);
    (window as any)[ID].close();
  };

  const handleSelect = (user: User) => {
    if (isSelected(user.id, selected)) {
      //pop it out
      const newSelected = selected.filter(({ id }) => id !== user.id);
      setSelected(newSelected);
    } else {
      //add it in
      setSelected([...selected, user]);
    }
  };

  const handleCancelInvite = async (id: string) => {
    const params: UpdateInviteParams = { status: INVITE_STATUS.CANCELLED };
    await updateInvite(id, params, token);
    const newPendingInvites = pendingInvites.filter(({ id: pendingInviteId }) => pendingInviteId != id);
    setPendingInvites(newPendingInvites);
  };

  const handleSendInvites = async () => {
    const invites = selected.map((user) => ({
      receiver_id: user.id,
    }));
    const params = {
      board_id: board.id,
      sender_id: user.id,
      invites,
    };
    await createInvites(params, token);
    setSelected([]);
    const response = await listInvitesByBoard(board.id, token, INVITE_STATUS.PENDING);
    const sentPendingInvites = response.result.filter((invite) => invite.sender_id == user.id);
    setPendingInvites(sentPendingInvites);
    toast.success('Invites sent.');
  };

  return (
    <div>
      <button className="btn btn-primary btn-sm" onClick={handleOpen}>
        <FaPlus />
        Invite
      </button>
      <dialog id={ID} className="modal">
        <div className="card bg-white w-[700px] shadow-md">
          <div className="card-body">
            <div className="flex justify-between items-center">
              <h3 className="text-xl font-bold">Invite members to this board</h3>
              <button type="reset" className="btn btn-ghost btn-sm" onClick={handleClose}>
                <FiX size={24} />
              </button>
            </div>
            <div className="flex space-x-4">
              <SearchPanel selected={selected} handleSelect={handleSelect} board={board} />
              <RightPanel
                selected={selected}
                handleSelect={handleSelect}
                pendingInvites={pendingInvites}
                handleCancelInvite={handleCancelInvite}
              />
            </div>
            <div className="flex justify-end">
              <button className="btn btn-primary btn-sm" disabled={selected.length == 0} onClick={handleSendInvites}>
                Send Invites
              </button>
            </div>
          </div>
        </div>
      </dialog>
    </div>
  );
}

type SearchPanelProps = {
  selected: User[];
  handleSelect: (user: User) => void;
  board: BoardWithMembers;
};

const SearchPanel = ({ selected, handleSelect, board }: SearchPanelProps) => {
  const [email, setEmail] = useState('');
  const [previousEmail, setPreviousEmail] = useState('');
  const [search, setSearch] = useState<User[]>([]);
  const memberIDs = board.members.map((member) => member.id);

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const newEmail = e.target.value;
    setEmail(newEmail);
  };

  useEffect(() => {
    if (email != previousEmail) {
      fetchAndSetSearchedUsers();
      setPreviousEmail(email);
    }
  }, [email]);

  async function fetchAndSetSearchedUsers() {
    if (email !== '') {
      const response = await listUsersByFuzzyEmail(email);
      const filteredResults = response.result.filter((user) => !memberIDs.includes(user.id) && user.is_verified == true);
      setSearch(filteredResults);
    }
  }

  return (
    <div className="flex-grow-[2] card card-compact">
      <div className="form-control">
        <input
          id="email-search"
          type="text"
          placeholder="Search for users by email"
          className="input input-bordered text-sm !py-0"
          onChange={handleChange}
          value={email}
        />
      </div>
      <div className="text-sm pt-2 font-bold">Suggested</div>
      <div className="min-h-[200px] max-h-[400px] overflow-auto">
        {search.map((user) => {
          return (
            <div
              key={`select-user-${user.id}`}
              className="flex justify-between items-center cursor-pointer hover:bg-gray-100 pr-2"
              onClick={() => handleSelect(user)}
            >
              <div className="flex items-center space-x-2 px-2 py-3">
                <div key={`search-${user.id}`}>
                  <Avatar id={user.id} />
                </div>
                <span className="font-bold text-gray-700 text-sm">{user.name}</span>
                <span className="text-gray-400 text-sm">{`(${user.email})`}</span>
              </div>
              <input
                type="checkbox"
                checked={isSelected(user.id, selected)}
                className="checkbox checkbox-xs"
                onClick={() => handleSelect(user)}
              />
            </div>
          );
        })}
      </div>
    </div>
  );
};

type RightPanelProps = {
  selected: User[];
  pendingInvites: InviteWithReceiver[];
  handleSelect: (user: User) => void;
  handleCancelInvite: (id: string) => void;
};

const RightPanel = ({ pendingInvites, selected, handleSelect, handleCancelInvite }: RightPanelProps) => {
  return (
    <div className="flex-grow-[1] bg-gray-50 min-w-[250px]">
      <div className="card card-compact">
        <div className="card-body">
          <div className="card-title text-sm">{`${selected.length} user(s) selected`}</div>
          {selected.map((user) => (
            <div className="flex justify-between items-center w-full">
              <span className="text-sm text-gray-700" style={{ overflowWrap: 'anywhere' }}>
                {user.email}
              </span>
              <button className="btn btn-ghost btn-xs" onClick={() => handleSelect(user)}>
                <FiX />
              </button>
            </div>
          ))}
        </div>
      </div>
      {pendingInvites.length ? (
        <div className="card card-compact">
          <div className="card-body">
            <div className="card-title text-sm">{`${pendingInvites.length} pending invite(s)`}</div>
            {pendingInvites.map((invite) => (
              <div className="flex justify-between items-center w-full">
                <span className="text-sm text-gray-700" style={{ overflowWrap: 'anywhere' }}>
                  {invite.receiver.email}
                </span>
                <button className="btn btn-ghost btn-xs" onClick={() => handleCancelInvite(invite.id)}>
                  <FiX />
                </button>
              </div>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
};

function isSelected(userID: string, selected: User[]): boolean {
  let isSelected = false;
  selected.some((user) => {
    if (user.id == userID) {
      isSelected = true;
      return;
    }
  });
  return isSelected;
}
