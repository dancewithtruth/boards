import { getBoard } from '@/api/board';
import { listPosts } from '@/api/post';
import { Board, PostMap } from '@/components/dnd/board';
import { CustomDragLayer } from '@/components/dnd/customDragLayer';
import { COOKIE_NAME_JWT_TOKEN, SIDEBAR_WIDTH } from '@/constants';
import { cookies } from 'next/headers';

export const metadata = {
  title: 'Board',
  description: 'Collaborate with your team',
};

async function fetchPostsData(boardId: string) {
  const cookieStore = cookies();
  const jwtToken = cookieStore.get(COOKIE_NAME_JWT_TOKEN);
  if (jwtToken) {
    const response = await listPosts(boardId, jwtToken.value);
    const posts = response.data.reduce((map, post) => {
      map[post.id] = post;
      return map;
    }, {} as PostMap);
    return posts;
  } else {
    throw new Error('Please log in.');
  }
}

async function fetchBoardData(boardId: string) {
  const cookieStore = cookies();
  const jwtToken = cookieStore.get(COOKIE_NAME_JWT_TOKEN);
  if (jwtToken) {
    const board = await getBoard(boardId, jwtToken.value);
    return board;
  } else {
    throw new Error('Please log in.');
  }
}

export default async function BoardPage({ params: { id: boardId } }: { params: { id: string } }) {
  const posts = await fetchPostsData(boardId);
  const board = await fetchBoardData(boardId);
  return (
    <div className="flex" style={{ paddingLeft: SIDEBAR_WIDTH }}>
      <Board snapToGrid={true} board={board} posts={posts} />
      <CustomDragLayer snapToGrid={true} />
    </div>
  );
}
