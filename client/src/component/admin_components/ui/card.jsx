import React from 'react';

export const Card = ({ children, className }) => (
  <div className={`card ${className}`}>{children}</div>
);

export const CardContent = ({ children, className }) => (
  <div className={`card-content ${className}`}>{children}</div>
);