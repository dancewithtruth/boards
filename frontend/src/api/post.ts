import { sendGetRequest } from './index';
import { BASE_URL } from '@/constants';

export type PostGroupWithPosts = {
  id: string
  board_id: string
  title: string
  pos_x: number
  pos_y: number
  z_index: number
  posts: Post[]
  created_at: string
  updated_at: string
}

export type Post = {
  id: string;
  board_id: string;
  user_id: string;
  content: string;
  pos_x: number;
  pos_y: number;
  color: string;
  height: number;
  z_index: number;
  created_at: string;
  updated_at: string;
};

export type ListPostGroupsResponse = { result: Array<PostGroupWithPosts> };

export async function listPostGroups(boardID: string, token: string): Promise<ListPostGroupsResponse> {
  const url = `${BASE_URL}/post-groups/?boardID=${boardID}`;
  return sendGetRequest<ListPostGroupsResponse>(url, token);
}
