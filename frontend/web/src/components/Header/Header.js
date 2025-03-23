import './Header.css';
import { useNavigate } from 'react-router-dom';

function Header() {
  const navigate = useNavigate();

  const handleLoginClick = () => {
    navigate('/login');
  };

  return (
    <header>
      <img src="logo.png" alt="Logo"></img>
      <button onClick={handleLoginClick}>Login</button>
    </header>
  );
}

export default Header;
