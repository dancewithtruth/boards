import { sendGetRequest, sendPostRequest } from './base';
import { API_BASE_URL, LOCAL_STORAGE_AUTH_TOKEN } from '../../constants';

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

export async function createBoard(params: CreateBoardParams): Promise<BoardResponse> {
  const jwtToken = localStorage.getItem(LOCAL_STORAGE_AUTH_TOKEN) || undefined;
  console.log(jwtToken)
  const url = `${API_BASE_URL}/boards`;
  return sendPostRequest<BoardResponse>(url, params, jwtToken);
}
