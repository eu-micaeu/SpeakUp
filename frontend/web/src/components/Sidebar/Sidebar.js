import DeleteIcon from '@mui/icons-material/Delete';
import AppsIcon from '@mui/icons-material/Apps';
import CloseIcon from '@mui/icons-material/Close';
import styles from './Sidebar.module.css';

function Sidebar({
    isVisible,
    toggleSidebar,
    chats,
    currentChatId,
    setCurrentChatId,
    loadChatMessages,
    handleDeleteChat,
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
                    <>
                    </>
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
                <AppsIcon onClick={goToHome} style={{ color: "#fff" }} />
            </div>
        </aside>
    );
}

export default Sidebar;
