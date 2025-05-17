import { useNavigate } from 'react-router-dom';
import { useState } from 'react';

function OnBoarding() {
    const navigate = useNavigate();
    const [level, setLevel] = useState('');
    const [language, setLanguage] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();
        // Aqui você pode salvar esses dados no backend, se quiser
        console.log('Level:', level);
        console.log('Language:', language);

        navigate('/home');
    };

    return (
        <form onSubmit={handleSubmit}>
            <h2>Vamos começar!</h2>
            <label>
                Seu nível:
                <select value={level} onChange={(e) => setLevel(e.target.value)} required>
                    <option value="">Selecione...</option>
                    <option value="beginner">Iniciante</option>
                    <option value="intermediate">Intermediário</option>
                    <option value="advanced">Avançado</option>
                </select>
            </label>
            <label>
                Idioma que deseja aprender:
                <input
                    type="text"
                    value={language}
                    onChange={(e) => setLanguage(e.target.value)}
                    placeholder="ex: Inglês, Espanhol..."
                    required
                />
            </label>
            <button type="submit">Continuar</button>
        </form>
    );
}

export default OnBoarding;
