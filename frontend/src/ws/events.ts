import { EVENTS } from '@/constants';
import { CreatePostParams, DeletePostParams, FocusPostParams, Send } from './types';
import { Post, PostGroupWithPosts } from '@/api/post';

export const buildMessageRequest = (event: string, params: object): string => {
  const payload = {
    event,
    params,
  };
  return JSON.stringify(payload);
};

export const authenticateUser = (jwtToken: string, send: Send) => {
  const message = buildMessageRequest(EVENTS.USER_AUTHENTICATE, { jwt: jwtToken });
  send(message);
};

export const connectBoard = (boardID: string, send: Send) => {
  const message = buildMessageRequest(EVENTS.BOARD_CONNECT, { board_id: boardID });
  send(message);
};

export const createPost = (params: CreatePostParams, send: Send) => {
  const message = buildMessageRequest(EVENTS.POST_CREATE, params);
  send(message);
};

export const updatePost = (params: Partial<Post>, send: Send) => {
  const message = buildMessageRequest(EVENTS.POST_UPDATE, params);
  send(message);
};

export const updatePostGroup = (params: Partial<PostGroupWithPosts>, send: Send) => {
  const message = buildMessageRequest(EVENTS.POST_GROUP_UPDATE, params);
  send(message);
};

export const deletePost = (params: DeletePostParams, send: Send) => {
  const message = buildMessageRequest(EVENTS.POST_DELETE, params);
  send(message);
};

export const focusPost = (params: FocusPostParams, send: Send) => {
  const message = buildMessageRequest(EVENTS.POST_FOCUS, params);
  send(message);
};

export const deletePostGroup = (id: string, send: Send) => {
  const message = buildMessageRequest(EVENTS.POST_GROUP_DELETE, { post_group_id: id });
  send(message);
};
