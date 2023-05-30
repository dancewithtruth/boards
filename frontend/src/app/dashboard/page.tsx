'use client';

import Board from '@/components/board';
import NewBoardForm from '@/components/forms/newboard';
import ConfiguredToastContainer from '@/components/toastcontainer';
import { FaPlus } from 'react-icons/fa';
import { BoardWithMembers, getBoards } from '../../../helpers/api/boards';
import Head from 'next/head';
import { useEffect, useState } from 'react';

const Dashboard = () => {
  const [ownedBoards, setOwnedBoards] = useState<BoardWithMembers[]>([]);
  const [sharedBoards, setSharedBoards] = useState<BoardWithMembers[]>([]);
  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    const response = await getBoards();
    setOwnedBoards(response.owned);
    setSharedBoards(response.shared);
    console.log(response);
    //set response
  };
  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
      <Head>
        <title>Dashboard</title>
      </Head>
      <ConfiguredToastContainer />
      <h1 className="text-4xl font-bold mt-10 mb-10">Dashboard</h1>
      <div>
        <div className="flex justify-between items-end">
          <h2 className="text-2xl font-bold">My Boards</h2>
          <label htmlFor="my-modal-4" className="btn btn-primary">
            <FaPlus className="mr-2" />
            New Board
          </label>
          <input type="checkbox" id="my-modal-4" className="modal-toggle" />
          <label htmlFor="my-modal-4" className="modal cursor-pointer">
            <label className="modal-box relative" htmlFor="">
              <NewBoardForm />
            </label>
          </label>
        </div>
        <div className="divider"></div>
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 place-items-center">
            {ownedBoards.map(({ id, name, description, created_at, members }) => {
              return (
                <Board
                  key={id}
                  id={id}
                  name={name}
                  description={description}
                  members={members}
                  createdAt={created_at}
                />
              );
            })}
          </div>
        </div>
      </div>
      <div className="mt-12">
        <h2 className="text-2xl font-bold">Shared Boards</h2>
        <div className="divider"></div>
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 place-items-center">
            {sharedBoards.map(({ id, name, description, created_at, members }) => {
              return (
                <Board
                  key={id}
                  id={id}
                  name={name}
                  description={description}
                  members={members}
                  createdAt={created_at}
                />
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
