import styles from './Header.module.css';
import { useNavigate } from 'react-router-dom';
import LoginIcon from '@mui/icons-material/Login';
import LogoutIcon from '@mui/icons-material/Logout';
import { useAuth } from '../../contexts/Auth';

function Header() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleNavigation = (path) => {
    navigate(path);
  };

  return (
    <header className={styles.headerContainer}>
      <img 
        src="logo.png" 
        alt="Logo"
        className={styles.logo}
      />

      {user ? (
        <LogoutIcon
          className={`${styles.icon} ${styles.logoutIcon}`}
          onClick={() => {
            logout();
            handleNavigation('/login');
          }}
        />
      ) : (
        <LoginIcon
          className={`${styles.icon} ${styles.loginIcon}`}
          onClick={() => handleNavigation('/login')}
        />
      )}
    </header>
  );
}

export default Header;
