'use client';

import React, { useState } from 'react';

interface Post {
  id: number;
  title: string;
}

type PostMap = {
  [key: number]: Post;
};
export default function Playground() {
  const [posts, setPosts] = useState<PostMap>({ 1: { id: 1, title: 'This is a post' } });
  const handleClick = () => {
    const newPostId = getMaxNumber(posts) + 1;
    const newMap = { ...posts };
    newMap[newPostId] = { id: newPostId, title: 'Some random title' };
    setPosts(newMap);
  };

  return (
    <div>
      <button className="btn btn-primary" onClick={handleClick}>
        Click me
      </button>
      {Object.entries(posts).map(([postId, post]) => (
        <Post key={postId} post={post} />
      ))}
    </div>
  );
}

const Post = React.memo(({ post }: { post: Post }) => {
  console.log(`Render Post ID: ${post.id}`);
  return <div>{`Post ${post.id}`}</div>;
});

function arePropsEqual(oldProps: { post: Post }, newProps: { post: Post }): boolean {
  return true;
}

function getMaxNumber(obj: PostMap) {
  let maxNumber = 0;

  for (const key in obj) {
    if (typeof obj[key].id === 'number' && obj[key].id > maxNumber) {
      maxNumber = obj[key].id;
    }
  }

  return maxNumber;
}
