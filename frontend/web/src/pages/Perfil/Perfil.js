import styles from './Perfil.module.css';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import PersonIcon from '@mui/icons-material/Person';
import { useState, useEffect } from 'react';
import { updateUser } from '../../utils/api';
import { toast, ToastContainer } from 'react-toastify';
import { getDecodedToken, isAuthTokenValid } from '../../utils/cookies';

function Perfil(){
    const [profileData, setProfileData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [levels, setLevels] = useState([]);
    const [editedData, setEditedData] = useState({
        name: '',
        email: '',
        language: '',
        level: ''
    });

    useEffect(() => {
        const fetchProfileData = async () => {
            try {

                console.log(getDecodedToken);

                const data = getDecodedToken();

                console.log('Dados do token decodificado:', data);
                
                setProfileData({
                    name: data.name,
                    email: data.email,
                    language: data.language,
                    level: data.level,
                    stats: {
                        conversations: 0,
                        messages: 0,
                        minutesPracticed: 0
                    }
                });

                setEditedData({
                    name: data.name,
                    email: data.email,
                    language: data.language,
                    level: data.level
                });

                if (data.language === 'english') {
                    setLevels(['A1', 'A2', 'B1', 'B2', 'C1', 'C2']);
                }

                setError(null);
            } catch (error) {
                console.error('Erro ao buscar dados do perfil:', error);
                setError("Falha ao carregar dados do perfil");
            } finally {
                setLoading(false);
            }
        };

        if (isAuthTokenValid) {
            fetchProfileData();
        }
    }, [isAuthTokenValid]);

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
            
            const userId = getDecodedToken().user.id;

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
