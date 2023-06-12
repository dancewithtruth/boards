import { User, sendGetRequest, sendPostRequest } from './index';
import { BASE_URL } from '@/constants';

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
  members: User[];
  created_at: string;
  updated_at: string;
};

export type Membership = {
  role: string;
  added_at: string;
  updated_at: string;
};

export type BoardsResponse = Array<BoardResponse>;

export async function createBoard(params: CreateBoardParams, token: string): Promise<BoardResponse> {
  const url = `${BASE_URL}/boards`;
  return sendPostRequest<BoardResponse>(url, params, token);
}

export async function getBoard(boardID: string, token: string): Promise<BoardWithMembers> {
  const url = `${BASE_URL}/boards/${boardID}`;
  return sendGetRequest<BoardWithMembers>(url, token);
}

export async function getBoards(token: string): Promise<GetBoardsResponse> {
  const url = `${BASE_URL}/boards`;
  return sendGetRequest<GetBoardsResponse>(url, token);
}
