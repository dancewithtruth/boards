import { getBoard } from '@/api/board';
import { listPostGroups } from '@/api/post';
import { Board, PostGroupMap } from '@/components/dnd/board';
import { CustomDragLayer } from '@/components/dnd/customDragLayer';
import { COOKIE_NAME_JWT_TOKEN, SIDEBAR_WIDTH_PX } from '@/constants';
import { cookies } from 'next/headers';

export const metadata = {
  title: 'Board',
  description: 'Collaborate with your team',
};

async function fetchPostGroupsData(boardID: string) {
  const cookieStore = cookies();
  const jwtToken = cookieStore.get(COOKIE_NAME_JWT_TOKEN);
  if (jwtToken) {
    const response = await listPostGroups(boardID, jwtToken.value);
    const posts = response.result.reduce((map, postGroup) => {
      map[postGroup.id] = postGroup;
      return map;
    }, {} as PostGroupMap);
    return posts;
  } else {
    throw new Error('Please log in.');
  }
}

async function fetchBoardData(boardID: string) {
  const cookieStore = cookies();
  const jwtToken = cookieStore.get(COOKIE_NAME_JWT_TOKEN);
  if (jwtToken) {
    const board = await getBoard(boardID, jwtToken.value);
    return board;
  } else {
    throw new Error('Please log in.');
  }
}

export default async function BoardPage({ params: { id: boardID } }: { params: { id: string } }) {
  const postGroups = await fetchPostGroupsData(boardID);
  const board = await fetchBoardData(boardID);
  return (
    <div className="flex" style={{ paddingLeft: SIDEBAR_WIDTH_PX }}>
      <Board snapToGrid={true} board={board} postGroups={postGroups} />
      <CustomDragLayer snapToGrid={true} />
    </div>
  );
}
