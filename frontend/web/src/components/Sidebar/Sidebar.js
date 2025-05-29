import { useNavigate } from 'react-router-dom';
import styles from './Sidebar.module.css';
import {
    Add,
    Delete,
    AccountCircle,
    MeetingRoom,
    Try as TryIcon,
    AutoStories as AutoStoriesIcon,
} from '@mui/icons-material';

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
    setSidebarOpen
}) {
    const navigate = useNavigate();

    const handleNewChat = () => {
        setCurrentChatId(null);
        setMessages([]);
        setInputMessage('');
        setSidebarOpen(false);
    };

    const goToPerfil = () => navigate('/perfil');
    const goToIndex = () => {
        removeAuthTokenFromCookies();
        navigate('/');
    };
    const goToPalavreco = () => navigate('/palavreco');
    const goToPlanoDeEstudos = () => navigate('/teaching-plan');

    return (
        <aside className={`${styles.sidebar} ${!sidebarOpen ? styles.sidebarClosed : ''}`}>
            <div>
                <div className={styles.header}>
                    <img src="/logo.png" alt="Logo" width={30} />
                    <h3>SpeakUp</h3>
                </div>

                <button className={styles.newChat} onClick={handleNewChat}>
                    <Add />
                </button>

                <nav className={styles.chatList}>
                    {chats.map(chat => (
                        <div
                            key={chat.id}
                            className={styles.navItem}
                            onClick={() => {
                                setCurrentChatId(chat.id);
                                setSidebarOpen(false);
                            }}
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
                <button onClick={goToPalavreco}>
                    <TryIcon fontSize="small" />
                    <span>Palavreco</span>
                </button>
                <button onClick={goToPerfil}>
                    <AccountCircle fontSize="small" />
                    <span>Perfil</span>
                </button>
                <button onClick={goToIndex}>
                    <MeetingRoom fontSize="small" />
                    <span>Sair</span>
                </button>
            </div>
        </aside>
    );
}