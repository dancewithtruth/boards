import React, { FC, ReactNode } from 'react';
import { ItemTypes } from '../../constants';
import { useDrag } from 'react-dnd';

export interface PostProps {
  id: any;
  left: number;
  top: number;
  content: string;
  hideSourceOnDrag?: boolean;
  children?: ReactNode;
}

const Post: FC<PostProps> = ({ id, left, top, hideSourceOnDrag, content, children }) => {
  const [{ isDragging }, drag] = useDrag(
    () => ({
      type: ItemTypes.POST,
      item: { id, left, top },
      collect: (monitor) => ({
        isDragging: monitor.isDragging(),
      }),
    }),
    [id, left, top]
  );

  if (isDragging) return <div ref={drag} />;
  return (
    <div
      ref={drag}
      className={`text-sm text-gray-700 bg-base-200 cursor-move h-[100px] w-[200px] border border-gray-700 absolute`}
      style={{ top: top, left: left }}
    >
      {content}
    </div>
  );
};

export default Post;
