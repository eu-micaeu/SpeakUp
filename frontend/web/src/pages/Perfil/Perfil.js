import styles from './Perfil.module.css';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { useState, useEffect } from 'react';
import { toast, ToastContainer } from 'react-toastify';
import { useNavigate } from 'react-router-dom';

// Api functions
import {
    updateUser,
    deleteUser
} from '../../utils/api';

// Utils
import { getDecodedToken, isAuthTokenValid, removeAuthTokenFromCookies } from '../../utils/cookies';

function Perfil() {
    const [profileData, setProfileData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [levels, setLevels] = useState([]);
    const [editedData, setEditedData] = useState({
        name: '',
        email: '',
        language: '',
        level: ''
    });
    const navigate = useNavigate();

    useEffect(() => {
        const fetchProfileData = async () => {
            try {

                const data = getDecodedToken();

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

            } catch (error) {
                console.error('Erro ao buscar dados do perfil:', error);
            } finally {
                setLoading(false);
            }
        };

        if (isAuthTokenValid) {
            fetchProfileData();
        }
    }, [isAuthTokenValid]);

    const goToIndex = () => {
        removeAuthTokenFromCookies();
        navigate('/');
    };

    const handleInputChange = (field, value) => {
        setEditedData(prev => ({
            ...prev,
            [field]: value
        }));

        if (field === 'language') {
            if (value === 'english') {
                setLevels(['A1', 'A2', 'B1', 'B2', 'C1', 'C2']);
            }
        }
    };

    const handleSave = async () => {
        try {

            const userId = getDecodedToken().user_id;

            await updateUser(userId, editedData);

            toast.success('Perfil atualizado com sucesso!');

            setProfileData(prev => ({
                ...prev,
                name: `${editedData.name}`,
                email: editedData.email,
                language: editedData.language,
                level: editedData.level
            }));

        } catch (error) {
            toast.error('Falha ao atualizar o perfil.');
        } finally {
            setLoading(false);
        }
    };

    const renderContent = () => {
        if (loading) {
            return <p>Carregando dados do perfil...</p>;
        }

        if (!profileData) {
            return <p>Nenhum dado de perfil encontrado</p>;
        }

        return (
            <div className={styles.profileContainer}>
                <div className={styles.profileHeader}>
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
                        <button
                            className={styles.deleteButton}
                            onClick={() => {
                                if (window.confirm('Tem certeza que deseja excluir sua conta?')) {
                                    const userId = getDecodedToken().user_id;
                                    deleteUser(userId)
                                        .then(() => {
                                            toast.success('Conta excluída com sucesso!');
                                            goToIndex();
                                        })
                                        .catch((error) => {
                                            console.error('Erro ao excluir conta:', error);
                                            toast.error('Falha ao excluir conta.');
                                        });
                                }
                            }}
                        >
                            Excluir Conta
                        </button>
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
