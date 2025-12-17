import { useTranslation } from 'react-i18next';

const Config = () => {
    const { t } = useTranslation();
    return (
        <div style={{ padding: '20px' }}>
            <h1>{t('config.title')}</h1>
            <p>{t('config.content')}</p>
        </div>
    );
};

export default Config;
