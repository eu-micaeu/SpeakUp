import styles from './Footer.module.css';

export function Footer() {
    return (
        <footer className={styles.footer}>
            <div className={styles.container}>

                <div className={styles.column}>
                    <h4>Links úteis</h4>
                    <ul>
                        <li><a href="/login">Login</a></li>
                        <li><a href="#features">Funcionalidades</a></li>
                        <li><a href="#">Sobre</a></li>
                    </ul>
                </div>

                <div className={styles.column}>
                    <h4>Contato</h4>
                    <p>Email: contato@speakup.com</p>
                    <p>Suporte: suporte@speakup.com</p>
                </div>
            </div>

            <div className={styles.bottomBar}>
                <p>&copy; {new Date().getFullYear()} SpeakUp. Todos os direitos reservados.</p>
            </div>
        </footer>
    );
}
