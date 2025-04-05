import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../../contexts/Auth';

const PrivateRoute = ({ children }) => {
  const { user } = useAuth();

  if (user === null) {
    return <div>Carregando...</div>;
  }

  return user ? children : <Navigate to="/login" />;
  
};

export default PrivateRoute;