'use client';

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
            <NewBoardForm />
          </div>
        </div>
      </dialog>
    </>
  );
}
