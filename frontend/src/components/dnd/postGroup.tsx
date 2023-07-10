'use client';

import { Post, PostGroupWithPosts } from '@/api/post';
import { PostAugmented } from './board';
import { BoardWithMembers, User } from '@/api';
import { Send } from '@/ws/types';
import { CSSProperties, ChangeEvent, memo, useEffect, useState } from 'react';
import { DragSourceMonitor, useDrag } from 'react-dnd';
import { ITEM_TYPES } from './itemTypes';
import { updatePostGroup } from '@/ws/events';
import PostUI from './post';
import { getEmptyImage } from 'react-dnd-html5-backend';

type PostGroupProps = {
  postGroup: PostGroupWithPosts;
  user: User;
  board: BoardWithMembers;
  send: Send;
  setColorSetting: (color: string) => void;
};

const PostGroup = ({ postGroup, user, board, send, setColorSetting }: PostGroupProps) => {
  const { id, board_id, title, posts, pos_x, pos_y, z_index } = postGroup;
  const [isHovered, setIsHovered] = useState(false);
  const [isTitleFocused, setTitleFocused] = useState(false);
  const [titleValue, setTitleValue] = useState(title);
  console.log('Postgroup: re-render');

  useEffect(() => {
    setTitleValue(title);
  }, [title]);

  const handleMouseEnter = () => {
    setIsHovered(true);
  };

  const handleMouseLeave = () => {
    setIsHovered(false);
  };

  // handleTitleChange updates the input value
  const handleTitleChange = (event: ChangeEvent<HTMLInputElement>) => {
    const { value } = event.target;
    setTitleValue(value);
  };

  const handleTitleFocus = () => {
    setTitleFocused(true);
  };

  const handleTitleBlur = () => {
    setTitleFocused(false);
    updatePostGroup({ id, board_id, title: titleValue }, send);
  };

  const [{ isDragging }, drag] = useDrag(() => {
    return {
      type: ITEM_TYPES.POST_GROUP,
      item: { postGroup, name: ITEM_TYPES.POST_GROUP },
      collect: (monitor: DragSourceMonitor) => ({
        isDragging: monitor.isDragging(),
      }),
      canDrag: !isTitleFocused,
    };
  }, [isTitleFocused, postGroup]);

  return (
    <div
      ref={drag}
      className={
        posts.length > 1 ? 'shadow-md border border-dashed border-black backdrop-blur-sm cursor-move rounded-sm' : ''
      }
      style={getStyles(pos_x, pos_y, z_index, isDragging, isHovered)}
      role="DraggableGroupPost"
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      {posts.length > 1 ? (
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
      {posts.map((post, index) => (
        <PostUI
          key={index}
          user={user}
          board={board}
          postGroup={postGroup}
          post={post as PostAugmented}
          send={send}
          setColorSetting={setColorSetting}
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

PostGroup.displayName = 'PostGroup';

export default memo(PostGroup);
