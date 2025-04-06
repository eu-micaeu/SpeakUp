import React, { useState, useEffect, useRef } from 'react';
import {
  getChatsByUserId,
  createChat,
  deleteChat,
  getMessagesByChatId,
  addMessageToChat,
  generateAIResponseDialog,
  generateAIResponseCorrection,
  generateAIResponseTranslation,
  generateAIResponseTopic
} from '../../utils/api';
import { useNavigate } from 'react-router-dom';
import SendIcon from '@mui/icons-material/Send';
import DeleteIcon from '@mui/icons-material/Delete';
import LogoutIcon from '@mui/icons-material/Logout';
import PersonIcon from '@mui/icons-material/Person';
import SettingsIcon from '@mui/icons-material/Settings';
import AppsIcon from '@mui/icons-material/Apps';
import { useAuth } from '../../contexts/Auth';
import styles from './Chat.module.css';

function Chat() {
  const [isSidebarVisible, setIsSidebarVisible] = useState(false);
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [currentChatId, setCurrentChatId] = useState(null);
  const [showSettingsModal, setShowSettingsModal] = useState(false);
  const [isLoadingMessages, setIsLoadingMessages] = useState(false);
  const [isSendingMessage, setIsSendingMessage] = useState(false);
  const { logout } = useAuth();
  const messagesEndRef = useRef(null);
  const navigate = useNavigate();

  const handleSettingsClick = () => {
    setShowSettingsModal(true);
  };

  const handleOptionSelect = (option) => {
    setShowSettingsModal(false);
  };

  const closeModal = () => {
    setShowSettingsModal(false);
  };

  const goToIndex = () => {
    logout();
    navigate('/');
  };

  const goToHome = () => {
    navigate('/home');
  };

  const clearCookies = () => {
    document.cookie.split(";").forEach((c) => {
      document.cookie = c
        .replace(/^ +/, "")
        .replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/");
    });
  };

  const toggleSidebar = () => {
    setIsSidebarVisible(!isSidebarVisible);
  };

  useEffect(() => {
    getChatsByUserId().then((response) => {
      setChats(response.chats || []);
    }).catch(error => {
      console.error("Erro ao buscar chats:", error);
      setChats([]);
    });
  }, []);

  useEffect(() => {
    if (currentChatId) {
      loadChatMessages(currentChatId);
    } else {
      setMessages([]);
    }
  }, [currentChatId]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  const loadChatMessages = async (chatId) => {
    try {
      setIsLoadingMessages(true);
      const response = await getMessagesByChatId(chatId);
      const formattedMessages = response.map(msg => ({
        id: msg.id,
        text: msg.content,
        sender: msg.sender,
        created_id: msg.created_at,
        chat_id: msg.chat_id,
        type: msg.type
      }));
      setMessages(formattedMessages);
    } catch (error) {
      console.error("Erro ao carregar mensagens:", error);
      setMessages([]);
    } finally {
      setIsLoadingMessages(false);
    }
  };

  const handleSendMessage = async () => {
    if (inputMessage.trim() === '' || isSendingMessage) return;

    setIsSendingMessage(true);

    const messageContent = inputMessage.trim();
    setInputMessage('');

    try {
      let chatId = currentChatId;

      if (!chatId) {
        const topicResponse = await generateAIResponseTopic(messageContent);
        const newChat = await createChat(topicResponse.response);
        chatId = newChat.id;
        setCurrentChatId(chatId);
        setChats(prevChats => [...prevChats, newChat]);
      }

      const aiCorrectionResponse = await generateAIResponseCorrection(messageContent);
      const combinedMessageContent = `${messageContent}\n\nCorreção: ${aiCorrectionResponse.response}`;

      const savedMessage = await addMessageToChat(chatId, combinedMessageContent, 'user', 'request');

      setMessages(prevMessages => [
        ...prevMessages,
        {
          id: savedMessage.id,
          text: savedMessage.content,
          sender: savedMessage.sender || 'user',
          timestamp: savedMessage.timestamp,
          chatId: savedMessage.chat_id,
          type: 'request'
        }
      ]);

      const aiResponseDialog = await generateAIResponseDialog(aiCorrectionResponse.response);
      const aiTranslation = await generateAIResponseTranslation(aiResponseDialog.response);

      const aiResponseWithTranslation = `${aiResponseDialog.response}\n\n[TRANSLATION]: ${aiTranslation.response}`;

      const savedAIMessage = await addMessageToChat(chatId, aiResponseWithTranslation, 'ai', 'response');

      setMessages(prevMessages => [
        ...prevMessages,
        {
          id: savedAIMessage.id,
          text: savedAIMessage.content,
          sender: savedAIMessage.sender || 'ai',
          timestamp: savedAIMessage.timestamp,
          chatId: savedAIMessage.chat_id,
          type: 'response'
        }
      ]);

    } catch (error) {
      console.error("Erro ao processar mensagem:", error);
      setInputMessage(messageContent);
    } finally {
      setIsSendingMessage(false);
    }
  };

  const handleDeleteChat = async (chatId) => {
    try {
      await deleteChat(chatId);
      setChats(prevChats => prevChats.filter(chat => chat.id !== chatId));
      if (currentChatId === chatId) {
        setCurrentChatId(null);
        setMessages([]);
        loadChatMessages(null);
      }
    } catch (error) {
      console.error("Erro ao excluir chat:", error);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleSendMessage();
    }
  };

  return (
    <div className={styles.pageHome}>
      <PersonIcon
        style={{ color: "#fff", cursor: "pointer", position: "absolute", top: "20px", right: "20px", fontSize: "30px" }}
        onMouseEnter={(e) => e.target.style.color = "rgb(187, 187, 187)"}
        onMouseLeave={(e) => e.target.style.color = "#fff"}
      />

      <aside className={`${styles.sidebar} ${!isSidebarVisible ? styles.sidebarHidden : ''}`}>
        <button
          className={`${styles.toggleButton} ${isSidebarVisible ? '' : styles.toggleButtonHidden}`}
          onClick={toggleSidebar}
        >
          &times;
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
                />
              </li>
            ))
          ) : (
            <ol>Crie um chat!</ol>
          )}
          <button
            className={styles.btCreateChat}
            onClick={() => { setCurrentChatId(null); setMessages([]); setIsSidebarVisible(false) }}
          >
            +
          </button>
        </ul>

        <div className={styles.actionsDiv}>
          <LogoutIcon
            onClick={() => { goToIndex(); clearCookies(); }}
            style={{ color: "#ff0000", cursor: "pointer", margin: "20px" }}
            onMouseEnter={(e) => e.target.style.color = "#ffffff"}
            onMouseLeave={(e) => e.target.style.color = "#ff0000"}
          />
          <AppsIcon
            onClick={goToHome}
            style={{ color: "#fff", cursor: "pointer", margin: "20px" }}
            onMouseEnter={(e) => e.target.style.color = "rgb(194, 194, 194)"}
            onMouseLeave={(e) => e.target.style.color = "#fff"}
          />
        </div>
      </aside>

      <main className={styles.mainContent}>
        <div
          className={`${styles.modalOverlay} ${showSettingsModal ? '' : styles.modalOverlayHidden}`}
          onClick={closeModal}
        >
          <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <button className={styles.closeButton} onClick={closeModal}>&times;</button>
            <h3 className={styles.modalTitle}>Escolha o Modelo</h3>
            <button className={styles.optionButton} onClick={() => handleOptionSelect('OpenAI')}>
              OpenAI
            </button>
            <button className={styles.optionButton} onClick={() => handleOptionSelect('Gemini')}>
              Gemini
            </button>
          </div>
        </div>

        {!isSidebarVisible && (
          <button
            className={`${styles.toggleButton} ${styles.toggleButtonMainContent} ${isSidebarVisible ? styles.toggleButtonMainContentSidebarVisible : ''}`}
            onClick={toggleSidebar}
          >
            &#9776;
          </button>
        )}

        <div className={styles.divSpeakUp}>
          <img src='./logo.png' width={50} alt="SpeakUp Logo" />
          <h2>SpeakUp</h2>
        </div>
        <div className={styles.chatContainer}>
          <div className={styles.messages}>
            {messages
              .filter(msg => msg.chatId === currentChatId || msg.chat_id === currentChatId)
              .map((message) => (
                <div
                  key={message.id}
                  className={`${styles.message} ${message.sender === 'user' ? styles.user : styles.ai} ${message.type === 'correction' ? styles.correction : ''}`}
                >
                  {(message.text || message.content).split(/\n{1,}/).map((line, index, lines) => {
                    if (line.startsWith('[TRANSLATION]: ')) {
                      return (
                        <div key={index} style={{
                          marginTop: "10px",
                          paddingTop: "10px",
                          borderTop: "1px solid #555",
                          color: "#aaa",
                        }}>
                          <strong>Tradução:</strong> {line.replace('[TRANSLATION]: ', '')}
                        </div>
                      );
                    }

                    if (index === 1 && line.startsWith('Correção: ')) {
                      return (
                        <React.Fragment key={index}>
                          <hr />
                          <span style={{ color: '#1eff00' }}>{line.replace('Correção: ', '')}</span>
                        </React.Fragment>
                      );
                    }

                    return (
                      <React.Fragment key={index}>
                        {line}
                        {index < lines.length - 1 && <br />}
                      </React.Fragment>
                    );
                  })}
                </div>
              ))}
            <div ref={messagesEndRef} />
          </div>

          <div className={styles.chatInput}>
            <SettingsIcon
              className={styles.styledSettingsIcon}
              onClick={handleSettingsClick}
            />

            <input
              type="text"
              placeholder="Digite sua mensagem..."
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={handleKeyPress}
              disabled={isSendingMessage}
            />

            {isSendingMessage ? (
              <div className={styles.loadingSpinner} style={{ margin: '0 10px' }} />
            ) : (
              <SendIcon
                onKeyDown={handleKeyPress}
                onClick={handleSendMessage}
                style={{ color: "#fff", cursor: "pointer" }}
                onMouseEnter={(e) => e.target.style.color = "rgb(187, 187, 187)"}
                onMouseLeave={(e) => e.target.style.color = "#fff"}
              />
            )}
          </div>
        </div>

        <p>O SpeakUp ainda está em Beta, poderá ter erros nas respostas e/ou correções.</p>
      </main>
    </div>
  );
}

export default Chat;