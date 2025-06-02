import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './Sidebar.module.css';
import {
    Add,
    Delete,
    AccountCircle,
    MeetingRoom,
    Try as TryIcon,
    AutoStories as AutoStoriesIcon,
    Settings,
    Info as InfoIcon
} from '@mui/icons-material';
import {
    Modal,
    Box,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
} from '@mui/material';

// Api functions
import {
    deleteChat
} from '../../utils/api';

// Utils
import { removeAuthTokenFromCookies } from '../../utils/cookies';

export default function Sidebar({
    chats,
    setChats,
    setCurrentChatId,
    setMessages,
    setInputMessage,
    sidebarOpen,
    setSidebarOpen,
    selectedChat,
    setSelectedChat
}) {
    const navigate = useNavigate();
    const [openSettings, setOpenSettings] = useState(false);

    const handleNewChat = () => {
        setCurrentChatId(null);
        setMessages([]);
        setInputMessage('');
        setSelectedChat(null);
        setSidebarOpen(false);
    };

    const handleChatClick = (chat) => {
        setCurrentChatId(chat.id);
        setSelectedChat(chat);
    };

    const goToPerfil = () => {
        navigate('/perfil');
        setOpenSettings(false);
    };
    const goToIndex = () => {
        removeAuthTokenFromCookies();
        navigate('/');
        setOpenSettings(false);
    };
    const goToPalavreco = () => navigate('/palavreco');
    const goToPlanoDeEstudos = () => navigate('/teaching-plan');

    const handleOpenSettings = () => {
        setOpenSettings(true);
    };

    const handleCloseSettings = () => {
        setOpenSettings(false);
    };

    return (
        <aside className={`${styles.sidebar} ${!sidebarOpen ? styles.sidebarClosed : ''}`}>
            <div className={styles.sidebarContent}>
                <div>
                    <div className={styles.header}>
                        <img src="/logo.png" alt="Logo" width={35} />
                        <h3>SpeakUp</h3>
                    </div>

                    <button className={styles.newChat} onClick={handleNewChat}>
                        <Add />
                    </button>

                    <nav className={styles.chatList}>
                        {chats.map(chat => (
                            <div
                                key={chat.id}
                                className={`${styles.navItem} ${selectedChat?.id === chat.id ? styles.selected : ''}`}
                                onClick={() => handleChatClick(chat)}
                            >
                                <span>{chat.topic}</span>
                                <button
                                    className={styles.deleteButton}
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        deleteChat(chat.id)
                                            .then(() => {
                                                setMessages([]);
                                                setCurrentChatId(null);
                                                setSelectedChat(null);
                                                setChats(chats.filter(c => c.id !== chat.id));
                                            })
                                            .catch(err => console.error('Error deleting chat:', err));
                                    }}
                                >
                                    <Delete fontSize="small" />
                                </button>
                            </div>
                        ))}
                    </nav>
                </div>

                <div className={styles.footer}>
                    <button onClick={goToPlanoDeEstudos}>
                        <AutoStoriesIcon fontSize="small" />
                        <span>Plano de Estudo</span>
                    </button>
                    {/* <button onClick={goToPalavreco}>
                        <TryIcon fontSize="small" />
                        <span>Palavreco</span>
                    </button> */}
                    <button onClick={handleOpenSettings}>
                        <Settings fontSize="small" />
                        <span>Configurações</span>
                    </button>
                </div>
            </div>

            <Modal
                open={openSettings}
                onClose={handleCloseSettings}
                aria-labelledby="settings-modal-title"
                aria-describedby="settings-modal-description"
            >
                <Box sx={{
                    position: 'absolute',
                    top: '50%',
                    left: '50%',
                    transform: 'translate(-50%, -50%)',
                    width: 200,
                    bgcolor: 'black',
                    boxShadow: 24,
                    p: 2,
                    borderRadius: 5
                }}>
                    <List>
                        <ListItem button onClick={goToPerfil} sx={{ borderRadius: "10px", '&:hover': { cursor: 'pointer', backgroundColor: "#2a2a2a" } }}>
                            <ListItemIcon>
                                <AccountCircle fontSize="small" sx={{color: "white"}} />
                            </ListItemIcon>
                            <ListItemText primary="Perfil" sx={{color: "white"}} />
                        </ListItem>
                        <ListItem button onClick={goToIndex} sx={{ borderRadius: "10px", '&:hover': { cursor: 'pointer', backgroundColor: "#2a2a2a" } }}>
                            <ListItemIcon>
                                <MeetingRoom fontSize="small" sx={{color: "red"}} />
                            </ListItemIcon>
                            <ListItemText primary="Sair" sx={{color: "red"}} />
                        </ListItem>
                        {/* <ListItem button onClick={() => navigate('/about')} sx={{ borderRadius: "10px", '&:hover': { cursor: 'pointer', backgroundColor: "#2a2a2a" } }}>
                            <ListItemIcon>
                                <InfoIcon fontSize="small" sx={{color: "white"}} />
                            </ListItemIcon>
                            <ListItemText primary="Sobre" sx={{color: "white"}} />
                        </ListItem> */}
                    </List>
                </Box>
            </Modal>
        </aside>
    );
}