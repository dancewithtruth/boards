import { Post } from '@/api/post';

export interface DragItem {
  id: string;
  pos_x: number;
  pos_y: number;
  single_post: Post | null;
}
