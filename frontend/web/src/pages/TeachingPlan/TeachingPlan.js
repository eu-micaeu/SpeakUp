import styles from './TeachingPlan.module.css';
import { getDecodedToken } from '../../utils/cookies';
import { useNavigate } from 'react-router-dom';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';

const plansByLevel = {
    N5: {
        title: 'Plano de Ensino - Japonês N5',
        weeks: [
            'Hiragana e katakana, cumprimentos e apresentações básicas.',
            'Números, dias da semana, expressões de tempo.',
            'Vocabulário do cotidiano, partículas は、が、の.',
            'Verbos básicos no presente, forma -masu, estrutura de frases simples.',
            'Família, perguntas simples com 何、どこ、だれ.',
            'Adjetivos i/na, frases descritivas com です e あります/います.'
        ]
    },
    N4: {
        title: 'Plano de Ensino - Japonês N4',
        weeks: [
            'Kanji básicos (~300), revisão do N5.',
            'Formas verbais: て-form, passado, negativo.',
            'Uso de から、ので、けど、が para conectar frases.',
            'Expressões do cotidiano: transportes, compras, escola.',
            'Verbo + たい (querer), verbo + ことがある (experiência).',
            'Estrutura ～ながら、～たり～たりする, prática de leitura simples.'
        ]
    },
    N3: {
        title: 'Plano de Ensino - Japonês N3',
        weeks: [
            'Kanji intermediários (~650), leitura de textos curtos.',
            'Formas verbais complexas: 〜そう、〜よう、〜らしい.',
            'Expressões causativas e passivas.',
            'Vocabulário relacionado a trabalho, sociedade, cotidiano urbano.',
            'Uso de ～べき、～はず、～かもしれない.',
            'Compreensão de textos de nível intermediário.'
        ]
    },
    N2: {
        title: 'Plano de Ensino - Japonês N2',
        weeks: [
            'Kanji avançados (~1000), leitura de artigos e reportagens.',
            'Expressões idiomáticas e de negócios.',
            'Estudo de gramática avançada: ～に違いない、～わけではない.',
            'Discurso indireto, linguagem polida e formal.',
            'Debates e opinião sobre temas sociais.',
            'Simulados de provas N2 e redação de textos argumentativos.'
        ]
    },
    N1: {
        title: 'Plano de Ensino - Japonês N1',
        weeks: [
            'Kanji complexos (~2000), leitura de textos acadêmicos.',
            'Gramática refinada: nuances e exceções.',
            'Compreensão de textos técnicos, políticos e científicos.',
            'Uso avançado da linguagem formal e informal.',
            'Discussão de temas abstratos e culturais.',
            'Preparação intensiva para o exame N1 com simulados e redação formal.'
        ]
    },
    A1: {
        title: 'Plano de Ensino - Inglês A1',
        weeks: [
            'Cumprimentos, apresentações, alfabeto.',
            'Números, cores, dias da semana.',
            'Objetos do dia a dia, artigos definidos/indefinidos.',
            'Família, possessivos, perguntas simples.',
            'Descrições físicas, adjetivos básicos.',
            'Verbos básicos no presente, frases simples.'
        ]
    },
    A2: {
        title: 'Plano de Ensino - Inglês A2',
        weeks: [
            'Rotina diária, present simple, advérbios de frequência.',
            'Lugares na cidade, direções, there is/are.',
            'Vocabulário de compras, preços, demonstrativos (this/that/these/those).',
            'Comida e bebidas, substantivos contáveis/incontáveis, some/any.',
            'Hobbies, like/love/hate + verbo-ing, atividades de lazer.',
            'Planos de viagem, futuro com "going to", vocabulário de viagem.'
        ]
    },
    B1: {
        title: 'Plano de Ensino - Inglês B1',
        weeks: [
            'Revisão de tempos verbais, conectores.',
            'Expressar opiniões, argumentos simples.',
            'Vocabulário sobre trabalho e carreira.',
            'Relatar experiências passadas (present perfect).',
            'Fazer comparações, expressões idiomáticas.',
            'Simulações de conversas do dia a dia.'
        ]
    },
    B2: {
        title: 'Plano de Ensino - Inglês B2',
        weeks: [
            'Debates sobre temas atuais, vocabulário avançado.',
            'Expressões idiomáticas e gírias.',
            'Escrita de textos formais e informais.',
            'Discussão de filmes/livros, vocabulário cultural.',
            'Preparação para exames de proficiência.',
            'Simulações de entrevistas de emprego.'
        ]
    },
    C1: {
        title: 'Plano de Ensino - Inglês C1',
        weeks: [
            'Discussão de temas complexos, vocabulário técnico.',
            'Análise crítica de textos e artigos.',
            'Produção escrita avançada, ensaios acadêmicos.',
            'Debates sobre ética e moralidade.',
            'Preparação para apresentações orais.',
            'Simulações de situações profissionais.'
        ]
    },
    C2: {
        title: 'Plano de Ensino - Inglês C2',
        weeks: [
            'Discussão de temas filosóficos e abstratos.',
            'Produção escrita criativa, contos e poesias.',
            'Análise de obras literárias e cinematográficas.',
            'Debates sobre política e sociedade.',
            'Preparação para exames de proficiência avançada.',
            'Simulações de negociações internacionais.'
        ]
    }
};

function TeachingPlan() {
    const navigate = useNavigate();
    const user = getDecodedToken();
    const levelPlan = plansByLevel[user.level];

    return (
        <div className={styles.pageTeachingPlan}>
            <ArrowBackIcon
                className={styles.backArrow}
                onClick={() => navigate(-1)}
            />

            {levelPlan ? (
                <div className={styles.teachingPlan}>
                    <h1 className={styles.title}>{levelPlan.title}</h1>

                    <h2 className={styles.sectionTitle}>📅 Plano Semanal (12 Semanas)</h2>
                    <ul className={styles.planList}>
                        {levelPlan.weeks.map((content, index) => (
                            <li key={index}>
                                <strong>Semana {index * 2 + 1}-{index * 2 + 2}:</strong> {content}
                            </li>
                        ))}
                    </ul>
                </div>
            ) : (
                <div className={styles.teachingPlan}>
                    <h1 className={styles.title}>Acesso Negado</h1>
                    <p>Você não tem permissão para acessar esta página.</p>
                </div>
            )}
        </div>
    );
}

export default TeachingPlan;
