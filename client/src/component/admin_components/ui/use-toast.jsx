import { useState } from 'react';

export const useToast = () => {
  const [toasts, setToasts] = useState([]);

  const toast = ({ title, description, variant }) => {
    setToasts([...toasts, { title, description, variant }]);
  };

  return { toast, toasts };
};