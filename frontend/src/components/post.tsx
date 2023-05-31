import React, { ChangeEvent, FC, ReactNode, useRef, useState } from 'react';
import { ItemTypes, POST_COLORS } from '../../constants';
import { useDrag } from 'react-dnd';
import { FaRegTrashAlt, FaVoteYea } from 'react-icons/fa';
import { User } from '../../api/users';
import Avatar from './avatar';

export interface PostData {
  id: any;
  left: number;
  top: number;
  content: string;
  color: string;
  zIndex: number;
  customHeight?: number;
  user: User;
}
export interface PostProps {
  data: PostData;
  updatePost: (data: PostData) => void;
  setColor: (color: string) => void;
  deletePost: (id: string) => void;
  hideSourceOnDrag?: boolean;
  children?: ReactNode;
}

const Post: FC<PostProps> = ({ data, updatePost, setColor, deletePost, hideSourceOnDrag, children }) => {
  const { id, left, top, content, color, zIndex, customHeight } = data;
  const [isHovered, setIsHovered] = useState(false);
  const [isFocused, setIsFocused] = useState(false);
  const [textareaHeight, setTextareaHeight] = useState(customHeight);
  const [textareaValue, setTextareaValue] = useState(content);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleChange = (event: ChangeEvent<HTMLTextAreaElement>) => {
    const { value } = event.target;
    setTextareaValue(value);
    const textarea = textareaRef.current;
    if (textarea) {
      const scrollHeight = textarea.scrollHeight;
      setTextareaHeight(scrollHeight);
    }
  };

  const handleFocus = () => {
    setIsFocused(true);
  };

  const handleBlur = (event: ChangeEvent<HTMLTextAreaElement>) => {
    setIsFocused(false);
    const payload = { content: textareaValue, customHeight: textareaHeight } as PostData;
    // send payload and update value to backend
  };

  const handleMouseEnter = () => {
    setIsHovered(true);
  };

  const handleMouseLeave = () => {
    setIsHovered(false);
  };

  const [{ isDragging }, drag] = useDrag(
    () => ({
      type: ItemTypes.POST,
      item: { id, left, top },
      canDrag: (_) => !isFocused,
      collect: (monitor) => {
        return {
          isDragging: monitor.isDragging(),
        };
      },
    }),
    [id, left, top, isFocused]
  );

  if (isDragging) return null;
  return (
    <div
      ref={drag}
      className={`card card-compact min-h-[100px] w-[275px] cursor-move shadow-md absolute`}
      style={{
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
          <div className="flex space-x-1 items-center">
            {Object.keys(POST_COLORS).map((key) => {
              const colorName = displayColor(key);
              const colorHex = POST_COLORS[key];
              const data = { color: colorHex } as PostData;
              return (
                <div key={key} data-tooltip-id="my-tooltip" data-tooltip-content={colorName}>
                  <button
                    className=" w-3 h-3 btn-square border border-gray-300"
                    style={{ backgroundColor: colorHex }}
                    onClick={() => {
                      updatePost(data);
                      setColor(colorHex);
                    }}
                  />
                </div>
              );
            })}
          </div>
          <button className="btn-xs text-gray-500 hover:text-gray-700" onClick={() => deletePost(data.id)}>
            <FaRegTrashAlt />
          </button>
        </div>
        <textarea
          ref={textareaRef}
          className="textarea textarea-ghost textarea-sm textarea-bordered leading-4"
          onFocus={handleFocus}
          onBlur={handleBlur}
          onChange={handleChange}
          value={textareaValue}
          style={{ ...(textareaHeight && { height: textareaHeight }), resize: 'none' }}
        />
        <div className="flex h-6 justify-between items-center">
          <div key={`author-${data.user.id}`} data-tooltip-id="my-tooltip" data-tooltip-content={data.user.name}>
            <Avatar id={data.user.id} size={4} />
          </div>
        </div>
      </div>
    </div>
  );
};

function displayColor(str: string): string {
  return str.toLowerCase().replace(/_/g, ' ');
}

export default Post;
