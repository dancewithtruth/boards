import { User, sendGetRequest, sendPatchRequest, sendPostRequest } from './index';
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
  members: UserWithMembership[];
  created_at: string;
  updated_at: string;
};

export type UserWithMembership = {
  membership: Membership;
} & User;

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
};

export type InviteWithReceiver = {
  id: string;
  board_id: string;
  sender_id: string;
  receiver: User;
  status: string;
  created_at: string;
  updated_at: string;
};

export type InviteWithBoardAndSender = {
  id: string;
  board: BoardResponse;
  sender: User;
  receiver_id: string;
  status: string;
  created_at: string;
  updated_at: string;
};

export type CreateInvitesParams = {
  board_id: string;
  sender_id: string;
  invites: { receiver_id: string }[];
};

export type CreateInvitesResponse = {
  result: BoardInvite[];
};

export type BoardsResponse = Array<BoardResponse>;

export type ListInvitesByBoardResponse = {
  result: InviteWithReceiver[];
};

export type ListInvitesByReceiverResponse = {
  result: InviteWithBoardAndSender[];
};

export type UpdateInviteParams = {
  status: string;
};

export type VerifyEmailResponse = {
  message: string;
};

export async function createBoard(params: CreateBoardParams, token: string): Promise<BoardResponse> {
  const url = `${BASE_URL}/boards`;
  return sendPostRequest<BoardResponse>(url, params, token);
}

export async function getBoard(boardID: string, token: string): Promise<BoardWithMembers> {
  const url = `${BASE_URL}/boards/${boardID}`;
  return sendGetRequest<BoardWithMembers>(url, token);
}

export async function listBoards(token: string): Promise<GetBoardsResponse> {
  const url = `${BASE_URL}/boards`;
  return sendGetRequest<GetBoardsResponse>(url, token);
}

export async function createInvites(params: CreateInvitesParams, token: string): Promise<CreateInvitesResponse> {
  const url = `${BASE_URL}/boards/${params.board_id}/invites`;
  return sendPostRequest<CreateInvitesResponse>(url, params, token);
}

export async function listInvitesByBoard(
  boardID: string,
  token: string,
  statusFilter?: string
): Promise<ListInvitesByBoardResponse> {
  let url = `${BASE_URL}/boards/${boardID}/invites`;
  if (statusFilter) {
    url += `?status=${statusFilter}`;
  }
  return sendGetRequest<ListInvitesByBoardResponse>(url, token);
}

export async function listInvitesByReceiver(
  token: string,
  statusFilter?: string
): Promise<ListInvitesByReceiverResponse> {
  let url = `${BASE_URL}/invites`;
  if (statusFilter) {
    url += `?status=${statusFilter}`;
  }
  return sendGetRequest<ListInvitesByReceiverResponse>(url, token);
}

export async function updateInvite(id: string, params: UpdateInviteParams, token: string) {
  const url = `${BASE_URL}/invites/${id}`;
  return sendPatchRequest(url, params, token);
}

export async function verifyEmail(code: string, token: string): Promise<VerifyEmailResponse> {
  const url = `${BASE_URL}/users/verify-email`;
  const params = {code: code}
  return sendPostRequest<VerifyEmailResponse>(url, params, token);
}
