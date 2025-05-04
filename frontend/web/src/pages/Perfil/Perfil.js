import styles from './Perfil.module.css';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { useState, useEffect } from 'react';
import { getUserById } from '../../utils/api';
import { useAuth } from '../../contexts/Auth';
import Cookies from 'js-cookie';

function Perfil(){
    const [profileData, setProfileData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const { user } = useAuth();
useEffect(() => {
    const fetchProfileData = async () => {
            try {
                if (!user) {
                    setLoading(false);
                    return;
                }

                // Usa prioritariamente o ID do cookie, já que é extraído do token JWT
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

setProfileData({
    name: `${data.user.first_name} ${data.user.last_name}`,
    email: data.user.email,
    language: data.user.language,
    level: data.user.level,
    stats: {
        conversations: 0,
        messages: 0,
        minutesPracticed: 0
    }
});
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

    const renderContent = () => {
        if (loading) {
            return <p>Carregando dados do perfil...</p>;
        }

        if (error) {
            return <p>Erro: {error}</p>;
        }

        if (!profileData) {
            return <p>Nenhum dado de perfil encontrado</p>;
        }

        return (
            <div className={styles.profileContainer}>
                <div className={styles.profileHeader}>
                    <div className={styles.avatarContainer}>
                        <img src="/avatar-placeholder.png" alt="Foto do perfil" className={styles.avatar} />
                    </div>
                    <div className={styles.userNames}>
                        <h1 className={styles.userName}>{profileData?.name}</h1>
                        <p className={styles.userEmail}>{profileData?.email}</p>
                    </div>
                </div>

                <div className={styles.profileInfo}>
                    <div className={styles.infoSection}>
                        <h2>Informações Pessoais</h2>
                        <div className={styles.infoField}>
                            <label>Nome Completo</label>
                            <input type="text" disabled value={profileData?.name} />
                        </div>
                        <div className={styles.infoField}>
                            <label>E-mail</label>
                            <input type="email" disabled value={profileData?.email} />
                        </div>
                        <div className={styles.infoField}>
                            <label>Idioma Principal</label>
                            <input type="text" disabled value={profileData?.language} />
                        </div>
                        <div className={styles.infoField}>
                            <label>Nivel</label>
                            <input type="text" disabled value={profileData?.level} />
                        </div>
                    </div>

                    <div className={styles.infoSection}>
                        <h2>Estatísticas de Uso</h2>
                        <div className={styles.statsGrid}>
                            <div className={styles.statItem}>
                                <span className={styles.statValue}>{profileData.stats?.conversations || 0}</span>
                                <span className={styles.statLabel}>Conversas</span>
                            </div>
                            <div className={styles.statItem}>
                                <span className={styles.statValue}>{profileData.stats?.messages || 0}</span>
                                <span className={styles.statLabel}>Mensagens</span>
                            </div>
                            <div className={styles.statItem}>
                                <span className={styles.statValue}>{profileData.stats?.minutesPracticed || 0}</span>
                                <span className={styles.statLabel}>Minutos Praticados</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    };

    return (
        <div className={styles.pagePerfil}>
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