import styles from './Header.module.css';
import { useNavigate } from 'react-router-dom';
import LoginIcon from '@mui/icons-material/Login';
import LogoutIcon from '@mui/icons-material/Logout';
import { useAuth } from '../../contexts/Auth';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Box,
  Typography
} from '@mui/material';
import { useState } from 'react';

function Header() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const [openDialog, setOpenDialog] = useState(false);

  const handleNavigation = (path) => {
    navigate(path);
  };

  const handleLogoutClick = () => {
    setOpenDialog(true);
  };

  const handleLogoutConfirm = () => {
    logout();
    handleNavigation('/login');
    setOpenDialog(false);
  };

  const handleDialogClose = () => {
    setOpenDialog(false);
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
          onClick={handleLogoutClick}
        />
      ) : (
        <LoginIcon
          className={`${styles.icon} ${styles.loginIcon}`}
          onClick={() => handleNavigation('/login')}
        />
      )}

      <Dialog
        open={openDialog}
        onClose={handleDialogClose}
        PaperProps={{
          sx: {
            background: 'linear-gradient(to bottom, #1a1a1a, #000)',
            color: 'white',
            textAlign: 'center',
            p: 4,
            borderRadius: 3,
            minWidth: 300,
          }
        }}
      >
        <DialogTitle sx={{ pb: 1 }}>
          <Typography variant="h5" sx={{ color: '#ff4d4f', fontWeight: 'bold' }}>
            Tem certeza?
          </Typography>
        </DialogTitle>

        <DialogContent>
          <Typography sx={{ color: '#ccc' }}>
            Deseja mesmo sair da sua conta?
          </Typography>
        </DialogContent>

        <DialogActions sx={{ justifyContent: 'center', mt: 3 }}>
          <Button
            onClick={handleDialogClose}
            sx={{
              border: '1px solid #888',
              color: '#ccc',
              mr: 2,
              '&:hover': {
                backgroundColor: '#333',
                color: 'white',
              }
            }}
          >
            Cancelar
          </Button>

          <Button
            onClick={handleLogoutConfirm}
            sx={{
              backgroundColor: '#ff4d4f',
              color: 'white',
              '&:hover': {
                backgroundColor: '#e60023',
              }
            }}
          >
            Sair
          </Button>
        </DialogActions>
      </Dialog>

    </header>
  );
}

export default Header;
