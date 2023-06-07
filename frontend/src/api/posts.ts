import { sendGetRequest, sendPostRequest } from './base';
import { BASE_URL, LOCAL_STORAGE_AUTH_TOKEN } from '@/constants';


export type PostResponse = {
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

export type ListPostsResponse = { data: Array<PostResponse>}

export async function listPosts(boardId: string): Promise<ListPostsResponse> {
  const token = localStorage.getItem(LOCAL_STORAGE_AUTH_TOKEN) || undefined;
  const url = `${BASE_URL}/posts/?boardId=${boardId}`;
  return sendGetRequest<ListPostsResponse>(url, token);
}
