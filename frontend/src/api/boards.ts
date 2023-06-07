import { sendGetRequest, sendPostRequest } from './base';
import { BASE_URL, LOCAL_STORAGE_AUTH_TOKEN } from '@/constants';

export type CreateBoardParams = {
  name?: string;
  description?: string;
};

export type BoardResponse = {
  id: string;
  name: string;
  description: string;
  user_id: boolean;
  created_at: string;
  updated_at: string;
};

export type GetBoardsResponse = {
  owned: BoardWithMembers[];
  shared: BoardWithMembers[];
};

export type BoardWithMembers = {
  id: string;
  name: string;
  description: string;
  user_id: boolean;
  members: Member[];
  created_at: string;
  updated_at: string;
};

export type Member = {
  id: string;
  name: string;
  email: string;
  membership: Membership;
  created_at: string;
  updated_at: string;
};

export type Membership = {
  role: string;
  added_at: string;
  updated_at: string;
};

export type BoardsResponse = Array<BoardResponse>;

export async function createBoard(params: CreateBoardParams): Promise<BoardResponse> {
  const token = localStorage.getItem(LOCAL_STORAGE_AUTH_TOKEN) || undefined;
  const url = `${BASE_URL}/boards`;
  return sendPostRequest<BoardResponse>(url, params, token);
}

export async function getBoard(boardId: string): Promise<GetBoardsResponse> {
  const token = localStorage.getItem(LOCAL_STORAGE_AUTH_TOKEN) || undefined;
  const url = `${BASE_URL}/boards/${boardId}`;
  return sendGetRequest<GetBoardsResponse>(url, token);
}

export async function getBoards(): Promise<GetBoardsResponse> {
  const token = localStorage.getItem(LOCAL_STORAGE_AUTH_TOKEN) || undefined;
  const url = `${BASE_URL}/boards`;
  return sendGetRequest<GetBoardsResponse>(url, token);
}
