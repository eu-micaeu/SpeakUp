import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styles from './Chat.module.css';
import {
  Add,
  Delete,
  AccountCircle,
  Mic,
  AttachFile,
  Send,
  Dehaze,
} from '@mui/icons-material';
import MeetingRoomIcon from '@mui/icons-material/MeetingRoom';
import TryIcon from '@mui/icons-material/Try';
import AutoStoriesIcon from '@mui/icons-material/AutoStories';
import { removeAuthTokenFromCookies } from '../../utils/cookies';

import {
  getChatsByUserId,
  createChat,
  getMessagesByChatId,
  addMessageToChat,
  generateAIResponseDialog,
  generateAIResponseCorrection,
  generateAIResponseTranslation,
  generateAIResponseTopic
} from '../../utils/api';

export default function Chat() {
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [currentChatId, setCurrentChatId] = useState(null);
  const [isSending, setIsSending] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(true);

  const navigate = useNavigate();

  function goToPerfil() {
    navigate('/perfil');
  }

  function goToIndex() {
    removeAuthTokenFromCookies();
    navigate('/');
  }

  function goToPalavreco() {
    navigate('/palavreco');
  }

  function goToPlanoDeEstudos() {
    navigate('/teaching-plan');
  }

  useEffect(() => {
    getChatsByUserId()
      .then(res => setChats(res.chats || []))
      .catch(() => setChats([]));
  }, []);

  useEffect(() => {
    if (currentChatId) {
      getMessagesByChatId(currentChatId)
        .then(res => {
          const msgs = res.map(m => ({
            id: m.id,
            text: m.content,
            sender: m.sender,
            type: m.type
          }));
          setMessages(msgs);
        })
        .catch(() => setMessages([]));
    }
  }, [currentChatId]);

  const handleSend = async () => {
    if (!inputMessage.trim() || isSending) return;

    setIsSending(true);
    const userInput = inputMessage.trim();
    setInputMessage('');

    try {
      let chatId = currentChatId;

      if (!chatId) {
        const topicRes = await generateAIResponseTopic(userInput);
        const newChat = await createChat(topicRes.response);
        chatId = newChat.id;
        setCurrentChatId(chatId);
        setChats(prev => [...prev, newChat]);
      }

      const correction = await generateAIResponseCorrection(userInput);
      const fullUserMsg = `${userInput}\n\nCorreção: ${correction.response}`;

      const savedUserMsg = await addMessageToChat(chatId, fullUserMsg, 'user', 'request');
      setMessages(prev => [...prev, {
        id: savedUserMsg.id,
        text: fullUserMsg,
        sender: 'user',
        type: 'request'
      }]);

      const dialogRes = await generateAIResponseDialog(correction.response, chatId);
      const translation = await generateAIResponseTranslation(dialogRes.response);
      const fullAIResponse = `${dialogRes.response}\n\n[TRANSLATION]: ${translation.response}`;

      const savedBotMsg = await addMessageToChat(chatId, fullAIResponse, 'ai', 'response');
      setMessages(prev => [...prev, {
        id: savedBotMsg.id,
        text: fullAIResponse,
        sender: 'ai',
        type: 'response'
      }]);

    } catch (err) {
      console.error(err);
      setInputMessage(userInput);
    } finally {
      setIsSending(false);
    }
  };

  return (
    <div className={styles.Chat}>
      <aside className={`${styles.sidebar} ${!sidebarOpen ? styles.sidebarClosed : ''}`}>

        <div>
          <div className={styles.header}>
            <img src="/logo.png" alt="Logo" width={30} />
            <h3>SpeakUp</h3>
          </div>

          <button className={styles.newChat} onClick={() => {
            setCurrentChatId(null);
            setMessages([]);
            setInputMessage('');
            setSidebarOpen(false);
          }}>
            <Add />
          </button>

          <nav>
            {chats.map(chat => (
              <a
                key={chat.id}
                className={styles.navItem}
                href="#"
                onClick={() => setCurrentChatId(chat.id)}
              >
                <span>{chat.topic}</span>
                <Delete />
              </a>
            ))}
          </nav>
        </div>

        <div className={styles.footer}>
          <button onClick={goToPlanoDeEstudos}><AutoStoriesIcon /> Plano de Estudo</button>
          <button onClick={goToPalavreco}><TryIcon /> Palavreco</button>
          <button onClick={goToPerfil}><AccountCircle /> Perfil</button>
          <button onClick={goToIndex}><MeetingRoomIcon /> Sair</button>
        </div>
      </aside>

      <main className={styles.main}>
        <header className={styles.chatHeader}>
          <Dehaze onClick={() => setSidebarOpen(prev => !prev)} className={`${styles.toggleIcon} ${!sidebarOpen ? styles.rotated : ''}`} />
        </header>

        <div className={styles.chatBody}>
          {messages.map(msg => (
            <div
              key={msg.id}
              className={msg.sender === 'user' ? styles.userMessage : styles.botMessage}
            >
              <div>
                <p>{msg.text.split('\n\n')[0]}</p>
                {msg.type === 'request' && (
                  <div className={styles.responseExtras}>
                    <p><strong>Correção:</strong> {msg.text.split('\n\n')[1]?.replace('Correção: ', '')}</p>
                  </div>
                )}
                {msg.type === 'response' && (
                  <div className={styles.responseExtras}>
                    <p><strong>Tradução:</strong> {msg.text.split('[TRANSLATION]: ')[1]}</p>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>

        <footer className={styles.chatFooter}>
          <div className={styles.inputBox}>
            <input
              type="text"
              placeholder="Digite sua mensagem aqui..."
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleSend()}
            />
            <button onClick={handleSend} disabled={isSending}>
              <Send />
            </button>
          </div>
        </footer>
      </main>
    </div>
  );
}
