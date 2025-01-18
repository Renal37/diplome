import React from 'react';

export const Select = ({ children, ...props }) => (
  <div className="select-container">
    <select {...props}>
      {children}
    </select>
  </div>
);

export const SelectContent = ({ children }) => (
  <>{children}</>
);

export const SelectItem = ({ value, children }) => (
  <option value={value}>{children}</option>
);

export const SelectTrigger = ({ children }) => (
  <>{children}</>
);

export const SelectValue = ({ children }) => (
  <>{children}</>
);