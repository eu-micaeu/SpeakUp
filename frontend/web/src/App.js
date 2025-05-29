import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

// Pages
import Index from './pages/Index/Index';
import Login from './pages/Login/Login';
import Register from './pages/Register/Register';
import Chat from './pages/Chat/Chat';
import TeachingPlan from './pages/TeachingPlan/TeachingPlan';
import Perfil from './pages/Perfil/Perfil';
import OnBoarding from './pages/OnBoarding/OnBoarding';
import Palavreco from './pages/Palavreco/Palavreco';

// Private Route
import PrivateRoute from './components/PrivateRoute/PrivateRoute';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Index />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/onboarding" element={<OnBoarding />} />
        <Route
          path="/chat"
          element={
            <PrivateRoute>
              <Chat />
            </PrivateRoute>
          }
        />
        <Route
          path="/teaching-plan"
          element={
            <PrivateRoute>
              <TeachingPlan />
            </PrivateRoute>
          }
        />
        <Route
          path="/perfil"
          element={
            <PrivateRoute>
              <Perfil />
            </PrivateRoute>
          }
        />
        <Route
          path="/palavreco"
          element={
            <PrivateRoute>
              <Palavreco />
            </PrivateRoute>
          }
        />
      </Routes>
    </Router>
  );
}

export default App;
