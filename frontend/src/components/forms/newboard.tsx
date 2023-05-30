'use client';
import { useRouter } from 'next/navigation';
import { ChangeEvent, FormEvent, useState } from 'react';
import { createBoard } from '../../../helpers/api/boards';
import { toast } from 'react-toastify';

const NewBoardForm = () => {
  const [name, setName] = useState<string | undefined>();
  const [description, setDescription] = useState<string | undefined>();
  const [isLoading, setIsLoading] = useState(false);

  const router = useRouter();

  const handleNameChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.value == '') {
      setName(undefined);
    } else {
      setName(e.target.value);
    }
  };

  const handleDescriptionChange = (e: ChangeEvent<HTMLTextAreaElement>) => {
    if (e.target.value == '') {
      setDescription(undefined);
    } else {
      setDescription(e.target.value);
    }
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsLoading(true);

    try {
      const board = await createBoard({ name, description });
      toast.success('Board created!');
      //TODO: Reload get boards and redirect to new board
      router.push(`/boards/${board.id}`);
      toast.info('Redirecting to new board...');
    } catch (error) {
      toast.error(String(error));
    } finally {
      setIsLoading(false);
    }
  };
  return (
    <>
      <div className="card-body">
        <h3 className="text-2xl font-bold">New Board</h3>
        <p className="text-gray-500">Invite collaborators to your brand new board.</p>
        <form onSubmit={handleSubmit} className="mt-4">
          <div className="form-control">
            <label className="label">
              <span className="label-text">Name</span>
              <span className="label-text-alt text-xs text-gray-300">optional</span>
            </label>
            <input
              type="text"
              id="name"
              className="input input-bordered w-full"
              placeholder="Board name"
              value={name}
              onChange={handleNameChange}
            />
          </div>
          <div className="form-control">
            <label className="label">
              <span className="label-text">Description</span>
              <span className="label-text-alt text-xs text-gray-300">optional</span>
            </label>
            <textarea
              className="textarea textarea-bordered"
              placeholder="Enter your board description"
              value={description}
              onChange={handleDescriptionChange}
            ></textarea>
          </div>
          <div className="form-control mt-6">
            <div className="flex flex-col w-full border-opacity-50">
              <button type="submit" className="btn btn-secondary">
                {isLoading ? 'Loading...' : 'Create board'}
              </button>
            </div>
          </div>
        </form>
      </div>
    </>
  );
};

export default NewBoardForm;
