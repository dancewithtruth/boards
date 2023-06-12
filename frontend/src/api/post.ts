import { sendGetRequest } from './index';
import { BASE_URL } from '@/constants';

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

export type ListPostsResponse = { data: Array<Post> };

export async function listPosts(boardID: string, token: string): Promise<ListPostsResponse> {
  const url = `${BASE_URL}/posts/?boardID=${boardID}`;
  return sendGetRequest<ListPostsResponse>(url, token);
}
