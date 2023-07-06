'use client';

import { Post, PostGroupWithPosts } from '@/api/post';
import { PostAugmented } from './board';
import { BoardWithMembers, User } from '@/api';
import { Send } from '@/ws/types';
import { PostUI as PostUI } from './post';
import { CSSProperties, ChangeEvent, useState } from 'react';
import { DragSourceMonitor, useDrag, useDrop } from 'react-dnd';
import { ITEM_TYPES } from './itemTypes';
import { PostGroupDragItem } from './interfaces';
import { deletePostGroup, updatePost, updatePostGroup } from '@/ws/events';

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
  const [isTitleFocused, setTitleFocused] = useState(false);
  const [titleValue, setTitleValue] = useState(postGroup.title);
  const { id, board_id, pos_x, pos_y, z_index } = postGroup;

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

  const [{ isOver }, drop] = useDrop(
    () => ({
      accept: [ITEM_TYPES.POST_GROUP, ITEM_TYPES.POST],
      drop(item: any, monitor) {
        if (item.name == ITEM_TYPES.POST_GROUP) {
          const { id: source_post_group_id, single_post } = item as PostGroupDragItem;
          if (single_post) {
            const target_post_group_id = id;
            if (source_post_group_id != target_post_group_id) {
              console.log('Moving post from post group ID ', source_post_group_id, ' to ', target_post_group_id);
              updatePost({ id: single_post.id, post_group_id: target_post_group_id }, send);
              // Unset post
              deletePostGroup(single_post.post_group_id, send);
            }
          }
        } else if (item.name == ITEM_TYPES.POST) {
          updatePost({ ...item.post, post_group_id: id }, send);
        }

        return undefined;
      },
      collect: (monitor) => ({
        isOver: monitor.isOver(),
      }),
    }),
    []
  );

  if (isDragging) {
    return null;
  }

  return (
    <div
      ref={(node) => drag(drop(node))}
      className={
        postGroup.posts.length > 1
          ? 'shadow-md border border-dashed border-black backdrop-blur-sm cursor-move rounded-sm'
          : ''
      }
      style={getStyles(pos_x, pos_y, z_index, isDragging, isHovered, isOver)}
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
  isHovered: boolean,
  isOver: boolean
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
    border: isOver ? '2px solid black' : '',
  };
}

export default PostGroup;
