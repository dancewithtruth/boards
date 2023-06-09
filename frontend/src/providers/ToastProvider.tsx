'use client';
import { ToastContainer } from 'react-toastify';

const ToastProvider = () => {
  return (
    <ToastContainer
      position="bottom-center"
      autoClose={3000}
      hideProgressBar={true}
      newestOnTop={false}
      closeOnClick={true}
      rtl={false}
      pauseOnFocusLoss={true}
      draggable={true}
      pauseOnHover={true}
      theme="dark"
    />
  );
};

export default ToastProvider;
