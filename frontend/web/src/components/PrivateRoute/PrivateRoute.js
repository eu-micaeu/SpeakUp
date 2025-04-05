import { useAuth } from "../../contexts/Auth";
import { useLocation, Navigate } from "react-router-dom";

function PrivateRoute({ children }) {
    const { user } = useAuth();
    const location = useLocation();

    console.log("PrivateRoute - User:", user);
    console.log("PrivateRoute - Location:", location);

    return user ? (
        children
    ) : (
        <Navigate to="/login" state={{ from: location }} replace />
    );
}

export default PrivateRoute;