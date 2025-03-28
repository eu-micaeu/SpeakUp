import { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { getChatsByUserId, createChat, getMessagesByChatId, addMessageToChat, generateAIResponseDialog, generateAIResponseCorrection } from '../../utils/api';
import { useNavigate } from 'react-router-dom';
import SendIcon from '@mui/icons-material/Send';

// Estilos com styled-components
const PageHome = styled.div`
  display: flex;
  height: 100vh;
`;

const Sidebar = styled.aside`
  width: ${props => props.$isVisible ? '300px' : '0'};
  background-color: #000000;
  padding: ${props => props.$isVisible ? '20px' : '0'};
  position: relative;
  transition: margin-left 0.3s ease;
  color: #ffffff;
  overflow: hidden;

  ol {
    padding: 0;
    text-align: center;
  }

  ul {
    margin: 50px 0;
    padding: 0;
  }

  ul li {
    list-style: none;
    margin: 20px 0 0 0;
    text-align: center;
    padding: 10px;
    background-color: #313131;
    border-radius: 10px;
    cursor: pointer;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;

    &:hover {
      background-color: #424242;
    }

    &.active-chat {
      background-color: #4CAF50;
    }
  }
`;

const MainContent = styled.main`
  padding: 20px;
  background-color: #313131;
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  align-items: center;
  color: #fff;

  h1 {
    margin: 25px;
  }
`;

const ToggleButton = styled.button`
  background-color: transparent;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #fff;
  position: fixed;
  top: 20px;
  left: 20px;
  z-index: 1000;
  display: ${props => props.$sidebarVisible ? 'block' : 'none'};

  ${MainContent} & {
    color: #fff;
    position: fixed;
    display: block;
    left: ${props => props.$sidebarVisible ? '320px' : '20px'};
  }
`;

const ChatContainer = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 60%;
`;

const Messages = styled.div`
  flex: 1;
  padding: 10px;
  border-radius: 10px;
  overflow-y: auto;
  overflow-x: hidden;
  max-height: 70vh;

  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-track {
    border-radius: 4px;
  }

  &::-webkit-scrollbar-thumb {
    background: #555;
    border-radius: 4px;
  }

  &::-webkit-scrollbar-thumb:hover {
    background: #777;
  }

`;

const Message = styled.div`
  padding: 20px;
  margin: 8px;
  border-radius: 4px;
  line-height: 2;

  &.user {
    color: #e3f2fd;
    background-color: #1d1d1d;
    margin-left: auto;
    width: 50%;
  }

  &.ai {
    background-color: #242424;
    margin-right: auto;
    width: 75%;
  }

  &.ai.correction {
    background-color: #1eff00;
    color: #000000;
    text-align: center;
    width: 100%;
    padding: 10px;
    margin: 10px 0;
  }

`;

const ChatInput = styled.div`
  display: flex;
  gap: 10px;
  width: 100%;

  input {
    flex: 1;
    padding: 20px;
    width: 100%;
    border: none;
    font-size: 1rem;
    border-radius: 4px;
    background-color: #3d3d3d;
    color: #fff;

    &:focus {
      outline: none;
    }
  }

  button {
    font-size: 0.8rem;
    padding: 20px;
    background-color: #4CAF50;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;

    &:hover {
      background-color: #45a049;
    }
  }
`;

const BtLogout = styled.button`
  border: none;
  border-top: 2px solid #ff0000;
  color: #ff0000;
  background-color: transparent;
  padding: 10px;
  font-size: 20px;
  cursor: pointer;
  font-weight: bold;
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  font-family: 'Karla', sans-serif;

  &:hover {
    background-color: #ff0000;
    color: #000000;
  }
`;

const DivSpeakUp = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: row;
  margin-bottom: 25px;
`;

// Componente Home
function Home() {
  const [isSidebarVisible, setIsSidebarVisible] = useState(false);
  const [chats, setChats] = useState([]);
  const [messages, setMessages] = useState([]);
  const [inputMessage, setInputMessage] = useState('');
  const [currentChatId, setCurrentChatId] = useState(null);
  const messagesEndRef = useRef(null);
  const navigate = useNavigate();

  const goToIndex = () => {
    navigate('/');
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
      const response = await getMessagesByChatId(chatId);
      const formattedMessages = response.map(msg => ({
        id: msg.id,
        text: msg.content,
        sender: msg.sender || 'user',
        timestamp: msg.timestamp,
        chatId: msg.chat_id,
        type: msg.type
      }));
      setMessages(formattedMessages);
    } catch (error) {
      console.error("Erro ao carregar mensagens:", error);
      setMessages([]);
    }
  };

  const handleSendMessage = async () => {
    if (inputMessage.trim() === '') return;

    const messageContent = inputMessage.trim();
    setInputMessage('');
    let tempMessage;
    let tempAIMessage;

    try {
      let chatId = currentChatId;

      if (!chatId) {
        const firstWord = messageContent.split(' ')[0];
        const topic = firstWord.length > 0 ? firstWord : "Novo Chat";

        const newChat = await createChat(topic);
        chatId = newChat.id;
        setCurrentChatId(chatId);
        setChats(prevChats => [...prevChats, newChat]);
      }

      // Mensagem do usuário
      tempMessage = {
        id: Date.now(),
        text: messageContent,
        sender: 'user',
        timestamp: new Date().toISOString(),
        chatId: chatId,
        type: 'request'
      };

      setMessages(prevMessages => [...prevMessages, tempMessage]);

      // Salva a mensagem do usuário no backend
      const savedMessage = await addMessageToChat(chatId, messageContent, 'user', 'request');

      setMessages(prevMessages => [
        ...prevMessages.filter(m => m.id !== tempMessage.id),
        {
          id: savedMessage.id,
          text: savedMessage.content,
          sender: savedMessage.sender || 'user',
          timestamp: savedMessage.timestamp,
          chatId: savedMessage.chat_id,
          type: 'request'
        }
      ]);

      // Resposta da IA - Correction
      const aiCorrectionResponse = await generateAIResponseCorrection(messageContent);
      const correctionMessage = {
        id: Date.now() + 2,
        text: aiCorrectionResponse.response,
        sender: 'ai',
        timestamp: new Date().toISOString(),
        chatId: chatId,
        type: 'correction'
      };
      setMessages(prevMessages => [...prevMessages, correctionMessage]);

      // Salva a correção da IA no backend com sender 'ai'
      const savedCorrectionMessage = await addMessageToChat(chatId, aiCorrectionResponse.response, 'ai', 'correction');
      setMessages(prevMessages => [
        ...prevMessages.filter(m => m.id !== correctionMessage.id),
        {
          id: savedCorrectionMessage.id,
          text: savedCorrectionMessage.content,
          sender: savedCorrectionMessage.sender || 'ai',
          timestamp: savedCorrectionMessage.timestamp,
          chatId: savedCorrectionMessage.chat_id,
          type: 'correction'
        }
      ]);

      // Resposta da IA - Dialog
      const aiResponseDialog = await generateAIResponseDialog(messageContent);

      tempAIMessage = {
        id: Date.now() + 1,
        text: aiResponseDialog.response,
        sender: 'ai',
        timestamp: new Date().toISOString(),
        chatId: chatId,
        type: 'response'
      };

      setMessages(prevMessages => [...prevMessages, tempAIMessage]);

      // Salva a resposta da IA no backend com sender 'ai'
      const savedAIMessage = await addMessageToChat(chatId, aiResponseDialog.response, 'ai', 'response');

      setMessages(prevMessages => [
        ...prevMessages.filter(m => m.id !== tempAIMessage.id),
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
      if (tempMessage) setMessages(prev => prev.filter(m => m.id !== tempMessage.id));
      if (tempAIMessage) setMessages(prev => prev.filter(m => m.id !== tempAIMessage.id));
      setInputMessage(messageContent);
    }
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter') {
      handleSendMessage();
    }
  };

  return (
    <PageHome>
      <Sidebar $isVisible={isSidebarVisible}>
        <ToggleButton onClick={toggleSidebar} $sidebarVisible={isSidebarVisible}>
          &times;
        </ToggleButton>
        <ul>
          {chats?.length > 0 ? (
            chats.map((chat) => (
              <li
                key={chat.id}
                className={chat.id === currentChatId ? 'active-chat' : ''}
                onClick={() => {
                  setCurrentChatId(chat.id);
                  loadChatMessages(chat.id);
                }}
              >
                {chat.topic || "Chat sem título"}
              </li>
            ))
          ) : (
            <ol>Sem chats, por ora</ol>
          )}
        </ul>
        <BtLogout onClick={goToIndex}>Sair</BtLogout>
      </Sidebar>

      <MainContent $sidebarVisible={isSidebarVisible}>
        {!isSidebarVisible && (
          <ToggleButton onClick={toggleSidebar} $sidebarVisible={isSidebarVisible}>
            &#9776;
          </ToggleButton>
        )}
        
        <DivSpeakUp>
          <img src='./logo.png' width={75} alt="SpeakUp Logo" />
          <h1>SpeakUp</h1>
        </DivSpeakUp>
        <ChatContainer>
          <Messages>
            {messages
              .filter(msg => msg.chatId === currentChatId || msg.chat_id === currentChatId)
              .map((message) => (
                <Message key={message.id} className={`${message.sender} ${message.type}`}>
                  {message.text || message.content}
                </Message>
              ))}
            <div ref={messagesEndRef} />
          </Messages>

          <ChatInput>
            <input
              type="text"
              placeholder="Digite sua mensagem..."
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyPress={handleKeyPress}
            />
            <SendIcon
              style={{ cursor: 'pointer', color: '#fff', marginTop: '20px' }}
              onClick={handleSendMessage}
            />
          </ChatInput>
        </ChatContainer>
      </MainContent>
    </PageHome>
  );
}

export default Home;
