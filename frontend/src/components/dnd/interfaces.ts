import { Post } from '@/api/post';

export interface PostGroupDragItem {
  name: 'post_group',
  id: string;
  pos_x: number;
  pos_y: number;
  posts: Post[];
}

export interface PostDragItem {
  name: 'post',
  post: Post
}
