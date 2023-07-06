'use client';

import { Post, PostGroupWithPosts } from '@/api/post';
import { PostAugmented } from './board';
import { BoardWithMembers, User } from '@/api';
import { Send } from '@/ws/types';
import { PostUI as PostUI } from './post';
import { CSSProperties, useState } from 'react';
import { DragSourceMonitor, useDrag, useDrop } from 'react-dnd';
import { ItemTypes } from './itemTypes';
import { DragItem } from './interfaces';
import { deletePostGroup, updatePost } from '@/ws/events';

type PostGroupProps = {
  postGroup: PostGroupWithPosts;
  user: User;
  board: BoardWithMembers;
  send: Send;
  setColorSetting: (color: string) => void;
  handleDeletePost: (post: Post) => void;
  setPost: (post: Post) => void;
  unsetPost: (post: Post) => void;
  unsetPostGroup: (id: string) => void;
};

const PostGroup = ({
  postGroup,
  user,
  board,
  send,
  setColorSetting,
  handleDeletePost,
  setPost,
  unsetPost,
  unsetPostGroup,
}: PostGroupProps) => {
  const [isHovered, setIsHovered] = useState(false);
  const { id, board_id, pos_x, pos_y, z_index } = postGroup;

  const handleMouseEnter = () => {
    setIsHovered(true);
  };

  const handleMouseLeave = () => {
    setIsHovered(false);
  };
  const [{ isDragging }, drag] = useDrag(() => {
    const single_post = postGroup.posts.length === 1 ? postGroup.posts[0] : null;
    return {
      type: ItemTypes.POST_GROUP,
      item: { id, pos_x, pos_y, single_post },
      collect: (monitor: DragSourceMonitor) => ({
        isDragging: monitor.isDragging(),
      }),
    };
  }, [id, pos_x, pos_y]);

  const [, drop] = useDrop(
    () => ({
      accept: ItemTypes.POST_GROUP,
      drop(item: DragItem, monitor) {
        const { id: source_post_group_id, single_post } = item;
        if (single_post) {
          const target_post_group_id = id;
          if (source_post_group_id != target_post_group_id) {
            console.log('Moving post from post group ID ', source_post_group_id, ' to ', target_post_group_id);
            updatePost({ id: single_post.id, board_id, post_group_id: target_post_group_id }, send);
            // Unset post
            deletePostGroup(single_post.post_group_id, send);
          }
        }
        return undefined;
      },
    }),
    []
  );

  return (
    <div
      ref={drag}
      className={
        postGroup.posts.length > 1
          ? 'shadow-md border border-dashed border-black backdrop-blur-sm cursor-move rounded-sm'
          : ''
      }
      style={getStyles(pos_x, pos_y, z_index, isDragging, isHovered)}
      role="DraggableGroupPost"
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      <div ref={drop}>
        {postGroup.posts.length > 1 ? (
          <div className="flex justify-between min-h-8">
            <span>{postGroup.title}</span>
          </div>
        ) : null}
        {postGroup.posts.map((post, index) => (
          <PostUI
            key={index}
            user={user}
            board={board}
            post={post as PostAugmented}
            send={send}
            setColorSetting={setColorSetting}
            handleDeletePost={handleDeletePost}
          />
        ))}
      </div>
    </div>
  );
};

function getStyles(
  pos_x: number,
  pos_y: number,
  z_index: number,
  isDragging: boolean,
  isHovered: boolean
): CSSProperties {
  const transform = `translate3d(${pos_x}px, ${pos_y}px, 0)`;
  return {
    position: 'absolute',
    transform,
    WebkitTransform: transform,
    // IE fallback: hide the real node using CSS when dragging
    // because IE will ignore our custom "empty image" drag preview.
    opacity: isDragging ? 0 : 1,
    height: isDragging ? 0 : '',
    zIndex: isHovered ? '10000' : z_index,
  };
}

export default PostGroup;
