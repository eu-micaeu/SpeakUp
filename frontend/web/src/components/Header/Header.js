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
            backgroundColor: 'black',
            color: 'white',
            textAlign: 'center',
            p: 3,
            borderRadius: 2,
            boxShadow: 'none',
          }
        }}
      >
        <DialogTitle>
          <Typography variant="h6" sx={{ color: 'white' }}>
            Confirmar Logout
          </Typography>
        </DialogTitle>

        <DialogContent>
          <Typography sx={{ color: 'white' }}>
            Tem certeza que deseja sair?
          </Typography>
        </DialogContent>

        <DialogActions sx={{ justifyContent: 'center', mt: 2 }}>
          <Button onClick={handleDialogClose} sx={{ color: 'white' }}>
            Cancelar
          </Button>
          <Button onClick={handleLogoutConfirm} sx={{ color: 'red' }}>
            Sair
          </Button>
        </DialogActions>
      </Dialog>
    </header>
  );
}

export default Header;
