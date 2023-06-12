'use client';

import { BoardWithMembers, User, listUsersByFuzzyEmail } from '@/api';
import { ChangeEvent, useEffect, useState } from 'react';
import { FaPlus } from 'react-icons/fa';
import { FiX } from 'react-icons/fi';
import Avatar from '../avatar';

export default function InviteMemberModal({ board }: { board: BoardWithMembers }) {
  const ID = 'modal_invite_member';
  const [selected, setSelected] = useState<User[]>([]);

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

  return (
    <div>
      <button className="btn btn-primary btn-sm" onClick={() => (window as any)[ID].showModal()}>
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
              <SelectedPanel selected={selected} handleSelect={handleSelect} />
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
    const response = await listUsersByFuzzyEmail(email);
    const filteredResults = response.result.filter((user) => !memberIDs.includes(user.id));
    setSearch(filteredResults);
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

type SelectedPanelProps = {
  selected: User[];
  handleSelect: (user: User) => void;
};

const SelectedPanel = ({ selected, handleSelect }: SelectedPanelProps) => {
  return (
    <div className="flex-grow-[1] card card-compact bg-gray-50 min-w-[250px]">
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
