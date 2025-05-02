import React, { useState, useEffect, useRef } from 'react';
import { useNavigate } from 'react-router-dom';
import SendIcon from '@mui/icons-material/Send';
import DeleteIcon from '@mui/icons-material/Delete';
import LogoutIcon from '@mui/icons-material/Logout';
import SettingsIcon from '@mui/icons-material/Settings';
import AppsIcon from '@mui/icons-material/Apps';
import { useAuth } from '../../contexts/Auth';
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
import styles from './Chat.module.css';

function Chat() {
  // State Variables
  const [isSidebarVisible, setIsSidebarVisible] = useState(false);
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [currentChatId, setCurrentChatId] = useState(null);
  const [showSettingsModal, setShowSettingsModal] = useState(false);
  const [isLoadingMessages, setIsLoadingMessages] = useState(false);
  const [isSendingMessage, setIsSendingMessage] = useState(false);
  const [showWelcomePopup, setShowWelcomePopup] = useState(true);

  // Refs and Navigators
  const messagesEndRef = useRef(null);
  const navigate = useNavigate();

  // Auth context
  const { logout } = useAuth();

  // Handlers
  const closeWelcomePopup = () => setShowWelcomePopup(false);
  const handleSettingsClick = () => setShowSettingsModal(true);
  const closeModal = () => setShowSettingsModal(false);
  const goToIndex = () => {
    logout();
    navigate('/');
  };
  const goToHome = () => navigate('/home');
  const toggleSidebar = () => setIsSidebarVisible(!isSidebarVisible);

  const handleOptionSelect = (option) => {
    setShowSettingsModal(false);
  };

  const clearCookies = () => {
    document.cookie.split(";").forEach(c => {
      document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/");
    });
  };

  // Side effect to load chats
  useEffect(() => {
    getChatsByUserId().then((response) => {
      setChats(response.chats || []);
    }).catch((error) => {
      console.error("Erro ao buscar chats:", error);
      setChats([]);
    });
  }, []);

  // Side effect to load messages for the current chat
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

  // Scroll to the bottom of the chat messages
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

      // Create new chat if no current chat
      if (!chatId) {
        const topicResponse = await generateAIResponseTopic(messageContent);
        const newChat = await createChat(topicResponse.response);
        chatId = newChat.id;
        setCurrentChatId(chatId);
        setChats(prevChats => [...prevChats, newChat]);
      }

      const aiCorrectionResponse = await generateAIResponseCorrection(messageContent);
      const combinedMessageContent = `${messageContent}\n\nCorreção: ${aiCorrectionResponse.response}`;

      // Save user message
      const savedMessage = await addMessageToChat(chatId, combinedMessageContent, 'user', 'request');
      setMessages(prevMessages => [...prevMessages, {
        id: savedMessage.id,
        text: savedMessage.content,
        sender: savedMessage.sender || 'user',
        timestamp: savedMessage.timestamp,
        chatId: savedMessage.chat_id,
        type: 'request'
      }]);

      // Get AI responses
      const aiResponseDialog = await generateAIResponseDialog(aiCorrectionResponse.response, chatId);
      const aiTranslation = await generateAIResponseTranslation(aiResponseDialog.response);
      const aiResponseWithTranslation = `${aiResponseDialog.response}\n\n[TRANSLATION]: ${aiTranslation.response}`;

      // Save AI message
      const savedAIMessage = await addMessageToChat(chatId, aiResponseWithTranslation, 'ai', 'response');
      setMessages(prevMessages => [...prevMessages, {
        id: savedAIMessage.id,
        text: savedAIMessage.content,
        sender: savedAIMessage.sender || 'ai',
        timestamp: savedAIMessage.timestamp,
        chatId: savedAIMessage.chat_id,
        type: 'response'
      }]);

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
      {showWelcomePopup && (
        <div className={styles.modalOverlay} onClick={closeWelcomePopup}>
          <div className={styles.modalContent} onClick={(e) => e.stopPropagation()} role="dialog" aria-modal="true">
            <h2 className={styles.modalTitle}>👋 Bem-vindo ao Chat de Prática!</h2>
            <p className={styles.modalDescription}>Aqui você pode praticar o idioma conversando com nossa IA.</p>
            <ul className={styles.modalList}>
              <li>✅ Correção gramatical</li>
              <li>✅ Tradução para o português</li>
              <li>✅ Diálogos simulados para praticar</li>
            </ul>
            <div className={styles.modalExample}>
              <strong>Exemplos:</strong><br /><br />
              <em>Entrada:</em> "I ned a car"<br />
              <em>Saída:</em> "I need a car"<br /><br />
              <em>Entrada:</em> "How are you doin?"<br />
              <em>Saída:</em> "How are you doing?"<br /><br />
              <em>Entrada:</em> "Let's go beach tomorrow?"<br />
              <em>Saída:</em> "Let's go to the beach tomorrow?"<br /><br />
              <em>Entrada:</em> "I don't know how say this."<br />
              <em>Saída:</em> "I don't know how to say this."
            </div>
            <button className={styles.optionButton} onClick={closeWelcomePopup}>Entendi!</button>
          </div>
        </div>
      )}

      <aside className={`${styles.sidebar} ${!isSidebarVisible ? styles.sidebarHidden : ''}`}>
        <button className={`${styles.toggleButton} ${isSidebarVisible ? '' : styles.toggleButtonHidden}`} onClick={toggleSidebar}>&times;</button>
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
                <DeleteIcon onClick={(e) => { e.stopPropagation(); handleDeleteChat(chat.id); }} />
              </li>
            ))
          ) : (
            <ol>Crie um chat!</ol>
          )}
          <button className={styles.btCreateChat} onClick={() => { setCurrentChatId(null); setMessages([]); setIsSidebarVisible(false) }}>+</button>
        </ul>
        <div className={styles.actionsDiv}>
          <LogoutIcon onClick={() => { goToIndex(); clearCookies(); }} style={{ color: "#ff0000" }} />
          <AppsIcon onClick={goToHome} style={{ color: "#fff" }} />
        </div>
      </aside>

      <main className={styles.mainContent}>
        <div className={`${styles.modalOverlay} ${showSettingsModal ? '' : styles.modalOverlayHidden}`} onClick={closeModal}>
          <div className={styles.modalContent} onClick={(e) => e.stopPropagation()}>
            <button className={styles.closeButton} onClick={closeModal}>&times;</button>
            <h3 className={styles.modalTitle}>Escolha o Modelo</h3>
            <button className={styles.optionButton} onClick={() => handleOptionSelect('OpenAI')}>OpenAI</button>
            <button className={styles.optionButton} onClick={() => handleOptionSelect('Gemini')}>Gemini</button>
          </div>
        </div>

        {!isSidebarVisible && (
          <button className={`${styles.toggleButton} ${styles.toggleButtonMainContent}`} onClick={toggleSidebar}>
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
                <div key={message.id} className={`${styles.message} ${message.sender === 'user' ? styles.user : styles.ai}`}>
                  {(message.text || message.content).split(/\n{1,}/).map((line, index, lines) => {
                    if (line.startsWith('[TRANSLATION]: ')) {
                      return (
                        <div key={index} style={{ marginTop: "10px", paddingTop: "10px", borderTop: "1px solid #555", color: "#aaa" }}>
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
            <SettingsIcon className={styles.styledSettingsIcon} onClick={handleSettingsClick} />
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
                onClick={handleSendMessage}
                style={{ color: "#fff", cursor: "pointer" }}
                onMouseEnter={(e) => e.target.style.color = "rgb(187, 187, 187)"}
                onMouseLeave={(e) => e.target.style.color = "#fff"}
              />
            )}
          </div>
        </div>

        <p>Converse com nosso assistente de IA para praticar esse idioma em tempo real, recebendo feedback instantâneo.</p>
      </main>
    </div>
  );
}

export default Chat;
