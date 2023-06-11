'use client';

import { FiX } from 'react-icons/fi';
import { FaPlus } from 'react-icons/fa';

import NewBoardForm from '../forms/newBoard';

export default function NewBoardModal() {
  const ID = 'modal_new_board';
  return (
    <>
      <button className="btn btn-primary" onClick={() => (window as any)[ID].showModal()}>
        <FaPlus className="mr-2" />
        New Board
      </button>
      <dialog id={ID} className="modal">
        <div className="card bg-white w-96 shadow-md">
          <div className="card-body">
            <div className="flex justify-between items-center">
              <h3 className="text-2xl font-bold">New Board</h3>
              <button type="reset"className="btn btn-ghost btn-sm" onClick={() => (window as any)[ID].close()}>
                <FiX size={24} />
              </button>
            </div>
            <NewBoardForm />
          </div>
        </div>
      </dialog>
    </>
  );
}
