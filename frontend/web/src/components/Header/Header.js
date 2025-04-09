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
  width: 60px;
  margin: 25px 50px;
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
      />

      {user ? (
        <LogoutIcon
          sx={{
            color: "rgb(255, 0, 0)",
            margin: "25px 50px",
            fontSize: "3rem",
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
            color: "rgb(255, 255, 255)",
            margin: "25px 50px",
            fontSize: "3rem",
            cursor: "pointer",
            '&:hover': {
              color: "rgb(182, 182, 182)"
            }
          }}
          onClick={() => handleNavigation('/login')}
        />
      )}
    </HeaderContainer>
  );
}

export default Header;