import styles from './Perfil.module.css';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import PersonIcon from '@mui/icons-material/Person';
import { useState, useEffect } from 'react';
import { getUserById, updateUser } from '../../utils/api';
import { useAuth } from '../../contexts/Auth';
import Cookies from 'js-cookie';
import { toast, ToastContainer } from 'react-toastify';

function Perfil(){
    const [profileData, setProfileData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    // const [successMessage, setSuccessMessage] = useState(''); // Remove o estado de mensagem de sucesso
    const [levels, setLevels] = useState([]);
    const { user } = useAuth();
    const [editedData, setEditedData] = useState({
        name: '',
        email: '',
        language: '',
        level: ''
    });

    useEffect(() => {
        const fetchProfileData = async () => {
            try {
                if (!user) {
                    setLoading(false);
                    return;
                }

                const userId = Cookies.get('userId');
                
                if (!userId) {
                    setError("Usuário não encontrado");
                    setLoading(false);
                    return;
                }
                const data = await getUserById(userId);

                if (!data) {
                    setError("Falha ao carregar dados do perfil");
                    setLoading(false);
                    return;
                }

                const fullName = `${data.user.name}`;
                
                setProfileData({
                    name: fullName,
                    email: data.user.email,
                    language: data.user.language,
                    level: data.user.level,
                    stats: {
                        conversations: 0,
                        messages: 0,
                        minutesPracticed: 0
                    }
                });

                setEditedData({
                    name: data.user.name,
                    email: data.user.email,
                    language: data.user.language,
                    level: data.user.level
                });

                if (data.user.language === 'english') {
                    setLevels(['A1', 'A2', 'B1', 'B2', 'C1', 'C2']);
                } else if (data.user.language === 'japanese') {
                    setLevels(['N5', 'N4', 'N3', 'N2', 'N1']);
                }

                setError(null);
            } catch (error) {
                console.error('Erro ao buscar dados do perfil:', error);
                setError("Falha ao carregar dados do perfil");
            } finally {
                setLoading(false);
            }
        };

        if (Cookies.get('userId')) {
            fetchProfileData();
        }
    }, [user]);

    const handleInputChange = (field, value) => {
        setEditedData(prev => ({
            ...prev,
            [field]: value
        }));
        
        if (field === 'language') {
            if (value === 'english') {
                setLevels(['A1', 'A2', 'B1', 'B2', 'C1', 'C2']);
            } else if (value === 'japanese') {
                setLevels(['N5', 'N4', 'N3', 'N2', 'N1']);
            } else {
                setLevels([]);
            }
        }
    };

    const handleSave = async () => {
        try {
            setLoading(true);
            setError(null);
            // setSuccessMessage(''); // Remove a limpeza da mensagem de sucesso
            
            const userId = Cookies.get('userId');
            if (!userId) {
                throw new Error('ID do usuário não encontrado');
            }

            await updateUser(userId, editedData);
            
            toast.success('Perfil atualizado com sucesso!'); // Adiciona a notificação de toast
            
            setProfileData(prev => ({
                ...prev,
                name: `${editedData.name}`,
                email: editedData.email,
                language: editedData.language,
                level: editedData.level
            }));

        } catch (error) {
            console.error('Erro ao atualizar perfil:', error);
            setError('Falha ao atualizar o perfil. Por favor, tente novamente.');
            toast.error('Falha ao atualizar o perfil.'); // Opcional: Adiciona toast de erro
        } finally {
            setLoading(false);
        }
    };

    const renderContent = () => {
        if (loading) {
            return <p>Carregando dados do perfil...</p>;
        }

        if (error) {
            // Mantém a exibição de erro no componente, mas o toast também pode ser usado
            return <p>Erro: {error}</p>; 
        }

        if (!profileData) {
            return <p>Nenhum dado de perfil encontrado</p>;
        }

        return (
            <div className={styles.profileContainer}>
                <div className={styles.profileHeader}>
                    <div className={styles.avatarContainer}>
                        <PersonIcon style={{ width: '100%', height: '100%', color: '#fff' }} />
                    </div>
                    <div className={styles.userNames}>
                        <h1 className={styles.userName}>{editedData.name}</h1>
                        <p className={styles.userEmail}>{editedData.email}</p>
                    </div>
                </div>

                {/* Remove a exibição da mensagem de sucesso anterior
                {successMessage && (
                    <div className={styles.successMessage}>
                        {successMessage}
                    </div>
                )} */}

                <div className={styles.profileInfo}>
                    <div className={styles.infoSection}>
                        <h2>Informações Pessoais</h2>
                        <div className={styles.infoField}>
                            <label>Nome</label>
                            <input
                                type="text"
                                value={editedData.name}
                                onChange={(e) => handleInputChange('name', e.target.value)}
                            />
                        </div>
                        <div className={styles.infoField}>
                            <label>E-mail</label>
                            <input
                                type="email"
                                value={editedData.email}
                                onChange={(e) => handleInputChange('email', e.target.value)}
                            />
                        </div>
                        <div className={styles.infoField}>
                            <label>Idioma Principal</label>
                            <select
                                className={styles.select}
                                value={editedData.language}
                                onChange={(e) => handleInputChange('language', e.target.value)}
                            >
                                <option value="" disabled>Selecione</option>
                                <option value="english">Inglês</option>
                                {/* <option value="japanese">Japonês</option> */}
                            </select>
                        </div>
                        <div className={styles.infoField}>
                            <label>Nível</label>
                            <select
                                className={styles.select}
                                value={editedData.level}
                                onChange={(e) => handleInputChange('level', e.target.value)}
                                disabled={!editedData.language}
                            >
                                <option value="" disabled>Selecione</option>
                                {levels.length > 0 ? (
                                    levels.map((level) => (
                                        <option key={level} value={level}>
                                            {level}
                                        </option>
                                    ))
                                ) : (
                                    editedData.language === 'english' ? (
                                        ['A1', 'A2', 'B1', 'B2', 'C1', 'C2'].map((level) => (
                                            <option key={level} value={level}>
                                                {level}
                                            </option>
                                        ))
                                    ) : editedData.language === 'japanese' ? (
                                        ['N5', 'N4', 'N3', 'N2', 'N1'].map((level) => (
                                            <option key={level} value={level}>
                                                {level}
                                            </option>
                                        ))
                                    ) : null
                                )}
                            </select>
                        </div>
                        <div className={styles.infoField}>
                            <button 
                                className={styles.saveButton} 
                                onClick={handleSave}
                                disabled={loading}
                            >
                                {loading ? 'Salvando...' : 'Salvar'}
                            </button>
                        </div>
                    </div>
                    
                </div>
            </div>
        );
    };

    return (
        <div className={styles.pagePerfil}>
            <ToastContainer />
            <ArrowBackIcon
                className={styles.arrowBack}
                onClick={() => window.history.back()}
                style={{ cursor: 'pointer' }}
            />
            {renderContent()}
        </div>
    );
}

export default Perfil;
