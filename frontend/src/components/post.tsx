import React, { ChangeEvent, FC, useRef, useState } from 'react';
import { ItemTypes, POST_COLORS, POST_HEIGHT, POST_WIDTH } from '@/constants';
import { useDrag } from 'react-dnd';
import { FaRegTrashAlt, FaVoteYea } from 'react-icons/fa';
import { User } from '@/api/users';
import Avatar from './avatar';

export interface Post {
  id: any;
  left: number;
  top: number;
  content: string;
  color: string;
  zIndex: number;
  customHeight?: number;
  user: User;
}
interface PostProps {
  post: Post;
  updatePost: (data: Post) => void;
  setColorSetting: (color: string) => void;
  deletePost: (id: string) => void;
}

const Post: FC<PostProps> = ({ post, updatePost, setColorSetting, deletePost }) => {
  const { id, left, top, content, color, zIndex, customHeight, user } = post;
  const [isHovered, setIsHovered] = useState(false);
  const [isFocused, setIsFocused] = useState(false);
  const [textareaValue, setTextareaValue] = useState(content);
  const [textareaHeight, setTextareaHeight] = useState(customHeight);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  // handleChange updates the textarea value and the textarea height
  const handleChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { value } = event.target;
    setTextareaValue(value);

    if (textareaRef.current) {
      const scrollHeight = textareaRef.current.scrollHeight;
      setTextareaHeight(scrollHeight);
    }
  };

  // isFocused is used to prevent the Post from being dragged when the textarea is focused
  const handleFocus = () => {
    setIsFocused(true);
  };

  const handleBlur = (event: ChangeEvent<HTMLTextAreaElement>) => {
    setIsFocused(false);
  };

  // isHovered is used to customize styles if a Post is hovered
  const handleMouseEnter = () => {
    setIsHovered(true);
  };

  const handleMouseLeave = () => {
    setIsHovered(false);
  };

  const [{ isDragging }, drag] = useDrag(
    () => ({
      type: ItemTypes.POST,
      item: { id, left, top , thing: '123'},
      canDrag: (_) => !isFocused,
      collect: (monitor) => {
        return {
          isDragging: monitor.isDragging(),
        };
      },
    }),
    [id, left, top, isFocused]
  );

  // If Post is being dragged, then hide the original Post from view
  if (isDragging) return null;

  return (
    <div
      ref={drag}
      className="card card-compact cursor-move shadow-md absolute"
      style={{
        minHeight: POST_HEIGHT,
        width: POST_WIDTH,
        top: top,
        left: left,
        background: color,
        ...(isHovered
          ? {
              zIndex: 10000,
              border: `1px solid black`,
            }
          : { zIndex, border: `1px solid gray` }),
      }}
      onMouseEnter={handleMouseEnter}
      onMouseLeave={handleMouseLeave}
    >
      <div className="card-body !py-1">
        <div className="card-actions justify-between">
          <ColorPicker updatePost={updatePost} setColorSetting={setColorSetting} />
          <button className="btn-xs text-gray-500 hover:text-gray-700" onClick={() => deletePost(id)}>
            <FaRegTrashAlt />
          </button>
        </div>
        <textarea
          ref={textareaRef}
          className="textarea textarea-ghost textarea-sm textarea-bordered leading-4 resize-none"
          style={{ ...(textareaHeight && { height: textareaHeight }) }}
          value={textareaValue}
          onChange={handleChange}
          onFocus={handleFocus}
          onBlur={handleBlur}
        />
        <div className="flex h-6 justify-between items-center">
          <div key={`author-${post.user.id}`} data-tooltip-id="my-tooltip" data-tooltip-content={user.name}>
            <Avatar id={user.id} size={6} />
          </div>
        </div>
      </div>
    </div>
  );
};

interface ColorPickerProps {
  updatePost: (data: Post) => void;
  setColorSetting: (color: string) => void;
}

const ColorPicker = ({ updatePost, setColorSetting }: ColorPickerProps) => {
  return (
    <div className="flex space-x-1 items-center">
      {Object.keys(POST_COLORS).map((key) => {
        const colorName = displayColor(key);
        const colorValue = POST_COLORS[key];
        const data = { color: colorValue } as Post;
        return (
          <div key={`color-${key}`} data-tooltip-id="my-tooltip" data-tooltip-content={colorName}>
            <button
              className="w-3 h-3 btn-square border border-gray-300"
              style={{ backgroundColor: colorValue }}
              onClick={() => {
                updatePost(data);
                setColorSetting(colorValue);
              }}
            />
          </div>
        );
      })}
    </div>
  );
};

function displayColor(str: string): string {
  return str.toLowerCase().replace(/_/g, ' ');
}

export default Post;
