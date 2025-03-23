import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

// Pages
import Index from './pages/Index/Index';
import Login from './pages/Login/Login';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Index />} />
        <Route path="/login" element={<Login />} />
      </Routes>
    </Router>
  );
}

export default App;
