'use client';

import { DndProvider } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';

const Layout = ({ children }: { children: React.ReactNode }) => {
  return <DndProvider backend={HTML5Backend}>{children}</DndProvider>;
};

export default Layout;
