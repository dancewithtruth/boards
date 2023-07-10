'use client';

import { Post, PostGroupWithPosts } from '@/api/post';
import { PostAugmented } from './board';
import { BoardWithMembers, User } from '@/api';
import { Send } from '@/ws/types';
import { PostUI as PostUI } from './post';
import { CSSProperties, ChangeEvent, useEffect, useState } from 'react';
import { DragSourceMonitor, useDrag } from 'react-dnd';
import { ITEM_TYPES } from './itemTypes';
import { updatePostGroup } from '@/ws/events';

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

const PostGroup = ({ postGroup, user, board, send, setColorSetting, handleDeletePost }: PostGroupProps) => {
  const [isHovered, setIsHovered] = useState(false);
  const [isTitleFocused, setTitleFocused] = useState(false);
  const [titleValue, setTitleValue] = useState(postGroup.title);
  const { id, board_id, pos_x, pos_y, z_index } = postGroup;

  useEffect(() => {
    setTitleValue(postGroup.title);
  }, [postGroup]);

  const handleMouseEnter = () => {
    setIsHovered(true);
  };

  const handleMouseLeave = () => {
    setIsHovered(false);
  };

  const handleTitleFocus = () => {
    setTitleFocused(true);
  };

  const handleTitleBlur = () => {
    setTitleFocused(false);
    updatePostGroup({ id, board_id, title: titleValue }, send);
  };

  // handleTitleChange updates the input value
  const handleTitleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { value } = event.target;
    setTitleValue(value);
  };

  const [{ isDragging }, drag] = useDrag(() => {
    const single_post = postGroup.posts.length === 1 ? postGroup.posts[0] : null;
    return {
      type: ITEM_TYPES.POST_GROUP,
      item: { id, pos_x, pos_y, single_post, name: ITEM_TYPES.POST_GROUP },
      collect: (monitor: DragSourceMonitor) => ({
        isDragging: monitor.isDragging(),
      }),
      canDrag: !isTitleFocused,
    };
  }, [id, pos_x, pos_y, isTitleFocused]);

  if (isDragging) {
    return null;
  }

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
      <div>
        {postGroup.posts.length > 1 ? (
          <div className="flex items-center min-h-8">
            <input
              type="text"
              placeholder={'Edit name'}
              className="input ml-1 h-5"
              onFocus={handleTitleFocus}
              onBlur={handleTitleBlur}
              value={titleValue}
              onChange={handleTitleChange}
            />
          </div>
        ) : null}
        {postGroup.posts.map((post, index) => (
          <PostUI
            key={index}
            user={user}
            board={board}
            postGroup={postGroup}
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
