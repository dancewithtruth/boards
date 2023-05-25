import Board from '@/components/board';
import NewBoardForm from '@/components/forms/newboard';
import ConfiguredToastContainer from '@/components/toastcontainer';
import { FaPlus } from 'react-icons/fa';

const Dashboard = () => {
  return (
    <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
      <ConfiguredToastContainer />
      <h1 className="text-4xl font-bold mt-10 mb-10">Dashboard</h1>
      <div>
        <div className="flex justify-between items-end">
          <h2 className="text-2xl font-bold">My Boards</h2>
          <label htmlFor="my-modal-4" className="btn btn-primary">
            <FaPlus className="mr-2" />
            New Board
          </label>
          <input type="checkbox" id="my-modal-4" className="modal-toggle" />
          <label htmlFor="my-modal-4" className="modal cursor-pointer">
            <label className="modal-box relative" htmlFor="">
              <NewBoardForm />
            </label>
          </label>
        </div>
        <div className="divider"></div>
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 place-items-center">
            <Board title={'First Board'} description={'My very first board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
          </div>
        </div>
      </div>
      <div className="mt-12">
        <h2 className="text-2xl font-bold">Shared Boards</h2>
        <div className="divider"></div>
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 place-items-center">
            <Board title={'First Board'} description={'My very first board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
            <Board title={'Second Board'} description={'My second board'} createdAt={'2023-05-19 18:22:03.515'} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
