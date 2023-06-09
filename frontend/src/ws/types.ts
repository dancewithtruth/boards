export type CreatePostParams = {
  board_id: string;
  content: string;
  pos_x: number;
  pos_y: number;
  color: string;
  z_index: number;
};

export type DeletePostParams = {
  post_id: string;
  board_id: string;
};

export type FocusPostParams = {
  id: string;
  board_id: string;
};

export type Send = (data: string) => void;
