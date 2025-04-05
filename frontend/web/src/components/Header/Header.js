import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';
import LoginIcon from '@mui/icons-material/Login';
import LogoutIcon from '@mui/icons-material/Logout';
import { useAuth } from '../../contexts/Auth';

const HeaderContainer = styled.header`
  background-color: #000;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 0;
  height: 10%;
`;

const Logo = styled.img`
  width: 40px;
  margin: 20px;
  cursor: pointer;
`;

function Header() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleNavigation = (path) => {
    navigate(path);
  };

  return (
    <HeaderContainer>
      <Logo 
        src="logo.png" 
        alt="Logo" 
        onClick={() => handleNavigation('/')} 
      />

      {user ? (
        <LogoutIcon
          sx={{
            color: "rgb(255, 0, 0)",
            margin: "20px",
            fontSize: "2rem",
            cursor: "pointer",
            '&:hover': {
              color: "#fff"
            }
          }}
          onClick={() => {
            logout();
            handleNavigation('/');
          }}
        />
      ) : (
        <LoginIcon
          sx={{
            color: "rgb(0, 255, 0)",
            margin: "20px",
            fontSize: "2rem",
            cursor: "pointer",
            '&:hover': {
              color: "#fff"
            }
          }}
          onClick={() => handleNavigation('/login')}
        />
      )}
    </HeaderContainer>
  );
}

export default Header;