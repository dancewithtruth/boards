'use client';

import update from 'immutability-helper';
import { FC, useEffect } from 'react';
import { useState } from 'react';
import { useDrop } from 'react-dnd';

import { DraggablePost } from './draggablePost';
import type { DragItem } from './interfaces';
import { ItemTypes } from './itemTypes';
import { snapToGrid as doSnapToGrid } from './snapToGrid';
import { Post } from '@/api/post';
import {
  BOARD_SPACE_ADD,
  COOKIE_NAME_JWT_TOKEN,
  EVENT_BOARD_CONNECT,
  EVENT_BOARD_DISCONNECT,
  EVENT_POST_CREATE,
  EVENT_POST_DELETE,
  EVENT_POST_FOCUS,
  EVENT_POST_UPDATE,
  EVENT_USER_AUTHENTICATE,
  NAVBAR_HEIGHT,
  POST_COLORS,
  POST_HEIGHT,
  POST_WIDTH,
  SIDEBAR_WIDTH,
  WS_URL,
} from '@/constants';
import { useWebSocket } from '@/hooks/useWebSocket';
import Cookies from 'universal-cookie';
import {
  authenticateUser as authenticateUserWS,
  connectBoard as connectBoardWS,
  createPost as createPostWS,
  updatePost as updatePostWS,
} from '@/ws/events';
import { Overlay } from '../overlay';
import { getMaxFieldFromObj } from '@/utils';
import { toast } from 'react-toastify';
import { BoardWithMembers } from '@/api/board';
import Sidebar from '../sidebar';
import { User } from '@/api';

export type PostUI = {
  typingBy: User | null;
  autoFocus: boolean;
} & Post;

export type PostMap = {
  [key: string]: Partial<PostUI>;
};
export interface BoardProps {
  snapToGrid: boolean;
  board: BoardWithMembers;
  posts: PostMap;
}

export const Board: FC<BoardProps> = ({ board, snapToGrid, posts: initialPosts }) => {
  const TEXT_CONNECTING = 'Connecting to board';
  const TEXT_NOT_CONNECTED = 'Not connected, try refreshing';
  const [posts, setPosts] = useState<PostMap>(initialPosts);
  const [overlayText, setOverlayText] = useState(TEXT_CONNECTING);
  const [showOverlay, setShowOverlay] = useState(true);
  const [user, setUser] = useState<User>();
  const [connectedUsers, setConnectedUsers] = useState<User[]>([]);
  const [boardDimension, setBoardDimension] = useState({ height: 0, width: 0 });
  const [highestZ, setHighestZ] = useState(getMaxFieldFromObj(initialPosts, 'z_index'));
  const [colorSetting, setColorSetting] = useState(pickColor(posts));
  const { messages, error, send, readyState } = useWebSocket(WS_URL);
  const cookies = new Cookies();

  useEffect(() => {
    // Scroll to the top left portion of the page
    window.scrollTo(0, 0);
  }, []);

  // Expands the board based on post locations
  useEffect(() => {
    const newWidth = getMaxFieldFromObj(posts, 'pos_x') + POST_WIDTH + BOARD_SPACE_ADD;
    const newHeight = getMaxFieldFromObj(posts, 'pos_y') + POST_HEIGHT + BOARD_SPACE_ADD;
    setBoardDimension({ height: newHeight, width: newWidth });
  }, [posts]);

  // Handles the different connection states. Will authenticate the user if connection is established
  // or will show an error overlay if trouble connecting.
  useEffect(() => {
    if (readyState == WebSocket.OPEN) {
      setOverlayText(TEXT_CONNECTING);
      setShowOverlay(true);
      const jwtToken = cookies.get(COOKIE_NAME_JWT_TOKEN);
      authenticateUserWS(jwtToken, send);
    }
    if (readyState == WebSocket.CLOSED || readyState == WebSocket.CLOSING) {
      setOverlayText(TEXT_NOT_CONNECTED);
      setShowOverlay(true);
    }
  }, [readyState]);

  // Handles all the different post events
  useEffect(() => {
    if (messages.length === 0) {
      return;
    }
    messages.forEach(({ event, result, success, error_message }) => {
      try {
        switch (event) {
          case EVENT_USER_AUTHENTICATE:
            setUser(result.user);
            connectBoardWS(board.id, send);
            break;
          case EVENT_BOARD_CONNECT:
            if (success) {
              setShowOverlay(false);
              setConnectedUsers(result.connected_users.concat([result.new_user]));
            } else {
              toast.error(error_message);
            }
            break;
          case EVENT_BOARD_DISCONNECT:
            if (success) {
              const userID = result.user_id;
              const newConnectedUsers = connectedUsers.filter((user) => user.id != userID);
              setConnectedUsers(newConnectedUsers);
            }
            break;
          case EVENT_POST_CREATE:
            if (success) {
              if (result.user_id == user?.id) {
                result.autoFocus = true;
              }
              addPost(result);
            } else {
              toast.error(error_message);
            }
            break;
          case EVENT_POST_UPDATE:
            if (success) {
              updatePost({ ...result, typingBy: null });
            } else {
              toast.error(error_message);
            }
            break;
          case EVENT_POST_DELETE:
            if (success) {
              deletePost(result.post_id);
            } else {
              toast.error(error_message);
            }
            break;
          case EVENT_POST_FOCUS:
            if (success) {
              if (result.user.id != user?.id) {
                updatePost({ id: result.id, typingBy: result.user });
              }
            } else {
              toast.error(error_message);
            }
            break;

          default:
            break;
        }
      } catch (e) {
        console.log(e);
      }
    });
  }, [messages]);

  // handleDoubleClick creates a new post
  const handleDoubleClick = (event: React.MouseEvent<HTMLDivElement>) => {
    if (event.target === event.currentTarget) {
      const { offsetX, offsetY } = event.nativeEvent;
      const newZIndex = highestZ + 1;
      const params = {
        board_id: board.id,
        content: '',
        pos_x: offsetX,
        pos_y: offsetY,
        color: colorSetting,
        z_index: highestZ + 1,
      };
      createPostWS(params, send);
      setHighestZ(newZIndex);
    }
  };

  const addPost = (post: PostUI) => {
    setPosts(
      update(posts, {
        [post.id]: {
          $set: post,
        },
      })
    );
  };

  const updatePost = (post: { id: string } & Partial<PostUI>) => {
    setPosts(
      update(posts, {
        [post.id]: {
          $merge: post,
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

  const sendUpdatePost = (post: Partial<PostUI>) => {
    updatePostWS(post, send);
  };

  const [, drop] = useDrop(
    () => ({
      accept: ItemTypes.POST,
      drop(item: DragItem, monitor) {
        const delta = monitor.getDifferenceFromInitialOffset() as {
          x: number;
          y: number;
        };

        let pos_x = Math.max(item.pos_x + delta.x, 0);
        let pos_y = Math.max(item.pos_y + delta.y, 0);
        if (snapToGrid) {
          [pos_x, pos_y] = doSnapToGrid(pos_x, pos_y);
        }
        const newZIndex = getMaxFieldFromObj(posts, 'z_index') + 1;
        const newParams = { id: item.id, board_id: board.id, z_index: newZIndex, pos_x, pos_y };
        // pre-emptively update post on frontend before waiting on websocket to smoothen out experience
        updatePost({ id: item.id, z_index: newZIndex, pos_x, pos_y });
        sendUpdatePost(newParams);
        return undefined;
      },
    }),
    [updatePost]
  );

  return (
    <div className="flex">
      <Overlay show={showOverlay || !user} text={overlayText} />
      {user ? <Sidebar board={board} width={SIDEBAR_WIDTH} user={user} connectedUsers={connectedUsers} /> : null}
      <div
        ref={drop}
        className="relative sketchbook-bg"
        style={{
          minHeight: `calc(100vh - ${NAVBAR_HEIGHT})`,
          minWidth: `calc(100vw - ${SIDEBAR_WIDTH})`,
          height: boardDimension.height,
          width: boardDimension.width,
        }}
        onDoubleClick={handleDoubleClick}
      >
        {user
          ? Object.keys(posts).map((key) => (
              <DraggablePost
                key={key}
                user={user}
                board={board}
                {...(posts[key] as PostUI)}
                send={send}
                setColorSetting={setColorSetting}
              />
            ))
          : null}
      </div>
    </div>
  );
};

// pickColor returns the first color that hasn't been picked yet among the board. If no
// colors are available, return a random color
const pickColor = (posts: PostMap) => {
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
