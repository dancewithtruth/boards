import AppNavbar from '@/components/appnavbar';
import Sidebar from '@/components/sidebar';

const Layout = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="sketchbook-bg">
      <AppNavbar />
      <div className="flex flex-col">
        <div className="h-16 w-full" />
        <div className="flex">
          <Sidebar />
          <div className="w-24 h-full" />
          {children}
        </div>
      </div>
    </div>
  );
};

export default Layout;
