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
  user_id: string;
  members: User[];
  created_at: string;
  updated_at: string;
};


export type Membership = {
  role: string;
  added_at: string;
  updated_at: string;
};

export type BoardInvite = {
  id: string;
  board_id: string;
  sender_id: string;
  receiver_id: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export type CreateInvitesParams = {
  board_id: string;
  sender_id: string;
  invites: {receiver_id: string}[]
}

export type CreateInvitesResponse = {
  result: BoardInvite[]
}

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

export async function createInvites(params: CreateInvitesParams, token: string): Promise<CreateInvitesResponse> {
  const url = `${BASE_URL}/boards/${params.board_id}/invites`;
  return sendPostRequest<CreateInvitesResponse>(url, params, token);
}
