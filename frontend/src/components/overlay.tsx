import React from 'react';
import { Transition } from '@headlessui/react';

export const Overlay: React.FC<{ show: boolean; text: string }> = ({ show, text }) => {
  return (
    <Transition
      show={show}
      enter="transition-opacity duration-300"
      enterFrom="opacity-0"
      enterTo="opacity-100"
      leave="transition-opacity duration-200"
      leaveFrom="opacity-100"
      leaveTo="opacity-0"
    >
      <div className="fixed inset-0 bg-white bg-opacity-75 flex items-center justify-center" style={{ zIndex: 10002 }}>
        <div className="flex flex-col items-center">
          <div className="mb-4 text-xl font-bold text-gray-700">{text}</div>
          <div className="w-16 h-16 border-t-2 border-gray-600 rounded-full animate-spin"></div>
        </div>
      </div>
    </Transition>
  );
};
