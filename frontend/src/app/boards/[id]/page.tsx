'use client';

import Head from 'next/head';
import { Metadata } from 'next';
import { useEffect, useState, useCallback } from 'react';
import { useBoard } from '@/providers/board';
import { getBoard } from '../../../../helpers/api/boards';
import Post from '@/components/post';
import update from 'immutability-helper';
import type { CSSProperties, FC } from 'react';
import type { XYCoord } from 'react-dnd';
import { useDrop } from 'react-dnd';

import { ItemTypes } from '../../../../constants';

// export const metadata: Metadata = {
//   title: 'Boards',
//   description: 'Boards is a live collaboration tool aimed to increase your productivity.',
// };

export interface DragItem {
  type: string;
  id: string;
  top: number;
  left: number;
}

const Board = ({ params: { id } }: { params: { id: string } }) => {
  const { dispatch } = useBoard();
  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    const response = await getBoard(id);
    dispatch({ type: 'set_board', payload: response });
  };

  const handleDoubleClick = (event: React.MouseEvent<HTMLDivElement>) => {
    const { offsetX, offsetY } = event.nativeEvent;
    addPost('123abc', offsetX, offsetY);
  };

  const [posts, setPosts] = useState<{
    [key: string]: {
      top: number;
      left: number;
      content: string;
    };
  }>({
    a: { top: 20, left: 80, content: 'Drag me around' },
    b: { top: 180, left: 20, content: 'Drag me too' },
  });

  const addPost = useCallback(
    (id: string, left: number, top: number) => {
      setPosts(update(posts, { $merge: { [id]: { left, top } } }));
    },
    [posts, setPosts]
  );

  const movePost = useCallback(
    (id: string, left: number, top: number) => {
      setPosts(
        update(posts, {
          [id]: {
            $merge: { left, top },
          },
        })
      );
    },
    [posts, setPosts]
  );

  const [, drop] = useDrop(
    () => ({
      accept: ItemTypes.POST,
      drop(item: DragItem, monitor) {
        const delta = monitor.getDifferenceFromInitialOffset() as XYCoord;
        const left = Math.round(item.left + delta.x);
        const top = Math.round(item.top + delta.y);
        movePost(item.id, left, top);
        return undefined;
      },
    }),
    [movePost]
  );
  return (
    <>
      <Head>
        <title>Boards</title>
      </Head>
      <div
        ref={drop}
        className="h-screen relative"
        style={{ width: `calc(100vw - 12rem)` }}
        onDoubleClick={handleDoubleClick}
      >
        {Object.keys(posts).map((key) => {
          const { left, top, content } = posts[key] as {
            top: number;
            left: number;
            content: string;
          };
          return <Post key={key} id={key} left={left} top={top} content={content} />;
        })}
      </div>
    </>
  );
};

export default Board;
