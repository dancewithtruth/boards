import Head from 'next/head';
import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Boards',
  description: 'Boards is a live collaboration tool aimed to increase your productivity.',
};

const Board = ({ params }: { params: { id: string } }) => {
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
