import styled from 'styled-components';
import { useNavigate } from 'react-router-dom';
import LoginIcon from '@mui/icons-material/Login';

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
`;

function Header() {
  const navigate = useNavigate();

  const handleLoginClick = () => {
    navigate('/login');
  };

  return (
    <HeaderContainer>
      <Logo src="logo.png" alt="Logo" />
      <LoginIcon sx={{ color: "rgb(0, 255, 0)", margin: "20px", fontSize: "2rem", cursor: "pointer" }}
      onMouseEnter={(e) => e.target.style.color = "#fff"}
      onMouseLeave={(e) => e.target.style.color = "rgb(0, 255, 0)"}
      onClick={handleLoginClick} />
    </HeaderContainer>
  );
}

export default Header;
