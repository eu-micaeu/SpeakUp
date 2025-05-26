import { Navigate } from 'react-router-dom';
import { isAuthTokenValid } from '../../utils/cookies';

const PrivateRoute = ({ children }) => {
  return isAuthTokenValid() ? children : <Navigate to="/login" />;
};

export default PrivateRoute;
