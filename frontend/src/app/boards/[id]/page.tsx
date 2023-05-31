'use client';

import Head from 'next/head';
import { Metadata } from 'next';
import { useEffect, useState, useCallback } from 'react';
import { useBoard } from '@/providers/board';
import { getBoard } from '../../../../api/boards';
import Post, { PostData } from '@/components/post';
import update from 'immutability-helper';
import type { CSSProperties, FC } from 'react';
import type { XYCoord } from 'react-dnd';
import { useDrop } from 'react-dnd';

import { API_BASE_URL, ItemTypes, LOCAL_STORAGE_AUTH_TOKEN, POST_COLORS, WS_BASE_URL } from '../../../../constants';
import WebSocketConnection from '../../../../websocket';

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

export interface Post {
  top: number;
  left: number;
  content: string;
  zIndex: number;
  color: string;
  customHeight?: number;
}

const getHighestZIndex = (posts: { [key: string]: Post }) => {
  const indexValues = Object.values(posts).map(({ zIndex }) => zIndex);
  return Math.max(...indexValues);
};

const getColor = (posts: { [key: string]: Post }) => {
  const chosenColors = Object.values(posts).map(({ color }) => color);
  const availableColors = Object.values(POST_COLORS);
  availableColors.forEach((color) => {
    if (!chosenColors.includes(color)) {
      return color;
    }
  });
  const randIndex = Math.floor(Math.random() * availableColors.length);
  return availableColors[randIndex];
};

const Board = ({ params: { id } }: { params: { id: string } }) => {
  const { dispatch } = useBoard();
  const [posts, setPosts] = useState<{
    [key: string]: Post;
  }>({});
  const [highestZ, setHighestZ] = useState(getHighestZIndex(posts));
  const [colorSetting, setColorSetting] = useState(getColor(posts));
  var ws;
  useEffect(() => {
    fetchData();

    // Establish ws connection
    // TODO: Handle if browser does not have websocket API
    if (window.WebSocket) {
      const token = localStorage.getItem(LOCAL_STORAGE_AUTH_TOKEN);
      const wsUrl = `${WS_BASE_URL}/ws?token=${token}&boardId=${id}`;
      ws = new WebSocketConnection(wsUrl);
    }
  }, []);

  const fetchData = async () => {
    const response = await getBoard(id);
    dispatch({ type: 'set_board', payload: response });
  };

  const handleDoubleClick = (event: React.MouseEvent<HTMLDivElement>) => {
    if (event.target === event.currentTarget) {
      const { offsetX, offsetY } = event.nativeEvent;
      //send request
      const randId = Math.random() * 1000000;
      const data = { left: offsetX, top: offsetY, color: colorSetting } as PostData;
      addPost(randId.toString(), data);
    }
  };

  const addPost = useCallback(
    (id: string, data: PostData) => {
      const { left, top, color } = data;
      setPosts(update(posts, { $merge: { [id]: { left, top, color } } }));
    },
    [posts, setPosts]
  );

  const updatePost = useCallback(
    (id: string, data: PostData) => {
      const { left, top, zIndex, color } = data;
      const mergeObject: Partial<PostData> = {};

      if (left !== undefined) {
        mergeObject['left'] = left;
      }

      if (top !== undefined) {
        mergeObject['top'] = top;
      }

      if (zIndex !== undefined) {
        mergeObject['zIndex'] = zIndex;
      }

      if (color !== undefined) {
        mergeObject['color'] = color;
      }
      setPosts(
        update(posts, {
          [id]: {
            $merge: mergeObject,
          },
        })
      );
    },
    [posts, setPosts]
  );

  const deletePost = useCallback(
    (id: string) => {
      setPosts(
        update(posts, {
          $unset: [id],
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
        const zIndex = highestZ + 1;
        setHighestZ(zIndex);
        const data = { left: Math.max(left, 0), top: Math.max(top, 0), zIndex } as PostData;
        updatePost(item.id, data);
        return undefined;
      },
    }),
    [updatePost]
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
          const post = posts[key] as Post;
          const data = { id: key, ...post };
          return (
            <Post
              key={key}
              data={data}
              updatePost={(data: PostData) => updatePost(key, data)}
              setColor={setColorSetting}
              deletePost={deletePost}
            />
          );
        })}
      </div>
    </>
  );
};

export default Board;
