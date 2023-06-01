'use client';

import Head from 'next/head';
import { useEffect, useState, useCallback } from 'react';
import { useBoard } from '@/providers/board';
import { getBoard } from '../../../api/boards';
import PostUI, { Post } from '@/components/post';
import update from 'immutability-helper';
import type { XYCoord } from 'react-dnd';
import { useDrop } from 'react-dnd';

import { ItemTypes, LOCAL_STORAGE_AUTH_TOKEN, POST_COLORS, WS_BASE_URL } from '@/constants';
import WebSocketConnection from '../../../../websocket';
import { useUser } from '@/providers/user';
import { User } from '../../../api/users';
export interface DragItem {
  type: string;
  id: string;
  top: number;
  left: number;
}

type PostsMap = {
  [key: string]: Post;
};

const Board = ({ params: { id } }: { params: { id: string } }) => {
  const {
    state: { user },
  } = useUser();
  const { dispatch } = useBoard();

  const [posts, setPosts] = useState<PostsMap>({});
  const [highestZ, setHighestZ] = useState(getHighestZIndex(posts));
  const [colorSetting, setColorSetting] = useState(pickColor(posts));
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
    // TODO: Implement loading UI
    const response = await getBoard(id);
    // TODO: Implement redirect if 401 error
    dispatch({ type: 'set_board', payload: response });
  };

  // handleDoubleClick creates a new post
  const handleDoubleClick = (event: React.MouseEvent<HTMLDivElement>) => {
    if (event.target === event.currentTarget) {
      const { offsetX, offsetY } = event.nativeEvent;
      const randId = Math.random() * 1000000;
      const data = { id: randId, left: offsetX, top: offsetY, color: colorSetting, user } as Post;
      // TODO: Instead of add post, call ws.send with relevant data
      addPost(randId.toString(), data);
    }
  };

  const addPost = useCallback(
    (id: string, data: Partial<Post>) => {
      setPosts(update(posts, { $merge: { [id]: data } }));
    },
    [posts, setPosts]
  );

  const updatePost = useCallback(
    (id: string, data: Partial<Post>) => {
      setPosts(
        update(posts, {
          [id]: {
            $merge: data,
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
        const newLeft = Math.max(item.left + delta.x, 0);
        const newTop = Math.max(item.top + delta.y, 0);
        const zIndex = highestZ + 1;
        const data = { id: item.id, left: newLeft, top: newTop, zIndex } as Post;
        // TODO: Instead of update post, call ws.send with relevant data
        updatePost(item.id, data);
        setHighestZ(zIndex);
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
          return (
            <PostUI
              key={key}
              post={post}
              updatePost={(data: Post) => updatePost(key, data)}
              setColorSetting={setColorSetting}
              deletePost={deletePost}
            />
          );
        })}
      </div>
    </>
  );
};

// getHighestZIndex returns the highest z index by scanning posts on a board. If a board has 0
// posts, then return 0
const getHighestZIndex = (posts: { [key: string]: Post }) => {
  const zIndexValues = Object.values(posts).map(({ zIndex }) => zIndex);
  if (zIndexValues.length == 0) {
    return 0;
  }
  return Math.max(...zIndexValues);
};

// pickColor returns the first color that hasn't been picked yet among the board. If no
// colors are available, return a random color
const pickColor = (posts: { [key: string]: Post }) => {
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

export default Board;
