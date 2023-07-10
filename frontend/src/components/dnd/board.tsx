'use client';

import update from 'immutability-helper';
import { FC, useEffect } from 'react';
import { useState } from 'react';
import { useDrop } from 'react-dnd';

import type { PostGroupDragItem } from './interfaces';
import { ITEM_TYPES } from './itemTypes';
import { snapToGrid } from './snapToGrid';
import { Post, PostGroupWithPosts } from '@/api/post';
import {
  BOARD_SPACE_ADD,
  COOKIE_NAME_JWT_TOKEN,
  EVENTS,
  NAVBAR_HEIGHT_PX,
  POST_COLORS,
  DEFAULT_POST_HEIGHT,
  POST_WIDTH,
  SIDEBAR_WIDTH_PX,
  WS_URL,
  SIDEBAR_WIDTH,
  NAVBAR_HEIGHT,
} from '@/constants';
import { useWebSocket } from '@/hooks/useWebSocket';
import Cookies from 'universal-cookie';
import {
  authenticateUser,
  connectBoard,
  createPost,
  updatePostGroup,
  deletePost,
  deletePostGroup,
  detachPost as detachPostWS,
} from '@/ws/events';
import { Overlay } from '../overlay';
import { getMaxFieldFromObj } from '@/utils';
import { toast } from 'react-toastify';
import { BoardWithMembers } from '@/api/board';
import Sidebar from '../sidebar';
import { User } from '@/api';
import PostGroup from './postGroup';

export type PostAugmented = {
  typingBy: User | null;
  autoFocus: boolean;
} & Post;

export type PostGroupMap = {
  [key: string]: PostGroupWithPosts;
};
export interface BoardProps {
  board: BoardWithMembers;
  postGroups: PostGroupMap;
}

export const Board: FC<BoardProps> = ({ board, postGroups: initialPostGroups }) => {
  const TEXT_CONNECTING = 'Connecting to board';
  const TEXT_NOT_CONNECTED = 'Not connected, try refreshing';
  const [data, setData] = useState<PostGroupMap>(initialPostGroups);
  const [overlayText, setOverlayText] = useState(TEXT_CONNECTING);
  const [showOverlay, setShowOverlay] = useState(true);
  const [user, setUser] = useState<User>();
  const [connectedUsers, setConnectedUsers] = useState<User[]>([]);
  const [boardDimension, setBoardDimension] = useState({ height: 0, width: 0 });
  const [colorSetting, setColorSetting] = useState(pickColor());
  const { messages, send, readyState } = useWebSocket(WS_URL);
  const cookies = new Cookies();

  useEffect(() => {
    // Scroll to the top left portion of the page
    window.scrollTo(0, 0);
    // Usability message
    toast.info('Double-click anywhere to create a post!', { autoClose: false });
  }, []);

  // Expands the board based on post locations
  useEffect(() => {
    const { height: oldHeight, width: oldWidth } = boardDimension;
    const newWidth = getMaxFieldFromObj(data, 'pos_x') + POST_WIDTH + BOARD_SPACE_ADD;
    const heights = Object.values(data).map((postGroup) => {
      let heightOffset = postGroup.pos_y;
      postGroup.posts.forEach((post) => {
        if (post.height) {
          heightOffset += post.height + 66;
        } else {
          heightOffset += DEFAULT_POST_HEIGHT;
        }
      });
      return heightOffset;
    });
    const newHeight = Math.max(...heights) + BOARD_SPACE_ADD;
    if (oldHeight != newHeight || oldWidth !== newWidth) {
      setBoardDimension({ height: newHeight, width: newWidth });
    }
  }, [data]);

  // Handles the different connection states. Will authenticate the user if connection is established
  // or will show an error overlay if trouble connecting.
  useEffect(() => {
    if (readyState == WebSocket.OPEN) {
      setOverlayText(TEXT_CONNECTING);
      setShowOverlay(true);
      const jwtToken = cookies.get(COOKIE_NAME_JWT_TOKEN);
      authenticateUser(jwtToken, send);
    }
    if (readyState == WebSocket.CLOSED || readyState == WebSocket.CLOSING) {
      setOverlayText(TEXT_NOT_CONNECTED);
      setShowOverlay(true);
    }
  }, [readyState]);

  // Handles WebSocket events
  useEffect(() => {
    if (messages.length === 0) {
      return;
    }
    messages.forEach(({ event, result, success, error_message }) => {
      if (success) {
        switch (event) {
          case EVENTS.USER_AUTHENTICATE:
            setUser(result.user);
            connectBoard(board.id, send);
            break;
          case EVENTS.BOARD_CONNECT:
            setShowOverlay(false);
            setConnectedUsers(result.connected_users.concat([result.new_user]));
            break;
          case EVENTS.BOARD_DISCONNECT:
            const userID = result.user_id;
            const newConnectedUsers = connectedUsers.filter((user) => user.id != userID);
            setConnectedUsers(newConnectedUsers);
            break;
          case EVENTS.POST_CREATE:
            if (result.post.user_id == user?.id) {
              result.post.autoFocus = true;
            }
            // If post group does not exist, create a new post group with post child.
            if (!data[result.post_group.id]) {
              const postGroup = result.post_group;
              postGroup.posts = [result.post];
              setPostGroup(postGroup);
              break;
            }
            pushPost(result.post);
            break;
          case EVENTS.POST_UPDATE:
            if (result.updated_post.post_group_id !== result.old_post.post_group_id) {
              transferPost(result.old_post, result.updated_post);
            } else {
              setPost(result.updated_post);
            }
            break;
          case EVENTS.POST_DETACH:
            detachPost(result.old_post, result.updated_post, result.post_group);
            break;
          case EVENTS.POST_DELETE:
            unsetPost(result);
            break;
          case EVENTS.POST_GROUP_UPDATE:
            mergePostGroup({ ...result, typingBy: null });
            break;
          case EVENTS.POST_GROUP_DELETE:
            unsetPostGroup(result.id);
            break;
          case EVENTS.POST_FOCUS:
            if (result.user.id != user?.id) {
              setPost({ ...result.post, typingBy: result.user });
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
    toast.dismiss();
    if (event.target === event.currentTarget) {
      const { offsetX, offsetY } = event.nativeEvent;
      const newZIndex = getMaxFieldFromObj(data, 'z_index') + 1;
      const params = {
        board_id: board.id,
        content: '',
        pos_x: offsetX,
        pos_y: offsetY,
        color: colorSetting,
        z_index: newZIndex,
      };
      createPost(params, send);
    }
  };

  const handleDeletePost = (post: Post) => {
    // Delete post only if more than 1 posts in post group
    if (data[post.post_group_id]?.posts.length >= 2) {
      const params = { post_id: post.id, board_id: board.id };
      deletePost(params, send);
      return;
    }
    // Delete post group if only 1 post in post group
    deletePostGroup(post.post_group_id, send);
  };

  const setPostGroup = (postGroup: PostGroupWithPosts) => {
    setData((prevData) => ({
      ...prevData,
      [postGroup.id]: postGroup,
    }));
  };

  const pushPost = (post: Post) => {
    setData(
      update(data, {
        [post.post_group_id]: {
          posts: {
            $push: [post],
          },
        },
      })
    );
  };

  const mergePostGroup = (postGroup: { id: string } & Partial<PostGroupWithPosts>) => {
    const updatedData = update(data, {
      [postGroup.id]: {
        $merge: postGroup,
      },
    });
    setData(updatedData);
  };

  // setPost will attept to set post by finding the existing post in the post group. If none if found, it will
  // insert based on post order.
  const setPost = (post: Post) => {
    const indexByID = data[post.post_group_id].posts.findIndex((elem) => elem.id == post.id);
    let updatedData = data;
    if (indexByID !== -1) {
      updatedData = update(data, {
        [post.post_group_id]: {
          posts: {
            $splice: [[indexByID, 1]],
          },
        },
      });
    }
    let indexByOrder = updatedData[post.post_group_id].posts.findIndex((elem) => elem.post_order > post.post_order);
    if (indexByOrder === -1) {
      indexByOrder = updatedData[post.post_group_id].posts.length;
    }
    updatedData = update(updatedData, {
      [post.post_group_id]: {
        posts: {
          $splice: [[indexByOrder, 0, post]],
        },
      },
    });
    setData(updatedData);
  };

  const unsetPostGroup = (id: string) => {
    setData(
      update(data, {
        $unset: [id],
      })
    );
  };

  const unsetPost = (post: Post) => {
    const currentPostGroup = data[post.post_group_id];
    if (!currentPostGroup) return;
    const index = currentPostGroup.posts.findIndex((elem) => elem.id === post.id);
    if (index === -1) return;
    const updatedPostGroups = update(data, {
      [post.post_group_id]: {
        posts: {
          $splice: [[index, 1]],
        },
      },
    });
    setData(updatedPostGroups);
  };

  // transferPost is used to remove the post from the old post group and to insert the post
  // into the new post group
  const transferPost = (oldPost: Post, updatedPost: Post) => {
    const oldPostGroup = data[oldPost.post_group_id];
    const newPostGroup = data[updatedPost.post_group_id];
    const oldIndex = oldPostGroup?.posts.findIndex((elem) => elem.id === oldPost.id);
    let updatedPostGroups = data;
    if (oldIndex !== -1) {
      updatedPostGroups = update(updatedPostGroups, {
        [oldPostGroup.id]: {
          posts: {
            $splice: [[oldIndex, 1]],
          },
        },
      });
    }
    let indexByOrder = newPostGroup?.posts.findIndex((elem) => elem.post_order > updatedPost.post_order);
    if (indexByOrder === -1) {
      indexByOrder = newPostGroup?.posts.length;
    }
    updatedPostGroups = update(updatedPostGroups, {
      [newPostGroup.id]: {
        posts: {
          $splice: [[indexByOrder, 0, updatedPost]],
        },
      },
    });
    setData(updatedPostGroups);
  };

  const detachPost = (oldPost: Post, updatedPost: Post, postGroup: PostGroupWithPosts) => {
    // Add updated post to new post group
    postGroup.posts = [updatedPost];
    let updatedMap = update(data, {
      [postGroup.id]: {
        $set: postGroup,
      },
    });
    // Remove post from old post group
    const currentPostGroup = updatedMap[oldPost.post_group_id];
    if (!currentPostGroup) return;
    const index = currentPostGroup.posts.findIndex((elem) => elem.id === oldPost.id);
    if (index === -1) return;
    updatedMap = update(updatedMap, {
      [oldPost.post_group_id]: {
        posts: {
          $splice: [[index, 1]],
        },
      },
    });
    setData(updatedMap);
  };

  const [, drop] = useDrop(
    () => ({
      accept: [ITEM_TYPES.POST_GROUP, ITEM_TYPES.POST],
      drop(item: any, monitor) {
        const newZIndex = getMaxFieldFromObj(data, 'z_index') + 1;
        if (item.name === ITEM_TYPES.POST_GROUP) {
          let { id: post_group_id, pos_x, pos_y } = item.postGroup as PostGroupDragItem;
          const delta = monitor.getDifferenceFromInitialOffset() as {
            x: number;
            y: number;
          };
          console.log('dropped post group', item);
          if (!delta) {
            return undefined;
          }
          pos_x = Math.max(pos_x + delta.x, 0);
          pos_y = Math.max(pos_y + delta.y, 0);
          [pos_x, pos_y] = snapToGrid(pos_x, pos_y);
          const newParams = { id: post_group_id, board_id: board.id, z_index: newZIndex, pos_x, pos_y };
          // preemptively update post on frontend before waiting on websocket to smoothen out experience
          mergePostGroup({ id: post_group_id, z_index: newZIndex, pos_x, pos_y });
          updatePostGroup(newParams, send);
          return undefined;
        } else if (item.name === ITEM_TYPES.POST) {
          const sourceClientOffset = monitor.getSourceClientOffset();
          if (sourceClientOffset) {
            let pos_x = sourceClientOffset.x - SIDEBAR_WIDTH;
            let pos_y = sourceClientOffset.y - NAVBAR_HEIGHT;
            [pos_x, pos_y] = snapToGrid(pos_x, pos_y);
            detachPostWS({ id: item.post.id, pos_x, pos_y, z_index: newZIndex }, send);
          }
        }
      },
    }),
    [mergePostGroup]
  );

  return (
    <div className="flex">
      <Overlay show={showOverlay || !user} text={overlayText} />
      {user ? <Sidebar board={board} width={SIDEBAR_WIDTH_PX} user={user} connectedUsers={connectedUsers} /> : null}
      <div
        ref={drop}
        className="relative sketchbook-bg"
        style={{
          minHeight: `calc(100vh - ${NAVBAR_HEIGHT_PX})`,
          minWidth: `calc(100vw - ${SIDEBAR_WIDTH_PX})`,
          height: boardDimension.height,
          width: boardDimension.width,
        }}
        onDoubleClick={handleDoubleClick}
      >
        {user
          ? Object.entries(data).map(([key, postGroup]) => (
              <PostGroup
                key={key}
                postGroup={postGroup}
                user={user}
                board={board}
                send={send}
                setColorSetting={setColorSetting}
              />
            ))
          : null}
      </div>
    </div>
  );
};

Board.whyDidYouRender = true;

const pickColor = () => {
  const availableColors = Object.values(POST_COLORS);
  const randIndex = Math.floor(Math.random() * availableColors.length);
  return availableColors[randIndex];
};
