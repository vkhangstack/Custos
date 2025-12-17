import { useTranslation } from 'react-i18next';

const Log = () => {
    const { t } = useTranslation();
    return (
        <div style={{ padding: '20px' }}>
            <h1>{t('log.title')}</h1>
            <p>{t('log.content')}</p>
        </div>
    );
};

export default Log;
