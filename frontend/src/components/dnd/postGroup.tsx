'use client';

import { Post, PostGroupWithPosts } from '@/api/post';
import { PostWithTypingBy } from './board';
import { BoardWithMembers, User } from '@/api';
import { Send } from '@/ws/types';
import { PostUI as PostUI } from './post';
import { CSSProperties, useState } from 'react';
import { DragSourceMonitor, useDrag } from 'react-dnd';
import { ItemTypes } from './itemTypes';

type PostGroupProps = {
  postGroup: PostGroupWithPosts;
  user: User;
  board: BoardWithMembers;
  send: Send;
  setColorSetting: (color: string) => void;
  handleDeletePost: (post: Post) => void;
};

const PostGroup = ({ postGroup, user, board, send, setColorSetting, handleDeletePost }: PostGroupProps) => {
  const [isHovered, setIsHovered] = useState(false);
  const { id, pos_x, pos_y, z_index } = postGroup;

  const handleMouseEnter = () => {
    setIsHovered(true);
  };

  const handleMouseLeave = () => {
    setIsHovered(false);
  };
  const [{ isDragging }, drag] = useDrag(
    () => ({
      type: ItemTypes.POST_GROUP,
      item: { id, pos_x, pos_y },
      collect: (monitor: DragSourceMonitor) => ({
        isDragging: monitor.isDragging(),
      }),
    }),
    [id, pos_x, pos_y]
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
          post={post as PostWithTypingBy}
          send={send}
          setColorSetting={setColorSetting}
          handleDeletePost={handleDeletePost}
        />
      ))}
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
