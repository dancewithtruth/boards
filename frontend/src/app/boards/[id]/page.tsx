'use client';

import Head from 'next/head';
import { Metadata } from 'next';
import { useEffect, useState } from 'react';
import { useBoard } from '@/providers/board';
import { getBoard } from '../../../../helpers/api/boards';

// export const metadata: Metadata = {
//   title: 'Boards',
//   description: 'Boards is a live collaboration tool aimed to increase your productivity.',
// };

const Board = ({ params: { id } }: { params: { id: string } }) => {
  const { dispatch } = useBoard();
  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    const response = await getBoard(id);
    dispatch({ type: 'set_board', payload: response });
  };
  return (
    <>
      <Head>
        <title>Boards</title>
      </Head>
      <div className="w-screen h-screen"></div>
    </>
  );
};

export default Board;
