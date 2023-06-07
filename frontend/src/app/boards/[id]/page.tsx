'use client';

import Head from 'next/head';
import { useEffect, useState, useCallback } from 'react';
import { useBoard } from '@/providers/board';
import { getBoard } from '../../../api/boards';
import PostUI, { Post } from '@/components/post';
import update from 'immutability-helper';
import type { XYCoord } from 'react-dnd';
import { useDrop } from 'react-dnd';

import { BASE_URL, ItemTypes, LOCAL_STORAGE_AUTH_TOKEN, POST_COLORS, WS_URL } from '@/constants';
import { useUser } from '@/providers/user';
import { PostResponse, listPosts } from '@/api/posts';
export interface DragItem {
  type: string;
  id: string;
  top: number;
  left: number;
}

type PostsMap = {
  [key: string]: Post;
};

const Board = ({ params: { id: boardId } }: { params: { id: string } }) => {
  const {
    state: { user },
  } = useUser();
  const { dispatch } = useBoard();
  const [ws, setWs] = useState<WebSocket>();
  const [posts, setPosts] = useState<PostsMap>({});
  const [highestZ, setHighestZ] = useState(0);
  const [colorSetting, setColorSetting] = useState(pickColor(posts));
  useEffect(() => {
    fetchData();
    // Establish ws connection
    // TODO: Handle if browser does not have websocket API
    if (window.WebSocket) {
      console.log('has websocket');
      const ws = new WebSocket(WS_URL);
      ws.addEventListener('open', (event) => {
        const jwtToken = localStorage.getItem(LOCAL_STORAGE_AUTH_TOKEN);
        const eventUserAuthenticate = {
          event: 'user.authenticate',
          params: { jwt: jwtToken },
        };
        ws.send(JSON.stringify(eventUserAuthenticate));
      });
      ws.addEventListener('message', (event) => {
        console.log('message', event);
        const data = JSON.parse(event.data);
        if (data.event == 'user.authenticate') {
          const eventBoardConnect = {
            event: 'board.connect',
            params: { board_id: boardId },
          };
          ws.send(JSON.stringify(eventBoardConnect));
        }

        if (data.event == 'board.connect') {
          // Handle connected users and such
        }

        if (data.event == 'post.create') {
          //TODO: Show user is online
          const { id: postId, pos_x: left, pos_y: top, user_id: userId, content, color, z_index: zIndex } = data.result;
          console.log(data.result);
          const createdPost: Post = {
            id: postId,
            left,
            top,
            content,
            color,
            zIndex,
            userId,
          };
          addPost(postId, createdPost);
        }
      });
      ws.addEventListener('error', (event) => {
        console.log('error', event);
      });
      ws.addEventListener('close', (event) => {
        console.log('close', event);
      });
      setWs(ws);
    }
  }, []);

  const fetchData = async () => {
    // TODO: Implement loading UI
    const response = await getBoard(boardId);
    const posts = await listPosts(boardId);
    const postsMap = posts.data.reduce((map: PostsMap, post: PostResponse) => {
      const { id: postId, pos_x: left, pos_y: top, user_id: userId, content, color, z_index: zIndex } = post;
      const postFormatted: Post = {
        id: postId,
        left,
        top,
        content,
        color,
        zIndex,
        userId,
      };
      map[post.id] = postFormatted;
      return map;
    }, {});
    // TODO: Implement redirect if 401 error
    dispatch({ type: 'set_board', payload: response });
    setHighestZ(getHighestZIndex(postsMap));
    setPosts(postsMap);
  };

  // handleDoubleClick creates a new post
  const handleDoubleClick = (event: React.MouseEvent<HTMLDivElement>) => {
    if (event.target === event.currentTarget) {
      const { offsetX, offsetY } = event.nativeEvent;
      const params = {
        board_id: boardId,
        content: '',
        pos_x: offsetX,
        pos_y: offsetY,
        color: colorSetting,
        z_index: highestZ + 1,
      };
      const msgPostCreate = {
        event: 'post.create',
        params,
      };
      ws?.send(JSON.stringify(msgPostCreate));
      setHighestZ(highestZ + 1);
      // const randId = Math.random() * 1000000;
      // const data = { id: randId, left: offsetX, top: offsetY, color: colorSetting, user } as Post;
      // // TODO: Instead of add post, call ws.send with relevant data
      // addPost(randId.toString(), data);
    }
  };

  const addPost = (id: string, data: Partial<Post>) => {
    // console.log('first we get posts', posts);
    // const updatedPosts = update(posts, { $merge: { [id]: data } });
    // console.log('then we get updatedposts', updatedPosts);
    // setPosts(updatedPosts);

    setPosts((prevPosts) => {
      const updatedPosts = update(prevPosts, { $merge: { [id]: data } });
      return updatedPosts;
    });
  };

  const updatePost = (id: string, data: Partial<Post>) => {
    setPosts(
      update(posts, {
        [id]: {
          $merge: data,
        },
      })
    );
  };

  const deletePost = (id: string) => {
    setPosts(
      update(posts, {
        $unset: [id],
      })
    );
  };

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
        {' '}
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
