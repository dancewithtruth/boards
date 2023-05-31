'use client';

import AppNavbar from '@/components/appnavbar';
import Sidebar from '@/components/sidebar';
import { DndProvider } from 'react-dnd';
import { HTML5Backend } from 'react-dnd-html5-backend';

const Layout = ({ children }: { children: React.ReactNode }) => {
  const SIDEBAR_WIDTH = 192;
  return (
    <DndProvider backend={HTML5Backend}>
      <div className="sketchbook-bg">
        <AppNavbar />
        <div className="flex flex-col">
          <div className="h-16 w-full" />
          <div className="flex">
            <Sidebar width={SIDEBAR_WIDTH} />
            <div className="h-full" style={{ width: SIDEBAR_WIDTH }} />
            {children}
          </div>
        </div>
      </div>
    </DndProvider>
  );
};

export default Layout;
