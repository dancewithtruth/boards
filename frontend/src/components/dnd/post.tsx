'use client';

import { BoardWithMembers, User } from '@/api';
import { POST_COLORS, POST_HEIGHT, POST_WIDTH } from '@/constants';
import { displayColor } from '@/utils';
import { deletePost, focusPost, updatePost } from '@/ws/events';
import { Send } from '@/ws/types';
import { CSSProperties, ChangeEvent, FC, useEffect, useRef, useState } from 'react';
import { memo } from 'react';
import { FaRegTrashAlt } from 'react-icons/fa';
import Avatar from '../avatar';
import { PostUI } from './board';

type PostProps = {
  user: User;
  board: BoardWithMembers;
  send: Send;
  setColorSetting: (color: string) => void;
} & PostUI;

export const Post: FC<PostProps> = memo(function Post({
  user,
  id,
  user_id,
  board,
  content,
  color,
  height,
  send,
  setColorSetting,
  typingBy,
  autoFocus,
}) {
  const [textareaValue, setTextareaValue] = useState(content);
  const [isHovered, setIsHovered] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const allMembers = board.members;
  const authorName = getName(user_id, allMembers) || 'Unknown';

  useEffect(() => {
    setTextareaValue(content);
  }, [content]);

  // isHovered is used to customize styles if a Post is hovered
  const handleMouseEnter = () => {
    setIsHovered(true);
  };

  const handleMouseLeave = () => {
    setIsHovered(false);
  };

  const handleFocus = () => {
    focusPost({ id, board_id: board.id }, send);
  };

  // handleChange updates the textarea value
  const handleChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { value } = event.target;
    setTextareaValue(value);
  };

  const handleBlur = () => {
    updatePost({ id, board_id: board.id, content: textareaValue }, send);
  };

  const handleDelete = () => {
    deletePost({ post_id: id, board_id: board.id }, send);
  };

  const handlePickColor = (color: string) => {
    updatePost({ id, board_id: board.id, color }, send);
    setColorSetting(color);
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
        <button className="text-gray-500 hover:text-gray-700" onClick={handleDelete}>
          <FaRegTrashAlt />
        </button>
      </div>
    );
  };
  return (
    <div
      className="card card-compact border border-gray-500 cursor-move shadow-md"
      style={{ ...getStyles(color) }}
      role="Post"
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
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
        />
        <div className="flex h-6 justify-between items-center">
          <div data-tooltip-id="my-tooltip" data-tooltip-content={authorName}>
            <Avatar id={user_id} size={16} />
          </div>
        </div>
      </div>
    </div>
  );
});

function getStyles(color: string): CSSProperties {
  return {
    minHeight: POST_HEIGHT,
    width: POST_WIDTH,
    background: color,
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
