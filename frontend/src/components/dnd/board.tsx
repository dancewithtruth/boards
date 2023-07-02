'use client';

import update from 'immutability-helper';
import { FC, useEffect } from 'react';
import { useState } from 'react';
import { useDrop } from 'react-dnd';

import type { DragItem } from './interfaces';
import { ItemTypes } from './itemTypes';
import { snapToGrid as doSnapToGrid } from './snapToGrid';
import { Post, PostGroupWithPosts } from '@/api/post';
import {
  BOARD_SPACE_ADD,
  COOKIE_NAME_JWT_TOKEN,
  EVENT_BOARD_CONNECT,
  EVENT_BOARD_DISCONNECT,
  EVENT_POST_CREATE,
  EVENT_POST_DELETE,
  EVENT_POST_FOCUS,
  EVENT_POST_GROUP_DELETE,
  EVENT_POST_GROUP_UPDATE,
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
  updatePostGroup as updatePostGroupWS,
  deletePost as deletePostWS,
  deletePostGroup as deletePostGroupWS,
} from '@/ws/events';
import { Overlay } from '../overlay';
import { getMaxFieldFromObj } from '@/utils';
import { toast } from 'react-toastify';
import { BoardWithMembers } from '@/api/board';
import Sidebar from '../sidebar';
import { User } from '@/api';
import PostGroup from './postGroup';

export type PostWithTypingBy = {
  typingBy: User | null;
  autoFocus: boolean;
} & Post;

export type PostGroupMap = {
  [key: string]: PostGroupWithPosts;
};
export interface BoardProps {
  snapToGrid: boolean;
  board: BoardWithMembers;
  postGroups: PostGroupMap;
}

export const Board: FC<BoardProps> = ({ board, snapToGrid, postGroups: initialPostGroups }) => {
  const TEXT_CONNECTING = 'Connecting to board';
  const TEXT_NOT_CONNECTED = 'Not connected, try refreshing';
  const [postGroups, setPostGroups] = useState<PostGroupMap>(initialPostGroups);
  const [overlayText, setOverlayText] = useState(TEXT_CONNECTING);
  const [showOverlay, setShowOverlay] = useState(true);
  const [user, setUser] = useState<User>();
  const [connectedUsers, setConnectedUsers] = useState<User[]>([]);
  const [boardDimension, setBoardDimension] = useState({ height: 0, width: 0 });
  const [highestZ, setHighestZ] = useState(getMaxFieldFromObj(initialPostGroups, 'z_index'));
  const [colorSetting, setColorSetting] = useState(pickColor(postGroups));
  const { messages, error, send, readyState } = useWebSocket(WS_URL);
  const cookies = new Cookies();

  useEffect(() => {
    // Scroll to the top left portion of the page
    window.scrollTo(0, 0);
  }, []);

  // Expands the board based on post locations
  useEffect(() => {
    const newWidth = getMaxFieldFromObj(postGroups, 'pos_x') + POST_WIDTH + BOARD_SPACE_ADD;
    const newHeight = getMaxFieldFromObj(postGroups, 'pos_y') + POST_HEIGHT + BOARD_SPACE_ADD;
    setBoardDimension({ height: newHeight, width: newWidth });
  }, [postGroups]);

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
      if (success) {
        switch (event) {
          case EVENT_USER_AUTHENTICATE:
            setUser(result.user);
            connectBoardWS(board.id, send);
            break;
          case EVENT_BOARD_CONNECT:
            setShowOverlay(false);
            setConnectedUsers(result.connected_users.concat([result.new_user]));
            break;
          case EVENT_BOARD_DISCONNECT:
            const userID = result.user_id;
            const newConnectedUsers = connectedUsers.filter((user) => user.id != userID);
            setConnectedUsers(newConnectedUsers);
            break;
          case EVENT_POST_CREATE:
            if (result.post.user_id == user?.id) {
              result.post.autoFocus = true;
            }
            // If post group does not exist, create a new post group with post child.
            if (!postGroups[result.post_group.id]) {
              const postGroup = result.post_group;
              postGroup.posts = [result.post];
              addPostGroup(postGroup);
              break;
            }
            addPost(result.post);
            break;
          case EVENT_POST_UPDATE:
            updatePost(result);
            break;
          case EVENT_POST_DELETE:
            deletePost(result);
            break;
          case EVENT_POST_GROUP_UPDATE:
            updatePostGroup({ ...result, typingBy: null });
            break;
          case EVENT_POST_GROUP_DELETE:
            deletePostGroup(result.id);
            break;
          case EVENT_POST_FOCUS:
            if (result.user.id != user?.id) {
              updatePost({ ...result.post, typingBy: result.user });
            }
            break;
          default:
            break;
        }
      }
      toast.error(error_message);
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

  const handleDeletePost = (post: Post) => {
    if (postGroups[post.post_group_id]?.posts.length >= 2) {
      // Delete post
      const params = { post_id: post.id, board_id: board.id };
      deletePostWS(params, send);
      return;
    }
    deletePostGroupWS(post.post_group_id, send);
  };

  const addPostGroup = (postGroup: PostGroupWithPosts) => {
    setPostGroups(
      update(postGroups, {
        [postGroup.id]: {
          $set: postGroup,
        },
      })
    );
  };

  const addPost = (post: Post) => {
    setPostGroups(
      update(postGroups, {
        [post.post_group_id]: {
          posts: {
            $push: [post],
          },
        },
      })
    );
  };

  const updatePostGroup = (postGroup: { id: string } & Partial<PostGroupWithPosts>) => {
    setPostGroups(
      update(postGroups, {
        [postGroup.id]: {
          $merge: postGroup,
        },
      })
    );
  };

  const updatePost = (post: Post) => {
    const index = postGroups[post.post_group_id].posts.findIndex((elem) => elem.id == post.id);
    setPostGroups(
      update(postGroups, {
        [post.post_group_id]: {
          posts: {
            $splice: [[index, 1, post]],
          },
        },
      })
    );
  };

  const deletePostGroup = (id: string) => {
    setPostGroups(
      update(postGroups, {
        $unset: [id],
      })
    );
  };

  const deletePost = (post: Post) => {
    const index = postGroups[post.post_group_id].posts.findIndex((elem) => elem.id == post.id);
    setPostGroups(
      update(postGroups, {
        [post.post_group_id]: {
          posts: {
            $splice: [[index, 1]],
          },
        },
      })
    );
  };

  const [, drop] = useDrop(
    () => ({
      accept: ItemTypes.POST_GROUP,
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
        const newZIndex = getMaxFieldFromObj(postGroups, 'z_index') + 1;
        const newParams = { id: item.id, board_id: board.id, z_index: newZIndex, pos_x, pos_y };
        // pre-emptively update post on frontend before waiting on websocket to smoothen out experience
        updatePostGroup({ id: item.id, z_index: newZIndex, pos_x, pos_y });
        updatePostGroupWS(newParams, send);
        return undefined;
      },
    }),
    [updatePostGroup]
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
          ? Object.entries(postGroups).map(([key, postGroup]) => (
              <PostGroup
                key={key}
                user={user}
                board={board}
                postGroup={postGroup}
                send={send}
                setColorSetting={setColorSetting}
                handleDeletePost={handleDeletePost}
              />
            ))
          : null}
      </div>
    </div>
  );
};

const pickColor = (postGroupMap: PostGroupMap) => {
  const availableColors = Object.values(POST_COLORS);
  const randIndex = Math.floor(Math.random() * availableColors.length);
  return availableColors[randIndex];
};
