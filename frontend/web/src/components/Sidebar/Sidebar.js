import React from 'react';
import DeleteIcon from '@mui/icons-material/Delete';
import LogoutIcon from '@mui/icons-material/Logout';
import AppsIcon from '@mui/icons-material/Apps';
import CloseIcon from '@mui/icons-material/Close';  // Adicionando o CloseIcon
import styles from './Sidebar.module.css';

function Sidebar({
    isVisible,
    toggleSidebar,
    chats,
    currentChatId,
    setCurrentChatId,
    loadChatMessages,
    handleDeleteChat,
    goToIndex,
    clearCookies,
    goToHome,
    setMessages
}) {
    return (
        <aside className={`${styles.sidebar} ${!isVisible ? styles.sidebarHidden : ''}`}>
            <button
                className={`${styles.toggleButton} ${isVisible ? '' : styles.toggleButtonHidden}`}
                onClick={toggleSidebar}
            >
                <CloseIcon style={{ color: "#fff" }} />  {/* Cor opcional */}
            </button>

            <ul>
                {chats?.length > 0 ? (
                    chats.map((chat) => (
                        <li
                            key={chat.id}
                            className={chat.id === currentChatId ? styles.activeChat : ''}
                            onClick={() => {
                                setCurrentChatId(chat.id);
                                loadChatMessages(chat.id);
                            }}
                        >
                            {chat.topic || "Chat sem título"}
                            <DeleteIcon
                                onClick={(e) => {
                                    e.stopPropagation();
                                    handleDeleteChat(chat.id);
                                }}
                                className={styles.deleteIcon}
                            />
                        </li>
                    ))
                ) : (
                    <ol>Crie um chat!</ol>
                )}
                <button
                    className={styles.btCreateChat}
                    onClick={() => {
                        setCurrentChatId(null);
                        setMessages([]);
                        toggleSidebar();
                    }}
                >
                    +
                </button>
            </ul>

            <div className={styles.actionsDiv}>
                <LogoutIcon onClick={() => { goToIndex(); clearCookies(); }} style={{ color: "#ff0000" }} />
                <AppsIcon onClick={goToHome} style={{ color: "#fff" }} />
            </div>
        </aside>
    );
}

export default Sidebar;
