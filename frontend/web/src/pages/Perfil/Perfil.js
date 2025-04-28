import styles from './Perfil.module.css';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';

function Perfil(){
    return (
        <div className={styles.pagePerfil}>
            <ArrowBackIcon
                className={styles.arrowBack}
                onClick={() => window.history.back()}
                style={{ cursor: 'pointer' }}
            />
            
            <div className={styles.profileContainer}>
                <div className={styles.profileHeader}>
                    <div className={styles.avatarContainer}>
                        <img src="/avatar-placeholder.png" alt="Foto do perfil" className={styles.avatar} />
                    </div>
                    <div className={styles.userNames}>
                        <h1 className={styles.userName}>Manoel Gomes</h1>
                        <p className={styles.userEmail}>canetaazul@gmail.com</p>
                    </div>
                    
                </div>

                <div className={styles.profileInfo}>
                    <div className={styles.infoSection}>
                        <h2>Informações Pessoais</h2>
                        <div className={styles.infoField}>
                            <label>Nome Completo</label>
                            <input type="text" disabled value="Manoel Gomes" />
                        </div>
                        <div className={styles.infoField}>
                            <label>E-mail</label>
                            <input type="email" disabled value="canetaazul@gmail.com" />
                        </div>
                        <div className={styles.infoField}>
                            <label>Idioma Principal</label>
                            <input type="text" disabled value="Ingles" />
                        </div>
                        <div className={styles.infoField}>
                            <label>Nivel</label>
                            <input type="text" disabled value="B1" />
                        </div>
                    </div>

                    <div className={styles.infoSection}>
                        <h2>Estatísticas de Uso</h2>
                        <div className={styles.statsGrid}>
                            <div className={styles.statItem}>
                                <span className={styles.statValue}>0</span>
                                <span className={styles.statLabel}>Conversas</span>
                            </div>
                            <div className={styles.statItem}>
                                <span className={styles.statValue}>0</span>
                                <span className={styles.statLabel}>Mensagens</span>
                            </div>
                            <div className={styles.statItem}>
                                <span className={styles.statValue}>0</span>
                                <span className={styles.statLabel}>Minutos Praticados</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Perfil;