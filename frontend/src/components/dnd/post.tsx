'use client';

import { BoardWithMembers, User } from '@/api';
import { POST_COLORS, DEFAULT_POST_HEIGHT, POST_WIDTH } from '@/constants';
import { displayColor } from '@/utils';
import { deletePost, deletePostGroup, focusPost, updatePost } from '@/ws/events';
import { Send } from '@/ws/types';
import { CSSProperties, ChangeEvent, FC, useEffect, useRef, useState } from 'react';
import { memo } from 'react';
import { FaRegTrashAlt } from 'react-icons/fa';
import Avatar from '../avatar';
import { PostAugmented } from './board';
import { DragSourceMonitor, DropTargetMonitor, XYCoord, useDrag, useDrop } from 'react-dnd';
import { ITEM_TYPES } from './itemTypes';
import { PostGroupWithPosts } from '@/api/post';
import { PostGroupDragItem } from './interfaces';

type PostProps = {
  user: User;
  board: BoardWithMembers;
  postGroup: PostGroupWithPosts;
  post: PostAugmented;
  send: Send;
  setColorSetting: (color: string) => void;
};

const PostUI: FC<PostProps> = memo(function Post({ user, board, postGroup, post, send, setColorSetting }) {
  const { id, user_id, color, content, height, typingBy, autoFocus, post_order } = post;
  const [textareaValue, setTextareaValue] = useState(content);
  const [textareaHeight, setTextareaHeight] = useState(height);
  const [isHovered, setIsHovered] = useState(false);
  const [isFocused, setIsFocused] = useState(false);
  const [borderStyle, setBorderStyle] = useState<any>();
  const ref = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const allMembers = board.members;
  const authorName = getName(user_id, allMembers) || 'Unknown';
  const posts = postGroup.posts;
  const hasSiblings = posts.length > 1;

  const [{ isDragging }, drag] = useDrag(
    () => ({
      type: ITEM_TYPES.POST,
      item: { name: ITEM_TYPES.POST, post },
      collect: (monitor: DragSourceMonitor) => ({
        isDragging: monitor.isDragging(),
      }),
      canDrag: !typingBy && !isFocused && hasSiblings,
    }),
    [post]
  );

  const [{ isOver }, drop] = useDrop(
    () => ({
      accept: [ITEM_TYPES.POST_GROUP, ITEM_TYPES.POST],
      drop(item: any, monitor) {
        if (!ref.current) {
          return;
        }
        const hoverAbove = isHoverAbove(monitor, ref.current);
        let post_order = 0.0;
        const index = posts.findIndex((post) => post.id == id);

        if (index === 0 && hoverAbove) {
          post_order = post.post_order / 2;
        } else if (index === posts.length - 1 && !hoverAbove) {
          post_order = post.post_order + 1;
        } else {
          const currentOrder = posts[index].post_order;
          if (hoverAbove) {
            const aboveOrder = posts[index - 1].post_order;
            post_order = (aboveOrder + currentOrder) / 2;
          } else {
            const belowOrder = posts[index + 1].post_order;
            post_order = (belowOrder + currentOrder) / 2;
          }
        }

        if (item.name == ITEM_TYPES.POST_GROUP) {
          const { posts } = item.postGroup as PostGroupDragItem;
          const single_post = posts.length === 1 ? posts[0] : null;
          if (single_post) {
            const target_post_group_id = post.post_group_id;
            updatePost({ id: single_post.id, post_group_id: target_post_group_id, post_order }, send);
            deletePostGroup(single_post.post_group_id, send);
          }
        } else if (item.name == ITEM_TYPES.POST) {
          updatePost({ ...item.post, post_group_id: post.post_group_id, post_order }, send);
        }
        return undefined;
      },
      collect: (monitor) => ({
        isOver: monitor.isOver(),
      }),
      hover: (item: any, monitor) => {
        if (!ref.current) {
          return;
        }
        if (isHoverAbove(monitor, ref.current)) {
          setBorderStyle({ borderTop: '2px solid black' });
        } else {
          setBorderStyle({ borderBottom: '2px solid black' });
        }
      },
    }),
    [posts]
  );

  useEffect(() => {
    setTextareaValue(content);
    setTextareaHeight(height);
  }, [content, height]);

  useEffect(() => {
    if (!isOver) {
      setBorderStyle({ border: 'none  ' });
    }
  }, [isOver]);

  // isHovered is used to customize styles if a Post is hovered
  const handleMouseEnter = () => {
    setIsHovered(true);
  };

  const handleMouseLeave = () => {
    setIsHovered(false);
  };

  const handleFocus = () => {
    setIsFocused(true);
    focusPost({ id, board_id: board.id }, send);
  };

  // handleChange updates the textarea value and the textarea height
  const handleChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { value } = event.target;
    setTextareaValue(value);

    if (textareaRef.current) {
      let scrollHeight = textareaRef.current.scrollHeight;
      if (scrollHeight != textareaHeight) {
        // Reset height to auto to calculate scroll height based on contents inside textarea
        textareaRef.current.style.height = 'auto';
        // Reassign scroll height based on content
        scrollHeight = textareaRef.current.scrollHeight;
        setTextareaHeight(scrollHeight);
      }
    }
  };

  const handleBlur = () => {
    setIsFocused(false);
    updatePost({ id, content: textareaValue, height: textareaHeight }, send);
  };

  const handlePickColor = (color: string) => {
    updatePost({ id, color }, send);
    setColorSetting(color);
  };

  const handleDeletePost = () => {
    if (posts.length === 1) {
      deletePostGroup(post.post_group_id, send);
      return;
    }
    deletePost({ post_id: id, board_id: board.id }, send);
  };

  const ColorPicker = () => {
    return (
      <div className="flex space-x-1 items-center">
        {Object.keys(POST_COLORS).map((key) => {
          const colorName = displayColor(key);
          const colorValue = POST_COLORS[key];
          return (
            <div key={`color-${key}`} data-tooltip-id="my-tooltip" data-tooltip-content={colorName}>
              <button
                className="w-3 h-3 btn-square border border-gray-300"
                style={{ backgroundColor: colorValue }}
                onClick={() => handlePickColor(colorValue)}
              />
            </div>
          );
        })}
      </div>
    );
  };

  const PostActions = () => {
    return (
      <div
        className="card-actions justify-between items-center"
        style={{ visibility: isHovered ? 'visible' : 'hidden' }}
      >
        <ColorPicker />
        <button className="text-gray-500 hover:text-gray-700" onClick={handleDeletePost}>
          <FaRegTrashAlt />
        </button>
      </div>
    );
  };

  drag(drop(ref));

  return (
    <div
      ref={ref}
      style={getStyles(isDragging, color, borderStyle)}
      role="DraggablePost"
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      <div className="card card-compact border border-gray-500 cursor-move !rounded-none">
        <div className="card-body !py-1">
          <div className="h-4">
            {typingBy ? (
              <div className="text-center text-xs text-gray-600"> {`${typingBy.name} is typing...`}</div>
            ) : (
              <PostActions />
            )}
          </div>
          <textarea
            ref={textareaRef}
            className="textarea textarea-ghost textarea-sm textarea-bordered leading-4 resize-none"
            value={textareaValue}
            onChange={handleChange}
            onBlur={handleBlur}
            onFocus={handleFocus}
            autoFocus={autoFocus}
            style={{ height: textareaHeight }}
          />
          <div className="flex h-6 justify-between items-center">
            <div data-tooltip-id="my-tooltip" data-tooltip-content={authorName}>
              <Avatar id={user_id} size={16} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
});

function getStyles(isDragging: boolean, color: string, borderStyle: object): CSSProperties {
  return {
    opacity: isDragging ? 0 : 1,
    height: isDragging ? 0 : '',
    minHeight: DEFAULT_POST_HEIGHT,
    width: POST_WIDTH,
    background: color,
    ...borderStyle,
  };
}

function getName(userID: string, boardMembers: User[]): string | undefined {
  let name;
  boardMembers.forEach((user) => {
    if (user.id == userID) {
      name = user.name;
    }
  });
  return name;
}

function isHoverAbove(monitor: DropTargetMonitor, ref: HTMLDivElement): boolean {
  const hoverBoundingRect = ref.getBoundingClientRect();

  // Get vertical middle
  const hoverMiddleY = (hoverBoundingRect.bottom - hoverBoundingRect.top) / 2;

  // Determine mouse position
  const clientOffset = monitor.getClientOffset();

  // Get pixels to the top
  const hoverClientY = (clientOffset as XYCoord).y - hoverBoundingRect.top;
  return hoverClientY < hoverMiddleY;
}

export default PostUI;
