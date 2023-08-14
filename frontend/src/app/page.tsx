import Footer from '@/components/footer';
import WidthContainer from '@/components/widthContainer';
import { NAVBAR_HEIGHT_PX } from '@/constants';
import Link from 'next/link';
import { FaPaperPlane, FaStickyNote, FaRegFileCode } from 'react-icons/fa';

export default function Page() {
  return (
    <div>
      <div className="hero bg-base-200" style={{ height: `calc(100vh - ${NAVBAR_HEIGHT_PX})` }}>
        <div className="hero-content flex-col lg:flex-row lg:space-x-6">
          <div>
            <h1 className="text-5xl font-bold">Boards</h1>
            <p className="py-6 max-w-sm text-xl text-gray-700">
              Collaborate with others and turn your ideas into actions
            </p>
            <div className="flex space-x-4">
              <Link href="/signup" className="btn btn-primary">
                Sign Up
              </Link>
              <Link href="/dashboard" className="btn btn-secondary btn-outline">
                Visit app
              </Link>
            </div>
          </div>
          <img
            src="/Hero.png"
            style={{ width: '700px' }}
            className="rounded-lg shadow-2xl"
            alt="Picture of the author"
          />
        </div>
      </div>
      <div className="bg-white h-[50vh]">
        <WidthContainer className="h-full">
          <div className="flex items-center justify-between h-full">
            <div className="card w-96 border !rounded-none">
              <div className="card-body">
                <FaPaperPlane size={30} />
                <h2 className="card-title">Invite</h2>
                <p className="text-gray-700">Invite other collaborators and work together in real-time</p>
              </div>
            </div>
            <div className="card w-96 border !rounded-none">
              <div className="card-body">
                <FaStickyNote size={30} />
                <h2 className="card-title">Organize</h2>
                <p className="text-gray-700">Capture ideas, tasks, and goals on digital sticky notes</p>
              </div>
            </div>
            <div className="card w-96 border !rounded-none">
              <div className="card-body">
                <FaRegFileCode size={30} />
                <h2 className="card-title">Automate</h2>
                <p className="text-gray-700">Automatically group similar posts together and summarize content</p>
              </div>
            </div>
          </div>
        </WidthContainer>
      </div>
      <div className="bg-base-200 h-[50vh]">
        <WidthContainer className="h-full ">
          <h2 className="text-2xl">Boards is built with</h2>
        </WidthContainer>
      </div>
      <Footer />
    </div>
  );
}
